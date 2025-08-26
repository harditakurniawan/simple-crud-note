package main

import (
	"fmt"
	"log"
	"simple-crud-notes/configs"
	databases "simple-crud-notes/databases/pgsql"
	"simple-crud-notes/databases/redis"
	"simple-crud-notes/handlers"
	"simple-crud-notes/middleware"
	"simple-crud-notes/routes"
	"simple-crud-notes/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	configs.LoadEnv(".env")

	appConfig := configs.LoadAppConfig()
	appName := appConfig.APP_NAME
	apiPrefix := appConfig.APP_PREFIX
	port := appConfig.APP_PORT

	tokenDuration := 24 * time.Hour
	jwtService, err := utils.NewJWTService(
		appConfig.JWT_PRIVATE_KEY_PATH,
		appConfig.JWT_PUBLIC_KEY_PATH,
		tokenDuration,
	)

	if err != nil {
		log.Fatal("Failed to initialize JWT service: ", err)
	}

	handlers.SetJWTService(jwtService)

	databases.DatabaseInit(appConfig)
	redis.RedisInit(appConfig)

	if !fiber.IsChild() {
		log.Println("Initializing database connection in master process")

		databases.Migration()
	}

	/*
	* Prefork mode adalah fitur yang memungkinkan aplikasi Fiber untuk berjalan di beberapa proses worker.
	* Ini meningkatkan performa aplikasi dengan memanfaatkan multi-core CPU.
	 */
	app := fiber.New(fiber.Config{
		AppName:           appName,
		BodyLimit:         4 * 1024 * 1024,
		EnablePrintRoutes: true,
		ServerHeader:      appName,
		Prefork:           true,
	})

	api := app.Group(fmt.Sprintf("/%v", apiPrefix))

	api.Get("/health-checks", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "OK",
		})
	})

	app.Use(middleware.Recovery())
	routes.InitRoutesV1(api, jwtService)

	log.Printf("App is running on port %s", port)

	log.Fatal(app.Listen(fmt.Sprintf(":%v", port)))
}

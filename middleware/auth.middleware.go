package middleware

import (
	databases "simple-crud-notes/databases/pgsql"
	"simple-crud-notes/databases/pgsql/entities"
	"simple-crud-notes/utils"
	"simple-crud-notes/utils/enum"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func AuthenticateRequest(jwtService utils.JWTService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse(
				fiber.StatusUnauthorized,
				c.OriginalURL(),
				"Authorization header is required",
			))
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse(
				fiber.StatusUnauthorized,
				c.OriginalURL(),
				"Authorization header must start with 'Bearer '",
			))
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwtService.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse(
				fiber.StatusUnauthorized,
				c.OriginalURL(),
				"Invalid or expired token",
			))
		}

		// var (
		// 	accessToken entities.AccessToken
		// 	user        entities.User
		// )

		// result := databases.DB.Preload("User").Where("token = ?", tokenString).First(&accessToken)
		// if result.Error != nil || result.RowsAffected == 0 {
		// 	return c.Status(fiber.StatusUnauthorized).JSON(
		// 		utils.ErrorResponse(
		// 			fiber.StatusUnauthorized,
		// 			c.OriginalURL(),
		// 			"Expired token",
		// 		),
		// 	)
		// }
		// user = accessToken.User

		// fmt.Printf("AccessToken: %+v\n", user)

		userInfo, err := jwtService.GetUserInfoFromToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse(
				fiber.StatusUnauthorized,
				c.OriginalURL(),
				"Invalid token claims",
			))
		}

		isTokenExists := CheckExistingToken(userInfo.UserID, tokenString)

		if !isTokenExists {
			return c.Status(fiber.StatusUnauthorized).JSON(
				utils.ErrorResponse(
					fiber.StatusUnauthorized,
					c.OriginalURL(),
					"Token expired",
				),
			)
		}

		c.Locals("userInfo", userInfo)

		return c.Next()
	}
}

func CheckExistingToken(userId string, token string) bool {
	redisKeyPattern := enum.KEY_TOKEN + "_" + userId + "_"
	key := redisKeyPattern + token

	if utils.GetCache(key, nil) == nil {
		var accessToken entities.AccessToken
		result := databases.DB.First(&accessToken, "token = ?", token)

		if result.Error != nil || result.RowsAffected == 0 {
			return false
		}

		utils.SetCache(key, true, 1*time.Hour)
		return true
	}

	return true
}

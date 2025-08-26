package routes

import (
	"simple-crud-notes/handlers"
	"simple-crud-notes/middleware"
	"simple-crud-notes/utils"

	"github.com/gofiber/fiber/v2"
)

func InitRoutesV1(router fiber.Router, jwtService utils.JWTService) {
	v1 := router.Group("/v1")

	v1.Use(middleware.Logging())

	// REGISTRATION & AUTH
	v1.Post("/registrations",
		middleware.ValidateRequest((*utils.RegistrationDto)(nil)),
		handlers.Registration,
	)
	v1.Post("/sign-in",
		middleware.ValidateRequest((*utils.SignInDto)(nil)),
		handlers.SignIn,
	)

	protectedRoute := v1.Group("/", middleware.AuthenticateRequest(jwtService))
	protectedRoute.Post("/sign-out", handlers.SignOut)

	// PROFILE
	protectedRoute.Get("/profiles", handlers.Profile)

	// NOTES
	protectedRoute.Post("/notes",
		middleware.ValidateRequest((*utils.CreateNoteDto)(nil)),
		handlers.CreateNote,
	)
	protectedRoute.Get("/notes",
		middleware.WithPagination(),
		handlers.GetNotes,
	)
	protectedRoute.Get("/notes/:id",
		handlers.DetailNote,
	)
	protectedRoute.Patch("/notes/:id",
		middleware.ValidateRequest((*utils.UpdateNoteDto)(nil)),
		handlers.UpdateNote,
	)
	protectedRoute.Delete("/notes/:id",
		handlers.DeleteNote,
	)
}

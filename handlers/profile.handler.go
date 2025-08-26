package handlers

import (
	"errors"
	databases "simple-crud-notes/databases/pgsql"
	"simple-crud-notes/databases/pgsql/entities"
	"simple-crud-notes/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Profile(c *fiber.Ctx) error {
	userInfo, ok := c.Locals("userInfo").(*utils.UserInfo)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(
			utils.ErrorResponse(
				fiber.StatusInternalServerError,
				c.OriginalURL(),
				"User info not found in context",
			),
		)
	}

	var user entities.User
	result := databases.DB.Select("id", "name", "email", "created_at").First(&user, "id = ?", userInfo.UserID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(
				utils.ErrorResponse(
					fiber.StatusNotFound,
					c.OriginalURL(),
					"User not found",
				),
			)
		}

		return c.Status(fiber.StatusInternalServerError).JSON(
			utils.ErrorResponse(
				fiber.StatusInternalServerError,
				c.OriginalURL(),
				"Failed to fetch user profile",
			),
		)
	}

	responseData := fiber.Map{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"created_at": user.CreatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(
		utils.SuccessResponse(
			fiber.StatusOK,
			c.OriginalURL(),
			responseData,
		),
	)
}

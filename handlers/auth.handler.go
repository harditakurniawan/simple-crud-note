package handlers

import (
	"errors"
	"fmt"
	databases "simple-crud-notes/databases/pgsql"
	"simple-crud-notes/databases/pgsql/entities"
	"simple-crud-notes/utils"
	"simple-crud-notes/utils/enum"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var jwtService utils.JWTService

func SetJWTService(service utils.JWTService) {
	jwtService = service
}

func Registration(c *fiber.Ctx) error {
	registrationRequest, ok := c.Locals("validatedDTO").(*utils.RegistrationDto)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(
			fiber.StatusInternalServerError,
			c.OriginalURL(),
			"Failed to get validated request",
		))
	}

	var existingUser entities.User
	if err := databases.DB.Where("email = ?", registrationRequest.Email).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			utils.ErrorResponse(
				fiber.StatusBadRequest,
				c.OriginalURL(),
				"Email already exists",
			),
		)
	}

	newUser := entities.User{
		Name:     registrationRequest.Name,
		Email:    registrationRequest.Email,
		Password: registrationRequest.Password,
	}

	result := databases.DB.Create(&newUser)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			utils.ErrorResponse(
				fiber.StatusInternalServerError,
				c.OriginalURL(),
				result.Error.Error(),
			),
		)
	}

	registrationResponse := map[string]interface{}{
		"id":         newUser.ID,
		"name":       newUser.Name,
		"email":      newUser.Email,
		"created_at": newUser.CreatedAt,
		"updated_at": newUser.UpdatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(
		utils.SuccessResponse(
			fiber.StatusCreated,
			c.OriginalURL(),
			registrationResponse,
		),
	)
}

func SignIn(c *fiber.Ctx) error {
	signInRequest, ok := c.Locals("validatedDTO").(*utils.SignInDto)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(
			fiber.StatusInternalServerError,
			c.OriginalURL(),
			"Failed to get validated request",
		))
	}

	var user entities.User
	result := databases.DB.Where("email = ?", signInRequest.Email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(
				utils.ErrorResponse(
					fiber.StatusBadRequest,
					c.OriginalURL(),
					"Email not found",
				),
			)
		}

		return c.Status(fiber.StatusInternalServerError).JSON(
			utils.ErrorResponse(
				fiber.StatusInternalServerError,
				c.OriginalURL(),
				"Database error",
			),
		)
	}

	if !utils.CheckPasswordHash(signInRequest.Password, user.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(
			utils.ErrorResponse(
				fiber.StatusBadRequest,
				c.OriginalURL(),
				"Invalid password",
			),
		)
	}

	token, err := jwtService.GenerateToken(fmt.Sprintf("%d", user.ID), user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			utils.ErrorResponse(
				fiber.StatusInternalServerError,
				c.OriginalURL(),
				"Failed to generate token",
			),
		)
	}

	accessToken := entities.AccessToken{
		UserID: user.ID,
		Token:  token,
	}

	upsertResult := databases.DB.
		Where(entities.AccessToken{UserID: user.ID}).
		Assign(entities.AccessToken{
			Token: token,
		}).
		FirstOrCreate(&accessToken)

	if upsertResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			utils.ErrorResponse(
				fiber.StatusInternalServerError,
				c.OriginalURL(),
				"Failed to save access token",
			),
		)
	}

	redisKeyPattern := enum.KEY_TOKEN + "_" + fmt.Sprintf("%d", user.ID) + "_"
	utils.DeleteTokensByPattern(redisKeyPattern + "*")
	utils.SetCache(redisKeyPattern+token, true, 1*time.Hour)

	responseData := fiber.Map{
		"access_token": token,
	}

	return c.Status(fiber.StatusOK).JSON(
		utils.SuccessResponse(
			fiber.StatusOK,
			c.OriginalURL(),
			responseData,
		),
	)
}

func SignOut(c *fiber.Ctx) error {
	userInfo, ok := c.Locals("userInfo").(*utils.UserInfo)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(
			utils.ErrorResponse(
				fiber.StatusInternalServerError,
				c.OriginalURL(),
				enum.USER_NOT_FOUND_IN_CONTEXT,
			),
		)
	}

	result := databases.DB.Delete(&entities.AccessToken{}, "user_id = ? ", userInfo.UserID)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			utils.ErrorResponse(
				fiber.StatusInternalServerError,
				c.OriginalURL(),
				"Failed to revoke access token",
			),
		)
	}

	authHeader := c.Get("Authorization")
	parts := strings.Split(authHeader, "Bearer ")

	var token string
	if len(parts) == 2 {
		token = strings.TrimSpace(parts[1])
	} else {
		token = ""
	}

	redisKeyPattern := enum.KEY_TOKEN + "_" + userInfo.UserID + "_"
	utils.DeleteCache(redisKeyPattern + token)

	responseData := fiber.Map{
		"message": "sign out success",
	}

	return c.Status(fiber.StatusOK).JSON(
		utils.SuccessResponse(
			fiber.StatusOK,
			c.OriginalURL(),
			responseData,
		),
	)
}

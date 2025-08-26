package middleware

import (
	"fmt"
	"reflect"
	"simple-crud-notes/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func init() {
	validate.RegisterValidation("ContainSpecialChar", utils.CheckContainSpecialChar)
}

func ValidateRequest(dtoType interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		dtoValue := reflect.New(reflect.TypeOf(dtoType).Elem())
		request := dtoValue.Interface()

		if err := c.BodyParser(request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(
				fiber.StatusBadRequest,
				c.OriginalURL(),
				fmt.Sprintf("Failed to parse request body: %v", err.Error()),
			))
		}

		if err := validate.Struct(request); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			errorDetails := make([]fiber.Map, 0, len(validateErrs))
			for _, e := range validateErrs {
				errorDetails = append(errorDetails, fiber.Map{
					"field": e.Field(),
					"error": fmt.Sprintf("field %s: wanted %s %s, got `%v`", e.Field(), e.Tag(), e.Param(), e.Value()),
				})
			}

			return c.Status(fiber.StatusBadRequest).JSON(
				utils.ErrorResponse(
					fiber.StatusBadRequest,
					c.OriginalURL(),
					errorDetails,
				),
			)
		}

		c.Locals("validatedDTO", request)
		return c.Next()
	}
}

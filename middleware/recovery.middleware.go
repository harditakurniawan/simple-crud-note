package middleware

import (
	"fmt"
	"runtime/debug"
	"simple-crud-notes/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Recovery() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}

				stack := debug.Stack()

				utils.LogAsync(utils.LogEntry{
					Timestamp:       time.Now(),
					StatusCode:      fiber.StatusInternalServerError,
					Method:          c.Method(),
					Endpoint:        c.OriginalURL(),
					RequestHeader:   string(c.Request().Header.Header()),
					RequestBody:     string(c.Body()),
					RequestParams:   string(c.Request().URI().QueryString()),
					ResponseMessage: err.Error() + "\n" + string(stack),
					ProcessTime:     time.Since(start).String(),
				})

				c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(
					fiber.StatusInternalServerError,
					c.OriginalURL(),
					"Internal Server Error",
				))
			}
		}()

		return c.Next()
	}
}

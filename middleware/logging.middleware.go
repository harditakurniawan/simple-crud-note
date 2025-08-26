package middleware

import (
	"net/http"
	"simple-crud-notes/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Logging() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		body := c.Body()

		err := c.Next()

		utils.LogAsync(utils.LogEntry{
			Timestamp:       time.Now(),
			StatusCode:      c.Response().StatusCode(),
			Method:          c.Method(),
			Endpoint:        c.OriginalURL(),
			RequestHeader:   string(c.Request().Header.Header()),
			RequestBody:     string(body),
			RequestParams:   string(c.Request().URI().QueryString()),
			ResponseMessage: http.StatusText(c.Response().StatusCode()),
			ProcessTime:     time.Since(start).String(),
		})

		return err
	}
}

func WithLoggingAndRecovery(handler fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		body := c.Body()

		defer func() {
			if r := recover(); r != nil {
				panic(r)
			}

			utils.LogAsync(utils.LogEntry{
				Timestamp:       time.Now(),
				StatusCode:      c.Response().StatusCode(),
				Method:          c.Method(),
				Endpoint:        c.OriginalURL(),
				RequestHeader:   string(c.Request().Header.Header()),
				RequestBody:     string(body),
				RequestParams:   string(c.Request().URI().QueryString()),
				ResponseMessage: http.StatusText(c.Response().StatusCode()),
				ProcessTime:     time.Since(start).String(),
			})
		}()

		return handler(c)
	}
}

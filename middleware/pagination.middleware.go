package middleware

import (
	"fmt"
	"simple-crud-notes/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func WithPagination() fiber.Handler {
	return func(c *fiber.Ctx) error {
		pageSize := 10
		pageNumber := 1
		orderBy := c.Query("order_by", "id")
		sortType := c.Query("sort_type", "DESC")
		search := c.Query("search", "")

		if ps := c.Query("page_size"); ps != "" {
			if val, err := strconv.Atoi(ps); err == nil && val > 0 {
				pageSize = val
			} else {
				return c.Status(fiber.StatusBadRequest).JSON(
					utils.ErrorResponse(
						fiber.StatusBadRequest,
						c.OriginalURL(),
						"Invalid page_size",
					),
				)
			}
		}

		if pn := c.Query("page_number"); pn != "" {
			if val, err := strconv.Atoi(pn); err == nil && val > 0 {
				pageNumber = val
			} else {
				return c.Status(fiber.StatusBadRequest).JSON(
					utils.ErrorResponse(
						fiber.StatusBadRequest,
						c.OriginalURL(),
						"Invalid page_number",
					),
				)
			}
		}

		if sortType != "ASC" && sortType != "DESC" {
			return c.Status(fiber.StatusBadRequest).JSON(
				utils.ErrorResponse(
					fiber.StatusBadRequest,
					c.OriginalURL(),
					"sort_type must be ASC or DESC",
				),
			)
		}

		if len(search) > 100 {
			return c.Status(fiber.StatusBadRequest).JSON(
				utils.ErrorResponse(
					fiber.StatusBadRequest,
					c.OriginalURL(),
					"search query too long",
				),
			)
		}

		offset := (pageNumber - 1) * pageSize

		c.Locals("order", fmt.Sprintf("%s %s", orderBy, sortType))
		c.Locals("offset", offset)
		c.Locals("limit", pageSize)
		c.Locals("pageNumber", pageNumber)
		c.Locals("search", search)
		return c.Next()
	}
}

package middleware

import (
	"fmt"
	"simple-crud-notes/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PaginationParams struct {
	PageSize   int
	PageNumber int
	OrderBy    string
	SortType   string
	Search     string
}

func WithPagination() fiber.Handler {
	return func(c *fiber.Ctx) error {
		params, err := extractAndValidatePaginationParams(c)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				utils.ErrorResponse(
					fiber.StatusBadRequest,
					c.OriginalURL(),
					err.Error(),
				),
			)
		}

		offset := (params.PageNumber - 1) * params.PageSize

		c.Locals("order", fmt.Sprintf("%s %s", params.OrderBy, params.SortType))
		c.Locals("offset", offset)
		c.Locals("limit", params.PageSize)
		c.Locals("pageNumber", params.PageNumber)
		c.Locals("search", params.Search)

		return c.Next()
	}
}

func extractAndValidatePaginationParams(c *fiber.Ctx) (*PaginationParams, error) {
	params := &PaginationParams{
		PageSize:   10,
		PageNumber: 1,
		OrderBy:    c.Query("order_by", "id"),
		SortType:   c.Query("sort_type", "DESC"),
		Search:     c.Query("search", ""),
	}

	// Validate page_size
	if err := validatePageSize(c, params); err != nil {
		return nil, err
	}

	// Validate page_number
	if err := validatePageNumber(c, params); err != nil {
		return nil, err
	}

	// Validate sort_type
	if err := validateSortType(params); err != nil {
		return nil, err
	}

	// Validate search
	if err := validateSearch(params); err != nil {
		return nil, err
	}

	return params, nil
}

func validatePageSize(c *fiber.Ctx, params *PaginationParams) error {
	if ps := c.Query("page_size"); ps != "" {
		if val, err := strconv.Atoi(ps); err != nil || val <= 0 {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid page_size")
		} else {
			params.PageSize = val
		}
	}
	return nil
}

func validatePageNumber(c *fiber.Ctx, params *PaginationParams) error {
	if pn := c.Query("page_number"); pn != "" {
		if val, err := strconv.Atoi(pn); err != nil || val <= 0 {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid page_number")
		} else {
			params.PageNumber = val
		}
	}
	return nil
}

func validateSortType(params *PaginationParams) error {
	if params.SortType != "ASC" && params.SortType != "DESC" {
		return fiber.NewError(fiber.StatusBadRequest, "sort_type must be ASC or DESC")
	}
	return nil
}

func validateSearch(params *PaginationParams) error {
	if len(params.Search) > 100 {
		return fiber.NewError(fiber.StatusBadRequest, "search query too long")
	}
	return nil
}

package utils

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type PaginationMeta struct {
	LAST_PAGE int
	PER_PAGE  int
	PAGE      int
	TOTAL     int64
}

type MetaResponse struct {
	STATUS_CODE uint   `json:"status_code"`
	STATUS      string `json:"status"`
	MESSAGE     any    `json:"message"`
	URL         string `json:"url"`
	LAST_PAGE   int    `json:"last_page"`
	PER_PAGE    int    `json:"per_page"`
	PAGE        int    `json:"page"`
	TOTAL       int64  `json:"total"`
}

func PaginationResponse(statusCode int, url string, data any, meta PaginationMeta) fiber.Map {
	return fiber.Map{
		"meta": MetaResponse{
			STATUS_CODE: uint(statusCode),
			STATUS:      http.StatusText(statusCode),
			MESSAGE:     nil,
			URL:         url,
			LAST_PAGE:   meta.LAST_PAGE,
			PER_PAGE:    meta.PER_PAGE,
			PAGE:        meta.PAGE,
			TOTAL:       meta.TOTAL,
		},
		"data": data,
	}
}

func SuccessResponse(statusCode int, url string, data any) fiber.Map {
	return fiber.Map{
		"meta": MetaResponse{
			STATUS_CODE: uint(statusCode),
			STATUS:      http.StatusText(statusCode),
			MESSAGE:     nil,
			URL:         url,
			LAST_PAGE:   0,
			PER_PAGE:    0,
			PAGE:        0,
			TOTAL:       0,
		},
		"data": data,
	}
}

func ErrorResponse(statusCode int, url string, message any) fiber.Map {
	var finalMessage []fiber.Map

	switch msg := message.(type) {
	case []fiber.Map:
		finalMessage = msg
	default:
		finalMessage = []fiber.Map{
			{
				"field": "general",
				"error": message,
			},
		}
	}

	return fiber.Map{
		"meta": MetaResponse{
			STATUS_CODE: uint(statusCode),
			STATUS:      http.StatusText(statusCode),
			MESSAGE:     finalMessage,
			URL:         url,
			LAST_PAGE:   0,
			PER_PAGE:    0,
			PAGE:        0,
			TOTAL:       0,
		},
		"data": nil,
	}
}

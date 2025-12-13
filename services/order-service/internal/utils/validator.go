package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

// ValidateStruct validate struct với validator tags
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// ErrorResponse định dạng error response
func ErrorResponse(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"error": message,
	})
}

// SuccessResponse định dạng success response
func SuccessResponse(c *fiber.Ctx, data interface{}) error {
	return c.JSON(fiber.Map{
		"data": data,
	})
}

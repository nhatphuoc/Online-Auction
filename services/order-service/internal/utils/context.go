package utils

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

var (
	ErrMissingUserID = errors.New("user ID not found in context")
)

// GetUserIDFromContext extracts user ID from fiber context (set by auth middleware)
func GetUserIDFromContext(c *fiber.Ctx) (int64, error) {
	userID := c.Locals("user_id")
	if userID == nil {
		return 0, ErrMissingUserID
	}
	
	// Try to convert to int64
	switch v := userID.(type) {
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case string:
		// Try to parse string to int64
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, ErrMissingUserID
		}
		return id, nil
	default:
		return 0, ErrMissingUserID
	}
}

// GetUserEmailFromContext extracts user email from fiber context
func GetUserEmailFromContext(c *fiber.Ctx) (string, error) {
	email := c.Locals("email")
	if email == nil {
		return "", errors.New("email not found in context")
	}
	
	if emailStr, ok := email.(string); ok {
		return emailStr, nil
	}
	
	return "", errors.New("invalid email format")
}

// GetUserRoleFromContext extracts user role from fiber context
func GetUserRoleFromContext(c *fiber.Ctx) (string, error) {
	role := c.Locals("role")
	if role == nil {
		return "", errors.New("role not found in context")
	}
	
	if roleStr, ok := role.(string); ok {
		return roleStr, nil
	}
	
	return "", errors.New("invalid role format")
}

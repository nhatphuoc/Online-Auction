package middleware

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware extracts user info from headers (set by API Gateway after authentication)
// API Gateway already validated JWT token and forwards user info via headers
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user info from headers (set by API Gateway)
		userIDStr := c.Get("X-User-ID")
		if userIDStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing user information",
			})
		}

		// Parse user ID
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid user ID format",
			})
		}

		// Get additional user info from headers
		email := c.Get("X-User-Email")
		role := c.Get("X-User-Role")

		// Set user info in context for handlers to use
		c.Locals("user_id", userID)
		c.Locals("email", email)
		c.Locals("role", role)

		return c.Next()
	}
}

// OptionalAuthMiddleware extracts user info if present but doesn't require it
func OptionalAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userIDStr := c.Get("X-User-ID")
		if userIDStr != "" {
			userID, err := strconv.ParseInt(userIDStr, 10, 64)
			if err == nil {
				c.Locals("user_id", userID)
				c.Locals("email", c.Get("X-User-Email"))
				c.Locals("role", c.Get("X-User-Role"))
			}
		}
		return c.Next()
	}
}

// RoleMiddleware checks if user has required role
func RoleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		userRole := role.(string)
		for _, allowedRole := range allowedRoles {
			if strings.EqualFold(userRole, allowedRole) {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Insufficient permissions",
		})
	}
}


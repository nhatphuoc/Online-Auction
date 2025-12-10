package middleware

import (
	"api_gateway/internal/config"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type ValidateResponse struct {
	Valid  bool   `json:"valid"`
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

// AuthMiddleware validates JWT token via auth service
func AuthMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from X-User-Token header
		token := c.Get("X-User-Token")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing X-User-Token header",
			})
		}

		// Call auth service to validate token
		req, err := http.NewRequest("POST", cfg.AuthServiceURL+"/auth/validate-jwt", nil)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create validation request",
			})
		}

		// Set headers for auth service
		req.Header.Set("X-User-Token", token)
		req.Header.Set("X-Api-Gateway", cfg.APIGatewaySecret)
		req.Header.Set("X-Auth-Internal-Service", cfg.AuthInternalSecret)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to validate token",
			})
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Invalid token",
				"details": string(body),
			})
		}

		// Parse response
		var validateResp ValidateResponse
		if err := json.NewDecoder(resp.Body).Decode(&validateResp); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to parse validation response",
			})
		}

		if !validateResp.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token is not valid",
			})
		}

		// Store user info in context
		c.Locals("userID", validateResp.UserID)
		c.Locals("email", validateResp.Email)
		c.Locals("role", validateResp.Role)
		c.Locals("token", token)

		return c.Next()
	}
}

// ProxyMiddleware adds required headers when proxying to internal services
func ProxyMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user info from context (set by AuthMiddleware)
		userID, _ := c.Locals("userID").(string)
		email, _ := c.Locals("email").(string)
		role, _ := c.Locals("role").(string)
		token, _ := c.Locals("token").(string)

		// Set headers for internal services
		c.Request().Header.Set("X-User-ID", userID)
		c.Request().Header.Set("X-User-Email", email)
		c.Request().Header.Set("X-User-Role", role)
		c.Request().Header.Set("X-User-Token", token)
		c.Request().Header.Set("X-Api-Gateway", cfg.APIGatewaySecret)
		c.Request().Header.Set("X-Auth-Internal-Service", cfg.AuthInternalSecret)

		return c.Next()
	}
}

// AdminMiddleware checks if user has admin role
func AdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(string)
		if !ok || role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Admin access required",
			})
		}
		return c.Next()
	}
}

// SellerMiddleware checks if user has seller role
func SellerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(string)
		if !ok || (role != "seller" && role != "admin") {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Seller access required",
			})
		}
		return c.Next()
	}
}

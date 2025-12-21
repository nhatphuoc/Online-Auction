package middleware

import (
	"category_service/internal/config"
	"category_service/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware kiểm tra JWT token
func AuthMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// authHeader := c.Get("Authorization")
		// if authHeader == "" {
		// 	return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Token không được cung cấp")
		// }

		// // Lấy token từ "Bearer <token>"
		// parts := strings.Split(authHeader, " ")
		// if len(parts) != 2 || parts[0] != "Bearer" {
		// 	return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Format token không đúng")
		// }

		// token := parts[1]
		// claims, err := utils.ValidateToken(token, cfg.JWTSecret)
		// if err != nil {
		// 	return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Token không hợp lệ")
		// }

		// // Lưu thông tin user vào context
		// c.Locals("userID", claims.UserID)
		// c.Locals("email", claims.Email)
		// c.Locals("role", claims.Role)

		c.Locals("userID", "test-user-id")

		return c.Next()
	}
}

// AdminMiddleware kiểm tra quyền admin
func AdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role").(string)
		if role != "admin" {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "Bạn không có quyền truy cập")
		}
		return c.Next()
	}
}

// ExtractUserInfo middleware: lấy thông tin user từ header và xác nhận X-Internal-JWT
func ExtractUserInfo() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Get("X-User-ID")
		email := c.Get("X-User-Email")
		role := c.Get("X-User-Role")
		internalJWT := c.Get("X-Internal-JWT")

		if internalJWT == "" {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Missing X-Internal-JWT header")
		}

		c.Locals("userID", userID)
		c.Locals("email", email)
		c.Locals("role", role)
		c.Locals("internalJWT", internalJWT)
		return c.Next()
	}
}

// RequireAdminRole middleware: chỉ cho phép ROLE_ADMIN
func RequireAdminRole() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role != "ROLE_ADMIN" {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "Admin role required")
		}
		return c.Next()
	}
}

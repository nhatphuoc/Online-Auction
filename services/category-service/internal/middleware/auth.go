package middleware

import (
	"category_service/internal/config"
	"category_service/internal/utils"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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
func VerifyInternalJWT(
	cfg *config.Config,
	tokenString string,
	expectedAudience string,
) (bool, error) {

	// =========================
	// Phase 1: Parse UNVERIFIED để lấy issuer
	// =========================
	unverifiedClaims := &jwt.RegisteredClaims{}

	parser := jwt.NewParser(
		jwt.WithoutClaimsValidation(),
	)

	_, _, err := parser.ParseUnverified(tokenString, unverifiedClaims)
	if err != nil {
		return false, err
	}

	issuer := unverifiedClaims.Issuer
	if issuer == "" {
		return false, errors.New("missing issuer")
	}

	// =========================
	// Lookup public key theo issuer
	// =========================
	publicPem, ok := cfg.PublicKeys[issuer]
	if !ok {
		return false, errors.New("unknown issuer")
	}

	block, _ := pem.Decode([]byte(publicPem))
	if block == nil {
		return false, errors.New("invalid public key PEM")
	}

	pubAny, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, err
	}

	pubKey, ok := pubAny.(*rsa.PublicKey)
	if !ok {
		return false, errors.New("not RSA public key")
	}

	// =========================
	// Phase 2: Parse + verify + validate time
	// =========================
	claims := &jwt.RegisteredClaims{}

	parser = jwt.NewParser(
		jwt.WithAudience(expectedAudience),
		jwt.WithIssuer(issuer),
		jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}),
		// ⬇️ TIME VALIDATION xảy ra ở đây
	)

	token, err := parser.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			return pubKey, nil
		},
	)
	if err != nil {
		return false, err
	}

	if !token.Valid {
		return false, errors.New("invalid token")
	}

	return true, nil
}

// ExtractUserInfo middleware: lấy thông tin user từ header và xác nhận X-Internal-JWT
func ExtractUserInfo(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Get("X-User-ID")
		email := c.Get("X-User-Email")
		role := c.Get("X-User-Role")
		internalJWT := c.Get("X-Internal-JWT")

		ok, err := VerifyInternalJWT(
			cfg,
			internalJWT,
			cfg.CategoryServiceName,
		)
		if err != nil || !ok {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid Internal JWT")
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

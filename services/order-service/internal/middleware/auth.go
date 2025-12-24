package middleware

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"order_service/internal/config"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware extracts user info from headers (set by API Gateway after authentication)
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userIDStr := c.Get("X-User-ID")
		if userIDStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing user information",
			})
		}

		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid user ID format",
			})
		}

		email := c.Get("X-User-Email")
		role := c.Get("X-User-Role")

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

// AdminMiddleware checks if user is admin
func AdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role == nil || role != "ROLE_ADMIN" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Admin access required",
			})
		}
		return c.Next()
	}
}

// VerifyInternalJWT verifies internal JWT from API Gateway
func VerifyInternalJWT(
	cfg *config.Config,
	tokenString string,
	expectedAudience string,
) (bool, error) {

	fmt.Println("Verifying Internal JWT...")
	fmt.Println(expectedAudience)
	// Parse UNVERIFIED to get issuer
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

	// Lookup public key by issuer
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

	// Parse + verify + validate time
	claims := &jwt.RegisteredClaims{}

	parser = jwt.NewParser(
		jwt.WithAudience(expectedAudience),
		jwt.WithIssuer(issuer),
		jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}),
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

// ExtractUserInfo middleware: gets user info from headers and verifies X-Internal-JWT
func ExtractUserInfo(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Get("X-User-ID")
		email := c.Get("X-User-Email")
		role := c.Get("X-User-Role")
		internalJWT := c.Get("X-Internal-JWT")

		ok, err := VerifyInternalJWT(
			cfg,
			internalJWT,
			cfg.OrderServiceName,
		)
		if err != nil || !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid Internal JWT",
			})
		}

		c.Locals("userID", userID)
		c.Locals("email", email)
		c.Locals("role", role)
		c.Locals("internalJWT", internalJWT)
		return c.Next()
	}
}

// RequireAdminRole middleware: requires ROLE_ADMIN
func RequireAdminRole() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role != "ROLE_ADMIN" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Admin role required",
			})
		}
		return c.Next()
	}
}

package middleware

import (
	"api_gateway/internal/config"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type ValidateResponse struct {
	Valid  bool   `json:"valid"`
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

// AuthMiddleware validates JWT token directly in API Gateway
func AuthMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("X-User-Token")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing X-User-Token header",
			})
		}

		// Lấy public key từ config
		pubKeyPEM := cfg.JWTPublicKeyAuthService
		if pubKeyPEM == "" {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Missing JWT public key in config",
			})
		}
		pubKey, err := parseRSAPublicKeyFromPEM([]byte(pubKeyPEM))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Invalid public key",
			})
		}

		type CustomClaims struct {
			Role  interface{} `json:"role"`
			Email string      `json:"email"`
			Type  string      `json:"type"`
			jwt.RegisteredClaims
		}

		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return pubKey, nil
		})
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		claims, ok := token.Claims.(*CustomClaims)
		if !ok || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token claims",
			})
		}

		if claims.ExpiresAt == nil || claims.ExpiresAt.Time.Before(time.Now()) || claims.Type != "access" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token expired or invalid type",
			})
		}

		userID := claims.Subject
		email := claims.Email
		role := ""
		switch v := claims.Role.(type) {
		case string:
			role = v
		case []interface{}:
			if len(v) > 0 {
				role = v[0].(string)
			}
		}

		c.Locals("userID", userID)
		c.Locals("email", email)
		c.Locals("role", role)
		c.Locals("token", tokenString)

		return c.Next()
	}
}

// parseRSAPublicKeyFromPEM parses PEM encoded PKCS1 or PKCS8 public key
func parseRSAPublicKeyFromPEM(pubPEM []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pubPEM)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the public key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		return nil, errors.New("not RSA public key")
	}
}

// ProxyMiddleware adds required headers when proxying to internal services
func ProxyMiddleware(cfg *config.Config, serviceName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user info from context (set by AuthMiddleware)
		userID, _ := c.Locals("userID").(string)
		email, _ := c.Locals("email").(string)
		role, _ := c.Locals("role").(string)
		token, _ := c.Locals("token").(string)

		// Tạo JWT nội bộ ký bằng private key của API Gateway
		internalJWT, err := generateInternalJWT(cfg, serviceName)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to sign internal JWT"})
		}

		// Set headers for internal services
		c.Request().Header.Set("X-User-ID", userID)
		c.Request().Header.Set("X-User-Email", email)
		c.Request().Header.Set("X-User-Role", role)
		c.Request().Header.Set("X-User-Token", token)
		c.Request().Header.Set("X-Api-Gateway", cfg.APIGatewaySecret)
		c.Request().Header.Set("X-Auth-Internal-Service", cfg.AuthInternalSecret)
		c.Request().Header.Set("X-Internal-JWT", internalJWT)

		return c.Next()
	}
}

// generateInternalJWT tạo JWT ký bằng private key của API Gateway
func generateInternalJWT(cfg *config.Config, aud string) (string, error) {
	privPem := cfg.JWTPrivateKey
	if privPem == "" {
		return "", errors.New("missing JWT_PRIVATE_KEY in config")
	}
	block, _ := pem.Decode([]byte(privPem))
	if block == nil {
		return "", errors.New("failed to parse PEM block for private key")
	}
	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	claims := jwt.RegisteredClaims{
		Issuer:    "api-gateway",
		Audience:  []string{aud},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privKey)
}

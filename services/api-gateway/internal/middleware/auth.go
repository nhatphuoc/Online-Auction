package middleware

import (
	"api_gateway/internal/config"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log/slog"
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
			return c.Next()
		}

		token, _, err := new(jwt.Parser).ParseUnverified(
			tokenString,
			jwt.MapClaims{},
		)
		if err != nil {
			slog.Warn("Invalid token format received",
				slog.String("error", err.Error()),
				slog.String("path", c.Path()),
				slog.String("ip", c.IP()),
			)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token format",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			slog.Warn("Invalid token claims",
				slog.String("path", c.Path()),
				slog.String("ip", c.IP()),
			)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token claims",
			})
		}

		// Kiểm tra type == "access"
		if t, ok := claims["type"].(string); !ok || t != "access" {
			slog.Warn("Token is not access token",
				slog.String("token_type", t),
				slog.String("path", c.Path()),
			)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token is not access token",
			})
		}

		// Lấy userId (subject), email, role
		userID := ""
		if sub, ok := claims["sub"].(string); ok {
			userID = sub
		} else if subf, ok := claims["sub"].(float64); ok {
			userID = fmt.Sprintf("%.0f", subf)
		}
		email, _ := claims["email"].(string)
		role := ""
		if r, ok := claims["role"].(string); ok {
			role = r
		} else if r, ok := claims["role"].(map[string]interface{}); ok {
			// Trường hợp role là object (enum)
			if name, ok := r["name"].(string); ok {
				role = name
			}
		}

		c.Locals("userID", userID)
		c.Locals("email", email)
		c.Locals("role", role)
		c.Locals("token", tokenString)

		slog.Debug("User authenticated successfully",
			slog.String("user_id", userID),
			slog.String("email", email),
			slog.String("role", role),
			slog.String("path", c.Path()),
		)

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
			slog.Error("Failed to generate internal JWT",
				slog.String("error", err.Error()),
				slog.String("service", serviceName),
			)
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		// Set headers for internal services
		c.Request().Header.Set("X-User-ID", userID)
		c.Request().Header.Set("X-User-Email", email)
		c.Request().Header.Set("X-User-Role", role)
		c.Request().Header.Set("X-User-Token", token)
		c.Request().Header.Set("X-Api-Gateway", cfg.APIGatewayPrivateKey)
		c.Request().Header.Set("X-Auth-Internal-Service", cfg.AuthInternalSecret)
		c.Request().Header.Set("X-Internal-JWT", internalJWT)

		slog.Debug("Proxy headers set for internal service",
			slog.String("service", serviceName),
			slog.String("user_id", userID),
			slog.String("email", email),
			slog.String("role", role),
		)
		
		return c.Next()
	}
}

func ProxyMiddlewareForAuthenService(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Request().Header.Set("X-Api-Gateway", cfg.APIGatewayPrivateKey)
		c.Request().Header.Set("X-Auth-Internal-Service", cfg.AuthInternalSecret)

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
	privKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	claims := jwt.RegisteredClaims{
		Issuer:    cfg.APIGatewayName,
		Audience:  []string{aud},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privKey)
}

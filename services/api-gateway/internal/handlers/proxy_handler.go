package handlers

import (
	"api_gateway/internal/config"
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ProxyHandler struct {
	cfg *config.Config
}

func NewProxyHandler(cfg *config.Config) *ProxyHandler {
	return &ProxyHandler{cfg: cfg}
}

// ProxyRequest forwards the request to the target service
func (h *ProxyHandler) ProxyRequest(targetURL string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		startTime := time.Now()
		
		// Build target URL
		path := strings.Trim(c.Params("*"), "/")

		fullURL := strings.TrimSuffix(targetURL, "/")
		if path != "" {
			fullURL += "/" + path
		}

		// Add query params
		if len(c.Context().QueryArgs().String()) > 0 {
			fullURL += "?" + c.Context().QueryArgs().String()
		}

		// Log proxy request
		logAttrs := []any{
			slog.String("target_url", fullURL),
			slog.String("method", c.Method()),
			slog.String("original_path", c.Path()),
		}
		
		if userID, ok := c.Locals("userID").(string); ok && userID != "" {
			logAttrs = append(logAttrs, slog.String("user_id", userID))
		}
		
		slog.Debug("Proxying request to service", logAttrs...)
		
		// Create request
		req, err := http.NewRequest(c.Method(), fullURL, bytes.NewReader(c.Body()))
		if err != nil {
			slog.Error("Failed to create proxy request", 
				slog.String("error", err.Error()),
				slog.String("target_url", fullURL),
			)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create request",
			})
		}

		// Copy headers from original request
		c.Request().Header.VisitAll(func(key, value []byte) {
			req.Header.Set(string(key), string(value))
		})

		// Make request
		client := &http.Client{
			Timeout: 30 * time.Second,
		}
		resp, err := client.Do(req)
		if err != nil {
			duration := time.Since(startTime)
			slog.Error("Failed to reach service",
				slog.String("error", err.Error()),
				slog.String("target_url", fullURL),
				slog.Duration("duration", duration),
			)
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
				"error":   "Failed to reach service",
				"details": err.Error(),
			})
		}
		defer resp.Body.Close()

		// Copy response headers
		for key, values := range resp.Header {
			for _, value := range values {
				c.Response().Header.Add(key, value)
			}
		}

		// Copy response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("Failed to read service response",
				slog.String("error", err.Error()),
				slog.String("target_url", fullURL),
			)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to read response",
			})
		}

		duration := time.Since(startTime)
		
		// Log successful proxy
		slog.Info("Proxy request completed",
			slog.String("target_url", fullURL),
			slog.Int("status", resp.StatusCode),
			slog.Duration("duration", duration),
			slog.Int("response_size", len(body)),
		)

		c.Status(resp.StatusCode)
		return c.Send(body)
	}
}

func (h *ProxyHandler) ProxyWebSocket(c *fiber.Ctx) error {
	internalJWT := c.Get("X-Internal-JWT")
	if internalJWT == "" {
		slog.Warn("WebSocket proxy request missing internal JWT",
			slog.String("path", c.Path()),
			slog.String("ip", c.IP()),
		)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing X-Internal-JWT header",
		})
	}

	slog.Info("WebSocket connection info provided",
		slog.String("service", "comment-service"),
		slog.String("path", c.Path()),
	)

	return c.JSON(fiber.Map{
		"comment_service_websocket_url": h.cfg.CommentServiceWebSocketURL,
		"internal_jwt":                  internalJWT,
	})
}

func (h *ProxyHandler) OrderProxyWebSocket(c *fiber.Ctx) error {
	internalJWT := c.Get("X-Internal-JWT")
	if internalJWT == "" {
		slog.Warn("WebSocket proxy request missing internal JWT",
			slog.String("path", c.Path()),
			slog.String("ip", c.IP()),
		)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing X-Internal-JWT header",
		})
	}

	slog.Info("WebSocket connection info provided",
		slog.String("service", "order-service"),
		slog.String("path", c.Path()),
	)

	return c.JSON(fiber.Map{
		"order_service_websocket_url": h.cfg.OrderServiceWebSocketURL,
		"internal_jwt":                internalJWT,
	})
}

// Health check
func (h *ProxyHandler) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"service": "api-gateway",
		"version": h.cfg.OTelServiceVersion,
	})
}

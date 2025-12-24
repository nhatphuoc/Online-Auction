package handlers

import (
	"api_gateway/internal/config"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

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
		// Build target URL
		path := c.Params("*")
		if path == "" {
			path = "/"
		}
		fullURL := strings.TrimSuffix(targetURL, "/") + "/" + strings.TrimPrefix(path, "/")

		// Add query params
		if len(c.Context().QueryArgs().String()) > 0 {
			fullURL += "?" + c.Context().QueryArgs().String()
		}

		fmt.Printf("Proxying request to: %s\n", fullURL)
		// Create request
		req, err := http.NewRequest(c.Method(), fullURL, bytes.NewReader(c.Body()))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create request",
			})
		}

		// Copy headers from original request
		c.Request().Header.VisitAll(func(key, value []byte) {
			req.Header.Set(string(key), string(value))
		})

		// Make request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
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
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to read response",
			})
		}

		c.Status(resp.StatusCode)
		return c.Send(body)
	}
}

func (h *ProxyHandler) ProxyWebSocket(c *fiber.Ctx) error {
	internalJWT := c.Get("X-Internal-JWT")
	if internalJWT == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing X-Internal-JWT header",
		})
	}

	return c.JSON(fiber.Map{
		"comment_service_websocket_url": h.cfg.CommentServiceWebSocketURL,
		"internal_jwt":                  internalJWT,
	})
}

func (h *ProxyHandler) OrderProxyWebSocket(c *fiber.Ctx) error {
	internalJWT := c.Get("X-Internal-JWT")
	if internalJWT == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing X-Internal-JWT header",
		})
	}

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

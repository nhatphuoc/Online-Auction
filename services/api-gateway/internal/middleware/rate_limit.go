package middleware

import (
	"api_gateway/internal/config"
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

// RateLimiter implements token bucket algorithm with Redis
type RateLimiter struct {
	client         *redis.Client
	requestsPerIP  int
	windowSeconds  int
	burstSize      int
	enabled        bool
}

// NewRateLimiter creates a new rate limiter instance
func NewRateLimiter(cfg *config.Config) (*RateLimiter, error) {
	if !cfg.RateLimitEnabled {
		slog.Info("Rate limiting is disabled")
		return &RateLimiter{enabled: false}, nil
	}

	// Initialize Redis client
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.RedisAddr,
		Password:     cfg.RedisPassword,
		DB:           cfg.RedisDB,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 5,
	})

	// Test Redis connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		slog.Error("Failed to connect to Redis", "error", err)
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	slog.Info("Rate limiter initialized",
		"requests_per_ip", cfg.RateLimitRequestsPerIP,
		"window_seconds", cfg.RateLimitWindow,
		"burst_size", cfg.RateLimitBurstSize,
	)

	return &RateLimiter{
		client:        client,
		requestsPerIP: cfg.RateLimitRequestsPerIP,
		windowSeconds: cfg.RateLimitWindow,
		burstSize:     cfg.RateLimitBurstSize,
		enabled:       true,
	}, nil
}

// Middleware returns a Fiber middleware handler for rate limiting
func (rl *RateLimiter) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !rl.enabled {
			return c.Next()
		}

		// Get client identifier (IP address)
		clientIP := c.IP()
		
		// Allow internal health checks
		if c.Path() == "/health" || c.Path() == "/metrics" {
			return c.Next()
		}

		// Check rate limit
		allowed, remaining, resetTime, err := rl.allowRequest(c.Context(), clientIP)
		if err != nil {
			slog.Error("Rate limiter error", "error", err, "ip", clientIP)
			// On error, allow the request to proceed (fail open)
			return c.Next()
		}

		// Set rate limit headers
		c.Set("X-RateLimit-Limit", strconv.Itoa(rl.requestsPerIP))
		c.Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Set("X-RateLimit-Reset", strconv.FormatInt(resetTime, 10))

		if !allowed {
			// Calculate retry-after in seconds
			retryAfter := resetTime - time.Now().Unix()
			if retryAfter < 0 {
				retryAfter = 1
			}
			c.Set("Retry-After", strconv.FormatInt(retryAfter, 10))

			slog.Warn("Rate limit exceeded",
				"ip", clientIP,
				"path", c.Path(),
				"method", c.Method(),
				"retry_after", retryAfter,
			)

			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":   "Rate limit exceeded",
				"message": fmt.Sprintf("Too many requests. Please try again in %d seconds", retryAfter),
				"limit":   rl.requestsPerIP,
				"window":  fmt.Sprintf("%ds", rl.windowSeconds),
			})
		}

		return c.Next()
	}
}

// allowRequest checks if a request from the given identifier is allowed
// Uses sliding window log algorithm for accurate rate limiting
func (rl *RateLimiter) allowRequest(ctx context.Context, identifier string) (bool, int, int64, error) {
	now := time.Now()
	windowStart := now.Add(-time.Duration(rl.windowSeconds) * time.Second)
	
	key := fmt.Sprintf("ratelimit:%s", identifier)
	
	// Use Redis pipeline for atomic operations
	pipe := rl.client.Pipeline()
	
	// Remove old entries outside the current window
	pipe.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(windowStart.UnixNano(), 10))
	
	// Count requests in current window
	countCmd := pipe.ZCard(ctx, key)
	
	// Execute pipeline
	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, 0, 0, err
	}
	
	currentCount := int(countCmd.Val())
	
	// Calculate remaining requests
	remaining := rl.requestsPerIP - currentCount
	if remaining < 0 {
		remaining = 0
	}
	
	// Calculate reset time (end of current window)
	resetTime := now.Add(time.Duration(rl.windowSeconds) * time.Second).Unix()
	
	// Check if request is allowed
	// Allow burst up to burstSize above the limit
	maxAllowed := rl.requestsPerIP + rl.burstSize
	
	if currentCount >= maxAllowed {
		return false, 0, resetTime, nil
	}
	
	// Add current request timestamp to sorted set
	score := now.UnixNano()
	member := fmt.Sprintf("%d:%s", score, generateRequestID())
	
	pipe2 := rl.client.Pipeline()
	pipe2.ZAdd(ctx, key, redis.Z{
		Score:  float64(score),
		Member: member,
	})
	
	// Set expiration on the key to prevent memory leak
	pipe2.Expire(ctx, key, time.Duration(rl.windowSeconds*2)*time.Second)
	
	_, err = pipe2.Exec(ctx)
	if err != nil {
		return false, 0, 0, err
	}
	
	// Update remaining count
	remaining = rl.requestsPerIP - (currentCount + 1)
	if remaining < 0 {
		remaining = 0
	}
	
	return true, remaining, resetTime, nil
}

// generateRequestID generates a simple unique identifier for each request
func generateRequestID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}

// Close closes the Redis connection
func (rl *RateLimiter) Close() error {
	if rl.client != nil {
		return rl.client.Close()
	}
	return nil
}

// ResetLimit resets the rate limit for a specific identifier (for testing or admin purposes)
func (rl *RateLimiter) ResetLimit(ctx context.Context, identifier string) error {
	if !rl.enabled {
		return nil
	}
	
	key := fmt.Sprintf("ratelimit:%s", identifier)
	return rl.client.Del(ctx, key).Err()
}

// GetCurrentCount returns the current request count for an identifier
func (rl *RateLimiter) GetCurrentCount(ctx context.Context, identifier string) (int64, error) {
	if !rl.enabled {
		return 0, nil
	}
	
	now := time.Now()
	windowStart := now.Add(-time.Duration(rl.windowSeconds) * time.Second)
	key := fmt.Sprintf("ratelimit:%s", identifier)
	
	// Remove old entries
	rl.client.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(windowStart.UnixNano(), 10))
	
	// Count current requests
	return rl.client.ZCard(ctx, key).Result()
}

package handlers

import (
	"context"
	"fmt"
	"order_service/internal/config"
	"order_service/internal/models"
	"strconv"

	"github.com/go-pg/pg/v10"
	"github.com/gofiber/fiber/v2"
)

type LikeHandler struct {
	db  *pg.DB
	cfg *config.Config
}

func NewLikeHandler(db *pg.DB, cfg *config.Config) *LikeHandler {
	return &LikeHandler{
		db:  db,
		cfg: cfg,
	}
}

// AddToWatchList godoc
// @Summary Add product to watch list
// @Description Add a product to user's watch list (favorites)
// @Tags WatchList
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.AddToWatchListRequest true "Add to watch list request"
// @Success 201 {object} models.WatchListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /watchlist [post]
func (h *LikeHandler) AddToWatchList(c *fiber.Ctx) error {
	ctx := context.Background()

	fmt.Println("AddToWatchList called")
	// Get user ID from context (set by middleware)
	userID := c.Locals("userID")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized - User ID not found",
		})
	}

	// Convert userID from string to int64
	userIDStr, ok := userID.(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	userIDInt, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Parse request
	var req models.AddToWatchListRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if req.ProductID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product_id",
		})
	}

	// Create watch list item
	watchItem := &models.WatchList{
		UserID:    userIDInt,
		ProductID: req.ProductID,
	}

	// Insert into database
	_, insertErr := h.db.ModelContext(ctx, watchItem).Insert()
	if insertErr != nil {
		// Check for unique constraint violation (product already in watch list)
		if pgErr, ok := insertErr.(pg.Error); ok {
			if pgErr.Field('C') == "23505" { // Unique violation error code
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{
					"error": "Product already in watch list",
				})
			}
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add product to watch list",
		})
	}

	// Fetch complete product information using JOIN
	var response models.WatchListResponse
	query := `
		SELECT 
			w.id,
			w.product_id,
			p.thumbnail_url,
			p.name,
			p.current_price,
			p.buy_now_price,
			p.created_at,
			p.end_at,
			p.bid_count,
			p.category_name
		FROM watch_list w
		INNER JOIN products p ON w.product_id = p.id
		WHERE w.id = ?
	`
	_, err = h.db.QueryOneContext(ctx, &response, query, watchItem.ID)
	if err != nil {
		// If we can't fetch product details, return basic info
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Product added to watch list successfully",
			"data": fiber.Map{
				"id":         watchItem.ID,
				"product_id": watchItem.ProductID,
			},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Product added to watch list successfully",
		"data":    response,
	})
}

// RemoveFromWatchList godoc
// @Summary Remove product from watch list
// @Description Remove a product from user's watch list
// @Tags WatchList
// @Produce json
// @Security BearerAuth
// @Param product_id path int true "Product ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /watchlist/{product_id} [delete]
func (h *LikeHandler) RemoveFromWatchList(c *fiber.Ctx) error {
	ctx := context.Background()

	// Get user ID from context
	userID := c.Locals("userID")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized - User ID not found",
		})
	}

	// Convert userID from string to int64
	userIDStr, ok := userID.(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	userIDInt, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Get product ID from params
	productIDStr := c.Params("product_id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil || productID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product_id",
		})
	}

	// Delete from watch list
	result, err := h.db.ModelContext(ctx, (*models.WatchList)(nil)).
		Where("user_id = ?", userIDInt).
		Where("product_id = ?", productID).
		Delete()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to remove product from watch list",
		})
	}

	if result.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found in watch list",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product removed from watch list successfully",
	})
}

// GetWatchList godoc
// @Summary Get user's watch list
// @Description Get all products in user's watch list
// @Tags WatchList
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 20)"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /watchlist [get]
func (h *LikeHandler) GetWatchList(c *fiber.Ctx) error {
	ctx := context.Background()

	// Get user ID from context
	userID := c.Locals("userID")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized - User ID not found",
		})
	}

	// Convert userID from string to int64
	userIDStr, ok := userID.(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	userIDInt, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Parse pagination params
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	// Get watch list items with product details using JOIN
	var response []models.WatchListResponse
	query := `
		SELECT 
			w.id,
			w.product_id,
			p.thumbnail_url,
			p.name,
			p.current_price,
			p.buy_now_price,
			p.created_at,
			p.end_at,
			p.bid_count,
			p.category_name
		FROM watch_list w
		INNER JOIN products p ON w.product_id = p.id
		WHERE w.user_id = ?
		ORDER BY w.created_at DESC
		LIMIT ? OFFSET ?
	`

	_, err = h.db.QueryContext(ctx, &response, query, userIDInt, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch watch list",
		})
	}

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM watch_list WHERE user_id = ?`
	_, err = h.db.QueryOneContext(ctx, pg.Scan(&total), countQuery, userIDInt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to count watch list items",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Watch list fetched successfully",
		"data":    response,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// CheckInWatchList godoc
// @Summary Check if product is in watch list
// @Description Check if a specific product is in user's watch list
// @Tags WatchList
// @Produce json
// @Security BearerAuth
// @Param product_id path int true "Product ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /watchlist/{product_id}/check [get]
func (h *LikeHandler) CheckInWatchList(c *fiber.Ctx) error {
	ctx := context.Background()

	// Get user ID from context
	userID := c.Locals("userID")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized - User ID not found",
		})
	}

	// Convert userID from string to int64
	userIDStr, ok := userID.(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	userIDInt, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Get product ID from params
	productIDStr := c.Params("product_id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil || productID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product_id",
		})
	}

	// Check if product exists in watch list
	exists, err := h.db.ModelContext(ctx, (*models.WatchList)(nil)).
		Where("user_id = ?", userIDInt).
		Where("product_id = ?", productID).
		Exists()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check watch list",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"is_in_watchlist": exists,
		"product_id":      productID,
	})
}

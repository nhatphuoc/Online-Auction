package handlers

import (
	"category_service/internal/models"
	"log/slog"
	"math"
	"strconv"

	"github.com/go-pg/pg/v10"
	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	db *pg.DB
}

func NewProductHandler(db *pg.DB) *ProductHandler {
	return &ProductHandler{db: db}
}

// GetProductsByCategory godoc
// @Summary Get products by category
// @Description Get paginated list of products filtered by category
// @Tags products
// @Produce json
// @Param category_id query int true "Category ID"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 20, max: 100)"
// @Param status query string false "Product status (ACTIVE, PENDING, FINISHED, REJECTED)"
// @Param sort_by query string false "Sort field (created_at, current_price, end_at)"
// @Param sort_order query string false "Sort order (asc, desc)"
// @Success 200 {object} models.ProductListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products [get]
func (h *ProductHandler) GetProductsByCategory(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse query parameters
	var params models.ProductQueryParams
	if err := c.QueryParser(&params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid query parameters",
		})
	}

	// Validate category_id is required
	if params.CategoryID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "category_id is required",
		})
	}

	// Set defaults
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = 20
	}
	if params.PageSize > 100 {
		params.PageSize = 100
	}
	if params.SortBy == "" {
		params.SortBy = "created_at"
	}
	if params.SortOrder == "" {
		params.SortOrder = "desc"
	}

	// Validate sort_by
	validSortFields := map[string]bool{
		"created_at":    true,
		"current_price": true,
		"end_at":        true,
		"name":          true,
	}
	if !validSortFields[params.SortBy] {
		params.SortBy = "created_at"
	}

	// Validate sort_order
	if params.SortOrder != "asc" && params.SortOrder != "desc" {
		params.SortOrder = "desc"
	}

	// Get category and its children IDs
	categoryIDs, err := h.getCategoryWithChildren(ctx, params.CategoryID)
	if err != nil {
		slog.Error("Error fetching category tree", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch category information",
		})
	}

	// Build query
	query := h.db.ModelContext(ctx, &[]*models.Product{}).
		WhereIn("category_id IN (?)", categoryIDs).
		Relation("Category").
		Relation("Images")

	// Filter by status if provided
	if params.Status != "" {
		validStatuses := []string{"ACTIVE", "PENDING", "FINISHED", "REJECTED"}
		isValid := false
		for _, s := range validStatuses {
			if params.Status == s {
				isValid = true
				break
			}
		}
		if isValid {
			query = query.Where("status = ?", params.Status)
		}
	}

	// Get total count
	total, err := query.Count()
	if err != nil {
		slog.Error("Error counting products", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to count products",
		})
	}

	// Apply sorting
	sortClause := params.SortBy + " " + params.SortOrder
	query = query.Order(sortClause)

	// Apply pagination
	offset := (params.Page - 1) * params.PageSize
	query = query.Limit(params.PageSize).Offset(offset)

	// Execute query
	var products []*models.Product
	err = query.Select()
	if err != nil {
		slog.Error("Error fetching products", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch products",
		})
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(params.PageSize)))

	response := models.ProductListResponse{
		Products:   products,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	}

	return c.JSON(response)
}

// GetProductByID godoc
// @Summary Get product by ID
// @Description Get detailed information about a product
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} models.Product
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products/{id} [get]
func (h *ProductHandler) GetProductByID(c *fiber.Ctx) error {
	ctx := c.Context()

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	var product models.Product
	err = h.db.ModelContext(ctx, &product).
		Where("id = ?", id).
		Relation("Category").
		Relation("Images").
		Select()

	if err != nil {
		if err == pg.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Product not found",
			})
		}
		slog.Error("Error fetching product", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	return c.JSON(product)
}

// Helper function to get category and all its children IDs
func (h *ProductHandler) getCategoryWithChildren(ctx interface{}, categoryID int64) ([]int64, error) {
	var categories []*models.Category
	
	// Get the category and all its potential children
	err := h.db.Model(&categories).
		WhereGroup(func(q *pg.Query) (*pg.Query, error) {
			q = q.WhereOr("id = ?", categoryID)
			q = q.WhereOr("parent_id = ?", categoryID)
			return q, nil
		}).
		Where("is_active = ?", true).
		Select()

	if err != nil {
		return nil, err
	}

	if len(categories) == 0 {
		return []int64{categoryID}, nil
	}

	// Collect all IDs
	ids := make([]int64, 0, len(categories))
	for _, cat := range categories {
		ids = append(ids, cat.ID)
	}

	return ids, nil
}

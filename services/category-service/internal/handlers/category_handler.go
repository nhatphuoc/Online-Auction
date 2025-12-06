package handlers

import (
	"category_service/internal/models"
	"category_service/internal/utils"
	"log/slog"
	"strconv"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/gofiber/fiber/v2"
)

type CategoryHandler struct {
	db *pg.DB
}

func NewCategoryHandler(db *pg.DB) *CategoryHandler {
	return &CategoryHandler{db: db}
}

// CreateCategory godoc
// @Summary Create a new category
// @Description Create a new product category
// @Tags categories
// @Accept json
// @Produce json
// @Param category body models.CreateCategoryRequest true "Category data"
// @Success 201 {object} models.Category
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /categories [post]
func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {
	ctx := c.Context()
	
	var req models.CreateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Determine level based on parent
	level := 1
	if req.ParentID != nil {
		var parent models.Category
		err := h.db.ModelContext(ctx, &parent).Where("id = ?", *req.ParentID).Select()
		if err != nil {
			if err == pg.ErrNoRows {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Parent category not found",
				})
			}
			slog.Error("Error fetching parent category", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Database error",
			})
		}
		level = parent.Level + 1
		
		// Limit to 2 levels only
		if level > 2 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Maximum category depth is 2 levels",
			})
		}
	}

	category := &models.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		ParentID:    req.ParentID,
		Level:       level,
		IsActive:    true,
		DisplayOrder: req.DisplayOrder,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err := h.db.ModelContext(ctx, category).Insert()
	if err != nil {
		slog.Error("Error creating category", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create category",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(category)
}

// GetCategories godoc
// @Summary Get all categories
// @Description Get hierarchical list of all categories
// @Tags categories
// @Produce json
// @Param parent_id query int false "Parent category ID (empty for root categories)"
// @Param level query int false "Category level (1 or 2)"
// @Success 200 {object} models.CategoryTreeResponse
// @Failure 500 {object} map[string]interface{}
// @Router /categories [get]
func (h *CategoryHandler) GetCategories(c *fiber.Ctx) error {
	ctx := c.Context()
	
	parentIDStr := c.Query("parent_id")
	levelStr := c.Query("level")

	query := h.db.ModelContext(ctx, &[]*models.Category{}).
		Where("is_active = ?", true).
		Order("display_order ASC", "name ASC")

	// Filter by parent_id if provided
	if parentIDStr != "" {
		if parentIDStr == "null" || parentIDStr == "0" {
			query = query.WhereGroup(func(q *pg.Query) (*pg.Query, error) {
				return q.WhereOr("parent_id IS NULL").WhereOr("parent_id = 0"), nil
			})
		} else {
			parentID, err := strconv.ParseInt(parentIDStr, 10, 64)
			if err == nil {
				query = query.Where("parent_id = ?", parentID)
			}
		}
	}

	// Filter by level if provided
	if levelStr != "" {
		level, err := strconv.Atoi(levelStr)
		if err == nil && (level == 1 || level == 2) {
			query = query.Where("level = ?", level)
		}
	}

	var categories []*models.Category
	err := query.Select()
	if err != nil {
		slog.Error("Error fetching categories", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch categories",
		})
	}

	// Build hierarchical response if no filters applied
	if parentIDStr == "" && levelStr == "" {
		tree := h.buildCategoryTree(ctx, categories)
		return c.JSON(models.CategoryTreeResponse{
			Categories: tree,
		})
	}

	// Return flat list if filters applied
	response := make([]*models.CategoryResponse, len(categories))
	for i, cat := range categories {
		response[i] = h.toCategoryResponse(cat)
	}

	return c.JSON(models.CategoryTreeResponse{
		Categories: response,
	})
}

// GetCategoryByID godoc
// @Summary Get category by ID
// @Description Get a single category by ID with its children
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} models.CategoryResponse
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetCategoryByID(c *fiber.Ctx) error {
	ctx := c.Context()
	
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid category ID",
		})
	}

	var category models.Category
	err = h.db.ModelContext(ctx, &category).
		Where("id = ?", id).
		Relation("Children", func(q *pg.Query) (*pg.Query, error) {
			return q.Where("is_active = ?", true).Order("display_order ASC"), nil
		}).
		Select()

	if err != nil {
		if err == pg.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Category not found",
			})
		}
		slog.Error("Error fetching category", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	response := h.toCategoryResponse(&category)
	return c.JSON(response)
}

// UpdateCategory godoc
// @Summary Update a category
// @Description Update an existing category
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body models.UpdateCategoryRequest true "Category data"
// @Success 200 {object} models.Category
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	ctx := c.Context()
	
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid category ID",
		})
	}

	var req models.UpdateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var category models.Category
	err = h.db.ModelContext(ctx, &category).Where("id = ?", id).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Category not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// Update fields
	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.Slug != nil {
		category.Slug = *req.Slug
	}
	if req.Description != nil {
		category.Description = *req.Description
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}
	if req.DisplayOrder != nil {
		category.DisplayOrder = *req.DisplayOrder
	}
	if req.ParentID != nil {
		// Validate parent exists
		var parent models.Category
		err := h.db.ModelContext(ctx, &parent).Where("id = ?", *req.ParentID).Select()
		if err != nil {
			if err == pg.ErrNoRows {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Parent category not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Database error",
			})
		}
		
		// Check level constraint
		newLevel := parent.Level + 1
		if newLevel > 2 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Maximum category depth is 2 levels",
			})
		}
		
		category.ParentID = req.ParentID
		category.Level = newLevel
	}

	category.UpdatedAt = time.Now()

	_, err = h.db.ModelContext(ctx, &category).WherePK().Update()
	if err != nil {
		slog.Error("Error updating category", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update category",
		})
	}

	return c.JSON(category)
}

// DeleteCategory godoc
// @Summary Delete a category
// @Description Soft delete a category (marks as inactive)
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	ctx := c.Context()
	
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid category ID",
		})
	}

	// Check if category has children
	count, err := h.db.ModelContext(ctx, &models.Category{}).
		Where("parent_id = ?", id).
		Where("is_active = ?", true).
		Count()

	if err != nil {
		slog.Error("Error checking category children", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	if count > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot delete category with active children",
		})
	}

	// Soft delete
	_, err = h.db.ModelContext(ctx, &models.Category{}).
		Set("is_active = ?", false).
		Set("updated_at = ?", time.Now()).
		Where("id = ?", id).
		Update()

	if err != nil {
		slog.Error("Error deleting category", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete category",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Category deleted successfully",
	})
}

// Helper functions
func (h *CategoryHandler) buildCategoryTree(ctx interface{}, categories []*models.Category) []*models.CategoryResponse {
	categoryMap := make(map[int64]*models.CategoryResponse)
	var roots []*models.CategoryResponse

	// First pass: create all response objects
	for _, cat := range categories {
		categoryMap[cat.ID] = h.toCategoryResponse(cat)
	}

	// Second pass: build tree structure
	for _, cat := range categories {
		response := categoryMap[cat.ID]
		if cat.ParentID == nil || *cat.ParentID == 0 {
			roots = append(roots, response)
		} else {
			if parent, ok := categoryMap[*cat.ParentID]; ok {
				parent.Children = append(parent.Children, response)
			}
		}
	}

	return roots
}

func (h *CategoryHandler) toCategoryResponse(cat *models.Category) *models.CategoryResponse {
	response := &models.CategoryResponse{
		ID:          cat.ID,
		Name:        cat.Name,
		Slug:        cat.Slug,
		Description: cat.Description,
		ParentID:    cat.ParentID,
		Level:       cat.Level,
		IsActive:    cat.IsActive,
		DisplayOrder: cat.DisplayOrder,
		CreatedAt:   cat.CreatedAt,
		UpdatedAt:   cat.UpdatedAt,
	}

	if len(cat.Children) > 0 {
		response.Children = make([]*models.CategoryResponse, len(cat.Children))
		for i, child := range cat.Children {
			response.Children[i] = h.toCategoryResponse(child)
		}
	}

	return response
}

// GetCategoriesByParent godoc
// @Summary Get categories by parent ID
// @Description Get all child categories of a parent category
// @Tags categories
// @Produce json
// @Param parent_id path int true "Parent Category ID"
// @Success 200 {object} []models.CategoryResponse
// @Failure 500 {object} map[string]interface{}
// @Router /categories/parent/{parent_id} [get]
func (h *CategoryHandler) GetCategoriesByParent(c *fiber.Ctx) error {
	ctx := c.Context()
	
	parentID, err := strconv.ParseInt(c.Params("parent_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid parent ID",
		})
	}

	var categories []*models.Category
	err = h.db.ModelContext(ctx, &categories).
		Where("parent_id = ?", parentID).
		Where("is_active = ?", true).
		Order("display_order ASC", "name ASC").
		Select()

	if err != nil {
		slog.Error("Error fetching categories", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch categories",
		})
	}

	response := make([]*models.CategoryResponse, len(categories))
	for i, cat := range categories {
		response[i] = h.toCategoryResponse(cat)
	}

	return c.JSON(response)
}

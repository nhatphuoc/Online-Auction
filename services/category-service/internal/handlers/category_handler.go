package handlers

import (
	"category_service/internal/models"
	"category_service/internal/utils"
	"context"
	"strconv"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/gofiber/fiber/v2"
)

var locVN = mustLoad("Asia/Ho_Chi_Minh")

func mustLoad(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}
	return loc
}

func FixedTimeNow() time.Time {
	nowVN := time.Now().In(locVN)
	return time.Date(
		nowVN.Year(), nowVN.Month(), nowVN.Day(),
		nowVN.Hour(), nowVN.Minute(), nowVN.Second(),
		nowVN.Nanosecond(),
		time.UTC,
	)
}

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
	ctx := context.Background()
	var req models.CreateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	level := 1
	if req.ParentID != nil {
		var parentLevel int
		_, err := h.db.QueryOneContext(ctx, pg.Scan(&parentLevel), "SELECT level FROM categories WHERE id = ? AND is_active = true", *req.ParentID)
		if err != nil {
			if err == pg.ErrNoRows {
				return utils.ErrorResponse(c, fiber.StatusBadRequest, "Parent category not found")
			}
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
		}
		level = parentLevel + 1
		if level > 2 {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Maximum category depth is 2 levels")
		}
	}
	var id int64
	query := `INSERT INTO categories (name, slug, description, parent_id, level, is_active, display_order, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`
	_, err := h.db.QueryOneContext(ctx, pg.Scan(&id), query,
		req.Name, req.Slug, req.Description, req.ParentID, level, true, req.DisplayOrder, FixedTimeNow(), FixedTimeNow())
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create category: "+err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":            id,
		"name":          req.Name,
		"slug":          req.Slug,
		"description":   req.Description,
		"parent_id":     req.ParentID,
		"level":         level,
		"is_active":     true,
		"display_order": req.DisplayOrder,
		"created_at":    FixedTimeNow(),
		"updated_at":    FixedTimeNow(),
	})
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
	ctx := context.Background()
	parentIDStr := c.Query("parent_id")
	levelStr := c.Query("level")
	var categories []*models.Category
	var query string
	var args []interface{}
	query = "SELECT id, name, slug, description, parent_id, level, is_active, display_order, created_at, updated_at FROM categories WHERE is_active = true"
	if parentIDStr != "" && parentIDStr != "null" && parentIDStr != "0" {
		query += " AND parent_id = ?"
		parentID, err := strconv.ParseInt(parentIDStr, 10, 64)
		if err == nil {
			args = append(args, parentID)
		}
	}
	if levelStr != "" {
		level, err := strconv.Atoi(levelStr)
		if err == nil && (level == 1 || level == 2) {
			query += " AND level = ?"
			args = append(args, level)
		}
	}
	query += " ORDER BY display_order ASC, name ASC"
	_, err := h.db.QueryContext(ctx, &categories, query, args...)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch categories")
	}
	if parentIDStr == "" && levelStr == "" {
		tree := h.buildCategoryTree(ctx, categories)
		return c.JSON(models.CategoryTreeResponse{
			Categories: tree,
		})
	}
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
	ctx := context.Background()
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid category ID")
	}
	var category models.Category
	_, err = h.db.QueryOneContext(ctx, &category,
		"SELECT id, name, slug, description, parent_id, level, is_active, display_order, created_at, updated_at FROM categories WHERE id = ? AND is_active = true", id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Category not found")
	}
	// Lấy children
	var children []*models.Category
	_, err = h.db.QueryContext(ctx, &children,
		"SELECT id, name, slug, description, parent_id, level, is_active, display_order, created_at, updated_at FROM categories WHERE parent_id = ? AND is_active = true ORDER BY display_order ASC", id)
	category.Children = children
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
	var req models.UpdateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	ctx := context.Background()
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid category ID")
	}

	// Start transaction for data consistency
	tx, err := h.db.Begin()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to start transaction")
	}
	defer tx.Rollback()

	// Get current category info
	var currentCategory models.Category
	_, err = tx.QueryOneContext(ctx, &currentCategory,
		"SELECT id, name, slug, description, parent_id, level, is_active, display_order, created_at, updated_at FROM categories WHERE id = ?", id)
	if err != nil {
		if err == pg.ErrNoRows {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Category not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
	}

	// Build update query dynamically
	setFields := ""
	args := []interface{}{}
	categoryNameChanged := false
	var newLevel int = currentCategory.Level

	if req.Name != nil && *req.Name != currentCategory.Name {
		setFields += "name = ?, "
		args = append(args, *req.Name)
		categoryNameChanged = true
	}
	if req.Slug != nil {
		setFields += "slug = ?, "
		args = append(args, *req.Slug)
	}
	if req.Description != nil {
		setFields += "description = ?, "
		args = append(args, *req.Description)
	}
	if req.IsActive != nil {
		setFields += "is_active = ?, "
		args = append(args, *req.IsActive)
	}
	if req.DisplayOrder != nil {
		setFields += "display_order = ?, "
		args = append(args, *req.DisplayOrder)
	}
	if req.ParentID != nil {
		var parentCategory models.Category
		_, err := tx.QueryOneContext(ctx, &parentCategory, 
			"SELECT id, name, level, parent_id FROM categories WHERE id = ? AND is_active = true", *req.ParentID)
		if err != nil {
			if err == pg.ErrNoRows {
				return utils.ErrorResponse(c, fiber.StatusBadRequest, "Parent category not found")
			}
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
		}
		if parentCategory.Level+1 > 2 {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Maximum category depth is 2 levels")
		}
		newLevel = parentCategory.Level + 1
		setFields += "parent_id = ?, level = ?, "
		args = append(args, *req.ParentID, newLevel)
		categoryNameChanged = true // Parent changed, need to update products
	}

	if setFields == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "No fields to update")
	}

	setFields += "updated_at = ?"
	args = append(args, FixedTimeNow())
	updateQuery := "UPDATE categories SET " + setFields + " WHERE id = ?"
	args = append(args, id)
	
	_, err = tx.ExecContext(ctx, updateQuery, args...)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update category: "+err.Error())
	}

	// Get updated category
	var updatedCategory models.Category
	_, err = tx.QueryOneContext(ctx, &updatedCategory,
		"SELECT id, name, slug, description, parent_id, level, is_active, display_order, created_at, updated_at FROM categories WHERE id = ?", id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch updated category")
	}

	// Update products table if category name or parent changed
	if categoryNameChanged {
		if err := h.updateProductsAfterCategoryChange(ctx, tx, id, &updatedCategory); err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update products: "+err.Error())
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to commit transaction")
	}

	return c.JSON(updatedCategory)
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
	ctx := context.Background()
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid category ID")
	}

	// Check if category exists
	var categoryExists int
	_, err = h.db.QueryOneContext(ctx, pg.Scan(&categoryExists), 
		"SELECT COUNT(*) FROM categories WHERE id = ? AND is_active = true", id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
	}
	if categoryExists == 0 {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Category not found")
	}

	// Check for active children categories
	var childrenCount int
	_, err = h.db.QueryOneContext(ctx, pg.Scan(&childrenCount), 
		"SELECT COUNT(*) FROM categories WHERE parent_id = ? AND is_active = true", id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
	}
	if childrenCount > 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Cannot delete category with active children categories")
	}

	// Check for products using this category
	var productsCount int
	_, err = h.db.QueryOneContext(ctx, pg.Scan(&productsCount), 
		"SELECT COUNT(*) FROM products WHERE category_id = ?", id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Database error checking products")
	}
	if productsCount > 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 
			"Cannot delete category with existing products. Please reassign or remove products first")
	}

	// Soft delete category
	_, err = h.db.ExecContext(ctx, 
		"UPDATE categories SET is_active = false, updated_at = ? WHERE id = ?", 
		FixedTimeNow(), id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete category: "+err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "Category deleted successfully",
	})
}

// Helper functions

// updateProductsAfterCategoryChange updates products table when category name or hierarchy changes
func (h *CategoryHandler) updateProductsAfterCategoryChange(ctx context.Context, tx *pg.Tx, categoryID int64, updatedCategory *models.Category) error {
	// Get parent category info if exists
	var parentCategoryName *string
	var parentCategoryID *int64
	
	if updatedCategory.ParentID != nil && *updatedCategory.ParentID != 0 {
		var parentCategory models.Category
		_, err := tx.QueryOneContext(ctx, &parentCategory,
			"SELECT id, name FROM categories WHERE id = ? AND is_active = true", *updatedCategory.ParentID)
		if err == nil {
			parentCategoryName = &parentCategory.Name
			parentCategoryID = &parentCategory.ID
		}
	}

	// Update products that belong to this category
	var updateQuery string
	var args []interface{}
	
	if updatedCategory.Level == 1 {
		// Level 1 category: update category_name and clear parent info
		updateQuery = `UPDATE products 
			SET category_name = ?, 
			    parent_category_id = NULL, 
			    parent_category_name = NULL 
			WHERE category_id = ?`
		args = []interface{}{updatedCategory.Name, categoryID}
	} else {
		// Level 2 category: update both category_name and parent info
		updateQuery = `UPDATE products 
			SET category_name = ?, 
			    parent_category_id = ?, 
			    parent_category_name = ? 
			WHERE category_id = ?`
		args = []interface{}{updatedCategory.Name, parentCategoryID, parentCategoryName, categoryID}
	}

	_, err := tx.ExecContext(ctx, updateQuery, args...)
	if err != nil {
		return err
	}

	// Also update products in child categories if this is a level 1 category
	if updatedCategory.Level == 1 {
		var childCategories []*models.Category
		_, err := tx.QueryContext(ctx, &childCategories,
			"SELECT id, name FROM categories WHERE parent_id = ? AND is_active = true", categoryID)
		if err != nil {
			return err
		}

		for _, child := range childCategories {
			_, err := tx.ExecContext(ctx,
				`UPDATE products 
				SET parent_category_id = ?, 
				    parent_category_name = ? 
				WHERE category_id = ?`,
				categoryID, updatedCategory.Name, child.ID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

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
		ID:           cat.ID,
		Name:         cat.Name,
		Slug:         cat.Slug,
		Description:  cat.Description,
		ParentID:     cat.ParentID,
		Level:        cat.Level,
		IsActive:     cat.IsActive,
		DisplayOrder: cat.DisplayOrder,
		CreatedAt:    cat.CreatedAt,
		UpdatedAt:    cat.UpdatedAt,
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
	ctx := context.Background()
	parentID, err := strconv.ParseInt(c.Params("parent_id"), 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid parent ID")
	}
	var categories []*models.Category
	_, err = h.db.QueryContext(ctx, &categories,
		"SELECT id, name, slug, description, parent_id, level, is_active, display_order, created_at, updated_at FROM categories WHERE parent_id = ? AND is_active = true ORDER BY display_order ASC, name ASC", parentID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch categories")
	}
	response := make([]*models.CategoryResponse, len(categories))
	for i, cat := range categories {
		response[i] = h.toCategoryResponse(cat)
	}
	return c.JSON(response)
}

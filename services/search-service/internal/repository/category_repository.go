package repository

import (
	"context"
	"search-service/internal/models"

	"github.com/go-pg/pg/v10"
)

type CategoryRepository struct {
	db *pg.DB
}

func NewCategoryRepository(db *pg.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// GetCategoryByID fetches a category by ID from PostgreSQL
func (r *CategoryRepository) GetCategoryByID(ctx context.Context, id int64) (*models.Category, error) {
	var category models.Category
	err := r.db.ModelContext(ctx, &category).
		Where("id = ?", id).
		Select()
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// GetAllCategories fetches all categories from PostgreSQL (for bulk indexing)
func (r *CategoryRepository) GetAllCategories(ctx context.Context) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.ModelContext(ctx, &categories).
		Select()
	if err != nil {
		return nil, err
	}
	return categories, nil
}

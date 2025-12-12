package repository

import (
	"context"
	"search-service/internal/models"

	"github.com/go-pg/pg/v10"
)

type ProductRepository struct {
	db *pg.DB
}

func NewProductRepository(db *pg.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// GetProductByID fetches a product by ID from PostgreSQL
func (r *ProductRepository) GetProductByID(ctx context.Context, id int64) (*models.Product, error) {
	var product models.Product
	err := r.db.ModelContext(ctx, &product).
		Where("id = ?", id).
		Select()
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// GetProductWithCategory fetches a product with its category information
func (r *ProductRepository) GetProductWithCategory(ctx context.Context, id int64) (*models.Product, *models.Category, error) {
	var product models.Product
	err := r.db.ModelContext(ctx, &product).
		Where("id = ?", id).
		Select()
	if err != nil {
		return nil, nil, err
	}

	var category models.Category
	err = r.db.ModelContext(ctx, &category).
		Where("id = ?", product.CategoryID).
		Select()
	if err != nil {
		return &product, nil, err
	}

	return &product, &category, nil
}

// GetAllProducts fetches all products from PostgreSQL (for bulk indexing)
func (r *ProductRepository) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	var products []models.Product
	err := r.db.ModelContext(ctx, &products).
		Select()
	if err != nil {
		return nil, err
	}
	return products, nil
}

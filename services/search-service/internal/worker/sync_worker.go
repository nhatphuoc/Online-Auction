package worker

import (
	"context"
	"log"
	"search-service/internal/elasticsearch"
	"search-service/internal/models"
	"search-service/internal/repository"
)

type SyncWorker struct {
	productRepo  *repository.ProductRepository
	categoryRepo *repository.CategoryRepository
	indexer      *elasticsearch.Indexer
}

func NewSyncWorker(
	productRepo *repository.ProductRepository,
	categoryRepo *repository.CategoryRepository,
	indexer *elasticsearch.Indexer,
) *SyncWorker {
	return &SyncWorker{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
		indexer:      indexer,
	}
}

func (w *SyncWorker) HandleEvent(ctx context.Context, event *models.Event) error {
	log.Printf("Processing event: type=%s, entity_id=%d", event.Type, event.EntityID)

	switch event.Type {
	case models.EventProductCreated, models.EventProductUpdated:
		return w.syncProduct(ctx, event.EntityID)
	case models.EventProductDeleted:
		return w.deleteProduct(ctx, event.EntityID)
	case models.EventCategoryCreated, models.EventCategoryUpdated:
		return w.syncCategory(ctx, event.EntityID)
	case models.EventCategoryDeleted:
		return w.deleteCategory(ctx, event.EntityID)
	default:
		log.Printf("Unknown event type: %s", event.Type)
		return nil
	}
}

func (w *SyncWorker) syncProduct(ctx context.Context, productID int64) error {
	product, category, err := w.productRepo.GetProductWithCategory(ctx, productID)
	if err != nil {
		log.Printf("Error fetching product %d from database: %v", productID, err)
		return err
	}

	esDoc := w.indexer.ConvertProductToESDocument(product, category)
	
	if err := w.indexer.IndexProduct(ctx, esDoc); err != nil {
		log.Printf("Error indexing product %d: %v", productID, err)
		return err
	}

	log.Printf("Successfully indexed product %d", productID)
	return nil
}

func (w *SyncWorker) deleteProduct(ctx context.Context, productID int64) error {
	if err := w.indexer.DeleteProduct(ctx, productID); err != nil {
		log.Printf("Error deleting product %d from index: %v", productID, err)
		return err
	}

	log.Printf("Successfully deleted product %d from index", productID)
	return nil
}

func (w *SyncWorker) syncCategory(ctx context.Context, categoryID int64) error {
	category, err := w.categoryRepo.GetCategoryByID(ctx, categoryID)
	if err != nil {
		log.Printf("Error fetching category %d from database: %v", categoryID, err)
		return err
	}

	esDoc := w.indexer.ConvertCategoryToESDocument(category)
	
	if err := w.indexer.IndexCategory(ctx, esDoc); err != nil {
		log.Printf("Error indexing category %d: %v", categoryID, err)
		return err
	}

	log.Printf("Successfully indexed category %d", categoryID)
	return nil
}

func (w *SyncWorker) deleteCategory(ctx context.Context, categoryID int64) error {
	if err := w.indexer.DeleteCategory(ctx, categoryID); err != nil {
		log.Printf("Error deleting category %d from index: %v", categoryID, err)
		return err
	}

	log.Printf("Successfully deleted category %d from index", categoryID)
	return nil
}

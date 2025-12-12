package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"search-service/internal/models"
	"search-service/internal/utils"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type Indexer struct {
	es                  *elasticsearch.Client
	productIndexName    string
	categoryIndexName   string
}

func NewIndexer(es *elasticsearch.Client, productIndex, categoryIndex string) *Indexer {
	return &Indexer{
		es:                es,
		productIndexName:  productIndex,
		categoryIndexName: categoryIndex,
	}
}

func (idx *Indexer) IndexProduct(ctx context.Context, product *models.ProductESDocument) error {
	data, err := json.Marshal(product)
	if err != nil {
		return fmt.Errorf("error marshaling product: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      idx.productIndexName,
		DocumentID: fmt.Sprintf("%d", product.ID),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, idx.es)
	if err != nil {
		return fmt.Errorf("error indexing product: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing product: %s", res.String())
	}

	return nil
}

func (idx *Indexer) DeleteProduct(ctx context.Context, productID int64) error {
	req := esapi.DeleteRequest{
		Index:      idx.productIndexName,
		DocumentID: fmt.Sprintf("%d", productID),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, idx.es)
	if err != nil {
		return fmt.Errorf("error deleting product: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 404 {
		return fmt.Errorf("error deleting product: %s", res.String())
	}

	return nil
}

func (idx *Indexer) IndexCategory(ctx context.Context, category *models.CategoryESDocument) error {
	data, err := json.Marshal(category)
	if err != nil {
		return fmt.Errorf("error marshaling category: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      idx.categoryIndexName,
		DocumentID: fmt.Sprintf("%d", category.ID),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, idx.es)
	if err != nil {
		return fmt.Errorf("error indexing category: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing category: %s", res.String())
	}

	return nil
}

func (idx *Indexer) DeleteCategory(ctx context.Context, categoryID int64) error {
	req := esapi.DeleteRequest{
		Index:      idx.categoryIndexName,
		DocumentID: fmt.Sprintf("%d", categoryID),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, idx.es)
	if err != nil {
		return fmt.Errorf("error deleting category: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 404 {
		return fmt.Errorf("error deleting category: %s", res.String())
	}

	return nil
}

func (idx *Indexer) ConvertProductToESDocument(product *models.Product, category *models.Category) *models.ProductESDocument {
	doc := &models.ProductESDocument{
		ID:            product.ID,
		Name:          product.Name,
		NameNoAccent:  utils.RemoveVietnameseAccents(product.Name),
		Description:   product.Description,
		DescriptionNoAccent: utils.RemoveVietnameseAccents(product.Description),
		CategoryID:    product.CategoryID,
		SellerID:      product.SellerID,
		StartingPrice: product.StartingPrice,
		CurrentPrice:  product.CurrentPrice,
		BuyNowPrice:   product.BuyNowPrice,
		StepPrice:     product.StepPrice,
		Status:        product.Status,
		ThumbnailURL:  product.ThumbnailURL,
		AutoExtend:    product.AutoExtend,
		CurrentBidder: product.CurrentBidder,
		EndAt:         product.EndAt,
		CreatedAt:     product.CreatedAt,
	}

	if category != nil {
		doc.CategoryName = category.Name
		doc.CategorySlug = category.Slug
	}

	return doc
}

func (idx *Indexer) ConvertCategoryToESDocument(category *models.Category) *models.CategoryESDocument {
	return &models.CategoryESDocument{
		ID:           category.ID,
		Name:         category.Name,
		NameNoAccent: utils.RemoveVietnameseAccents(category.Name),
		Slug:         category.Slug,
		Description:  category.Description,
		ParentID:     category.ParentID,
		Level:        category.Level,
		IsActive:     category.IsActive,
		DisplayOrder: category.DisplayOrder,
		CreatedAt:    category.CreatedAt,
		UpdatedAt:    category.UpdatedAt,
	}
}

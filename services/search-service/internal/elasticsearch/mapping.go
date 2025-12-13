package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

func CreateProductIndex(ctx context.Context, es *elasticsearch.Client, indexName string) error {
	mapping := map[string]interface{}{
		"settings": map[string]interface{}{
			"analysis": map[string]interface{}{
				"analyzer": map[string]interface{}{
					"vietnamese_analyzer": map[string]interface{}{
						"type":      "custom",
						"tokenizer": "standard",
						"filter": []string{
							"lowercase",
							"asciifolding",
						},
					},
				},
			},
			"number_of_shards":   1,
			"number_of_replicas": 1,
		},
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"id":                    map[string]string{"type": "long"},
				"name":                  map[string]string{"type": "text", "analyzer": "standard"},
				"name_no_accent":        map[string]string{"type": "text", "analyzer": "vietnamese_analyzer"},
				"description":           map[string]string{"type": "text", "analyzer": "standard"},
				"description_no_accent": map[string]string{"type": "text", "analyzer": "vietnamese_analyzer"},
				"category_id":           map[string]string{"type": "long"},
				"category_name":         map[string]string{"type": "text", "analyzer": "vietnamese_analyzer"},
				"category_slug":         map[string]string{"type": "keyword"},
				"seller_id":             map[string]string{"type": "long"},
				"starting_price":        map[string]string{"type": "double"},
				"current_price":         map[string]string{"type": "double"},
				"buy_now_price":         map[string]string{"type": "double"},
				"step_price":            map[string]string{"type": "double"},
				"status":                map[string]string{"type": "keyword"},
				"thumbnail_url":         map[string]string{"type": "keyword"},
				"auto_extend":           map[string]string{"type": "boolean"},
				"current_bidder":        map[string]string{"type": "long"},
				"current_bid_count":     map[string]string{"type": "integer"},
				"end_at":                map[string]string{"type": "date"},
				"created_at":            map[string]string{"type": "date"},
				"current_bidder_info": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"user_id":    map[string]string{"type": "long"},
						"username":   map[string]string{"type": "keyword"},
						"avatar_url": map[string]string{"type": "keyword"},
					},
				},
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(mapping); err != nil {
		return fmt.Errorf("error encoding mapping: %w", err)
	}

	res, err := es.Indices.Create(
		indexName,
		es.Indices.Create.WithContext(ctx),
		es.Indices.Create.WithBody(&buf),
	)
	if err != nil {
		return fmt.Errorf("error creating index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error creating index: %s", res.String())
	}

	log.Printf("Index '%s' created successfully", indexName)
	return nil
}

func CreateCategoryIndex(ctx context.Context, es *elasticsearch.Client, indexName string) error {
	mapping := map[string]interface{}{
		"settings": map[string]interface{}{
			"analysis": map[string]interface{}{
				"analyzer": map[string]interface{}{
					"vietnamese_analyzer": map[string]interface{}{
						"type":      "custom",
						"tokenizer": "standard",
						"filter": []string{
							"lowercase",
							"asciifolding",
						},
					},
				},
			},
			"number_of_shards":   1,
			"number_of_replicas": 1,
		},
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"id":             map[string]string{"type": "long"},
				"name":           map[string]string{"type": "text", "analyzer": "standard"},
				"name_no_accent": map[string]string{"type": "text", "analyzer": "vietnamese_analyzer"},
				"slug":           map[string]string{"type": "keyword"},
				"description":    map[string]string{"type": "text"},
				"parent_id":      map[string]string{"type": "long"},
				"level":          map[string]string{"type": "integer"},
				"is_active":      map[string]string{"type": "boolean"},
				"display_order":  map[string]string{"type": "integer"},
				"created_at":     map[string]string{"type": "date"},
				"updated_at":     map[string]string{"type": "date"},
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(mapping); err != nil {
		return fmt.Errorf("error encoding mapping: %w", err)
	}

	res, err := es.Indices.Create(
		indexName,
		es.Indices.Create.WithContext(ctx),
		es.Indices.Create.WithBody(&buf),
	)
	if err != nil {
		return fmt.Errorf("error creating index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error creating index: %s", res.String())
	}

	log.Printf("Index '%s' created successfully", indexName)
	return nil
}

func IndexExists(ctx context.Context, es *elasticsearch.Client, indexName string) (bool, error) {
	res, err := es.Indices.Exists(
		[]string{indexName},
		es.Indices.Exists.WithContext(ctx),
	)
	if err != nil {
		return false, fmt.Errorf("error checking index existence: %w", err)
	}
	defer res.Body.Close()

	return res.StatusCode == 200, nil
}

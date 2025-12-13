package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"search-service/internal/models"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
)

type Searcher struct {
	es               *elasticsearch.Client
	productIndexName string
	boostMinutes     int
	boostScore       float64
}

func NewSearcher(es *elasticsearch.Client, productIndex string, boostMinutes int, boostScore float64) *Searcher {
	return &Searcher{
		es:               es,
		productIndexName: productIndex,
		boostMinutes:     boostMinutes,
		boostScore:       boostScore,
	}
}

func (s *Searcher) SearchProducts(ctx context.Context, req *models.SearchRequest) (*models.SearchResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	from := (req.Page - 1) * req.PageSize
	
	query := s.buildQuery(req)
	
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("error encoding query: %w", err)
	}

	res, err := s.es.Search(
		s.es.Search.WithContext(ctx),
		s.es.Search.WithIndex(s.productIndexName),
		s.es.Search.WithBody(&buf),
		s.es.Search.WithFrom(from),
		s.es.Search.WithSize(req.PageSize),
	)
	if err != nil {
		return nil, fmt.Errorf("error searching products: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error searching products: %s", res.String())
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var searchResult map[string]interface{}
	if err := json.Unmarshal(body, &searchResult); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	hits := searchResult["hits"].(map[string]interface{})
	total := int64(hits["total"].(map[string]interface{})["value"].(float64))
	
	var products []models.ProductListingResponse
	for _, hit := range hits["hits"].([]interface{}) {
		source := hit.(map[string]interface{})["_source"]
		sourceBytes, _ := json.Marshal(source)
		
		var esDoc models.ProductESDocument
		if err := json.Unmarshal(sourceBytes, &esDoc); err != nil {
			continue
		}
		
		// Convert to ProductListingResponse with calculated fields
		listingProduct := s.convertToListingResponse(&esDoc)
		products = append(products, listingProduct)
	}

	totalPages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	
	return &models.SearchResponse{
		Products:   products,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *Searcher) buildQuery(req *models.SearchRequest) map[string]interface{} {
	query := make(map[string]interface{})
	
	// Build function_score query for boosting recent products
	functionScore := map[string]interface{}{
		"query": s.buildBaseQuery(req),
		"functions": []map[string]interface{}{
			{
				"gauss": map[string]interface{}{
					"created_at": map[string]interface{}{
						"origin": time.Now().Format(time.RFC3339),
						"scale":  fmt.Sprintf("%dm", s.boostMinutes),
						"decay":  0.5,
					},
				},
				"weight": s.boostScore,
			},
		},
		"score_mode": "sum",
		"boost_mode": "multiply",
	}
	
	query["query"] = map[string]interface{}{
		"function_score": functionScore,
	}
	
	// Add sorting
	if req.SortBy != "" {
		query["sort"] = s.buildSort(req)
	}
	
	return query
}

func (s *Searcher) buildBaseQuery(req *models.SearchRequest) map[string]interface{} {
	must := []map[string]interface{}{}
	filter := []map[string]interface{}{}
	
	// Text search on name and description (with Vietnamese no accent)
	if req.Query != "" {
		must = append(must, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  req.Query,
				"fields": []string{"name^3", "name_no_accent^3", "description", "description_no_accent"},
				"type":   "best_fields",
				"fuzziness": "AUTO",
			},
		})
	}
	
	// Filter by category
	if req.CategoryID != nil {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"category_id": *req.CategoryID,
			},
		})
	}
	
	// Filter by status
	if req.Status != "" {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"status": req.Status,
			},
		})
	}
	
	// Filter by price range
	if req.MinPrice != nil || req.MaxPrice != nil {
		priceRange := make(map[string]interface{})
		if req.MinPrice != nil {
			priceRange["gte"] = *req.MinPrice
		}
		if req.MaxPrice != nil {
			priceRange["lte"] = *req.MaxPrice
		}
		filter = append(filter, map[string]interface{}{
			"range": map[string]interface{}{
				"current_price": priceRange,
			},
		})
	}
	
	boolQuery := map[string]interface{}{}
	if len(must) > 0 {
		boolQuery["must"] = must
	} else {
		boolQuery["must"] = []map[string]interface{}{
			{"match_all": map[string]interface{}{}},
		}
	}
	if len(filter) > 0 {
		boolQuery["filter"] = filter
	}
	
	return map[string]interface{}{
		"bool": boolQuery,
	}
}

func (s *Searcher) buildSort(req *models.SearchRequest) []map[string]interface{} {
	sort := []map[string]interface{}{}
	
	sortOrder := "desc"
	if req.SortOrder == "asc" {
		sortOrder = "asc"
	}
	
	switch req.SortBy {
	case "price":
		sort = append(sort, map[string]interface{}{
			"current_price": map[string]string{"order": sortOrder},
		})
	case "end_at":
		sort = append(sort, map[string]interface{}{
			"end_at": map[string]string{"order": sortOrder},
		})
	case "created_at":
		sort = append(sort, map[string]interface{}{
			"created_at": map[string]string{"order": sortOrder},
		})
	default:
		sort = append(sort, map[string]interface{}{
			"_score": map[string]string{"order": "desc"},
		})
	}
	
	return sort
}

// convertToListingResponse converts ProductESDocument to ProductListingResponse
func (s *Searcher) convertToListingResponse(doc *models.ProductESDocument) models.ProductListingResponse {
	return models.ProductListingResponse{
		ID:                doc.ID,
		Name:              doc.Name,
		ThumbnailURL:      doc.ThumbnailURL,
		CurrentPrice:      doc.CurrentPrice,
		BuyNowPrice:       doc.BuyNowPrice,
		CurrentBidderInfo: doc.CurrentBidderInfo,
		CreatedAt:         doc.CreatedAt,
		EndAt:             doc.EndAt,
		TimeRemaining:     calculateTimeRemaining(doc.EndAt),
		CurrentBidCount:   doc.CurrentBidCount,
		Status:            doc.Status,
		CategoryName:      doc.CategoryName,
	}
}

// calculateTimeRemaining calculates human-readable time remaining
func calculateTimeRemaining(endAt time.Time) string {
	now := time.Now()
	
	// If auction has ended
	if endAt.Before(now) {
		return "Đã kết thúc"
	}
	
	duration := endAt.Sub(now)
	
	// Calculate days, hours, minutes
	days := int(duration.Hours() / 24)
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60
	
	// Format output based on time remaining
	if days > 0 {
		return fmt.Sprintf("%d ngày %d giờ", days, hours)
	} else if hours > 0 {
		return fmt.Sprintf("%d giờ %d phút", hours, minutes)
	} else if minutes > 0 {
		return fmt.Sprintf("%d phút", minutes)
	} else {
		seconds := int(duration.Seconds())
		return fmt.Sprintf("%d giây", seconds)
	}
}

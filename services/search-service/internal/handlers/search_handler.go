package handlers

import (
	"search-service/internal/elasticsearch"
	"search-service/internal/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	searcher *elasticsearch.Searcher
}

func NewProductHandler(searcher *elasticsearch.Searcher) *ProductHandler {
	return &ProductHandler{
		searcher: searcher,
	}
}

func (h *ProductHandler) SearchProducts(c *fiber.Ctx) error {
	ctx := c.Context()

	var req models.SearchRequest
	req.Query = c.Query("query", "")
	req.Status = c.Query("status", "")
	req.SortBy = c.Query("sort_by", "")
	req.SortOrder = c.Query("sort_order", "desc")

	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
		if err == nil {
			req.CategoryID = &categoryID
		}
	}

	if minPriceStr := c.Query("min_price"); minPriceStr != "" {
		minPrice, err := strconv.ParseFloat(minPriceStr, 64)
		if err == nil {
			req.MinPrice = &minPrice
		}
	}

	if maxPriceStr := c.Query("max_price"); maxPriceStr != "" {
		maxPrice, err := strconv.ParseFloat(maxPriceStr, 64)
		if err == nil {
			req.MaxPrice = &maxPrice
		}
	}

	if pageStr := c.Query("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err == nil {
			req.Page = page
		}
	} else {
		req.Page = 1
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		pageSize, err := strconv.Atoi(pageSizeStr)
		if err == nil {
			req.PageSize = pageSize
		}
	} else {
		req.PageSize = 20
	}

	result, err := h.searcher.SearchProducts(ctx, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":  "Failed to search products",
			"detail": err.Error(),
		})
	}

	return c.JSON(result)
}

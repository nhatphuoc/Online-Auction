package handlers

import (
	"auto-bidding-service/internal/models"
	"auto-bidding-service/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// AutoBidHandler xử lý HTTP requests cho auto-bidding
type AutoBidHandler struct {
	service *service.AutoBidService
}

// NewAutoBidHandler tạo handler mới
func NewAutoBidHandler(service *service.AutoBidService) *AutoBidHandler {
	return &AutoBidHandler{service: service}
}

// CreateAutoBid godoc
// @Summary Tạo auto-bid mới
// @Description Tạo một lệnh đấu giá tự động cho sản phẩm
// @Tags auto-bidding
// @Accept json
// @Produce json
// @Param request body models.CreateAutoBidRequest true "Auto-bid request"
// @Param X-User-ID header int true "User ID từ JWT"
// @Param X-User-Token header string true "JWT token"
// @Success 200 {object} models.AutoBidResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/auto-bids [post]
// @Security BearerAuth
func (h *AutoBidHandler) CreateAutoBid(c *fiber.Ctx) error {
	// Lấy user ID từ header (đã được API Gateway inject)
	userIDStr := c.Get("X-User-ID")
	if userIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "User ID not found in request",
		})
	}

	bidderID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid user ID",
		})
	}

	// Lấy JWT token
	userToken := c.Get("X-User-Token")
	if userToken == "" {
		userToken = c.Get("Authorization")
	}

	// Parse request body
	var req models.CreateAutoBidRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Validate
	if req.ProductID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid product ID",
		})
	}

	if req.MaxAmount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Max amount must be greater than 0",
		})
	}

	// Tạo auto-bid
	autoBid, err := h.service.CreateAutoBid(c.Context(), bidderID, req.ProductID, req.MaxAmount, userToken)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Auto-bid created successfully",
		"data":    autoBid,
	})
}

// TriggerAutoBidding godoc
// @Summary Trigger auto-bidding khi có bid mới
// @Description Được gọi bởi bidding-service khi có bid mới để trigger auto-bidding
// @Tags auto-bidding
// @Accept json
// @Produce json
// @Param request body models.TriggerAutoBidRequest true "Trigger request"
// @Param X-Internal-Key header string true "Internal service key"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/auto-bids/trigger [post]
func (h *AutoBidHandler) TriggerAutoBidding(c *fiber.Ctx) error {
	// Kiểm tra internal key (để bảo mật, chỉ cho phép bidding-service gọi)
	internalKey := c.Get("X-Internal-Key")
	// TODO: Validate internal key với environment variable

	var req models.TriggerAutoBidRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Trigger auto-bidding trong background
	// Lấy user token nếu có
	userToken := c.Get("X-User-Token")
	
	go h.service.TriggerAutoBidding(
		c.Context(),
		req.ProductID,
		req.CurrentPrice,
		req.BidIncrement,
		req.NewBidderID,
		req.NewBidAmount,
		userToken,
	)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Auto-bidding triggered",
	})
}

// GetMyAutoBids godoc
// @Summary Lấy danh sách auto-bid của user
// @Description Lấy tất cả auto-bid của user hiện tại
// @Tags auto-bidding
// @Produce json
// @Param X-User-ID header int true "User ID từ JWT"
// @Success 200 {array} models.AutoBidResponse
// @Failure 400 {object} map[string]interface{}
// @Router /api/auto-bids/my [get]
// @Security BearerAuth
func (h *AutoBidHandler) GetMyAutoBids(c *fiber.Ctx) error {
	userIDStr := c.Get("X-User-ID")
	if userIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "User ID not found in request",
		})
	}

	bidderID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid user ID",
		})
	}

	autoBids, err := h.service.GetAutoBidsByBidder(c.Context(), bidderID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get auto-bids",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    autoBids,
	})
}

// GetAutoBidByID godoc
// @Summary Lấy thông tin auto-bid theo ID
// @Description Lấy chi tiết một auto-bid
// @Tags auto-bidding
// @Produce json
// @Param id path int true "Auto-bid ID"
// @Success 200 {object} models.AutoBidResponse
// @Failure 404 {object} map[string]interface{}
// @Router /api/auto-bids/{id} [get]
// @Security BearerAuth
func (h *AutoBidHandler) GetAutoBidByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid auto-bid ID",
		})
	}

	autoBid, err := h.service.GetAutoBidByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    autoBid,
	})
}

// CancelAutoBid godoc
// @Summary Hủy auto-bid
// @Description Hủy một lệnh đấu giá tự động
// @Tags auto-bidding
// @Produce json
// @Param id path int true "Auto-bid ID"
// @Param X-User-ID header int true "User ID từ JWT"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/auto-bids/{id}/cancel [post]
// @Security BearerAuth
func (h *AutoBidHandler) CancelAutoBid(c *fiber.Ctx) error {
	userIDStr := c.Get("X-User-ID")
	if userIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "User ID not found in request",
		})
	}

	bidderID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid user ID",
		})
	}

	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid auto-bid ID",
		})
	}

	if err := h.service.CancelAutoBid(c.Context(), id, bidderID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Auto-bid cancelled successfully",
	})
}

package handlers

import (
	"order_service/internal/models"
	"order_service/internal/service"
	"order_service/internal/utils"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	service   service.OrderService
	validator *validator.Validate
}

func NewOrderHandler(service service.OrderService) *OrderHandler {
	return &OrderHandler{
		service:   service,
		validator: validator.New(),
	}
}

// CreateOrder creates a new order after auction ends
// @Summary Create a new order
// @Description Create a new order after auction ends (typically called by auction-service)
// @Tags orders
// @Accept json
// @Produce json
// @Param order body models.CreateOrderRequest true "Order data"
// @Success 201 {object} models.Order
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	req := new(models.CreateOrderRequest)
	
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	
	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	order, err := h.service.CreateOrder(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	return c.Status(fiber.StatusCreated).JSON(order)
}

// GetOrderByID retrieves order by ID
// @Summary Get order by ID
// @Description Get order details by ID (only for buyer or seller)
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Security BearerAuth
// @Success 200 {object} models.Order
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /orders/{id} [get]
func (h *OrderHandler) GetOrderByID(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID",
		})
	}
	
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	
	order, err := h.service.GetOrderByID(id, userID)
	if err != nil {
		if err == service.ErrOrderNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Order not found",
			})
		}
		if err == service.ErrUnauthorized {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized access",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	return c.JSON(order)
}

// GetUserOrders retrieves all orders for current user
// @Summary Get user orders
// @Description Get all orders for current user (as buyer or seller)
// @Tags orders
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Security BearerAuth
// @Success 200 {object} models.OrderListResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /orders [get]
func (h *OrderHandler) GetUserOrders(c *fiber.Ctx) error {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	
	response, err := h.service.GetUserOrders(userID, page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	return c.JSON(response)
}

// UpdateOrderStatus updates order status
// @Summary Update order status
// @Description Update order status (buyer or seller depending on status)
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param status body models.UpdateOrderStatusRequest true "Status update data"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /orders/{id}/status [patch]
func (h *OrderHandler) UpdateOrderStatus(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID",
		})
	}
	
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	
	req := new(models.UpdateOrderStatusRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	
	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	err = h.service.UpdateOrderStatus(id, userID, req)
	if err != nil {
		if err == service.ErrOrderNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Order not found",
			})
		}
		if err == service.ErrUnauthorized {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized access",
			})
		}
		if err == service.ErrInvalidStatus {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid status transition",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	return c.JSON(fiber.Map{
		"message": "Order status updated successfully",
	})
}

// CancelOrder cancels an order
// @Summary Cancel order
// @Description Cancel an order (buyer or seller can cancel)
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param data body map[string]string true "Cancel reason"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /orders/{id}/cancel [post]
func (h *OrderHandler) CancelOrder(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID",
		})
	}
	
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	
	var body struct {
		Reason string `json:"reason"`
	}
	
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	
	err = h.service.CancelOrder(id, userID, body.Reason)
	if err != nil {
		if err == service.ErrOrderNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Order not found",
			})
		}
		if err == service.ErrUnauthorized {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized access",
			})
		}
		if err == service.ErrInvalidStatus {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot cancel order in current status",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	return c.JSON(fiber.Map{
		"message": "Order cancelled successfully",
	})
}

// SendMessage sends a message in order chat
// @Summary Send message
// @Description Send a message in order chat (buyer or seller)
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param message body models.SendMessageRequest true "Message data"
// @Security BearerAuth
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /orders/{id}/messages [post]
func (h *OrderHandler) SendMessage(c *fiber.Ctx) error {
	orderID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID",
		})
	}
	
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	
	req := new(models.SendMessageRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	
	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	err = h.service.SendMessage(orderID, userID, req.Message)
	if err != nil {
		if err == service.ErrOrderNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Order not found",
			})
		}
		if err == service.ErrUnauthorized {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized access",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Message sent successfully",
	})
}

// GetMessages retrieves all messages for an order
// @Summary Get order messages
// @Description Get all messages for an order chat
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Security BearerAuth
// @Success 200 {array} models.OrderMessage
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /orders/{id}/messages [get]
func (h *OrderHandler) GetMessages(c *fiber.Ctx) error {
	orderID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID",
		})
	}
	
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	
	messages, err := h.service.GetMessages(orderID, userID)
	if err != nil {
		if err == service.ErrOrderNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Order not found",
			})
		}
		if err == service.ErrUnauthorized {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized access",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	return c.JSON(messages)
}

// RateOrder rates an order
// @Summary Rate order
// @Description Rate an order (buyer rates seller or seller rates buyer)
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param rating body models.RateOrderRequest true "Rating data"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /orders/{id}/rate [post]
func (h *OrderHandler) RateOrder(c *fiber.Ctx) error {
	orderID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID",
		})
	}
	
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	
	req := new(models.RateOrderRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	
	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	err = h.service.RateOrder(orderID, userID, req)
	if err != nil {
		if err == service.ErrOrderNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Order not found",
			})
		}
		if err == service.ErrUnauthorized {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized access",
			})
		}
		if err == service.ErrOrderNotCompleted {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Order must be completed or delivered to rate",
			})
		}
		if err == service.ErrAlreadyRated {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "You have already rated this order",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	return c.JSON(fiber.Map{
		"message": "Rating submitted successfully",
	})
}

// GetRating retrieves rating for an order
// @Summary Get order rating
// @Description Get rating for an order
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} models.OrderRating
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /orders/{id}/rating [get]
func (h *OrderHandler) GetRating(c *fiber.Ctx) error {
	orderID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID",
		})
	}
	
	rating, err := h.service.GetRating(orderID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	if rating == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Rating not found",
		})
	}
	
	return c.JSON(rating)
}

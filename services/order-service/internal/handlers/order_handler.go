package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"order_service/internal/config"
	"order_service/internal/middleware"
	"order_service/internal/models"
	"order_service/internal/utils"
	"strconv"
	"sync"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Client represents a connected WebSocket client
type Client struct {
	Conn    *websocket.Conn
	UserID  int64
	OrderID int64
	Send    chan []byte
}

// Hub manages WebSocket connections and message broadcasting
type Hub struct {
	clients    map[int64]map[*Client]bool // orderID -> clients
	broadcast  chan *BroadcastMessage
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

type BroadcastMessage struct {
	OrderID int64
	Message []byte
}

var hub *Hub

func init() {
	hub = &Hub{
		clients:    make(map[int64]map[*Client]bool),
		broadcast:  make(chan *BroadcastMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	go hub.Run()
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if _, ok := h.clients[client.OrderID]; !ok {
				h.clients[client.OrderID] = make(map[*Client]bool)
			}
			h.clients[client.OrderID][client] = true
			h.mu.Unlock()
			slog.Info("Client registered", "userID", client.UserID, "orderID", client.OrderID)

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.clients[client.OrderID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.Send)
					if len(clients) == 0 {
						delete(h.clients, client.OrderID)
					}
				}
			}
			h.mu.Unlock()
			slog.Info("Client unregistered", "userID", client.UserID, "orderID", client.OrderID)

		case message := <-h.broadcast:
			h.mu.RLock()
			clients := h.clients[message.OrderID]
			h.mu.RUnlock()

			for client := range clients {
				select {
				case client.Send <- message.Message:
				default:
					close(client.Send)
					h.mu.Lock()
					delete(h.clients[message.OrderID], client)
					h.mu.Unlock()
				}
			}
		}
	}
}

type OrderHandler struct {
	db        *pg.DB
	validator *validator.Validate
	cfg       *config.Config
}

func NewOrderHandler(db *pg.DB, cfg *config.Config) *OrderHandler {
	return &OrderHandler{
		db:        db,
		validator: validator.New(),
		cfg:       cfg,
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
	ctx := context.Background()
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

	// Create order
	order := &models.Order{
		AuctionID:  req.AuctionID,
		WinnerID:   req.WinnerID,
		SellerID:   req.SellerID,
		FinalPrice: req.FinalPrice,
		Status:     models.OrderStatusPendingPayment,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	_, err := h.db.ModelContext(ctx, order).Insert()
	if err != nil {
		slog.Error("Failed to create order", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create order",
		})
	}

	// Create rating record
	rating := &models.OrderRating{
		OrderID:   order.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err = h.db.ModelContext(ctx, rating).Insert()
	if err != nil {
		slog.Error("Failed to create rating record", "error", err)
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
// @Router /orders/{id} [get]
func (h *OrderHandler) GetOrderByID(c *fiber.Ctx) error {
	ctx := context.Background()
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

	// Get order
	order := &models.Order{}
	err = h.db.ModelContext(ctx, order).
		Where("id = ?", id).
		Relation("Rating").
		Select()

	if err != nil {
		if err == pg.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Order not found",
			})
		}
		slog.Error("Failed to get order", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get order",
		})
	}

	// Check if user is buyer or seller
	if order.WinnerID != userID && order.SellerID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied",
		})
	}

	return c.JSON(order)
}

// GetUserOrders retrieves all orders for current user (as buyer or seller)
// @Summary Get user orders
// @Description Get all orders for current user (as buyer or seller)
// @Tags orders
// @Accept json
// @Produce json
// @Param role query string false "Filter by role: buyer or seller"
// @Param status query string false "Filter by status"
// @Security BearerAuth
// @Success 200 {array} models.Order
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /orders [get]
func (h *OrderHandler) GetUserOrders(c *fiber.Ctx) error {
	ctx := context.Background()
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	role := c.Query("role")
	status := c.Query("status")

	query := h.db.ModelContext(ctx, &[]models.Order{}).Relation("Rating")

	// Filter by role
	if role == "buyer" {
		query = query.Where("winner_id = ?", userID)
	} else if role == "seller" {
		query = query.Where("seller_id = ?", userID)
	} else {
		// Get all orders where user is either buyer or seller
		query = query.WhereGroup(func(q *pg.Query) (*pg.Query, error) {
			q = q.WhereOr("winner_id = ?", userID).
				WhereOr("seller_id = ?", userID)
			return q, nil
		})
	}

	// Filter by status
	if status != "" {
		query = query.Where("status = ?", status)
	}

	orders := []models.Order{}
	err = query.Order("created_at DESC").Select(&orders)
	if err != nil {
		slog.Error("Failed to get orders", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get orders",
		})
	}

	return c.JSON(orders)
}

// PayOrder handles payment for order
// @Summary Pay for order
// @Description Buyer pays for the order
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param payment body models.PaymentRequest true "Payment data"
// @Security BearerAuth
// @Success 200 {object} models.Order
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /orders/{id}/pay [post]
func (h *OrderHandler) PayOrder(c *fiber.Ctx) error {
	ctx := context.Background()
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

	req := new(models.PaymentRequest)
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

	// Get order
	order := &models.Order{}
	err = h.db.ModelContext(ctx, order).Where("id = ?", id).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Order not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get order",
		})
	}

	// Check if user is buyer
	if order.WinnerID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Only buyer can pay for order",
		})
	}

	// Check order status
	if order.Status != models.OrderStatusPendingPayment {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Cannot pay order with status: %s", order.Status),
		})
	}

	// Update order with payment info (Mock payment - always successful)
	now := time.Now()
	order.PaymentMethod = req.PaymentMethod
	order.PaymentProof = req.PaymentProof
	order.Status = models.OrderStatusPaid
	order.PaidAt = &now
	order.UpdatedAt = now

	_, err = h.db.ModelContext(ctx, order).
		Column("payment_method", "payment_proof", "status", "paid_at", "updated_at").
		WherePK().
		Update()

	if err != nil {
		slog.Error("Failed to update order", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update order",
		})
	}

	return c.JSON(order)
}

// ProvideShippingAddress handles providing shipping address
// @Summary Provide shipping address
// @Description Buyer provides shipping address
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param address body models.ShippingAddressRequest true "Shipping address data"
// @Security BearerAuth
// @Success 200 {object} models.Order
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /orders/{id}/shipping-address [post]
func (h *OrderHandler) ProvideShippingAddress(c *fiber.Ctx) error {
	ctx := context.Background()
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

	req := new(models.ShippingAddressRequest)
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

	// Get order
	order := &models.Order{}
	err = h.db.ModelContext(ctx, order).Where("id = ?", id).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Order not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get order",
		})
	}

	// Check if user is buyer
	if order.WinnerID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Only buyer can provide shipping address",
		})
	}

	// Check order status (must be paid)
	if order.Status != models.OrderStatusPaid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Cannot provide address for order with status: %s", order.Status),
		})
	}

	// Update order
	order.ShippingAddress = req.ShippingAddress
	order.ShippingPhone = req.ShippingPhone
	order.Status = models.OrderStatusAddressProvided
	order.UpdatedAt = time.Now()

	_, err = h.db.ModelContext(ctx, order).
		Column("shipping_address", "shipping_phone", "status", "updated_at").
		WherePK().
		Update()

	if err != nil {
		slog.Error("Failed to update order", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update order",
		})
	}

	return c.JSON(order)
}

// SendShippingInvoice handles sending shipping invoice
// @Summary Send shipping invoice
// @Description Seller sends shipping invoice and tracking number
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param invoice body models.ShippingInvoiceRequest true "Shipping invoice data"
// @Security BearerAuth
// @Success 200 {object} models.Order
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /orders/{id}/shipping-invoice [post]
func (h *OrderHandler) SendShippingInvoice(c *fiber.Ctx) error {
	ctx := context.Background()
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

	req := new(models.ShippingInvoiceRequest)
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

	// Get order
	order := &models.Order{}
	err = h.db.ModelContext(ctx, order).Where("id = ?", id).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Order not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get order",
		})
	}

	// Check if user is seller
	if order.SellerID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Only seller can send shipping invoice",
		})
	}

	// Check order status (must have address provided)
	if order.Status != models.OrderStatusAddressProvided {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Cannot send invoice for order with status: %s", order.Status),
		})
	}

	// Update order
	order.TrackingNumber = req.TrackingNumber
	order.ShippingInvoice = req.ShippingInvoice
	order.Status = models.OrderStatusShipping
	order.UpdatedAt = time.Now()

	_, err = h.db.ModelContext(ctx, order).
		Column("tracking_number", "shipping_invoice", "status", "updated_at").
		WherePK().
		Update()

	if err != nil {
		slog.Error("Failed to update order", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update order",
		})
	}

	return c.JSON(order)
}

// ConfirmDelivery handles buyer confirming delivery
// @Summary Confirm delivery
// @Description Buyer confirms that they received the product
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Security BearerAuth
// @Success 200 {object} models.Order
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /orders/{id}/confirm-delivery [post]
func (h *OrderHandler) ConfirmDelivery(c *fiber.Ctx) error {
	ctx := context.Background()
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

	// Get order
	order := &models.Order{}
	err = h.db.ModelContext(ctx, order).Where("id = ?", id).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Order not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get order",
		})
	}

	// Check if user is buyer
	if order.WinnerID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Only buyer can confirm delivery",
		})
	}

	// Check order status
	if order.Status != models.OrderStatusShipping {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Cannot confirm delivery for order with status: %s", order.Status),
		})
	}

	// Update order
	now := time.Now()
	order.Status = models.OrderStatusDelivered
	order.DeliveredAt = &now
	order.UpdatedAt = now

	_, err = h.db.ModelContext(ctx, order).
		Column("status", "delivered_at", "updated_at").
		WherePK().
		Update()

	if err != nil {
		slog.Error("Failed to update order", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update order",
		})
	}

	return c.JSON(order)
}

// CancelOrder handles canceling order
// @Summary Cancel order
// @Description Seller can cancel order at any time before completion
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param cancel body models.CancelOrderRequest true "Cancel data"
// @Security BearerAuth
// @Success 200 {object} models.Order
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /orders/{id}/cancel [post]
func (h *OrderHandler) CancelOrder(c *fiber.Ctx) error {
	ctx := context.Background()
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

	req := new(models.CancelOrderRequest)
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

	// Get order
	order := &models.Order{}
	err = h.db.ModelContext(ctx, order).Where("id = ?", id).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Order not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get order",
		})
	}

	// Check if user is seller
	if order.SellerID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Only seller can cancel order",
		})
	}

	// Check order status (cannot cancel completed or already cancelled orders)
	if order.Status == models.OrderStatusCompleted || order.Status == models.OrderStatusCancelled {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Cannot cancel order with status: %s", order.Status),
		})
	}

	// Update order
	now := time.Now()
	order.Status = models.OrderStatusCancelled
	order.CancelReason = req.CancelReason
	order.CancelledAt = &now
	order.UpdatedAt = now

	_, err = h.db.ModelContext(ctx, order).
		Column("status", "cancel_reason", "cancelled_at", "updated_at").
		WherePK().
		Update()

	if err != nil {
		slog.Error("Failed to update order", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update order",
		})
	}

	// Auto-rate buyer with -1 when seller cancels
	rating := &models.OrderRating{}
	err = h.db.ModelContext(ctx, rating).Where("order_id = ?", order.ID).Select()
	if err == nil {
		negativeOne := -1
		now := time.Now()
		rating.SellerRating = &negativeOne
		rating.SellerComment = fmt.Sprintf("Order cancelled by seller. Reason: %s", req.CancelReason)
		rating.SellerRatedAt = &now
		rating.UpdatedAt = now

		_, err = h.db.ModelContext(ctx, rating).
			Column("seller_rating", "seller_comment", "seller_rated_at", "updated_at").
			WherePK().
			Update()

		if err != nil {
			slog.Error("Failed to update rating", "error", err)
		} else {
			// Update buyer's rating stats
			h.updateUserRating(ctx, order.WinnerID, -1)
		}
	}

	return c.JSON(order)
}

// SendMessage sends a chat message in order
// @Summary Send message
// @Description Send a chat message between buyer and seller
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param message body models.SendMessageRequest true "Message data"
// @Security BearerAuth
// @Success 201 {object} models.OrderMessage
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /orders/{id}/messages [post]
func (h *OrderHandler) SendMessage(c *fiber.Ctx) error {
	ctx := context.Background()
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

	// Get order
	order := &models.Order{}
	err = h.db.ModelContext(ctx, order).Where("id = ?", id).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Order not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get order",
		})
	}

	// Check if user is buyer or seller
	if order.WinnerID != userID && order.SellerID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied",
		})
	}

	// Create message
	message := &models.OrderMessage{
		OrderID:   id,
		SenderID:  userID,
		Message:   req.Message,
		CreatedAt: time.Now(),
	}

	_, err = h.db.ModelContext(ctx, message).Insert()
	if err != nil {
		slog.Error("Failed to save message", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save message",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(message)
}

// GetMessages retrieves all messages for an order
// @Summary Get messages
// @Description Get all chat messages for an order
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Security BearerAuth
// @Success 200 {array} models.OrderMessage
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /orders/{id}/messages [get]
func (h *OrderHandler) GetMessages(c *fiber.Ctx) error {
	ctx := context.Background()
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

	limit := c.QueryInt("limit", 50)
	offset := c.QueryInt("offset", 0)

	// Get order to check permissions
	order := &models.Order{}
	err = h.db.ModelContext(ctx, order).Where("id = ?", id).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Order not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get order",
		})
	}

	// Check if user is buyer or seller
	if order.WinnerID != userID && order.SellerID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied",
		})
	}

	// Get messages
	messages := []models.OrderMessage{}
	err = h.db.ModelContext(ctx, &messages).
		Where("order_id = ?", id).
		Order("created_at ASC").
		Limit(limit).
		Offset(offset).
		Select()

	if err != nil && err != pg.ErrNoRows {
		slog.Error("Failed to get messages", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get messages",
		})
	}

	// Nếu không có message nào thì trả về mảng rỗng [] và status 200
	return c.JSON(messages)
}

// RateOrder rates an order
// @Summary Rate order
// @Description Buyer or seller rates the transaction
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param rating body models.RateOrderRequest true "Rating data"
// @Security BearerAuth
// @Success 200 {object} models.OrderRating
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /orders/{id}/rate [post]
func (h *OrderHandler) RateOrder(c *fiber.Ctx) error {
	ctx := context.Background()
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

	// Get order
	order := &models.Order{}
	err = h.db.ModelContext(ctx, order).Where("id = ?", id).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Order not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get order",
		})
	}

	// Check if user is buyer or seller
	isBuyer := order.WinnerID == userID
	isSeller := order.SellerID == userID

	if !isBuyer && !isSeller {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied",
		})
	}

	// Get rating record
	rating := &models.OrderRating{}
	err = h.db.ModelContext(ctx, rating).Where("order_id = ?", id).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Rating record not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get rating",
		})
	}

	now := time.Now()
	var targetUserID int64
	var oldRating *int

	// Update rating based on user role
	if isBuyer {
		// Buyer rates seller
		oldRating = rating.BuyerRating
		rating.BuyerRating = &req.Rating
		rating.BuyerComment = req.Comment
		rating.BuyerRatedAt = &now
		targetUserID = order.SellerID
	} else {
		// Seller rates buyer
		oldRating = rating.SellerRating
		rating.SellerRating = &req.Rating
		rating.SellerComment = req.Comment
		rating.SellerRatedAt = &now
		targetUserID = order.WinnerID
	}

	rating.UpdatedAt = now

	// Update rating record
	_, err = h.db.ModelContext(ctx, rating).
		Column("buyer_rating", "buyer_comment", "buyer_rated_at", "seller_rating", "seller_comment", "seller_rated_at", "updated_at").
		WherePK().
		Update()

	if err != nil {
		slog.Error("Failed to update rating", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update rating",
		})
	}

	// Update user rating stats
	if oldRating != nil {
		// Remove old rating first
		h.updateUserRating(ctx, targetUserID, -*oldRating)
	}
	// Add new rating
	h.updateUserRating(ctx, targetUserID, req.Rating)

	// Check if both parties have rated, if yes, mark order as completed
	if rating.BuyerRating != nil && rating.SellerRating != nil && order.Status == models.OrderStatusDelivered {
		now := time.Now()
		order.Status = models.OrderStatusCompleted
		order.CompletedAt = &now
		order.UpdatedAt = now

		_, err = h.db.ModelContext(ctx, order).
			Column("status", "completed_at", "updated_at").
			WherePK().
			Update()

		if err != nil {
			slog.Error("Failed to complete order", "error", err)
		}
	}

	return c.JSON(rating)
}

// GetRating retrieves rating for an order
// @Summary Get rating
// @Description Get rating information for an order
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} models.OrderRating
// @Failure 404 {object} map[string]interface{}
// @Router /orders/{id}/rating [get]
func (h *OrderHandler) GetRating(c *fiber.Ctx) error {
	ctx := context.Background()
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID",
		})
	}

	rating := &models.OrderRating{}
	err = h.db.ModelContext(ctx, rating).Where("order_id = ?", id).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Rating not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get rating",
		})
	}

	return c.JSON(rating)
}

// GetUserRating retrieves a user's rating statistics
// @Summary Get user rating
// @Description Get a user's rating statistics (total reviews and good reviews)
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /users/{id}/rating [get]
func (h *OrderHandler) GetUserRating(c *fiber.Ctx) error {
	ctx := context.Background()
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	user := &models.User{}
	err = h.db.ModelContext(ctx, user).
		Column("id", "total_number_good_reviews", "total_number_reviews").
		Where("id = ?", id).
		Select()

	if err != nil {
		if err == pg.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user rating",
		})
	}

	rating := 0.0
	if user.TotalNumberReviews > 0 {
		rating = float64(user.TotalNumberGoodReviews) / float64(user.TotalNumberReviews) * 100
	}

	return c.JSON(fiber.Map{
		"user_id":                   user.ID,
		"total_number_reviews":      user.TotalNumberReviews,
		"total_number_good_reviews": user.TotalNumberGoodReviews,
		"rating_percentage":         rating,
	})
}

// GetAllOrders retrieves all orders (admin only)
// @Summary Get all orders
// @Description Get all orders in the system (admin only)
// @Tags orders
// @Accept json
// @Produce json
// @Param status query string false "Filter by status"
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Security BearerAuth
// @Success 200 {array} models.Order
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /admin/orders [get]
func (h *OrderHandler) GetAllOrders(c *fiber.Ctx) error {
	ctx := context.Background()

	// Check if user is admin
	role := c.Locals("role")
	if role != "ROLE_ADMIN" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Admin access required",
		})
	}

	status := c.Query("status")
	limit := c.QueryInt("limit", 50)
	offset := c.QueryInt("offset", 0)

	query := h.db.ModelContext(ctx, &[]models.Order{}).Relation("Rating")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	orders := []models.Order{}
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Select(&orders)
	if err != nil {
		slog.Error("Failed to get orders", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get orders",
		})
	}

	return c.JSON(orders)
}

// updateUserRating updates user's rating statistics in database
func (h *OrderHandler) updateUserRating(ctx context.Context, userID int64, rating int) error {
	user := &models.User{}
	err := h.db.ModelContext(ctx, user).Where("id = ?", userID).Select()
	if err != nil {
		slog.Error("Failed to get user for rating update", "error", err, "userID", userID)
		return err
	}

	// Update counts
	user.TotalNumberReviews++
	if rating == 1 {
		user.TotalNumberGoodReviews++
	}

	_, err = h.db.ModelContext(ctx, user).
		Column("total_number_reviews", "total_number_good_reviews").
		WherePK().
		Update()

	if err != nil {
		slog.Error("Failed to update user rating", "error", err, "userID", userID)
		return err
	}

	return nil
}

// HandleWebSocket handles WebSocket connections for order chat
func (h *OrderHandler) HandleWebSocket(c *websocket.Conn) {
	// Get query params
	orderIDStr := c.Query("orderId")
	tokenString := c.Query("X-User-Token")
	internalJWT := c.Query("X-Internal-JWT")
	slog.Error("We-------")

	if orderIDStr == "" || internalJWT == "" {
		slog.Error("Missing required parameters")
		c.Close()
		return
	}

	// Verify internal JWT
	ok, err := middleware.VerifyInternalJWT(
		h.cfg,
		internalJWT,
		h.cfg.OrderServiceName,
	)

	if err != nil || !ok {
		slog.Error("Invalid Internal JWT", "error", err)
		c.Close()
		return
	} else {
		slog.Error("Internal JWT verified")
	}

	if tokenString == "" {
		slog.Error("Missing X-User-Token header")
		c.Close()
		return
	}

	token, _, err := new(jwt.Parser).ParseUnverified(
		tokenString,
		jwt.MapClaims{},
	)
	if err != nil {
		slog.Error("Invalid token format")
		c.Close()
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		slog.Error("Invalid token claims")
		c.Close()
		return
	}

	// Check type == "access"
	if t, ok := claims["type"].(string); !ok || t != "access" {
		slog.Error("Token is not access token")
		c.Close()
		return
	}

	// Get userID
	var userID int64
	if sub, ok := claims["sub"].(float64); ok {
		userID = int64(sub)
	} else if subStr, ok := claims["sub"].(string); ok {
		if parsed, err := strconv.ParseInt(subStr, 10, 64); err == nil {
			userID = parsed
		}
	}

	if userID == 0 {
		slog.Error("Invalid user ID in token")
		c.Close()
		return
	}

	var orderID int64
	fmt.Sscanf(orderIDStr, "%d", &orderID)

	// Verify user has access to this order
	ctx := context.Background()
	order := &models.Order{}
	err = h.db.ModelContext(ctx, order).Where("id = ?", orderID).Select()
	if err != nil {
		slog.Error("Order not found", "orderID", orderID)
		c.Close()
		return
	}

	if order.WinnerID != userID && order.SellerID != userID {
		slog.Error("User does not have access to this order", "userID", userID, "orderID", orderID)
		c.Close()
		return
	}

	client := &Client{
		Conn:    c,
		UserID:  userID,
		OrderID: orderID,
		Send:    make(chan []byte, 256),
	}

	hub.register <- client

	// Start goroutines for reading and writing
	go h.writePump(client)
	h.readPump(client)

	return
}

// readPump reads messages from WebSocket connection
func (h *OrderHandler) readPump(client *Client) {
	defer func() {
		hub.unregister <- client
		client.Conn.Close()
	}()

	for {
		var wsMsg models.WebSocketMessage
		err := client.Conn.ReadJSON(&wsMsg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("WebSocket error", "error", err)
			}
			break
		}

		// Handle different message types
		switch wsMsg.Type {
		case "message":
			// Save message to database
			message := &models.OrderMessage{
				OrderID:   client.OrderID,
				SenderID:  client.UserID,
				Message:   wsMsg.Content,
				CreatedAt: time.Now(),
			}

			ctx := context.Background()
			_, err := h.db.ModelContext(ctx, message).Insert()
			if err != nil {
				slog.Error("Failed to save message", "error", err)
				continue
			}

			// Broadcast to all clients connected to this order
			responseMsg := models.WebSocketMessage{
				Type:    "message",
				OrderID: client.OrderID,
				Data:    message,
			}

			msgBytes, _ := json.Marshal(responseMsg)
			hub.broadcast <- &BroadcastMessage{
				OrderID: client.OrderID,
				Message: msgBytes,
			}

		case "typing":
			// Broadcast typing indicator
			typingMsg := models.WebSocketMessage{
				Type:    "typing",
				OrderID: client.OrderID,
				Data: fiber.Map{
					"userId": client.UserID,
				},
			}
			msgBytes, _ := json.Marshal(typingMsg)
			hub.broadcast <- &BroadcastMessage{
				OrderID: client.OrderID,
				Message: msgBytes,
			}
		}
	}
}

// writePump writes messages to WebSocket connection
func (h *OrderHandler) writePump(client *Client) {
	defer func() {
		client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := client.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				return
			}
		}
	}
}

package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"order_service/internal/config"
	"order_service/internal/middleware"
	"order_service/internal/models"
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

var locVN = mustLoad("Asia/Ho_Chi_Minh")

func mustLoad(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}
	return loc
}

func FixedTimeNow() time.Time {
	nowVN := time.Now().In(locVN)
	return time.Date(
		nowVN.Year(), nowVN.Month(), nowVN.Day(),
		nowVN.Hour(), nowVN.Minute(), nowVN.Second(),
		nowVN.Nanosecond(),
		time.UTC,
	)
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

	// Create order using raw query
	var orderID int64
	now := FixedTimeNow()
	query := `INSERT INTO orders (auction_id, winner_id, seller_id, final_price, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING id`
	_, err := h.db.QueryOneContext(ctx, pg.Scan(&orderID), query,
		req.AuctionID, req.WinnerID, req.SellerID, req.FinalPrice, models.OrderStatusPendingPayment, now, now)
	if err != nil {
		slog.Error("Failed to create order", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create order",
		})
	}

	// Create rating record
	ratingQuery := `INSERT INTO order_ratings (order_id, created_at, updated_at) VALUES (?, ?, ?)`
	_, err = h.db.ExecContext(ctx, ratingQuery, orderID, now, now)
	if err != nil {
		slog.Error("Failed to create rating record", "error", err)
	}

	// Return created order
	order := &models.Order{
		ID:         orderID,
		AuctionID:  req.AuctionID,
		WinnerID:   req.WinnerID,
		SellerID:   req.SellerID,
		FinalPrice: req.FinalPrice,
		Status:     models.OrderStatusPendingPayment,
		CreatedAt:  now,
		UpdatedAt:  now,
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

	userID := c.Locals("userID").(string)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Get order using raw query with user names
	var order models.Order
	orderQuery := `SELECT o.id, o.auction_id, o.winner_id, o.seller_id, o.final_price, o.status, o.payment_method, o.payment_proof, 
		o.shipping_address, o.shipping_phone, o.tracking_number, o.shipping_invoice, o.paid_at, o.delivered_at, 
		o.completed_at, o.cancelled_at, o.cancel_reason, o.created_at, o.updated_at,
		buyer.full_name AS buyer_name, seller.full_name AS seller_name
		FROM orders o
		LEFT JOIN users buyer ON o.winner_id = buyer.id
		LEFT JOIN users seller ON o.seller_id = seller.id
		WHERE o.id = ?`
	_, err = h.db.QueryOneContext(ctx, &order, orderQuery, id)
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

	userIDInt64, _ := strconv.ParseInt(userID, 10, 64)

	fmt.Println("UserID:", userIDInt64, "WinnerID:", order.WinnerID, "SellerID:", order.SellerID)
	// Check if user is buyer or seller
	if order.WinnerID != userIDInt64 && order.SellerID != userIDInt64 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied",
		})
	}

	// Get rating if exists
	var rating models.OrderRating
	ratingQuery := `SELECT id, order_id, buyer_rating, buyer_comment, seller_rating, seller_comment, 
		buyer_rated_at, seller_rated_at, created_at, updated_at 
		FROM order_ratings WHERE order_id = ?`
	_, err = h.db.QueryOneContext(ctx, &rating, ratingQuery, id)
	if err == nil {
		order.Rating = &rating
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
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /orders [get]
func (h *OrderHandler) GetUserOrders(c *fiber.Ctx) error {
	fmt.Println("---------------------")
	ctx := context.Background()
	userID := c.Locals("userID").(string)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	role := c.Query("role")
	status := c.Query("status")

	// Pagination parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Build query based on role filter
	var whereClause string
	var args []interface{}

	// Filter by role (accepts both "buyer"/"seller" and "ROLE_BIDDER"/"ROLE_SELLER")
	if role == "ROLE_BIDDER" || role == "buyer" {
		whereClause = "winner_id = ?"
		args = append(args, userID)
	} else if role == "ROLE_SELLER" || role == "seller" {
		whereClause = "seller_id = ?"
		args = append(args, userID)
	} else {
		// Get all orders where user is either buyer or seller
		whereClause = "(winner_id = ? OR seller_id = ?)"
		args = append(args, userID, userID)
	}

	// Filter by status
	if status != "" {
		whereClause += " AND status = ?"
		args = append(args, status)
	}

	// Get total count
	countQuery := "SELECT COUNT(*) FROM orders WHERE " + whereClause
	var total int
	_, err := h.db.QueryOneContext(ctx, pg.Scan(&total), countQuery, args...)
	if err != nil {
		slog.Error("Failed to count orders", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to count orders",
		})
	}

	// Get orders with pagination and user names
	ordersQuery := `SELECT o.id, o.auction_id, o.winner_id, o.seller_id, o.final_price, o.status, o.payment_method, o.payment_proof, 
		o.shipping_address, o.shipping_phone, o.tracking_number, o.shipping_invoice, o.paid_at, o.delivered_at, 
		o.completed_at, o.cancelled_at, o.cancel_reason, o.created_at, o.updated_at,
		buyer.full_name AS buyer_name, seller.full_name AS seller_name
		FROM orders o
		LEFT JOIN users buyer ON o.winner_id = buyer.id
		LEFT JOIN users seller ON o.seller_id = seller.id
		WHERE ` + whereClause + ` ORDER BY o.created_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	var orders []models.Order
	_, err = h.db.QueryContext(ctx, &orders, ordersQuery, args...)
	if err != nil {
		slog.Error("Failed to get orders", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get orders",
		})
	}

	// Get ratings for each order
	for i := range orders {
		var rating models.OrderRating
		ratingQuery := `SELECT id, order_id, buyer_rating, buyer_comment, seller_rating, seller_comment, 
			buyer_rated_at, seller_rated_at, created_at, updated_at 
			FROM order_ratings WHERE order_id = ?`
		_, err = h.db.QueryOneContext(ctx, &rating, ratingQuery, orders[i].ID)
		if err == nil {
			orders[i].Rating = &rating
		}
	}

	totalPages := (total + limit - 1) / limit

	return c.JSON(fiber.Map{
		"data": orders,
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
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

	userID := c.Locals("userID").(string)
	if userID == "" {
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

	// Get order using raw query
	var order models.Order
	orderQuery := `SELECT id, auction_id, winner_id, seller_id, final_price, status, payment_method, payment_proof, 
		shipping_address, shipping_phone, tracking_number, shipping_invoice, paid_at, delivered_at, 
		completed_at, cancelled_at, cancel_reason, created_at, updated_at 
		FROM orders WHERE id = ?`
	_, err = h.db.QueryOneContext(ctx, &order, orderQuery, id)
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

	userIDInt64, _ := strconv.ParseInt(userID, 10, 64)
	// Check if user is buyer
	if order.WinnerID != userIDInt64 {
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
	now := FixedTimeNow()
	updateQuery := `UPDATE orders SET payment_method = ?, payment_proof = ?, status = ?, paid_at = ?, updated_at = ? 
		WHERE id = ?`
	_, err = h.db.ExecContext(ctx, updateQuery, req.PaymentMethod, req.PaymentProof, models.OrderStatusPaid, now, now, id)
	if err != nil {
		slog.Error("Failed to update order", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update order",
		})
	}

	// Update local object for response
	order.PaymentMethod = req.PaymentMethod
	order.PaymentProof = req.PaymentProof
	order.Status = models.OrderStatusPaid
	order.PaidAt = &now
	order.UpdatedAt = now

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

	userID := c.Locals("userID").(string)
	if userID == "" {
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

	// Get order using raw query
	var order models.Order
	orderQuery := `SELECT id, auction_id, winner_id, seller_id, final_price, status, payment_method, payment_proof, 
		shipping_address, shipping_phone, tracking_number, shipping_invoice, paid_at, delivered_at, 
		completed_at, cancelled_at, cancel_reason, created_at, updated_at 
		FROM orders WHERE id = ?`
	_, err = h.db.QueryOneContext(ctx, &order, orderQuery, id)
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

	userIDInt64, _ := strconv.ParseInt(userID, 10, 64)
	// Check if user is buyer
	if order.WinnerID != userIDInt64 {
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
	now := FixedTimeNow()
	updateQuery := `UPDATE orders SET shipping_address = ?, shipping_phone = ?, status = ?, updated_at = ? WHERE id = ?`
	_, err = h.db.ExecContext(ctx, updateQuery, req.ShippingAddress, req.ShippingPhone, models.OrderStatusAddressProvided, now, id)
	if err != nil {
		slog.Error("Failed to update order", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update order",
		})
	}

	// Update local object for response
	order.ShippingAddress = req.ShippingAddress
	order.ShippingPhone = req.ShippingPhone
	order.Status = models.OrderStatusAddressProvided
	order.UpdatedAt = now

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

	userID := c.Locals("userID").(string)
	if userID == "" {
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

	// Get order using raw query
	var order models.Order
	orderQuery := `SELECT id, auction_id, winner_id, seller_id, final_price, status, payment_method, payment_proof, 
		shipping_address, shipping_phone, tracking_number, shipping_invoice, paid_at, delivered_at, 
		completed_at, cancelled_at, cancel_reason, created_at, updated_at 
		FROM orders WHERE id = ?`
	_, err = h.db.QueryOneContext(ctx, &order, orderQuery, id)
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

	userIDInt64, _ := strconv.ParseInt(userID, 10, 64)

	// Check if user is seller
	if order.SellerID != userIDInt64 {
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
	now := FixedTimeNow()
	updateQuery := `UPDATE orders SET tracking_number = ?, shipping_invoice = ?, status = ?, updated_at = ? WHERE id = ?`
	_, err = h.db.ExecContext(ctx, updateQuery, req.TrackingNumber, req.ShippingInvoice, models.OrderStatusShipping, now, id)
	if err != nil {
		slog.Error("Failed to update order", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update order",
		})
	}

	// Update local object for response
	order.TrackingNumber = req.TrackingNumber
	order.ShippingInvoice = req.ShippingInvoice
	order.Status = models.OrderStatusShipping
	order.UpdatedAt = now

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

	userID := c.Locals("userID").(string)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Get order using raw query
	var order models.Order
	orderQuery := `SELECT id, auction_id, winner_id, seller_id, final_price, status, payment_method, payment_proof, 
		shipping_address, shipping_phone, tracking_number, shipping_invoice, paid_at, delivered_at, 
		completed_at, cancelled_at, cancel_reason, created_at, updated_at 
		FROM orders WHERE id = ?`
	_, err = h.db.QueryOneContext(ctx, &order, orderQuery, id)
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

	userIDInt64, _ := strconv.ParseInt(userID, 10, 64)
	// Check if user is buyer
	if order.WinnerID != userIDInt64 {
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
	now := FixedTimeNow()
	updateQuery := `UPDATE orders SET status = ?, delivered_at = ?, updated_at = ? WHERE id = ?`
	_, err = h.db.ExecContext(ctx, updateQuery, models.OrderStatusDelivered, now, now, id)
	if err != nil {
		slog.Error("Failed to update order", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update order",
		})
	}

	// Update local object for response
	order.Status = models.OrderStatusDelivered
	order.DeliveredAt = &now
	order.UpdatedAt = now

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

	userID := c.Locals("userID").(string)
	if userID == "" {
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

	// Get order using raw query
	var order models.Order
	orderQuery := `SELECT id, auction_id, winner_id, seller_id, final_price, status, payment_method, payment_proof, 
		shipping_address, shipping_phone, tracking_number, shipping_invoice, paid_at, delivered_at, 
		completed_at, cancelled_at, cancel_reason, created_at, updated_at 
		FROM orders WHERE id = ?`
	_, err = h.db.QueryOneContext(ctx, &order, orderQuery, id)
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

	userIDInt64, _ := strconv.ParseInt(userID, 10, 64)
	// Check if user is seller
	if order.SellerID != userIDInt64 {
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
	now := FixedTimeNow()
	updateQuery := `UPDATE orders SET status = ?, cancel_reason = ?, cancelled_at = ?, updated_at = ? WHERE id = ?`
	_, err = h.db.ExecContext(ctx, updateQuery, models.OrderStatusCancelled, req.CancelReason, now, now, id)
	if err != nil {
		slog.Error("Failed to update order", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update order",
		})
	}

	// Update local object for response
	order.Status = models.OrderStatusCancelled
	order.CancelReason = req.CancelReason
	order.CancelledAt = &now
	order.UpdatedAt = now

	// Auto-rate buyer with -1 when seller cancels (create rating record if not exists)
	var rating models.OrderRating
	ratingQuery := `SELECT id, order_id, buyer_rating, buyer_comment, seller_rating, seller_comment, 
		buyer_rated_at, seller_rated_at, created_at, updated_at 
		FROM order_ratings WHERE order_id = ?`
	_, err = h.db.QueryOneContext(ctx, &rating, ratingQuery, order.ID)

	// If rating record doesn't exist, create it
	if err == pg.ErrNoRows {
		insertQuery := `INSERT INTO order_ratings (order_id, created_at, updated_at) 
			VALUES (?, ?, ?) RETURNING id`
		_, err = h.db.QueryOneContext(ctx, pg.Scan(&rating.ID), insertQuery, order.ID, now, now)
		if err != nil {
			slog.Error("Failed to create rating record for cancelled order", "error", err)
			// Don't return error, cancellation was successful
		} else {
			rating.OrderID = order.ID
			rating.CreatedAt = now
			rating.UpdatedAt = now
		}
	} else if err != nil {
		slog.Error("Failed to get rating record", "error", err)
	}

	// If we have a rating record (existing or newly created), update it
	if rating.ID > 0 {
		negativeOne := -1
		updateRatingQuery := `UPDATE order_ratings SET seller_rating = ?, seller_comment = ?, seller_rated_at = ?, updated_at = ? WHERE id = ?`
		comment := fmt.Sprintf("Order cancelled by seller. Reason: %s", req.CancelReason)
		_, err = h.db.ExecContext(ctx, updateRatingQuery, negativeOne, comment, now, now, rating.ID)
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

	userID := c.Locals("userID").(string)
	if userID == "" {
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

	// Get order using raw query
	var order models.Order
	orderQuery := `SELECT id, auction_id, winner_id, seller_id, final_price, status 
		FROM orders WHERE id = ?`
	_, err = h.db.QueryOneContext(ctx, &order, orderQuery, id)
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

	userIDInt64, _ := strconv.ParseInt(userID, 10, 64)

	// Check if user is buyer or seller
	if order.WinnerID != userIDInt64 && order.SellerID != userIDInt64 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied",
		})
	}

	// Create message using raw query
	now := FixedTimeNow()
	var messageID int64
	insertQuery := `INSERT INTO order_messages (order_id, sender_id, message, created_at) 
		VALUES (?, ?, ?, ?) RETURNING id`
	_, err = h.db.QueryOneContext(ctx, pg.Scan(&messageID), insertQuery, id, userIDInt64, req.Message, now)
	if err != nil {
		slog.Error("Failed to save message", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save message",
		})
	}

	// Return created message
	message := &models.OrderMessage{
		ID:        messageID,
		OrderID:   id,
		SenderID:  userIDInt64,
		Message:   req.Message,
		CreatedAt: now,
	}

	return c.Status(fiber.StatusCreated).JSON(message)
}

// GetMessages retrieves all messages for an order (Chat History)
// @Summary Get chat history
// @Description Get all chat messages for an order with pagination
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /orders/product/{id}/messages [get]
func (h *OrderHandler) GetMessages(c *fiber.Ctx) error {
	ctx := context.Background()
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID",
		})
	}

	userID := c.Locals("userID").(string)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	limit := c.QueryInt("limit", 50)
	offset := c.QueryInt("offset", 0)

	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	// Get order to check permissions using raw query
	var order models.Order
	orderQuery := `SELECT id, auction_id, winner_id, seller_id FROM orders WHERE id = ?`
	_, err = h.db.QueryOneContext(ctx, &order, orderQuery, id)
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

	userIDInt64, _ := strconv.ParseInt(userID, 10, 64)
	// Check if user is buyer or seller
	if order.WinnerID != userIDInt64 && order.SellerID != userIDInt64 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied",
		})
	}

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM order_messages WHERE order_id = ?`
	_, err = h.db.QueryOneContext(ctx, pg.Scan(&total), countQuery, id)
	if err != nil {
		slog.Error("Failed to count messages", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to count messages",
		})
	}

	// Get messages using raw query - ordered DESC to get newest first, then reverse
	messagesQuery := `SELECT id, order_id, sender_id, message, created_at 
		FROM order_messages WHERE order_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`
	var messages []models.OrderMessage
	_, err = h.db.QueryContext(ctx, &messages, messagesQuery, id, limit, offset)
	if err != nil && err != pg.ErrNoRows {
		slog.Error("Failed to get messages", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get messages",
		})
	}

	// Reverse messages to show oldest first (chronological order)
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	// Return with pagination metadata
	return c.JSON(fiber.Map{
		"data": messages,
		"pagination": fiber.Map{
			"total":  total,
			"limit":  limit,
			"offset": offset,
		},
	})
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

	userID := c.Locals("userID").(string)
	if userID == "" {
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

	// Get order using raw query
	var order models.Order
	orderQuery := `SELECT id, auction_id, winner_id, seller_id, status FROM orders WHERE id = ?`
	_, err = h.db.QueryOneContext(ctx, &order, orderQuery, id)
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
	userIDInt64, _ := strconv.ParseInt(userID, 10, 64)
	isBuyer := order.WinnerID == userIDInt64
	isSeller := order.SellerID == userIDInt64

	if !isBuyer && !isSeller {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied",
		})
	}

	// Get or create rating record using raw query
	var rating models.OrderRating
	ratingQuery := `SELECT id, order_id, buyer_rating, buyer_comment, seller_rating, seller_comment, 
		buyer_rated_at, seller_rated_at, created_at, updated_at 
		FROM order_ratings WHERE order_id = ?`
	_, err = h.db.QueryOneContext(ctx, &rating, ratingQuery, id)

	// If rating record doesn't exist, create it
	if err == pg.ErrNoRows {
		now := FixedTimeNow()
		insertQuery := `INSERT INTO order_ratings (order_id, created_at, updated_at) 
			VALUES (?, ?, ?) RETURNING id`
		_, err = h.db.QueryOneContext(ctx, pg.Scan(&rating.ID), insertQuery, id, now, now)
		if err != nil {
			slog.Error("Failed to create rating record", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create rating record",
			})
		}
		rating.OrderID = id
		rating.CreatedAt = now
		rating.UpdatedAt = now
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get rating",
		})
	}

	now := FixedTimeNow()
	var targetUserID int64
	var oldRating *int

	// Update rating based on user role
	if isBuyer {
		// Buyer rates seller
		oldRating = rating.BuyerRating
		targetUserID = order.SellerID
		// Update using raw query
		updateQuery := `UPDATE order_ratings SET buyer_rating = ?, buyer_comment = ?, buyer_rated_at = ?, updated_at = ? WHERE id = ?`
		_, err = h.db.ExecContext(ctx, updateQuery, req.Rating, req.Comment, now, now, rating.ID)
	} else {
		// Seller rates buyer
		oldRating = rating.SellerRating
		targetUserID = order.WinnerID
		// Update using raw query
		updateQuery := `UPDATE order_ratings SET seller_rating = ?, seller_comment = ?, seller_rated_at = ?, updated_at = ? WHERE id = ?`
		_, err = h.db.ExecContext(ctx, updateQuery, req.Rating, req.Comment, now, now, rating.ID)
	}

	if err != nil {
		slog.Error("Failed to update rating", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update rating",
		})
	}

	// Update local object for response
	if isBuyer {
		rating.BuyerRating = &req.Rating
		rating.BuyerComment = req.Comment
		rating.BuyerRatedAt = &now
	} else {
		rating.SellerRating = &req.Rating
		rating.SellerComment = req.Comment
		rating.SellerRatedAt = &now
	}
	rating.UpdatedAt = now

	// Update user rating stats
	if oldRating != nil {
		// Remove old rating first
		if err := h.removeUserRating(ctx, targetUserID, *oldRating); err != nil {
			slog.Error("Failed to remove old rating", "error", err)
		}
	}
	// Add new rating
	if err := h.addUserRating(ctx, targetUserID, req.Rating); err != nil {
		slog.Error("Failed to add new rating", "error", err)
	}

	// Check if both parties have rated, if yes, mark order as completed
	if rating.BuyerRating != nil && rating.SellerRating != nil && order.Status == models.OrderStatusDelivered {
		completeQuery := `UPDATE orders SET status = ?, completed_at = ?, updated_at = ? WHERE id = ?`
		_, err = h.db.ExecContext(ctx, completeQuery, models.OrderStatusCompleted, now, now, id)
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

	var rating models.OrderRating
	ratingQuery := `SELECT id, order_id, buyer_rating, buyer_comment, seller_rating, seller_comment, 
		buyer_rated_at, seller_rated_at, created_at, updated_at 
		FROM order_ratings WHERE order_id = ?`
	_, err = h.db.QueryOneContext(ctx, &rating, ratingQuery, id)
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

	var user models.User
	userQuery := `SELECT id, total_number_good_reviews, total_number_reviews FROM users WHERE id = ?`
	_, err = h.db.QueryOneContext(ctx, &user, userQuery, id)
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

	// Build query
	var whereClause string
	var args []interface{}

	ordersQuery := `SELECT id, auction_id, winner_id, seller_id, final_price, status, payment_method, payment_proof, 
		shipping_address, shipping_phone, tracking_number, shipping_invoice, paid_at, delivered_at, 
		completed_at, cancelled_at, cancel_reason, created_at, updated_at 
		FROM orders`

	if status != "" {
		whereClause = " WHERE status = ?"
		args = append(args, status)
	}

	ordersQuery += whereClause + " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	var orders []models.Order
	_, err := h.db.QueryContext(ctx, &orders, ordersQuery, args...)
	if err != nil {
		slog.Error("Failed to get orders", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get orders",
		})
	}

	// Get ratings for each order
	for i := range orders {
		var rating models.OrderRating
		ratingQuery := `SELECT id, order_id, buyer_rating, buyer_comment, seller_rating, seller_comment, 
			buyer_rated_at, seller_rated_at, created_at, updated_at 
			FROM order_ratings WHERE order_id = ?`
		_, err = h.db.QueryOneContext(ctx, &rating, ratingQuery, orders[i].ID)
		if err == nil {
			orders[i].Rating = &rating
		}
	}

	return c.JSON(orders)
}

// updateUserRating updates user's rating statistics in database
// addUserRating adds a new rating to user's stats
func (h *OrderHandler) addUserRating(ctx context.Context, userID int64, rating int) error {
	var user models.User
	userQuery := `SELECT id, total_number_reviews, total_number_good_reviews FROM users WHERE id = ?`
	_, err := h.db.QueryOneContext(ctx, &user, userQuery, userID)
	if err != nil {
		slog.Error("Failed to get user for rating update", "error", err, "userID", userID)
		return err
	}

	// Update counts
	totalReviews := user.TotalNumberReviews + 1
	totalGoodReviews := user.TotalNumberGoodReviews
	if rating == 1 {
		totalGoodReviews++
	}

	updateQuery := `UPDATE users SET total_number_reviews = ?, total_number_good_reviews = ? WHERE id = ?`
	_, err = h.db.ExecContext(ctx, updateQuery, totalReviews, totalGoodReviews, userID)

	if err != nil {
		slog.Error("Failed to update user rating", "error", err, "userID", userID)
		return err
	}

	return nil
}

// removeUserRating removes an existing rating from user's stats
func (h *OrderHandler) removeUserRating(ctx context.Context, userID int64, rating int) error {
	var user models.User
	userQuery := `SELECT id, total_number_reviews, total_number_good_reviews FROM users WHERE id = ?`
	_, err := h.db.QueryOneContext(ctx, &user, userQuery, userID)
	if err != nil {
		slog.Error("Failed to get user for rating removal", "error", err, "userID", userID)
		return err
	}

	// Decrease counts (ensure they don't go negative)
	totalReviews := user.TotalNumberReviews - 1
	if totalReviews < 0 {
		totalReviews = 0
	}

	totalGoodReviews := user.TotalNumberGoodReviews
	if rating == 1 && totalGoodReviews > 0 {
		totalGoodReviews--
	}

	updateQuery := `UPDATE users SET total_number_reviews = ?, total_number_good_reviews = ? WHERE id = ?`
	_, err = h.db.ExecContext(ctx, updateQuery, totalReviews, totalGoodReviews, userID)

	if err != nil {
		slog.Error("Failed to remove user rating", "error", err, "userID", userID)
		return err
	}

	return nil
}

// updateUserRating is a legacy wrapper - prefer addUserRating
func (h *OrderHandler) updateUserRating(ctx context.Context, userID int64, rating int) error {
	return h.addUserRating(ctx, userID, rating)
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
	fmt.Println(ok)

	fmt.Println(err)

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
	var order models.Order
	orderQuery := `SELECT id, winner_id, seller_id FROM orders WHERE id = ?`
	_, err = h.db.QueryOneContext(ctx, &order, orderQuery, orderID)
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
			// Save message to database using raw query
			now := time.Now().UTC()
			var messageID int64
			ctx := context.Background()
			insertQuery := `INSERT INTO order_messages (order_id, sender_id, message, created_at) 
				VALUES (?, ?, ?, ?) RETURNING id`
			_, err := h.db.QueryOneContext(ctx, pg.Scan(&messageID), insertQuery,
				client.OrderID, client.UserID, wsMsg.Content, now)
			if err != nil {
				slog.Error("Failed to save message", "error", err)
				continue
			}

			message := &models.OrderMessage{
				ID:        messageID,
				OrderID:   client.OrderID,
				SenderID:  client.UserID,
				Message:   wsMsg.Content,
				CreatedAt: now,
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

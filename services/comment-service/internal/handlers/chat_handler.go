package handlers

import (
	"comment_service/internal/config"
	"comment_service/internal/middleware"
	"comment_service/internal/models"
	"comment_service/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"sync"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Client represents a connected WebSocket client
type Client struct {
	Conn      *websocket.Conn
	UserID    int
	ProductID int
	Send      chan []byte
}

// Hub manages WebSocket connections and message broadcasting
type Hub struct {
	clients    map[int]map[*Client]bool // productID -> clients
	broadcast  chan *BroadcastMessage
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

type BroadcastMessage struct {
	ProductID int
	Message   []byte
}

var hub *Hub

func init() {
	hub = &Hub{
		clients:    make(map[int]map[*Client]bool),
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
			if _, ok := h.clients[client.ProductID]; !ok {
				h.clients[client.ProductID] = make(map[*Client]bool)
			}
			h.clients[client.ProductID][client] = true
			h.mu.Unlock()
			slog.Info("Client registered", "userID", client.UserID, "productID", client.ProductID)

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.clients[client.ProductID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.Send)
					if len(clients) == 0 {
						delete(h.clients, client.ProductID)
					}
				}
			}
			h.mu.Unlock()
			slog.Info("Client unregistered", "userID", client.UserID, "productID", client.ProductID)

		case message := <-h.broadcast:
			h.mu.RLock()
			clients := h.clients[message.ProductID]
			h.mu.RUnlock()

			for client := range clients {
				select {
				case client.Send <- message.Message:
				default:
					close(client.Send)
					h.mu.Lock()
					delete(h.clients[message.ProductID], client)
					h.mu.Unlock()
				}
			}
		}
	}
}

type CommentHandler struct {
	db *pg.DB
}

func NewCommentHandler(db *pg.DB) *CommentHandler {
	return &CommentHandler{db: db}
}

// GetProductComments godoc
// @Summary Lấy lịch sử bình luận
// @Description Lấy danh sách bình luận của sản phẩm
// @Tags Comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param productId path int true "Product ID"
// @Param limit query int false "Số lượng bình luận" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} models.CommentResponse
// @Failure 500 {object} map[string]interface{}
// @Router /api/comments/products/{productId} [get]
func (h *CommentHandler) GetProductComments(c *fiber.Ctx) error {
	ctx := context.Background()
	productID, err := c.ParamsInt("productId")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid product ID")
	}
	limit := c.QueryInt("limit", 50)
	offset := c.QueryInt("offset", 0)

	var comments []models.Comment
	err = h.db.ModelContext(ctx, &comments).
		Where("product_id = ?", productID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Select()

	if err != nil {
		slog.Error("Failed to get comments", "error", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Lỗi lấy bình luận")
	}

	return c.JSON(comments)
}

// HandleWebSocket handles WebSocket connections
func (h *CommentHandler) HandleWebSocket(c *websocket.Conn) {
	// Get query params
	productIDStr := c.Query("productId")
	tokenString := c.Query("X-User-Token")
	internalJWT := c.Query("X-Internal-JWT")

	if productIDStr == "" || internalJWT == "" {
		slog.Error("Missing required parameters")
		c.Close()
		return
	}

	// You may want to load config here if needed, or inject via handler struct
	cfg := config.LoadConfig()
	ok, err := middleware.VerifyInternalJWT(
		cfg,
		internalJWT,
		cfg.CommentServiceName,
	)
	if err != nil || !ok {
		slog.Error(err.Error())
		slog.Error("Invalid Internal JWT")
		c.Close()
		return
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

	// Kiểm tra type == "access"
	if t, ok := claims["type"].(string); !ok || t != "access" {
		slog.Error("Token is not access token")
		c.Close()
		return
	}

	// Lấy userId (subject), email, role
	var userID int
	if sub, ok := claims["sub"].(float64); ok {
		userID = int(sub)
	} else if subStr, ok := claims["sub"].(string); ok {
		// Try to parse string to int
		if parsed, err := strconv.Atoi(subStr); err == nil {
			userID = parsed
		}
	}
	role := ""
	if r, ok := claims["role"].(string); ok {
		role = r
	} else if r, ok := claims["role"].(map[string]interface{}); ok {
		// Trường hợp role là object (enum)
		if name, ok := r["name"].(string); ok {
			role = name
		}
	}
	if role == "" {
		slog.Error("User role not found in token")
		c.Close()
		return
	}

	var productID int
	fmt.Sscanf(productIDStr, "%d", &productID)

	client := &Client{
		Conn:      c,
		UserID:    userID,
		ProductID: productID,
		Send:      make(chan []byte, 256),
	}

	hub.register <- client

	// Start goroutines for reading and writing
	go h.writePump(client)
	h.readPump(client)
	return
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

// readPump reads messages from WebSocket connection
func (h *CommentHandler) readPump(client *Client) {
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
		case "comment":
			// Save comment to database
			comment := &models.Comment{
				ProductID: client.ProductID,
				SenderID:  client.UserID,
				Content:   wsMsg.Content,
				CreatedAt: FixedTimeNow(),
			}

			ctx := context.Background()
			_, err := h.db.ModelContext(ctx, comment).Insert()
			if err != nil {
				slog.Error("Failed to save comment", "error", err)
				continue
			}

			// Broadcast to all clients connected to this product
			responseMsg := models.WebSocketMessage{
				Type:      "comment",
				ProductID: client.ProductID,
				Data:      comment,
			}

			msgBytes, _ := json.Marshal(responseMsg)
			hub.broadcast <- &BroadcastMessage{
				ProductID: client.ProductID,
				Message:   msgBytes,
			}

		case "typing":
			// Broadcast typing indicator
			typingMsg := models.WebSocketMessage{
				Type:      "typing",
				ProductID: client.ProductID,
				Data: fiber.Map{
					"userId": client.UserID,
				},
			}
			msgBytes, _ := json.Marshal(typingMsg)
			hub.broadcast <- &BroadcastMessage{
				ProductID: client.ProductID,
				Message:   msgBytes,
			}
		}
	}
}

// writePump writes messages to WebSocket connection
func (h *CommentHandler) writePump(client *Client) {
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

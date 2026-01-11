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
	UserName  string // Masked username for display
	ProductID int
	Send      chan []byte
	Token     string // X-User-Token for user-service calls
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
	db  *pg.DB
	cfg *config.Config
}

func NewCommentHandler(db *pg.DB, cfg *config.Config) *CommentHandler {
	return &CommentHandler{
		db:  db,
		cfg: cfg,
	}
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

	// Struct map thêm full_name
	type CommentWithUser struct {
		models.Comment
		FullName string `pg:"full_name"`
	}

	var commentsWithUsers []CommentWithUser

	err = h.db.ModelContext(ctx, (*models.Comment)(nil)).
		Column("comment.*").
		ColumnExpr("users.full_name AS full_name").
		Join("LEFT JOIN users ON users.id = comment.sender_id").
		Where("comment.product_id = ?", productID).
		Order("comment.created_at DESC"). // lấy mới nhất trước
		Limit(limit).
		Offset(offset).
		Select(&commentsWithUsers)

	if err != nil {
		slog.Error("Failed to get comments", "error", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Lỗi lấy bình luận")
	}

	for i, j := 0, len(commentsWithUsers)-1; i < j; i, j = i+1, j-1 {
		commentsWithUsers[i], commentsWithUsers[j] =
			commentsWithUsers[j], commentsWithUsers[i]
	}

	// Build response
	responses := make([]models.CommentResponse, 0, len(commentsWithUsers))

	for _, cw := range commentsWithUsers {
		userName := ""
		if cw.FullName != "" {
			userName = utils.MaskUserName(cw.FullName)
		} else {
			userName = fmt.Sprintf("Người dùng #%d", cw.SenderID)
		}

		responses = append(responses, models.CommentResponse{
			ID:         cw.ID,
			ProductID:  cw.ProductID,
			SenderID:   cw.SenderID,
			SenderName: userName,
			Content:    cw.Content,
			CreatedAt:  cw.CreatedAt,
		})
	}

	return utils.SuccessResponse(c, responses)
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

// HandleWebSocket handles WebSocket connections
func (h *CommentHandler) HandleWebSocket(c *websocket.Conn) {
	// Get query params
	productIDStr := c.Query("productId")
	tokenString := c.Query("X-User-Token")
	internalJWT := c.Query("X-Internal-JWT")

	if productIDStr == "" || internalJWT == "" {
		slog.Error("Missing required parameters")
		c.Close()
		// (không cần return ở đây)
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

	// Fetch user's name and review info from database (not user-service)
	userName := ""
	var user models.User
	err = h.db.Model(&user).Where("id = ?", userID).Select()
	if err != nil {
		slog.Warn("Failed to get user from database", "userID", userID, "error", err)
		userName = fmt.Sprintf("Người dùng #%d", userID)
	} else {
		userName = utils.MaskUserName(user.FullName)
	}

	// Đánh dấu client có bị chặn comment không
	commentBlocked := false
	if user.TotalNumberReviews > 0 {
		ratio := float64(user.TotalNumberGoodReviews) / float64(user.TotalNumberReviews)
		if ratio < 0.8 {
			commentBlocked = true
		}
	}

	client := &Client{
		Conn:      c,
		UserID:    userID,
		UserName:  userName,
		ProductID: productID,
		Send:      make(chan []byte, 256),
		Token:     tokenString,
	}
	// Truyền trạng thái chặn comment qua context
	ctx := context.WithValue(context.Background(), "commentBlocked", commentBlocked)

	hub.register <- client

	// Start goroutines for reading and writing
	go h.writePump(client)
	h.readPumpWithContext(client, ctx)
	return
}

// readPumpWithContext reads messages from WebSocket connection, cho phép truyền trạng thái chặn comment
func (h *CommentHandler) readPumpWithContext(client *Client, ctx context.Context) {
	defer func() {
		hub.unregister <- client
		client.Conn.Close()
	}()

	commentBlocked, _ := ctx.Value("commentBlocked").(bool)

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
			if commentBlocked {
				// Gửi lỗi qua websocket, không insert comment
				errMsg := models.WebSocketMessage{
					Type: "error",
					Data: fiber.Map{
						"message": "Tỷ lệ tích cực dưới 80% - Không thể comment",
					},
				}
				msgBytes, _ := json.Marshal(errMsg)
				client.Conn.WriteMessage(websocket.TextMessage, msgBytes)
				continue
			}
			// Save comment to database
			comment := &models.Comment{
				ProductID: client.ProductID,
				SenderID:  client.UserID,
				Content:   wsMsg.Content,
				CreatedAt: FixedTimeNow(),
			}

			_, err := h.db.ModelContext(ctx, comment).Insert()
			if err != nil {
				slog.Error("Failed to save comment", "error", err)
				continue
			}

			// Broadcast to all clients with full comment response including username
			commentResp := models.CommentResponse{
				ID:         comment.ID,
				ProductID:  comment.ProductID,
				SenderID:   comment.SenderID,
				SenderName: client.UserName, // Use cached username from client
				Content:    comment.Content,
				CreatedAt:  comment.CreatedAt,
			}

			responseMsg := models.WebSocketMessage{
				Type:      "comment",
				ProductID: client.ProductID,
				Data:      commentResp,
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

	for message := range client.Send {
		if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
	client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
}

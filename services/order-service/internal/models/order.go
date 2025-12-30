package models

import "time"

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPendingPayment  OrderStatus = "PENDING_PAYMENT"  // Chờ người mua thanh toán
	OrderStatusPaid            OrderStatus = "PAID"             // Đã thanh toán
	OrderStatusAddressProvided OrderStatus = "ADDRESS_PROVIDED" // Người mua đã gửi địa chỉ
	OrderStatusShipping        OrderStatus = "SHIPPING"         // Đang vận chuyển
	OrderStatusDelivered       OrderStatus = "DELIVERED"        // Đã giao hàng
	OrderStatusCompleted       OrderStatus = "COMPLETED"        // Hoàn thành
	OrderStatusCancelled       OrderStatus = "CANCELLED"        // Đã hủy
)

// Order represents an order after auction ends
type Order struct {
	tableName struct{} `pg:"orders"`

	ID              int64       `json:"id" pg:"id,pk"`
	AuctionID       int64       `json:"auction_id" pg:"auction_id,notnull"`     // ID sản phẩm đấu giá
	WinnerID        int64       `json:"winner_id" pg:"winner_id,notnull"`       // Người thắng (buyer)
	SellerID        int64       `json:"seller_id" pg:"seller_id,notnull"`       // Người bán
	FinalPrice      float64     `json:"final_price" pg:"final_price,notnull"`   // Giá cuối cùng
	Status          OrderStatus `json:"status" pg:"status,notnull"`             // Trạng thái đơn hàng
	PaymentMethod   string      `json:"payment_method" pg:"payment_method"`     // Phương thức thanh toán
	PaymentProof    string      `json:"payment_proof" pg:"payment_proof"`       // Ảnh chứng từ thanh toán
	ShippingAddress string      `json:"shipping_address" pg:"shipping_address"` // Địa chỉ giao hàng
	ShippingPhone   string      `json:"shipping_phone" pg:"shipping_phone"`     // SĐT nhận hàng
	TrackingNumber  string      `json:"tracking_number" pg:"tracking_number"`   // Mã vận đơn
	ShippingInvoice string      `json:"shipping_invoice" pg:"shipping_invoice"` // Hóa đơn vận chuyển
	PaidAt          *time.Time  `json:"paid_at" pg:"paid_at"`                   // Thời gian thanh toán
	DeliveredAt     *time.Time  `json:"delivered_at" pg:"delivered_at"`         // Thời gian giao hàng
	CompletedAt     *time.Time  `json:"completed_at" pg:"completed_at"`         // Thời gian hoàn thành
	CancelledAt     *time.Time  `json:"cancelled_at" pg:"cancelled_at"`         // Thời gian hủy
	CancelReason    string      `json:"cancel_reason" pg:"cancel_reason"`       // Lý do hủy
	CreatedAt       time.Time   `json:"created_at" pg:"created_at,default:now()"`
	UpdatedAt       time.Time   `json:"updated_at" pg:"updated_at,default:now()"`

	// User names from JOIN
	BuyerName  string `json:"buyer_name,omitempty" pg:"-"`  // Tên người mua (from users table)
	SellerName string `json:"seller_name,omitempty" pg:"-"` // Tên người bán (from users table)

	// Relations
	Messages []*OrderMessage `json:"messages,omitempty" pg:"rel:has-many"`
	Rating   *OrderRating    `json:"rating,omitempty" pg:"rel:has-one"`
}

// OrderMessage represents chat messages between buyer and seller
type OrderMessage struct {
	tableName struct{} `pg:"order_messages"`

	ID        int64     `json:"id" pg:"id,pk"`
	OrderID   int64     `json:"order_id" pg:"order_id,notnull"`
	SenderID  int64     `json:"sender_id" pg:"sender_id,notnull"` // ID người gửi
	Message   string    `json:"message" pg:"message,notnull"`     // Nội dung tin nhắn
	CreatedAt time.Time `json:"created_at" pg:"created_at,default:now()"`
}

// OrderRating represents rating between buyer and seller
type OrderRating struct {
	tableName struct{} `pg:"order_ratings"`

	ID            int64      `json:"id" pg:"id,pk"`
	OrderID       int64      `json:"order_id" pg:"order_id,notnull,unique"`
	BuyerRating   *int       `json:"buyer_rating" pg:"buyer_rating"`     // +1 hoặc -1, null nếu chưa đánh giá
	BuyerComment  string     `json:"buyer_comment" pg:"buyer_comment"`   // Nhận xét của buyer về seller
	SellerRating  *int       `json:"seller_rating" pg:"seller_rating"`   // +1 hoặc -1, null nếu chưa đánh giá
	SellerComment string     `json:"seller_comment" pg:"seller_comment"` // Nhận xét của seller về buyer
	BuyerRatedAt  *time.Time `json:"buyer_rated_at" pg:"buyer_rated_at"`
	SellerRatedAt *time.Time `json:"seller_rated_at" pg:"seller_rated_at"`
	CreatedAt     time.Time  `json:"created_at" pg:"created_at,default:now()"`
	UpdatedAt     time.Time  `json:"updated_at" pg:"updated_at,default:now()"`
}

// User represents user information (for rating updates)
type User struct {
	tableName struct{} `pg:"users"`

	ID                     int64 `json:"id" pg:"id,pk"`
	TotalNumberGoodReviews int   `json:"total_number_good_reviews" pg:"total_number_good_reviews,notnull"`
	TotalNumberReviews     int   `json:"total_number_reviews" pg:"total_number_reviews,notnull"`
}

// CreateOrderRequest represents request to create order after auction ends
type CreateOrderRequest struct {
	AuctionID  int64   `json:"auction_id" validate:"required"`
	WinnerID   int64   `json:"winner_id" validate:"required"`
	SellerID   int64   `json:"seller_id" validate:"required"`
	FinalPrice float64 `json:"final_price" validate:"required,gt=0"`
}

// PaymentRequest represents request to pay for order
type PaymentRequest struct {
	PaymentMethod string `json:"payment_method" validate:"required,oneof=MOMO ZALOPAY VNPAY STRIPE PAYPAL"`
	PaymentProof  string `json:"payment_proof"` // Optional: URL to payment proof image
}

// ShippingAddressRequest represents request to provide shipping address
type ShippingAddressRequest struct {
	ShippingAddress string `json:"shipping_address" validate:"required,min=10,max=500"`
	ShippingPhone   string `json:"shipping_phone" validate:"required,min=10,max=20"`
}

// ShippingInvoiceRequest represents request to send shipping invoice
type ShippingInvoiceRequest struct {
	TrackingNumber  string `json:"tracking_number" validate:"required,min=5,max=100"`
	ShippingInvoice string `json:"shipping_invoice"` // Optional: URL to shipping invoice document
}

// CancelOrderRequest represents request to cancel order
type CancelOrderRequest struct {
	CancelReason string `json:"cancel_reason" validate:"required,min=10,max=500"`
}

// SendMessageRequest represents request to send a message
type SendMessageRequest struct {
	Message string `json:"message" validate:"required,min=1,max=1000"`
}

// RateOrderRequest represents request to rate order
type RateOrderRequest struct {
	Rating  int    `json:"rating" validate:"required,oneof=-1 1"` // +1 or -1
	Comment string `json:"comment" validate:"max=500"`            // Optional comment
}

// WebSocketMessage represents WebSocket message structure
type WebSocketMessage struct {
	Type    string      `json:"type"`    // "message", "typing", "notification"
	OrderID int64       `json:"orderId"` // Order ID
	Data    interface{} `json:"data"`    // Message data
	Content string      `json:"content"` // Message content (for backward compatibility)
}

// OrderListResponse represents paginated order list
type OrderListResponse struct {
	Orders     []*Order `json:"orders"`
	Total      int      `json:"total"`
	Page       int      `json:"page"`
	PageSize   int      `json:"page_size"`
	TotalPages int      `json:"total_pages"`
}

// OrderQueryParams represents query parameters for listing orders
type OrderQueryParams struct {
	UserID   int64       `query:"user_id"` // Filter by buyer or seller
	Status   OrderStatus `query:"status"`  // Filter by status
	Page     int         `query:"page"`
	PageSize int         `query:"page_size"`
}

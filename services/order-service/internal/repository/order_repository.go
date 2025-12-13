package repository

import (
	"order_service/internal/models"
	"time"

	"github.com/go-pg/pg/v10"
)

type OrderRepository interface {
	CreateOrder(order *models.Order) error
	GetOrderByID(id int64) (*models.Order, error)
	GetOrdersByUser(userID int64, page, pageSize int) ([]*models.Order, int, error)
	GetOrdersByStatus(status models.OrderStatus, page, pageSize int) ([]*models.Order, int, error)
	UpdateOrderStatus(id int64, req *models.UpdateOrderStatusRequest) error
	CancelOrder(id int64, reason string, cancelledBy int64) error
	
	// Messages
	CreateMessage(message *models.OrderMessage) error
	GetMessagesByOrder(orderID int64) ([]*models.OrderMessage, error)
	
	// Ratings
	CreateOrUpdateRating(rating *models.OrderRating) error
	GetRatingByOrderID(orderID int64) (*models.OrderRating, error)
	GetUserRatings(userID int64) ([]*models.OrderRating, error)
}

type orderRepository struct {
	db *pg.DB
}

func NewOrderRepository(db *pg.DB) OrderRepository {
	return &orderRepository{db: db}
}

// CreateOrder creates a new order
func (r *orderRepository) CreateOrder(order *models.Order) error {
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	order.Status = models.OrderStatusPendingPayment

	_, err := r.db.Model(order).Insert()
	return err
}

// GetOrderByID retrieves order by ID with relations
func (r *orderRepository) GetOrderByID(id int64) (*models.Order, error) {
	order := &models.Order{}
	err := r.db.Model(order).
		Where("id = ?", id).
		Relation("Messages").
		Relation("Rating").
		Select()
	
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return order, nil
}

// GetOrdersByUser retrieves orders by user (buyer or seller)
func (r *orderRepository) GetOrdersByUser(userID int64, page, pageSize int) ([]*models.Order, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	
	offset := (page - 1) * pageSize
	
	orders := []*models.Order{}
	
	query := r.db.Model(&orders).
		WhereGroup(func(q *pg.Query) (*pg.Query, error) {
			q = q.WhereOr("winner_id = ?", userID).
				WhereOr("seller_id = ?", userID)
			return q, nil
		}).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset)
	
	total, err := query.SelectAndCount()
	if err != nil {
		return nil, 0, err
	}
	
	return orders, total, nil
}

// GetOrdersByStatus retrieves orders by status
func (r *orderRepository) GetOrdersByStatus(status models.OrderStatus, page, pageSize int) ([]*models.Order, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	
	offset := (page - 1) * pageSize
	
	orders := []*models.Order{}
	
	query := r.db.Model(&orders).
		Where("status = ?", status).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset)
	
	total, err := query.SelectAndCount()
	if err != nil {
		return nil, 0, err
	}
	
	return orders, total, nil
}

// UpdateOrderStatus updates order status and related fields
func (r *orderRepository) UpdateOrderStatus(id int64, req *models.UpdateOrderStatusRequest) error {
	order := &models.Order{ID: id}
	
	updates := map[string]interface{}{
		"status":     req.Status,
		"updated_at": time.Now(),
	}
	
	// Update specific fields based on status
	switch req.Status {
	case models.OrderStatusPaymentConfirmed:
		if req.PaymentMethod != "" {
			updates["payment_method"] = req.PaymentMethod
		}
		if req.PaymentProof != "" {
			updates["payment_proof"] = req.PaymentProof
		}
	case models.OrderStatusAddressProvided:
		if req.ShippingAddress != "" {
			updates["shipping_address"] = req.ShippingAddress
		}
		if req.ShippingPhone != "" {
			updates["shipping_phone"] = req.ShippingPhone
		}
	case models.OrderStatusInvoiceSent:
		if req.TrackingNumber != "" {
			updates["tracking_number"] = req.TrackingNumber
		}
		if req.ShippingInvoice != "" {
			updates["shipping_invoice"] = req.ShippingInvoice
		}
	case models.OrderStatusDelivered:
		now := time.Now()
		updates["delivered_at"] = &now
	case models.OrderStatusCompleted:
		now := time.Now()
		updates["completed_at"] = &now
	case models.OrderStatusCancelled:
		now := time.Now()
		updates["cancelled_at"] = &now
		if req.CancelReason != "" {
			updates["cancel_reason"] = req.CancelReason
		}
	}
	
	_, err := r.db.Model(order).
		Where("id = ?", id).
		Set("status = ?status").
		Set("updated_at = ?updated_at").
		Update()
	
	// Update additional fields
	for key, value := range updates {
		if key != "status" && key != "updated_at" {
			_, err = r.db.Model(order).
				Where("id = ?", id).
				Set(key+" = ?", value).
				Update()
			if err != nil {
				return err
			}
		}
	}
	
	return err
}

// CancelOrder cancels an order
func (r *orderRepository) CancelOrder(id int64, reason string, cancelledBy int64) error {
	now := time.Now()
	
	_, err := r.db.Model(&models.Order{}).
		Where("id = ?", id).
		Set("status = ?", models.OrderStatusCancelled).
		Set("cancel_reason = ?", reason).
		Set("cancelled_at = ?", &now).
		Set("updated_at = ?", now).
		Update()
	
	return err
}

// CreateMessage creates a new message in order chat
func (r *orderRepository) CreateMessage(message *models.OrderMessage) error {
	message.CreatedAt = time.Now()
	_, err := r.db.Model(message).Insert()
	return err
}

// GetMessagesByOrder retrieves all messages for an order
func (r *orderRepository) GetMessagesByOrder(orderID int64) ([]*models.OrderMessage, error) {
	messages := []*models.OrderMessage{}
	
	err := r.db.Model(&messages).
		Where("order_id = ?", orderID).
		Order("created_at ASC").
		Select()
	
	if err != nil {
		return nil, err
	}
	
	return messages, nil
}

// CreateOrUpdateRating creates or updates rating for an order
func (r *orderRepository) CreateOrUpdateRating(rating *models.OrderRating) error {
	existing := &models.OrderRating{}
	err := r.db.Model(existing).
		Where("order_id = ?", rating.OrderID).
		Select()
	
	if err == pg.ErrNoRows {
		// Create new rating
		rating.CreatedAt = time.Now()
		rating.UpdatedAt = time.Now()
		_, err = r.db.Model(rating).Insert()
		return err
	}
	
	if err != nil {
		return err
	}
	
	// Update existing rating
	rating.ID = existing.ID
	rating.UpdatedAt = time.Now()
	
	_, err = r.db.Model(rating).
		Where("id = ?", existing.ID).
		UpdateNotZero()
	
	return err
}

// GetRatingByOrderID retrieves rating for an order
func (r *orderRepository) GetRatingByOrderID(orderID int64) (*models.OrderRating, error) {
	rating := &models.OrderRating{}
	
	err := r.db.Model(rating).
		Where("order_id = ?", orderID).
		Select()
	
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return rating, nil
}

// GetUserRatings retrieves all ratings for a user (as buyer or seller)
func (r *orderRepository) GetUserRatings(userID int64) ([]*models.OrderRating, error) {
	ratings := []*models.OrderRating{}
	
	// Get all orders where user is buyer or seller
	orders := []*models.Order{}
	err := r.db.Model(&orders).
		WhereGroup(func(q *pg.Query) (*pg.Query, error) {
			q = q.WhereOr("winner_id = ?", userID).
				WhereOr("seller_id = ?", userID)
			return q, nil
		}).
		Select()
	
	if err != nil {
		return nil, err
	}
	
	orderIDs := make([]int64, len(orders))
	for i, order := range orders {
		orderIDs[i] = order.ID
	}
	
	if len(orderIDs) == 0 {
		return ratings, nil
	}
	
	err = r.db.Model(&ratings).
		Where("order_id IN (?)", pg.In(orderIDs)).
		Select()
	
	if err != nil {
		return nil, err
	}
	
	return ratings, nil
}

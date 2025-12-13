package service

import (
	"errors"
	"order_service/internal/models"
	"order_service/internal/repository"
)

var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrUnauthorized       = errors.New("unauthorized access")
	ErrInvalidStatus      = errors.New("invalid status transition")
	ErrAlreadyRated       = errors.New("already rated")
	ErrOrderNotCompleted  = errors.New("order must be completed to rate")
)

type OrderService interface {
	CreateOrder(req *models.CreateOrderRequest) (*models.Order, error)
	GetOrderByID(id, userID int64) (*models.Order, error)
	GetUserOrders(userID int64, page, pageSize int) (*models.OrderListResponse, error)
	UpdateOrderStatus(id, userID int64, req *models.UpdateOrderStatusRequest) error
	CancelOrder(id, userID int64, reason string) error
	
	SendMessage(orderID, senderID int64, message string) error
	GetMessages(orderID, userID int64) ([]*models.OrderMessage, error)
	
	RateOrder(orderID, userID int64, req *models.RateOrderRequest) error
	GetRating(orderID int64) (*models.OrderRating, error)
}

type orderService struct {
	repo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) OrderService {
	return &orderService{repo: repo}
}

// CreateOrder creates a new order after auction ends
func (s *orderService) CreateOrder(req *models.CreateOrderRequest) (*models.Order, error) {
	order := &models.Order{
		AuctionID:  req.AuctionID,
		WinnerID:   req.WinnerID,
		SellerID:   req.SellerID,
		FinalPrice: req.FinalPrice,
		Status:     models.OrderStatusPendingPayment,
	}
	
	err := s.repo.CreateOrder(order)
	if err != nil {
		return nil, err
	}
	
	return order, nil
}

// GetOrderByID retrieves order by ID (only if user is buyer or seller)
func (s *orderService) GetOrderByID(id, userID int64) (*models.Order, error) {
	order, err := s.repo.GetOrderByID(id)
	if err != nil {
		return nil, err
	}
	
	if order == nil {
		return nil, ErrOrderNotFound
	}
	
	// Check if user is authorized to view this order
	if order.WinnerID != userID && order.SellerID != userID {
		return nil, ErrUnauthorized
	}
	
	return order, nil
}

// GetUserOrders retrieves all orders for a user
func (s *orderService) GetUserOrders(userID int64, page, pageSize int) (*models.OrderListResponse, error) {
	orders, total, err := s.repo.GetOrdersByUser(userID, page, pageSize)
	if err != nil {
		return nil, err
	}
	
	totalPages := (total + pageSize - 1) / pageSize
	
	return &models.OrderListResponse{
		Orders:     orders,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// UpdateOrderStatus updates order status with validation
func (s *orderService) UpdateOrderStatus(id, userID int64, req *models.UpdateOrderStatusRequest) error {
	order, err := s.repo.GetOrderByID(id)
	if err != nil {
		return err
	}
	
	if order == nil {
		return ErrOrderNotFound
	}
	
	// Validate status transition and authorization
	if err := s.validateStatusTransition(order, userID, req.Status); err != nil {
		return err
	}
	
	return s.repo.UpdateOrderStatus(id, req)
}

// validateStatusTransition validates if status transition is valid and user is authorized
func (s *orderService) validateStatusTransition(order *models.Order, userID int64, newStatus models.OrderStatus) error {
	currentStatus := order.Status
	
	// Define valid transitions
	validTransitions := map[models.OrderStatus][]models.OrderStatus{
		models.OrderStatusPendingPayment: {
			models.OrderStatusPaymentConfirmed,
			models.OrderStatusCancelled,
		},
		models.OrderStatusPaymentConfirmed: {
			models.OrderStatusAddressProvided,
			models.OrderStatusCancelled,
		},
		models.OrderStatusAddressProvided: {
			models.OrderStatusInvoiceSent,
			models.OrderStatusCancelled,
		},
		models.OrderStatusInvoiceSent: {
			models.OrderStatusDelivered,
			models.OrderStatusCancelled,
		},
		models.OrderStatusDelivered: {
			models.OrderStatusCompleted,
			models.OrderStatusCancelled,
		},
	}
	
	// Check if transition is valid
	allowedStatuses, ok := validTransitions[currentStatus]
	if !ok {
		return ErrInvalidStatus
	}
	
	isValid := false
	for _, status := range allowedStatuses {
		if status == newStatus {
			isValid = true
			break
		}
	}
	
	if !isValid {
		return ErrInvalidStatus
	}
	
	// Check authorization based on status
	switch newStatus {
	case models.OrderStatusPaymentConfirmed, models.OrderStatusAddressProvided, models.OrderStatusDelivered:
		// Only buyer can update these statuses
		if order.WinnerID != userID {
			return ErrUnauthorized
		}
	case models.OrderStatusInvoiceSent:
		// Only seller can update this status
		if order.SellerID != userID {
			return ErrUnauthorized
		}
	case models.OrderStatusCompleted:
		// Only buyer can mark as completed after delivery
		if order.WinnerID != userID {
			return ErrUnauthorized
		}
	case models.OrderStatusCancelled:
		// Both buyer and seller can cancel
		if order.WinnerID != userID && order.SellerID != userID {
			return ErrUnauthorized
		}
	}
	
	return nil
}

// CancelOrder cancels an order
func (s *orderService) CancelOrder(id, userID int64, reason string) error {
	order, err := s.repo.GetOrderByID(id)
	if err != nil {
		return err
	}
	
	if order == nil {
		return ErrOrderNotFound
	}
	
	// Check if user is authorized
	if order.WinnerID != userID && order.SellerID != userID {
		return ErrUnauthorized
	}
	
	// Cannot cancel if already completed or cancelled
	if order.Status == models.OrderStatusCompleted || order.Status == models.OrderStatusCancelled {
		return ErrInvalidStatus
	}
	
	return s.repo.CancelOrder(id, reason, userID)
}

// SendMessage sends a message in order chat
func (s *orderService) SendMessage(orderID, senderID int64, message string) error {
	order, err := s.repo.GetOrderByID(orderID)
	if err != nil {
		return err
	}
	
	if order == nil {
		return ErrOrderNotFound
	}
	
	// Check if user is buyer or seller
	if order.WinnerID != senderID && order.SellerID != senderID {
		return ErrUnauthorized
	}
	
	msg := &models.OrderMessage{
		OrderID:  orderID,
		SenderID: senderID,
		Message:  message,
	}
	
	return s.repo.CreateMessage(msg)
}

// GetMessages retrieves all messages for an order
func (s *orderService) GetMessages(orderID, userID int64) ([]*models.OrderMessage, error) {
	order, err := s.repo.GetOrderByID(orderID)
	if err != nil {
		return nil, err
	}
	
	if order == nil {
		return nil, ErrOrderNotFound
	}
	
	// Check if user is buyer or seller
	if order.WinnerID != userID && order.SellerID != userID {
		return nil, ErrUnauthorized
	}
	
	return s.repo.GetMessagesByOrder(orderID)
}

// RateOrder allows buyer or seller to rate the order
func (s *orderService) RateOrder(orderID, userID int64, req *models.RateOrderRequest) error {
	order, err := s.repo.GetOrderByID(orderID)
	if err != nil {
		return err
	}
	
	if order == nil {
		return ErrOrderNotFound
	}
	
	// Check if user is buyer or seller
	isBuyer := order.WinnerID == userID
	isSeller := order.SellerID == userID
	
	if !isBuyer && !isSeller {
		return ErrUnauthorized
	}
	
	// Order must be completed or delivered to rate
	if order.Status != models.OrderStatusCompleted && order.Status != models.OrderStatusDelivered {
		return ErrOrderNotCompleted
	}
	
	// Get existing rating
	rating, err := s.repo.GetRatingByOrderID(orderID)
	if err != nil {
		return err
	}
	
	if rating == nil {
		rating = &models.OrderRating{
			OrderID: orderID,
		}
	}
	
	// Update rating based on who is rating
	if isBuyer {
		// Buyer rates seller
		if rating.BuyerRating != nil {
			return ErrAlreadyRated
		}
		rating.BuyerRating = &req.Rating
		rating.BuyerComment = req.Comment
	} else {
		// Seller rates buyer
		if rating.SellerRating != nil {
			return ErrAlreadyRated
		}
		rating.SellerRating = &req.Rating
		rating.SellerComment = req.Comment
	}
	
	return s.repo.CreateOrUpdateRating(rating)
}

// GetRating retrieves rating for an order
func (s *orderService) GetRating(orderID int64) (*models.OrderRating, error) {
	return s.repo.GetRatingByOrderID(orderID)
}

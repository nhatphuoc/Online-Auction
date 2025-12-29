package config

import (
	"log"
	"order_service/internal/models"
	"time"

	"github.com/go-pg/pg/v10"
)

// SeedData creates sample orders for testing (seller 18, bidder 17, products 2,3,4,5,6)
func SeedData(db *pg.DB) error {
	log.Println("Starting database seeding for orders...")

	// Create sample orders for seller 18 and bidder 17
	now := time.Now()
	paidTime := now.Add(-24 * time.Hour)
	deliveredTime := now.Add(-6 * time.Hour)
	completedTime := now.Add(-1 * time.Hour)
	cancelledTime := now.Add(-3 * time.Hour)

	// Order 1: Pending Payment - Product 2
	order1 := &models.Order{
		AuctionID:  2,
		WinnerID:   17, // bidder
		SellerID:   18, // seller
		FinalPrice: 1500000.00,
		Status:     models.OrderStatusPendingPayment,
		CreatedAt:  now.Add(-48 * time.Hour),
		UpdatedAt:  now.Add(-48 * time.Hour),
	}

	// Order 2: Paid - Waiting for shipping address - Product 3
	order2 := &models.Order{
		AuctionID:     3,
		WinnerID:      17,
		SellerID:      18,
		FinalPrice:    2500000.00,
		Status:        models.OrderStatusPaid,
		PaymentMethod: "MOMO",
		PaymentProof:  "https://example.com/payment-proof-1.jpg",
		PaidAt:        &paidTime,
		CreatedAt:     now.Add(-72 * time.Hour),
		UpdatedAt:     now.Add(-24 * time.Hour),
	}

	// Order 3: Address Provided - Waiting for shipping - Product 4
	order3 := &models.Order{
		AuctionID:       4,
		WinnerID:        17,
		SellerID:        18,
		FinalPrice:      3200000.00,
		Status:          models.OrderStatusAddressProvided,
		PaymentMethod:   "VNPAY",
		PaymentProof:    "https://example.com/payment-proof-2.jpg",
		ShippingAddress: "123 Đường Lê Lợi, Phường Bến Thành, Quận 1, TP.HCM",
		ShippingPhone:   "0901234567",
		PaidAt:          &paidTime,
		CreatedAt:       now.Add(-96 * time.Hour),
		UpdatedAt:       now.Add(-12 * time.Hour),
	}

	// Order 4: Shipping - In transit - Product 5
	order4 := &models.Order{
		AuctionID:       5,
		WinnerID:        17,
		SellerID:        18,
		FinalPrice:      1800000.00,
		Status:          models.OrderStatusShipping,
		PaymentMethod:   "ZALOPAY",
		ShippingAddress: "456 Nguyễn Huệ, Phường Bến Nghé, Quận 1, TP.HCM",
		ShippingPhone:   "0912345678",
		TrackingNumber:  "VN123456789",
		ShippingInvoice: "https://example.com/shipping-invoice-1.pdf",
		PaidAt:          &paidTime,
		CreatedAt:       now.Add(-120 * time.Hour),
		UpdatedAt:       now.Add(-6 * time.Hour),
	}

	// Order 5: Delivered - Waiting for confirmation - Product 6
	order5 := &models.Order{
		AuctionID:       6,
		WinnerID:        17,
		SellerID:        18,
		FinalPrice:      4500000.00,
		Status:          models.OrderStatusDelivered,
		PaymentMethod:   "STRIPE",
		ShippingAddress: "789 Hai Bà Trưng, Phường Đa Kao, Quận 1, TP.HCM",
		ShippingPhone:   "0923456789",
		TrackingNumber:  "VN987654321",
		DeliveredAt:     &deliveredTime,
		PaidAt:          &paidTime,
		CreatedAt:       now.Add(-144 * time.Hour),
		UpdatedAt:       now.Add(-12 * time.Hour),
	}

	// Order 6: Completed with rating - Product 7
	order6 := &models.Order{
		AuctionID:       7,
		WinnerID:        17,
		SellerID:        18,
		FinalPrice:      5200000.00,
		Status:          models.OrderStatusCompleted,
		PaymentMethod:   "PAYPAL",
		ShippingAddress: "321 Pasteur, Phường Bến Nghé, Quận 1, TP.HCM",
		ShippingPhone:   "0934567890",
		TrackingNumber:  "VN111222333",
		DeliveredAt:     &deliveredTime,
		CompletedAt:     &completedTime,
		PaidAt:          &paidTime,
		CreatedAt:       now.Add(-168 * time.Hour),
		UpdatedAt:       now.Add(-1 * time.Hour),
	}

	// Order 7: Cancelled - Product 8
	order7 := &models.Order{
		AuctionID:    8,
		WinnerID:     17,
		SellerID:     18,
		FinalPrice:   2100000.00,
		Status:       models.OrderStatusCancelled,
		CancelReason: "Người mua không muốn mua nữa",
		CancelledAt:  &cancelledTime,
		CreatedAt:    now.Add(-192 * time.Hour),
		UpdatedAt:    now.Add(-3 * time.Hour),
	}

	// Insert orders
	orders := []*models.Order{order1, order2, order3, order4, order5, order6, order7}
	_, err := db.Model(&orders).Insert()
	if err != nil {
		return err
	}

	log.Printf("Created %d sample orders\n", len(orders))

	// Create sample messages for some orders
	messages := []*models.OrderMessage{
		{
			OrderID:   3,
			SenderID:  17,
			Message:   "Xin chào, tôi đã gửi địa chỉ giao hàng. Khi nào bạn gửi hàng?",
			CreatedAt: now.Add(-10 * time.Hour),
		},
		{
			OrderID:   3,
			SenderID:  18,
			Message:   "Cảm ơn bạn! Tôi sẽ gửi hàng trong 1-2 ngày tới.",
			CreatedAt: now.Add(-9 * time.Hour),
		},
		{
			OrderID:   4,
			SenderID:  17,
			Message:   "Hàng đến khi nào vậy ạ?",
			CreatedAt: now.Add(-5 * time.Hour),
		},
		{
			OrderID:   4,
			SenderID:  18,
			Message:   "Hàng đang trên đường giao, dự kiến 2-3 ngày nữa nhé!",
			CreatedAt: now.Add(-4 * time.Hour),
		},
		{
			OrderID:   5,
			SenderID:  17,
			Message:   "Tôi đã nhận được hàng, cảm ơn shop!",
			CreatedAt: now.Add(-11 * time.Hour),
		},
	}

	_, err = db.Model(&messages).Insert()
	if err != nil {
		log.Printf("Warning: Error creating messages: %v\n", err)
	} else {
		log.Printf("Created %d sample messages\n", len(messages))
	}

	// Create sample rating for completed order
	buyerRating := 1
	sellerRating := 1
	ratedTime := now.Add(-30 * time.Minute)

	rating := &models.OrderRating{
		OrderID:       6,
		BuyerRating:   &buyerRating,
		BuyerComment:  "Người bán rất nhiệt tình, hàng đẹp đúng như mô tả!",
		SellerRating:  &sellerRating,
		SellerComment: "Người mua thanh toán nhanh, giao dịch tốt!",
		BuyerRatedAt:  &ratedTime,
		SellerRatedAt: &ratedTime,
		CreatedAt:     ratedTime,
		UpdatedAt:     ratedTime,
	}

	_, err = db.Model(rating).Insert()
	if err != nil {
		log.Printf("Warning: Error creating rating: %v\n", err)
	} else {
		log.Println("Created sample rating")
	}

	log.Println("Database seeding completed successfully!")
	return nil
}

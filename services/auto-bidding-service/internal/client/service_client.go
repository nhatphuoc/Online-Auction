package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// BiddingServiceClient là client để gọi API của bidding-service
type BiddingServiceClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewBiddingServiceClient tạo client mới
func NewBiddingServiceClient(baseURL string) *BiddingServiceClient {
	return &BiddingServiceClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// BidRequest là request để đặt giá
type BidRequest struct {
	ProductID int64   `json:"product_id"`
	Amount    float64 `json:"amount"`
	RequestID string  `json:"request_id"`
}

// BidResponse là response từ bidding-service
type BidResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PlaceBid gọi API bidding-service để đặt giá
func (c *BiddingServiceClient) PlaceBid(productID int64, bidderID int64, amount float64, requestID string, userToken string) (*BidResponse, error) {
	reqBody := BidRequest{
		ProductID: productID,
		Amount:    amount,
		RequestID: requestID,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/bids", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Token", userToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var bidResp BidResponse
	if err := json.Unmarshal(body, &bidResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &bidResp, nil
}

// ProductServiceClient là client để gọi API của product-service
type ProductServiceClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewProductServiceClient tạo client mới
func NewProductServiceClient(baseURL string) *ProductServiceClient {
	return &ProductServiceClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ProductInfo là thông tin sản phẩm từ product-service
type ProductInfo struct {
	ID            int64   `json:"id"`
	Name          string  `json:"name"`
	CurrentPrice  float64 `json:"current_price"`
	StepPrice     float64 `json:"step_price"`
	HighestBidder int64   `json:"highest_bidder"`
	Status        string  `json:"status"`
}

// ProductResponse là response từ product-service
type ProductResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    *ProductInfo `json:"data,omitempty"`
}

// GetProduct lấy thông tin sản phẩm
func (c *ProductServiceClient) GetProduct(productID int64) (*ProductInfo, error) {
	url := fmt.Sprintf("%s/products/%d", c.baseURL, productID)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var prodResp ProductResponse
	if err := json.Unmarshal(body, &prodResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !prodResp.Success {
		return nil, fmt.Errorf("product service error: %s", prodResp.Message)
	}

	return prodResp.Data, nil
}

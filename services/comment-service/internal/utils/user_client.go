package utils

import (
	"comment_service/internal/config"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type UserSimpleResponse struct {
	Success bool `json:"success"`
	Data    struct {
		ID       int    `json:"id"`
		FullName string `json:"fullName"`
		Email    string `json:"email"`
	} `json:"data"`
	Message string `json:"message"`
}

var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

// GetUserName fetches user's full name from user-service
func GetUserName(cfg *config.Config, userID int, token string) (string, error) {
	url := fmt.Sprintf("%s/api/users/%d/simple", cfg.ServiceURLs["user-service"], userID)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		slog.Error("Failed to create request", "error", err)
		return "", err
	}

	// Add auth token
	if token != "" {
		req.Header.Set("X-User-Token", token)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		slog.Error("Failed to fetch user info", "error", err, "userID", userID)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		slog.Error("User service returned error", "status", resp.StatusCode, "body", string(body))
		return "", fmt.Errorf("failed to fetch user: status %d", resp.StatusCode)
	}

	var userResp UserSimpleResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		slog.Error("Failed to decode user response", "error", err)
		return "", err
	}

	if !userResp.Success {
		return "", fmt.Errorf("user service failed: %s", userResp.Message)
	}

	return userResp.Data.FullName, nil
}

// MaskUserName masks username for privacy (Người dùng ***xyz)
func MaskUserName(fullName string) string {
	if fullName == "" {
		return "Người dùng ẩn danh"
	}

	// Get last 3 characters
	runes := []rune(fullName)
	nameLen := len(runes)
	
	if nameLen <= 3 {
		return "***" + fullName
	}

	// Show last 3 chars, mask the rest
	maskedPart := "***"
	visiblePart := string(runes[nameLen-3:])
	
	return maskedPart + visiblePart
}

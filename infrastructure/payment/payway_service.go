package payment

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// PaywayConfig holds Payway configuration
type PaywayConfig struct {
	MerchantID   string
	APIKey       string
	APIUsername  string
	BaseURL      string // e.g., https://api-sandbox.payway.com.kh or https://api.payway.com.kh
	ReturnURL    string // Frontend URL to return after payment
	ContinueURL  string // Continue shopping URL
	CallbackURL  string // Backend webhook URL
}

// PaymentMethod represents the payment method type
type PaymentMethod string

const (
	PaymentMethodQR   PaymentMethod = "qr"   // KHQR payment
	PaymentMethodCard PaymentMethod = "card" // Credit/Debit card
)

// PaymentRequest represents a payment request to Payway
type PaymentRequest struct {
	TransactionID  string        `json:"tran_id"`
	Amount         string        `json:"amount"`
	Currency       string        `json:"currency"` // "USD" or "KHR"
	PaymentMethod  PaymentMethod `json:"payment_option"`
	ReturnURL      string        `json:"return_url"`
	ContinueURL    string        `json:"continue_success_url"`
	CustomerName   string        `json:"firstname"`
	CustomerEmail  string        `json:"email"`
	CustomerPhone  string        `json:"phone"`
	Items          string        `json:"items"`
	ShippingMethod string        `json:"shipping"`
}

// PaymentResponse represents Payway's response
type PaymentResponse struct {
	Status        int    `json:"status"`
	Message       string `json:"message"`
	TransactionID string `json:"tran_id"`
	PaymentURL    string `json:"payment_url"` // URL to redirect user for payment
	QRCode        string `json:"qr_code"`     // Base64 QR code for KHQR payments
	Hash          string `json:"hash"`
}

// PaymentCallbackData represents data received from Payway webhook
type PaymentCallbackData struct {
	TransactionID  string `json:"tran_id"`
	Status         string `json:"status"` // "success", "failed", "pending"
	Amount         string `json:"amount"`
	Currency       string `json:"currency"`
	Hash           string `json:"hash"`
	PaymentMethod  string `json:"payment_option"`
	CardNumber     string `json:"card_number_mask"` // Masked card number
	ApprovalCode   string `json:"approval_code"`
	ResponseCode   string `json:"response_code"`
	PaymentDate    string `json:"payment_date"`
	BankName       string `json:"bank"`
}

// PaywayService handles Payway payment operations
type PaywayService interface {
	InitiateQRPayment(transactionID, amount, currency, customerName, customerEmail, customerPhone, items string) (*PaymentResponse, error)
	InitiateCardPayment(transactionID, amount, currency, customerName, customerEmail, customerPhone, items string) (*PaymentResponse, error)
	CheckTransaction(transactionID string) (*TransactionStatusResponse, error)
	VerifyCallback(data *PaymentCallbackData) bool
	GenerateHash(data map[string]string) string
}

type paywayService struct {
	config *PaywayConfig
	client *http.Client
}

// NewPaywayService creates a new Payway service instance
func NewPaywayService(config *PaywayConfig) PaywayService {
	return &paywayService{
		config: config,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// InitiateQRPayment initiates a KHQR payment (for donations)
func (s *paywayService) InitiateQRPayment(transactionID, amount, currency, customerName, customerEmail, customerPhone, items string) (*PaymentResponse, error) {
	req := &PaymentRequest{
		TransactionID:  transactionID,
		Amount:         amount,
		Currency:       currency,
		PaymentMethod:  PaymentMethodQR,
		ReturnURL:      s.config.ReturnURL,
		ContinueURL:    s.config.ContinueURL,
		CustomerName:   customerName,
		CustomerEmail:  customerEmail,
		CustomerPhone:  customerPhone,
		Items:          items,
		ShippingMethod: "NA",
	}

	return s.sendPaymentRequest(req)
}

// InitiateCardPayment initiates a card payment (for sponsorships)
func (s *paywayService) InitiateCardPayment(transactionID, amount, currency, customerName, customerEmail, customerPhone, items string) (*PaymentResponse, error) {
	req := &PaymentRequest{
		TransactionID:  transactionID,
		Amount:         amount,
		Currency:       currency,
		PaymentMethod:  PaymentMethodCard,
		ReturnURL:      s.config.ReturnURL,
		ContinueURL:    s.config.ContinueURL,
		CustomerName:   customerName,
		CustomerEmail:  customerEmail,
		CustomerPhone:  customerPhone,
		Items:          items,
		ShippingMethod: "NA",
	}

	return s.sendPaymentRequest(req)
}

// sendPaymentRequest sends a payment request to Payway API
func (s *paywayService) sendPaymentRequest(req *PaymentRequest) (*PaymentResponse, error) {
	// Build request data for hash
	reqData := map[string]string{
		"merchant_id":          s.config.MerchantID,
		"tran_id":              req.TransactionID,
		"amount":               req.Amount,
		"currency":             req.Currency,
		"payment_option":       string(req.PaymentMethod),
		"return_url":           req.ReturnURL,
		"continue_success_url": req.ContinueURL,
		"firstname":            req.CustomerName,
		"email":                req.CustomerEmail,
		"phone":                req.CustomerPhone,
		"items":                req.Items,
		"shipping":             req.ShippingMethod,
	}

	// Generate hash
	hash := s.GenerateHash(reqData)
	reqData["hash"] = hash
	reqData["type"] = "purchase"

	// Encode request
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", s.config.BaseURL+"/api/payment-gateway/v1/payments/purchase", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.SetBasicAuth(s.config.APIUsername, s.config.APIKey)

	// Send request
	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var paymentResp PaymentResponse
	if err := json.Unmarshal(body, &paymentResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for errors
	if paymentResp.Status != 200 && paymentResp.Status != 0 {
		return nil, fmt.Errorf("payway error: %s (status: %d)", paymentResp.Message, paymentResp.Status)
	}

	return &paymentResp, nil
}

// TransactionStatusResponse represents transaction status check response
type TransactionStatusResponse struct {
	TransactionID string `json:"tran_id"`
	Status        string `json:"status"` // "success", "failed", "pending"
	Amount        string `json:"amount"`
	Currency      string `json:"currency"`
	PaymentMethod string `json:"payment_option"`
	PaymentDate   string `json:"payment_date"`
	ResponseCode  string `json:"response_code"`
	Message       string `json:"message"`
}

// CheckTransaction checks the status of a transaction
func (s *paywayService) CheckTransaction(transactionID string) (*TransactionStatusResponse, error) {
	reqData := map[string]string{
		"merchant_id": s.config.MerchantID,
		"tran_id":     transactionID,
	}

	// Generate hash
	hash := s.GenerateHash(reqData)
	reqData["hash"] = hash

	// Encode request
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", s.config.BaseURL+"/api/payment-gateway/v1/payments/check-transaction", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.SetBasicAuth(s.config.APIUsername, s.config.APIKey)

	// Send request
	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var statusResp TransactionStatusResponse
	if err := json.Unmarshal(body, &statusResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &statusResp, nil
}

// VerifyCallback verifies the authenticity of a callback from Payway
func (s *paywayService) VerifyCallback(data *PaymentCallbackData) bool {
	reqData := map[string]string{
		"tran_id":        data.TransactionID,
		"status":         data.Status,
		"amount":         data.Amount,
		"currency":       data.Currency,
		"payment_option": data.PaymentMethod,
	}

	expectedHash := s.GenerateHash(reqData)
	return expectedHash == data.Hash
}

// GenerateHash generates HMAC SHA-512 hash for Payway requests
func (s *paywayService) GenerateHash(data map[string]string) string {
	// Build string to hash (order matters!)
	// Format: merchant_id + tran_id + amount + items + shipping + ... (check Payway docs)
	var hashString string
	
	// Common fields in order
	if val, ok := data["merchant_id"]; ok {
		hashString += val
	}
	if val, ok := data["tran_id"]; ok {
		hashString += val
	}
	if val, ok := data["amount"]; ok {
		hashString += val
	}
	if val, ok := data["currency"]; ok {
		hashString += val
	}
	if val, ok := data["status"]; ok {
		hashString += val
	}
	if val, ok := data["payment_option"]; ok {
		hashString += val
	}

	// Create HMAC
	h := hmac.New(sha512.New, []byte(s.config.APIKey))
	h.Write([]byte(hashString))
	
	// Return base64 encoded hash
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// Helper function to format amount for Payway (must be in cents/smallest unit)
func FormatAmount(amount float64, currency string) string {
	if currency == "KHR" {
		// KHR doesn't use decimals
		return fmt.Sprintf("%.0f", amount)
	}
	// USD uses cents
	return fmt.Sprintf("%.0f", amount*100)
}

// Helper function to parse amount from Payway response
func ParseAmount(amountStr string, currency string) (float64, error) {
	var amount float64
	_, err := fmt.Sscanf(amountStr, "%f", &amount)
	if err != nil {
		return 0, errors.New("invalid amount format")
	}

	if currency == "USD" {
		// Convert cents to dollars
		return amount / 100, nil
	}
	// KHR is already in full units
	return amount, nil
}


package handler

import (
	"net/http"
	"strconv"

	"github.com/suk-chanthea/ezra/domain/dto"
	"github.com/suk-chanthea/ezra/usecase"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type DonationHandler struct {
	donationUseCase usecase.DonationUseCase
}

func NewDonationHandler(uc usecase.DonationUseCase) *DonationHandler {
	return &DonationHandler{donationUseCase: uc}
}

func (h *DonationHandler) Create(c *gin.Context) {
	var req dto.CreateDonationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			e := validationErrors[0]
			var message string
			switch e.Tag() {
			case "required":
				message = e.Field() + " is required"
			case "gt":
				message = e.Field() + " must be greater than 0"
			case "oneof":
				message = "invalid " + e.Field() + " value"
			default:
				message = "invalid " + e.Field()
			}
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Errors: message})
			return
		}
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid input"})
		return
	}

	// Additional validation for company donations
	if req.DonorType == "company" {
		if req.CompanyName == "" || req.CompanyEmail == "" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "company name and email are required for company donations"})
			return
		}
	}

	// Get user_id from JWT middleware (optional for company donations)
	var userID *uint
	if userIDValue, exists := c.Get("user_id"); exists {
		uid := userIDValue.(uint)
		userID = &uid
	}

	// For user donations, authentication is required
	if req.DonorType == "user" && userID == nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user must be authenticated for user donations"})
		return
	}

	// Create donation
	donation, err := h.donationUseCase.CreateDonation(&req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	// Build response message
	message := "donation created successfully"
	if req.InitiatePayment && donation.PaymentInfo != nil {
		if donation.Type == "donate" {
			message = "donation created successfully. Please scan the QR code to complete payment"
		} else {
			message = "donation created successfully. Please complete payment via the provided URL"
		}
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse{
		Message: message,
		Data:    donation,
	})
}

func (h *DonationHandler) GetAll(c *gin.Context) {
	// Parse filter and pagination parameters
	var filter dto.DonationFilterRequest
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid query parameters"})
		return
	}

	donations, pagination, err := h.donationUseCase.GetAllDonations(&filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.PaginatedResponse{
		Data:       donations,
		Pagination: pagination,
	})
}

func (h *DonationHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	donation, err := h.donationUseCase.GetDonationByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "donation not found"})
		return
	}

	c.JSON(http.StatusOK, donation)
}

func (h *DonationHandler) GetByUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
		return
	}

	// Parse pagination parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	donations, err := h.donationUseCase.GetDonationsByUserID(userID.(uint), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, donations)
}

func (h *DonationHandler) GetByType(c *gin.Context) {
	donationType := c.Param("type")
	if donationType != "donate" && donationType != "sponsor" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid donation type, must be 'donate' or 'sponsor'"})
		return
	}

	// Parse pagination parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	donations, err := h.donationUseCase.GetDonationsByType(donationType, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, donations)
}

func (h *DonationHandler) GetByEvent(c *gin.Context) {
	eventID, err := strconv.ParseUint(c.Param("event_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid event id"})
		return
	}

	// Parse pagination parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	donations, err := h.donationUseCase.GetDonationsByEventID(uint(eventID), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, donations)
}

func (h *DonationHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	var req dto.UpdateDonationStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			e := validationErrors[0]
			var message string
			switch e.Tag() {
			case "required":
				message = e.Field() + " is required"
			case "oneof":
				message = e.Field() + " must be one of: pending, completed, failed, refunded"
			default:
				message = "invalid " + e.Field()
			}
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Errors: message})
			return
		}
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid input"})
		return
	}

	if err := h.donationUseCase.UpdateDonationStatus(uint(id), &req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "donation status updated successfully"})
}

func (h *DonationHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
		return
	}

	if err := h.donationUseCase.DeleteDonation(uint(id), userID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "donation deleted successfully"})
}

func (h *DonationHandler) GetStats(c *gin.Context) {
	stats, err := h.donationUseCase.GetDonationStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *DonationHandler) GetStatsByEvent(c *gin.Context) {
	eventID, err := strconv.ParseUint(c.Param("event_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid event id"})
		return
	}

	stats, err := h.donationUseCase.GetDonationStatsByEventID(uint(eventID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *DonationHandler) InitiatePayment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	paymentResp, err := h.donationUseCase.InitiatePayment(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, paymentResp)
}

// PaywayWebhookRequest represents the webhook payload from Payway
type PaywayWebhookRequest struct {
	TransactionID string `json:"tran_id"`
	Status        string `json:"status"`
	ApprovalCode  string `json:"approval_code"`
	PaymentMethod string `json:"payment_option"`
	Hash          string `json:"hash"`
}

func (h *DonationHandler) HandlePaywayWebhook(c *gin.Context) {
	var webhookData PaywayWebhookRequest
	if err := c.ShouldBindJSON(&webhookData); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid webhook data"})
		return
	}

	// TODO: Verify webhook signature/hash here
	// For now, we'll process it directly

	// Handle the payment callback
	err := h.donationUseCase.HandlePaymentCallback(
		webhookData.TransactionID,
		webhookData.Status,
		webhookData.ApprovalCode,
		webhookData.PaymentMethod,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "webhook processed successfully"})
}

func (h *DonationHandler) CheckQRStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	donation, err := h.donationUseCase.GetDonationByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "donation not found"})
		return
	}

	// Check if it's a QR payment
	if donation.Type != "donate" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "not a QR payment"})
		return
	}

	// Return status with expiration info
	c.JSON(http.StatusOK, gin.H{
		"donation_id": donation.ID,
		"status":      donation.Status,
		"expired":     donation.Status == "pending", // Will be checked by frontend
	})
}

func (h *DonationHandler) RegenerateQR(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	// Regenerate payment (will create new QR with new expiration)
	paymentResp, err := h.donationUseCase.InitiatePayment(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "QR code regenerated successfully",
		Data:    paymentResp,
	})
}


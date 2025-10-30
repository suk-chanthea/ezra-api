package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/suk-chanthea/ezra/domain/dto"
	"github.com/suk-chanthea/ezra/usecase"
)

type SupporterHandler struct {
	supporterUseCase usecase.SupporterUseCase
}

func NewSupporterHandler(supporterUseCase usecase.SupporterUseCase) *SupporterHandler {
	return &SupporterHandler{
		supporterUseCase: supporterUseCase,
	}
}

// Create creates a new supporter
func (h *SupporterHandler) Create(c *gin.Context) {
	var req dto.CreateSupporterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	// Get user ID from context (optional)
	var userID *uint
	if userIDVal, exists := c.Get("user_id"); exists {
		if uid, ok := userIDVal.(uint); ok {
			userID = &uid
		}
	}

	supporter, err := h.supporterUseCase.CreateSupporter(&req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse{
		Message: "supporter created successfully",
		Data:    supporter,
	})
}

// GetAll retrieves all supporters with pagination
func (h *SupporterHandler) GetAll(c *gin.Context) {
	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	supporters, err := h.supporterUseCase.GetAllSupporters(req.GetPage(), req.GetPageSize())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, supporters)
}

// GetByID retrieves a supporter by ID
func (h *SupporterHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid supporter ID"})
		return
	}

	supporter, err := h.supporterUseCase.GetSupporterByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, supporter)
}

// GetByEmail retrieves a supporter by email
func (h *SupporterHandler) GetByEmail(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "email parameter is required"})
		return
	}

	supporter, err := h.supporterUseCase.GetSupporterByEmail(email)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, supporter)
}

// GetByType retrieves supporters by type (company, organization, or church)
func (h *SupporterHandler) GetByType(c *gin.Context) {
	supporterType := c.Param("type")
	if supporterType != "company" && supporterType != "organization" && supporterType != "church" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid supporter type"})
		return
	}

	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	supporters, err := h.supporterUseCase.GetSupportersByType(supporterType, req.GetPage(), req.GetPageSize())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, supporters)
}

// GetByUser retrieves supporters created by a user
func (h *SupporterHandler) GetByUser(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
		return
	}

	userID, ok := userIDVal.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "invalid user ID"})
		return
	}

	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	supporters, err := h.supporterUseCase.GetSupportersByUser(userID, req.GetPage(), req.GetPageSize())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, supporters)
}

// Update updates a supporter
func (h *SupporterHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid supporter ID"})
		return
	}

	var req dto.UpdateSupporterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	// Get user ID from context
	var userID *uint
	if userIDVal, exists := c.Get("user_id"); exists {
		if uid, ok := userIDVal.(uint); ok {
			userID = &uid
		}
	}

	supporter, err := h.supporterUseCase.UpdateSupporter(uint(id), &req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "supporter updated successfully",
		Data:    supporter,
	})
}

// Delete deletes a supporter
func (h *SupporterHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid supporter ID"})
		return
	}

	// Get user ID from context
	var userID *uint
	if userIDVal, exists := c.Get("user_id"); exists {
		if uid, ok := userIDVal.(uint); ok {
			userID = &uid
		}
	}

	if err := h.supporterUseCase.DeleteSupporter(uint(id), userID); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "supporter deleted successfully",
	})
}

// GetStats retrieves donation statistics for a supporter
func (h *SupporterHandler) GetStats(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid supporter ID"})
		return
	}

	stats, err := h.supporterUseCase.GetSupporterStats(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}


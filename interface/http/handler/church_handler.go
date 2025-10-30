package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/suk-chanthea/ezra/domain/dto"
	"github.com/suk-chanthea/ezra/usecase"
)

type ChurchHandler struct {
	churchUseCase usecase.ChurchUseCase
}

func NewChurchHandler(churchUseCase usecase.ChurchUseCase) *ChurchHandler {
	return &ChurchHandler{
		churchUseCase: churchUseCase,
	}
}

// Create creates a new church
func (h *ChurchHandler) Create(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	var req dto.CreateChurchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	church, err := h.churchUseCase.CreateChurch(&req, userID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse{
		Message: "church created successfully",
		Data:    church,
	})
}

// GetAll retrieves all churches with pagination
func (h *ChurchHandler) GetAll(c *gin.Context) {
	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	churches, err := h.churchUseCase.GetAllChurches(req.GetPage(), req.GetPageSize())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, churches)
}

// GetByID retrieves a church by ID
func (h *ChurchHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid church ID"})
		return
	}

	church, err := h.churchUseCase.GetChurchByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, church)
}

// GetByDenomination retrieves churches by denomination
func (h *ChurchHandler) GetByDenomination(c *gin.Context) {
	denomination := c.Param("denomination")
	if denomination == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "denomination parameter is required"})
		return
	}

	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	churches, err := h.churchUseCase.GetChurchesByDenomination(denomination, req.GetPage(), req.GetPageSize())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, churches)
}

// Update updates a church
func (h *ChurchHandler) Update(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid church ID"})
		return
	}

	var req dto.UpdateChurchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	church, err := h.churchUseCase.UpdateChurch(uint(id), &req, userID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "church updated successfully",
		Data:    church,
	})
}

// Delete deletes a church
func (h *ChurchHandler) Delete(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid church ID"})
		return
	}

	if err := h.churchUseCase.DeleteChurch(uint(id), userID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "church deleted successfully",
	})
}

// JoinChurch handles user requests to join a church
func (h *ChurchHandler) JoinChurch(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	var req dto.JoinChurchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.churchUseCase.JoinChurch(userID.(uint), req.ChurchID); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "church join request submitted. Waiting for owner approval",
	})
}

// LeaveChurch handles user requests to leave their current church
func (h *ChurchHandler) LeaveChurch(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	if err := h.churchUseCase.LeaveChurch(userID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "left church successfully",
	})
}

// GetPendingMembers retrieves pending membership requests for a church
func (h *ChurchHandler) GetPendingMembers(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid church ID"})
		return
	}

	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	members, err := h.churchUseCase.GetPendingMembers(uint(id), userID.(uint), req.GetPage(), req.GetPageSize())
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, members)
}

// GetMembers retrieves approved members of a church
func (h *ChurchHandler) GetMembers(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid church ID"})
		return
	}

	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	members, err := h.churchUseCase.GetApprovedMembers(uint(id), req.GetPage(), req.GetPageSize())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, members)
}

// ApproveMember handles approval/rejection of membership requests
func (h *ChurchHandler) ApproveMember(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid church ID"})
		return
	}

	var req dto.ApproveChurchMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.churchUseCase.ApproveMember(uint(id), userID.(uint), req.UserID, req.Status); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	message := "member approved successfully"
	if req.Status == "rejected" {
		message = "member rejected successfully"
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: message,
	})
}


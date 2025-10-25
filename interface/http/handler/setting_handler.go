package handler

import (
	"net/http"

	"github.com/suk-chanthea/ezra/domain/dto"
	"github.com/suk-chanthea/ezra/usecase"

	"github.com/gin-gonic/gin"
)

type SettingHandler struct {
	settingUseCase usecase.SettingUseCase
}

func NewSettingHandler(uc usecase.SettingUseCase) *SettingHandler {
	return &SettingHandler{settingUseCase: uc}
}

// GetSettings gets the current user's settings
func (h *SettingHandler) GetSettings(c *gin.Context) {
	// Get user_id from JWT middleware
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
		return
	}

	response, err := h.settingUseCase.GetUserSettings(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateSettings updates the current user's settings
func (h *SettingHandler) UpdateSettings(c *gin.Context) {
	// Get user_id from JWT middleware
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
		return
	}

	var req dto.UpdateSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
		return
	}

	response, err := h.settingUseCase.UpdateSettings(userID.(uint), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// ResetSettings resets the current user's settings to default values
func (h *SettingHandler) ResetSettings(c *gin.Context) {
	// Get user_id from JWT middleware
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
		return
	}

	response, err := h.settingUseCase.ResetToDefaults(userID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}


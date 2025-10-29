package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suk-chanthea/ezra/domain/entity"
	"github.com/suk-chanthea/ezra/domain/repository"
)

type DeviceTokenHandler struct {
	tokenRepo repository.DeviceTokenRepository
}

func NewDeviceTokenHandler(tokenRepo repository.DeviceTokenRepository) *DeviceTokenHandler {
	return &DeviceTokenHandler{
		tokenRepo: tokenRepo,
	}
}

type RegisterTokenRequest struct {
	Token    string `json:"token" binding:"required"`
	Platform string `json:"platform" binding:"required,oneof=ios android web"`
}

type UnregisterTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

// RegisterToken registers a device token for push notifications
// @Summary Register device token
// @Description Register a device token for receiving push notifications
// @Tags Device Tokens
// @Accept json
// @Produce json
// @Param request body RegisterTokenRequest true "Token registration details"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/device-tokens/register [post]
func (h *DeviceTokenHandler) RegisterToken(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req RegisterTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create device token entity
	deviceToken := entity.NewDeviceToken(userID, req.Token, req.Platform)

	// Validate
	if !deviceToken.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device token data"})
		return
	}

	// Save token (will upsert if exists)
	if err := h.tokenRepo.Save(c.Request.Context(), deviceToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register device token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Device token registered successfully",
		"token_id": deviceToken.ID,
	})
}

// UnregisterToken removes a device token
// @Summary Unregister device token
// @Description Remove a device token (e.g., on logout)
// @Tags Device Tokens
// @Accept json
// @Produce json
// @Param request body UnregisterTokenRequest true "Token to unregister"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/device-tokens/unregister [post]
func (h *DeviceTokenHandler) UnregisterToken(c *gin.Context) {
	var req UnregisterTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Delete token
	if err := h.tokenRepo.DeleteToken(c.Request.Context(), req.Token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unregister device token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Device token unregistered successfully",
	})
}

// DeleteAllTokens removes all device tokens for the current user
// @Summary Delete all device tokens
// @Description Remove all device tokens for the current user (e.g., on account settings)
// @Tags Device Tokens
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/device-tokens/clear [delete]
func (h *DeviceTokenHandler) DeleteAllTokens(c *gin.Context) {
	userID := c.GetUint("user_id")

	if err := h.tokenRepo.DeleteUserTokens(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete device tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "All device tokens deleted successfully",
	})
}


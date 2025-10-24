package handler

import (
	"net/http"
	"strconv"

	"github.com/suk-chanthea/ezra/domain/dto"
	"github.com/suk-chanthea/ezra/usecase"

	"github.com/gin-gonic/gin"
)

type BandHandler struct {
	bandUseCase usecase.BandUseCase
}

func NewBandHandler(uc usecase.BandUseCase) *BandHandler {
	return &BandHandler{bandUseCase: uc}
}

// Create creates a new band
func (h *BandHandler) Create(c *gin.Context) {
	var req dto.CreateBandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid input"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
		return
	}

	if err := h.bandUseCase.CreateBand(c.Request.Context(), &req, userID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse{Message: "band created successfully"})
}

// GetAll gets all bands with optional pagination
func (h *BandHandler) GetAll(c *gin.Context) {
	// Parse pagination parameters
	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid pagination parameters"})
		return
	}

	// If pagination parameters are provided, use paginated query
	if pagination.Page > 0 || pagination.PageSize > 0 {
		bands, meta, err := h.bandUseCase.GetAllBandsPaginated(c.Request.Context(), pagination.GetPage(), pagination.GetPageSize())
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusOK, dto.PaginatedResponse{
			Data:       bands,
			Pagination: meta,
		})
		return
	}

	// Otherwise, return all results
	bands, err := h.bandUseCase.GetAllBands(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, bands)
}

// GetByID gets a band by ID
func (h *BandHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	band, err := h.bandUseCase.GetBandByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "band not found"})
		return
	}

	c.JSON(http.StatusOK, band)
}

// GetByUser gets bands created by the authenticated user
func (h *BandHandler) GetByUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
		return
	}

	bands, err := h.bandUseCase.GetBandsByUserID(c.Request.Context(), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, bands)
}

// GetPublic gets all public bands with optional pagination
func (h *BandHandler) GetPublic(c *gin.Context) {
	// Parse pagination parameters
	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid pagination parameters"})
		return
	}

	// If pagination parameters are provided, use paginated query
	if pagination.Page > 0 || pagination.PageSize > 0 {
		bands, meta, err := h.bandUseCase.GetPublicBandsPaginated(c.Request.Context(), pagination.GetPage(), pagination.GetPageSize())
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusOK, dto.PaginatedResponse{
			Data:       bands,
			Pagination: meta,
		})
		return
	}

	// Otherwise, return all results
	bands, err := h.bandUseCase.GetPublicBands(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, bands)
}

// Update updates a band
func (h *BandHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	var req dto.UpdateBandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid input"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
		return
	}

	if err := h.bandUseCase.UpdateBand(c.Request.Context(), uint(id), &req, userID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "band updated successfully"})
}

// Delete deletes a band
func (h *BandHandler) Delete(c *gin.Context) {
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

	if err := h.bandUseCase.DeleteBand(c.Request.Context(), uint(id), userID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "band deleted successfully"})
}

// AddMusics adds music to a band
func (h *BandHandler) AddMusics(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	var req dto.AddMusicsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid input"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
		return
	}

	if err := h.bandUseCase.AddMusicsToBand(c.Request.Context(), uint(id), req.MusicIDs, userID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse{Message: "music added to band successfully"})
}

// RemoveMusic removes music from a band
func (h *BandHandler) RemoveMusic(c *gin.Context) {
	bandID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid band id"})
		return
	}

	musicID, err := strconv.ParseUint(c.Param("music_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid music id"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
		return
	}

	if err := h.bandUseCase.RemoveMusicFromBand(c.Request.Context(), uint(bandID), uint(musicID), userID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "music removed from band successfully"})
}

// GetMusics gets all music in a band
func (h *BandHandler) GetMusics(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	musics, err := h.bandUseCase.GetBandMusics(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, musics)
}

// ReorderMusics reorders music in a band
func (h *BandHandler) ReorderMusics(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	var req dto.ReorderMusicsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid input"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
		return
	}

	if err := h.bandUseCase.ReorderBandMusics(c.Request.Context(), uint(id), req.MusicOrders, userID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "music order updated successfully"})
}

// GetMembers gets all members of a band
func (h *BandHandler) GetMembers(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	members, err := h.bandUseCase.GetBandMembers(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, members)
}


package handler

import (
	"net/http"
	"strconv"

	"github.com/suk-chanthea/ezra/domain/dto"
	"github.com/suk-chanthea/ezra/usecase"

	"github.com/gin-gonic/gin"
)

type MusicHandler struct {
	musicUseCase usecase.MusicUseCase
}

func NewMusicHandler(uc usecase.MusicUseCase) *MusicHandler {
	return &MusicHandler{musicUseCase: uc}
}

type CreateMusicRequest struct {
	Title       string `json:"title" binding:"required,min=1,max=255"`
	Artist      string `json:"artist" binding:"omitempty,max=255"`
	Album       string `json:"album" binding:"omitempty,max=255"`
	Genre       string `json:"genre" binding:"omitempty,max=100"`
	Duration    int    `json:"duration" binding:"omitempty,min=0"`
	BPM         int    `json:"bpm" binding:"omitempty,min=0"`
	Key         string `json:"key" binding:"omitempty,max=10"`
	Cover       string `json:"cover" binding:"omitempty,max=255"`
	Lyrics      string `json:"lyrics" binding:"omitempty"`
	Description string `json:"description" binding:"omitempty"`
}

func (h *MusicHandler) Create(c *gin.Context) {
	var req CreateMusicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid input"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
		return
	}

	if err := h.musicUseCase.CreateMusic(req.Title, req.Cover, userID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse{Message: "music created successfully"})
}

func (h *MusicHandler) GetAll(c *gin.Context) {
	// Parse pagination parameters
	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid pagination parameters"})
		return
	}

	// If pagination parameters are provided, use paginated query
	if pagination.Page > 0 || pagination.PageSize > 0 {
		musics, meta, err := h.musicUseCase.GetAllMusicsPaginated(pagination.GetPage(), pagination.GetPageSize())
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusOK, dto.PaginatedResponse{
			Data:       musics,
			Pagination: meta,
		})
		return
	}

	// Otherwise, return all results
	musics, err := h.musicUseCase.GetAllMusics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, musics)
}

func (h *MusicHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	music, err := h.musicUseCase.GetMusicByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "music not found"})
		return
	}

	c.JSON(http.StatusOK, music)
}

func (h *MusicHandler) GetByUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
		return
	}

	musics, err := h.musicUseCase.GetMusicsByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, musics)
}

func (h *MusicHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	var req CreateMusicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid input"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
		return
	}

	if err := h.musicUseCase.UpdateMusic(uint(id), req.Title, req.Cover, userID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "music updated successfully"})
}

func (h *MusicHandler) Delete(c *gin.Context) {
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

	if err := h.musicUseCase.DeleteMusic(uint(id), userID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "music deleted successfully"})
}
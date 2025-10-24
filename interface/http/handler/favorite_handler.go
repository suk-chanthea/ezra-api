package handler

import (
	"net/http"
	"strconv"

	"github.com/suk-chanthea/ezra/domain/dto"
	"github.com/suk-chanthea/ezra/usecase"

	"github.com/gin-gonic/gin"
)

type FavoriteHandler struct {
	favoriteUseCase usecase.FavoriteUseCase
}

func NewFavoriteHandler(uc usecase.FavoriteUseCase) *FavoriteHandler {
	return &FavoriteHandler{favoriteUseCase: uc}
}

// AddFavorite adds a music to user's favorites
func (h *FavoriteHandler) AddFavorite(c *gin.Context) {
	musicID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid music id"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
		return
	}

	if err := h.favoriteUseCase.AddFavorite(c.Request.Context(), userID.(uint), uint(musicID)); err != nil {
		if err.Error() == "music not found" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
			return
		}
		if err.Error() == "already favorited" {
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse{Message: "music added to favorites"})
}

// RemoveFavorite removes a music from user's favorites
func (h *FavoriteHandler) RemoveFavorite(c *gin.Context) {
	musicID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid music id"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
		return
	}

	if err := h.favoriteUseCase.RemoveFavorite(c.Request.Context(), userID.(uint), uint(musicID)); err != nil {
		if err.Error() == "favorite not found" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "music removed from favorites"})
}

// GetUserFavorites gets all favorites for the authenticated user
func (h *FavoriteHandler) GetUserFavorites(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
		return
	}

	// Parse pagination parameters
	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid pagination parameters"})
		return
	}

	// If pagination parameters are provided, use paginated query
	if pagination.Page > 0 || pagination.PageSize > 0 {
		favorites, meta, err := h.favoriteUseCase.GetUserFavoritesPaginated(c.Request.Context(), userID.(uint), pagination.GetPage(), pagination.GetPageSize())
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusOK, dto.PaginatedResponse{
			Data:       favorites,
			Pagination: meta,
		})
		return
	}

	// Otherwise, return all results
	favorites, err := h.favoriteUseCase.GetUserFavorites(c.Request.Context(), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, favorites)
}

// IsFavorite checks if a music is favorited by the authenticated user
func (h *FavoriteHandler) IsFavorite(c *gin.Context) {
	musicID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid music id"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
		return
	}

	isFav, err := h.favoriteUseCase.IsFavorite(c.Request.Context(), userID.(uint), uint(musicID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"is_favorite": isFav})
}

// GetFavoriteCount gets the favorite count for a specific music
func (h *FavoriteHandler) GetFavoriteCount(c *gin.Context) {
	musicID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid music id"})
		return
	}

	count, err := h.favoriteUseCase.GetFavoriteCount(c.Request.Context(), uint(musicID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}


package handler

import (
	"net/http"

	"github.com/suk-chanthea/ezra/domain/dto"
	"github.com/suk-chanthea/ezra/usecase"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type EventHandler struct {
	eventUseCase usecase.EventUseCase
}

func NewEventHandler(uc usecase.EventUseCase) *EventHandler {
	return &EventHandler{eventUseCase: uc}
}

func (h *EventHandler) Create(c *gin.Context) {
	var req dto.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			e := validationErrors[0]
			var message string
			switch e.Tag() {
			case "required":
				message = e.Field() + " is required"
			case "min":
				message = e.Field() + " is too short"
			case "max":
				message = e.Field() + " is too long"
			case "gtfield":
				message = e.Field() + " must be after " + e.Param()
			default:
				message = "invalid " + e.Field()
			}
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Errors: message})
			return
		}
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid input"})
		return
	}

	// Get user_id from JWT middleware
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
		return
	}

	// Create event
	if err := h.eventUseCase.CreateEvent(&req, userID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse{Message: "event created successfully"})
}

func (h *EventHandler) GetAll(c *gin.Context) {
	events, err := h.eventUseCase.GetAllEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}
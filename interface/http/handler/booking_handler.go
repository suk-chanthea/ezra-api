package handler

import (
	"net/http"
	"strconv"

	"github.com/suk-chanthea/ezra/domain/dto"
	"github.com/suk-chanthea/ezra/usecase"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BookingHandler struct {
	bookingUseCase usecase.BookingUseCase
}

func NewBookingHandler(uc usecase.BookingUseCase) *BookingHandler {
	return &BookingHandler{bookingUseCase: uc}
}

func (h *BookingHandler) Create(c *gin.Context) {
	var req dto.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			e := validationErrors[0]
			var message string
			switch e.Tag() {
			case "required":
				message = e.Field() + " is required"
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

	// Create booking
	if err := h.bookingUseCase.CreateBooking(&req, userID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse{Message: "booking created successfully"})
}

func (h *BookingHandler) GetAll(c *gin.Context) {
	bookings, err := h.bookingUseCase.GetAllBookings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

func (h *BookingHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	booking, err := h.bookingUseCase.GetBookingByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "booking not found"})
		return
	}

	c.JSON(http.StatusOK, booking)
}

func (h *BookingHandler) GetByEvent(c *gin.Context) {
	eventID, err := strconv.ParseUint(c.Param("event_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid event id"})
		return
	}

	bookings, err := h.bookingUseCase.GetBookingsByEventID(uint(eventID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

func (h *BookingHandler) GetByUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
		return
	}

	bookings, err := h.bookingUseCase.GetBookingsByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

func (h *BookingHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	var req dto.UpdateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			e := validationErrors[0]
			var message string
			switch e.Tag() {
			case "required":
				message = e.Field() + " is required"
			case "oneof":
				message = e.Field() + " must be one of: pending, confirmed, cancelled"
			default:
				message = "invalid " + e.Field()
			}
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Errors: message})
			return
		}
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid input"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
		return
	}

	if err := h.bookingUseCase.UpdateBooking(uint(id), &req, userID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "booking updated successfully"})
}

func (h *BookingHandler) Delete(c *gin.Context) {
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

	if err := h.bookingUseCase.DeleteBooking(uint(id), userID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "booking deleted successfully"})
}


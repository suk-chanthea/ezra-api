package controller

import (
	"github.com/suk-chanthea/ezra/domain"
	"github.com/suk-chanthea/ezra/service"
    "net/http"

	"github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
)

type EventController struct {
    service service.EventService
}

func NewEventController(s service.EventService) *EventController {
    return &EventController{s}
}

func (c *EventController) Create(ctx *gin.Context) {
    var input domain.Event
    if err := ctx.ShouldBindJSON(&input); err != nil {
        // Check if it's a validation error
        if validationErrors, ok := err.(validator.ValidationErrors); ok {
            // Only return the first error as a string
            e := validationErrors[0] // first error
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
            ctx.JSON(http.StatusBadRequest, gin.H{"errors": message})
            return
        }
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
        return
    }
    

    // get user_id from JWT
    userID, exists := ctx.Get("user_id")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
        return
    }
    input.UserID = userID.(uint) // assign correct FK

    // save event (timestamps handled by GORM)
    if err := c.service.CreateEvent(input); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusCreated, gin.H{"message": "event created successfully"})
}

func (c *EventController) GetAll(ctx *gin.Context) {
    events, err := c.service.GetEvents()
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, events)
}

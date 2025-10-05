package controller

import (
    "github.com/suk-chanthea/ezra/domain"
    "github.com/suk-chanthea/ezra/service"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
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
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    input.CreatedAt = time.Now()
    input.UpdatedAt = time.Now()

    if err := c.service.CreateEvent(input); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusCreated, input)
}

func (c *EventController) GetAll(ctx *gin.Context) {
    events, err := c.service.GetEvents()
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, events)
}

package bootstrap

import (
	"github.com/suk-chanthea/ezra/api/middleware"
	"github.com/suk-chanthea/ezra/api/controller"
	"github.com/suk-chanthea/ezra/api/route"
	"github.com/suk-chanthea/ezra/repository"
	"github.com/suk-chanthea/ezra/service"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func SetupRouter(db *gorm.DB, secret string) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	//router.SetTrustedProxies([]string{"127.0.0.1"}) // only trust localhost
	// Public route
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "api work"})
	})

	// API group with auth middleware
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware(db, secret))

	// Event feature
	eventRepo := repository.NewEventRepository(db)
	eventService := service.NewEventService(eventRepo)
	eventController := controller.NewEventController(eventService)
	route.EventRoutes(api, eventController)

	return router
}

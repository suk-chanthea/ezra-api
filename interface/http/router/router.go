package router

import (
	"github.com/suk-chanthea/ezra/interface/http/handler"
	"github.com/suk-chanthea/ezra/interface/http/middleware"
	"github.com/suk-chanthea/ezra/usecase"

	"github.com/gin-gonic/gin"
)

type Router struct {
	authHandler  *handler.AuthHandler
	eventHandler *handler.EventHandler
	authUseCase  usecase.AuthUseCase
}

func NewRouter(
	authHandler *handler.AuthHandler,
	eventHandler *handler.EventHandler,
	authUseCase usecase.AuthUseCase,
) *Router {
	return &Router{
		authHandler:  authHandler,
		eventHandler: eventHandler,
		authUseCase:  authUseCase,
	}
}

func (r *Router) Setup() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Health check
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "api work..."})
	})

	// Public routes
	router.POST("/register", r.authHandler.Register)
	router.POST("/login", r.authHandler.Login)

	// Protected API group
	api := router.Group("/api")
	api.Use(middleware.JWTMiddleware(r.authUseCase))
	{
		// Event routes
		events := api.Group("/events")
		{
			events.POST("/", r.eventHandler.Create)
			events.GET("/", r.eventHandler.GetAll)
		}
	}

	return router
}
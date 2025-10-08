package router

import (
	"github.com/suk-chanthea/ezra/interface/http/handler"
	"github.com/suk-chanthea/ezra/interface/http/middleware"
	"github.com/suk-chanthea/ezra/usecase"

	"github.com/gin-gonic/gin"
)

type Router struct {
	authHandler  *handler.AuthHandler
	musicHandler *handler.MusicHandler
	eventHandler *handler.EventHandler
	authUseCase  usecase.AuthUseCase
}

func NewRouter(
	authHandler *handler.AuthHandler,
	musicHandler *handler.MusicHandler,
	eventHandler *handler.EventHandler,
	authUseCase usecase.AuthUseCase,
) *Router {
	return &Router{
		authHandler:  authHandler,
		musicHandler: musicHandler,
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
		// Music routes
		musics := api.Group("/musics")
		{
			musics.POST("/", r.musicHandler.Create)
			musics.GET("/", r.musicHandler.GetAll)
			musics.GET("/user", r.musicHandler.GetByUser)
			musics.GET("/:id", r.musicHandler.GetByID)
			musics.PUT("/:id", r.musicHandler.Update)
			musics.DELETE("/:id", r.musicHandler.Delete)
		}

		// Event routes
		events := api.Group("/events")
		{
			events.POST("/", r.eventHandler.Create)
			events.GET("/", r.eventHandler.GetAll)
		}
	}

	return router
}
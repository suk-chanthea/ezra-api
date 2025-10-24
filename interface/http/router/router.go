package router

import (
	"github.com/suk-chanthea/ezra/interface/http/handler"
	"github.com/suk-chanthea/ezra/interface/http/middleware"
	"github.com/suk-chanthea/ezra/usecase"

	"github.com/gin-gonic/gin"
)

type Router struct {
	authHandler     *handler.AuthHandler
	musicHandler    *handler.MusicHandler
	eventHandler    *handler.EventHandler
	bookingHandler  *handler.BookingHandler
	favoriteHandler *handler.FavoriteHandler
	authUseCase     usecase.AuthUseCase
}

func NewRouter(
	authHandler *handler.AuthHandler,
	musicHandler *handler.MusicHandler,
	eventHandler *handler.EventHandler,
	bookingHandler *handler.BookingHandler,
	favoriteHandler *handler.FavoriteHandler,
	authUseCase usecase.AuthUseCase,
) *Router {
	return &Router{
		authHandler:     authHandler,
		musicHandler:    musicHandler,
		eventHandler:    eventHandler,
		bookingHandler:  bookingHandler,
		favoriteHandler: favoriteHandler,
		authUseCase:     authUseCase,
	}
}

func (r *Router) Setup() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Health check
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "api work..."})
	})

	// Public routes (no authentication required)
	router.POST("/register", r.authHandler.Register)
	router.POST("/login", r.authHandler.Login)
	router.POST("/auth/google", r.authHandler.GoogleLogin)

	// Protected API group (authentication required)
	api := router.Group("/api")
	api.Use(middleware.JWTMiddleware(r.authUseCase))
	{
		// User/Auth routes
		api.POST("/logout", r.authHandler.Logout)
		api.DELETE("/user", r.authHandler.DeleteUser)

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
			events.GET("/user", r.eventHandler.GetByUser)
			events.GET("/:id", r.eventHandler.GetByID)
			events.PUT("/:id", r.eventHandler.Update)
			events.DELETE("/:id", r.eventHandler.Delete)
		}

		// Booking routes
		bookings := api.Group("/bookings")
		{
			bookings.POST("/", r.bookingHandler.Create)
			bookings.GET("/", r.bookingHandler.GetAll)
			bookings.GET("/user", r.bookingHandler.GetByUser)
			bookings.GET("/event/:event_id", r.bookingHandler.GetByEvent)
			bookings.GET("/:id", r.bookingHandler.GetByID)
			bookings.PUT("/:id", r.bookingHandler.Update)
			bookings.DELETE("/:id", r.bookingHandler.Delete)
		}

		// Favorite routes
		favorites := api.Group("/favorites")
		{
			favorites.GET("/", r.favoriteHandler.GetUserFavorites)
			favorites.POST("/music/:id", r.favoriteHandler.AddFavorite)
			favorites.DELETE("/music/:id", r.favoriteHandler.RemoveFavorite)
			favorites.GET("/music/:id/check", r.favoriteHandler.IsFavorite)
			favorites.GET("/music/:id/count", r.favoriteHandler.GetFavoriteCount)
		}
	}

	return router
}
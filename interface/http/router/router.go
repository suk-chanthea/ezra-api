package router

import (
	"github.com/suk-chanthea/ezra/interface/http/handler"
	"github.com/suk-chanthea/ezra/interface/http/middleware"
	"github.com/suk-chanthea/ezra/usecase"

	"github.com/gin-gonic/gin"
)

type Router struct {
	authHandler         *handler.AuthHandler
	musicHandler        *handler.MusicHandler
	eventHandler        *handler.EventHandler
	bookingHandler      *handler.BookingHandler
	favoriteHandler     *handler.FavoriteHandler
	bandHandler         *handler.BandHandler
	settingHandler      *handler.SettingHandler
	notificationHandler *handler.NotificationHandler
	deviceTokenHandler  *handler.DeviceTokenHandler
	donationHandler     *handler.DonationHandler
	supporterHandler    *handler.SupporterHandler
	churchHandler       *handler.ChurchHandler
	otpHandler          *handler.OTPHandler
	authUseCase         usecase.AuthUseCase
}

func NewRouter(
	authHandler *handler.AuthHandler,
	musicHandler *handler.MusicHandler,
	eventHandler *handler.EventHandler,
	bookingHandler *handler.BookingHandler,
	favoriteHandler *handler.FavoriteHandler,
	bandHandler *handler.BandHandler,
	settingHandler *handler.SettingHandler,
	notificationHandler *handler.NotificationHandler,
	deviceTokenHandler *handler.DeviceTokenHandler,
	donationHandler *handler.DonationHandler,
	supporterHandler *handler.SupporterHandler,
	churchHandler *handler.ChurchHandler,
	otpHandler *handler.OTPHandler,
	authUseCase usecase.AuthUseCase,
) *Router {
	return &Router{
		authHandler:         authHandler,
		musicHandler:        musicHandler,
		eventHandler:        eventHandler,
		bookingHandler:      bookingHandler,
		favoriteHandler:     favoriteHandler,
		bandHandler:         bandHandler,
		settingHandler:      settingHandler,
		notificationHandler: notificationHandler,
		deviceTokenHandler:  deviceTokenHandler,
		donationHandler:     donationHandler,
		supporterHandler:    supporterHandler,
		churchHandler:       churchHandler,
		otpHandler:          otpHandler,
		authUseCase:         authUseCase,
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

	// OTP routes (public - for email verification)
	router.POST("/otp/send", r.otpHandler.SendOTP)
	router.POST("/otp/verify", r.otpHandler.VerifyOTP)

	// Password reset route (public)
	router.POST("/auth/reset-password", r.authHandler.ResetPassword)

	// Public donation routes (companies can donate without auth)
	// router.POST("/donations", r.donationHandler.Create)
	router.POST("/donations/:id/pay", r.donationHandler.InitiatePayment)
	router.POST("/donations/:id/regenerate-qr", r.donationHandler.RegenerateQR)
	router.GET("/donations", r.donationHandler.GetAll)
	router.GET("/donations/stats", r.donationHandler.GetStats)
	router.GET("/donations/:id", r.donationHandler.GetByID)
	router.GET("/donations/:id/status", r.donationHandler.CheckQRStatus)
	router.GET("/donations/type/:type", r.donationHandler.GetByType)
	router.GET("/donations/event/:event_id", r.donationHandler.GetByEvent)
	router.GET("/donations/event/:event_id/stats", r.donationHandler.GetStatsByEvent)
	
	// Payway webhook endpoint (public - called by Payway)
	router.POST("/webhooks/payway", r.donationHandler.HandlePaywayWebhook)

	// Public supporter routes (for viewing supporters)
	router.GET("/supporters", r.supporterHandler.GetAll)
	router.GET("/supporters/:id", r.supporterHandler.GetByID)
	router.GET("/supporters/type/:type", r.supporterHandler.GetByType)
	router.GET("/supporters/search", r.supporterHandler.GetByEmail)

	// Public church routes (for viewing churches)
	router.GET("/churches", r.churchHandler.GetAll)
	router.GET("/churches/:id", r.churchHandler.GetByID)
	router.GET("/churches/:id/members", r.churchHandler.GetMembers) // Get approved members
	router.GET("/churches/denomination/:denomination", r.churchHandler.GetByDenomination)

	// Protected API group (authentication required)
	api := router.Group("/api")
	api.Use(middleware.JWTMiddleware(r.authUseCase))
	{
		// User/Auth routes
		api.GET("/me", r.authHandler.GetMe)
		api.PUT("/me", r.authHandler.UpdateMe)
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

		// Band routes
		bands := api.Group("/bands")
		{
			bands.POST("/", r.bandHandler.Create)
			bands.GET("/", r.bandHandler.GetAll)
			bands.GET("/user", r.bandHandler.GetByUser)
			bands.GET("/public", r.bandHandler.GetPublic)
			bands.GET("/:id", r.bandHandler.GetByID)
			bands.PUT("/:id", r.bandHandler.Update)
			bands.DELETE("/:id", r.bandHandler.Delete)
			
			// Band music management
			bands.GET("/:id/musics", r.bandHandler.GetMusics)
			bands.POST("/:id/musics", r.bandHandler.AddMusics)
			bands.DELETE("/:id/musics/:music_id", r.bandHandler.RemoveMusic)
			bands.PUT("/:id/musics/reorder", r.bandHandler.ReorderMusics)
			
			// Band member management
			bands.GET("/:id/members", r.bandHandler.GetMembers)
		}

		// Setting routes
		settings := api.Group("/settings")
		{
			settings.GET("/", r.settingHandler.GetSettings)
			settings.PUT("/", r.settingHandler.UpdateSettings)
			settings.POST("/reset", r.settingHandler.ResetSettings)
		}

		// Notification routes
		notifications := api.Group("/notifications")
		{
			notifications.POST("/", r.notificationHandler.Create)                                   // Send to specific user
			notifications.POST("/band/:band_id", r.notificationHandler.CreateBandNotification)     // Send to band/team
			notifications.POST("/broadcast", r.notificationHandler.CreateBroadcast)                // Send to all users
			notifications.GET("/", r.notificationHandler.GetAll)
			notifications.GET("/unread", r.notificationHandler.GetUnread)
			notifications.GET("/unread/count", r.notificationHandler.GetUnreadCount)
			notifications.GET("/:id", r.notificationHandler.GetByID)
			notifications.PUT("/:id/read", r.notificationHandler.MarkAsRead)
			notifications.PUT("/read-all", r.notificationHandler.MarkAllAsRead)
			notifications.DELETE("/:id", r.notificationHandler.Delete)
			notifications.DELETE("/", r.notificationHandler.DeleteAll)
		}

		// Device Token routes (FCM push notifications)
		deviceTokens := api.Group("/device-tokens")
		{
			deviceTokens.POST("/register", r.deviceTokenHandler.RegisterToken)     // Register FCM token
			deviceTokens.POST("/unregister", r.deviceTokenHandler.UnregisterToken) // Unregister FCM token
			deviceTokens.DELETE("/clear", r.deviceTokenHandler.DeleteAllTokens)    // Clear all tokens
		}

		// Protected Donation routes (require authentication)
		donations := api.Group("/donations")
		{
			donations.POST("/", r.donationHandler.Create)
			donations.GET("/user", r.donationHandler.GetByUser)
			donations.PUT("/:id/status", r.donationHandler.UpdateStatus)
			donations.DELETE("/:id", r.donationHandler.Delete)
		}

		// Protected Supporter routes (require authentication)
		supporters := api.Group("/supporters")
		{
			supporters.POST("/", r.supporterHandler.Create)
			supporters.GET("/user", r.supporterHandler.GetByUser)
			supporters.GET("/:id/stats", r.supporterHandler.GetStats)
			supporters.PUT("/:id", r.supporterHandler.Update)
			supporters.DELETE("/:id", r.supporterHandler.Delete)
		}

		// Protected Church routes (require authentication)
		churches := api.Group("/churches")
		{
			churches.POST("/", r.churchHandler.Create)              // Create church (user becomes owner)
			churches.PUT("/:id", r.churchHandler.Update)            // Update church (owner only)
			churches.DELETE("/:id", r.churchHandler.Delete)         // Delete church (owner only)
			churches.POST("/join", r.churchHandler.JoinChurch)      // Request to join a church
			churches.POST("/leave", r.churchHandler.LeaveChurch)    // Leave current church
			churches.GET("/:id/pending", r.churchHandler.GetPendingMembers) // Get pending members (owner only)
			churches.POST("/:id/approve", r.churchHandler.ApproveMember)    // Approve/reject member (owner only)
		}
	}

	return router
}
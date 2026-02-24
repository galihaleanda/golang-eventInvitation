package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/galihaleanda/event-invitation/internal/config"
	"github.com/galihaleanda/event-invitation/internal/infrastructure/cache"
	"github.com/galihaleanda/event-invitation/internal/infrastructure/database"
	"github.com/galihaleanda/event-invitation/internal/infrastructure/storage"
	"github.com/galihaleanda/event-invitation/internal/middleware"
	"github.com/galihaleanda/event-invitation/internal/repository"
	"github.com/galihaleanda/event-invitation/internal/service"
	"github.com/gin-gonic/gin"

	handler "github.com/galihaleanda/event-invitation/internal/handler/http"
)

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	fmt.Println("DB_USER:", "["+cfg.Database.User+"]")
	fmt.Println("DB_PASSWORD:", "["+cfg.Database.Password+"]")

	// Connect database
	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("✓ Connected to PostgreSQL")

	// Connect Redis (optional, warn if not available)
	_, err = cache.NewRedis(cfg)
	if err != nil {
		log.Printf("⚠ Redis not available: %v", err)
	} else {
		log.Println("✓ Connected to Redis")
	}

	// Ensure storage directories
	if err := storage.EnsureStorageDirs(cfg); err != nil {
		log.Fatalf("failed to create storage dirs: %v", err)
	}

	// Repositories
	userRepo := repository.NewUserRepository(db)
	templateRepo := repository.NewTemplateRepository(db)
	eventRepo := repository.NewEventRepository(db)
	guestRepo := repository.NewGuestRepository(db)
	mediaRepo := repository.NewMediaRepository(db)

	// Services
	authSvc := service.NewAuthService(userRepo, cfg)
	templateSvc := service.NewTemplateService(templateRepo)
	eventSvc := service.NewEventService(eventRepo, templateRepo, mediaRepo)
	rsvpSvc := service.NewRSVPService(guestRepo, eventRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authSvc)
	templateHandler := handler.NewTemplateHandler(templateSvc)
	eventHandler := handler.NewEventHandler(eventSvc)
	rsvpHandler := handler.NewRSVPHandler(rsvpSvc)
	mediaHandler := handler.NewMediaHandler(mediaRepo, eventRepo, cfg)

	// Gin setup
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())
	r.Use(gin.Recovery())

	// Serve uploaded files
	r.Static("/uploads", cfg.Storage.BasePath)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API v1
	v1 := r.Group("/api/v1")
	{
		// Auth (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Templates (public)
		templates := v1.Group("/templates")
		{
			templates.GET("", templateHandler.GetAll)
			templates.GET("/:id", templateHandler.GetByID)
		}

		// Public event page
		v1.GET("/e/:slug", eventHandler.GetPublic)

		// Public RSVP submission
		v1.POST("/events/:id/rsvp", rsvpHandler.Submit)

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg))
		{
			// Events
			events := protected.Group("/events")
			{
				events.POST("", eventHandler.Create)
				events.GET("", eventHandler.GetMyEvents)
				events.GET("/:id", eventHandler.GetByID)
				events.PATCH("/:id", eventHandler.Update)
				events.DELETE("/:id", eventHandler.Delete)
				events.PATCH("/:id/publish", eventHandler.Publish)
				events.PUT("/:id/theme", eventHandler.UpdateTheme)
				events.PATCH("/:id/sections/:sectionId", eventHandler.UpdateSection)

				// Guests (owner only)
				events.GET("/:id/guests", rsvpHandler.GetGuests)

				// Media
				events.POST("/:id/media", mediaHandler.Upload)
				events.GET("/:id/media", mediaHandler.GetByEvent)
				events.DELETE("/:id/media/:mediaId", mediaHandler.Delete)
			}
		}
	}

	addr := fmt.Sprintf(":%s", cfg.App.Port)
	log.Printf("🚀 Server running on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

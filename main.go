package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/mkamrul9/spotsync-api/config"
	"github.com/mkamrul9/spotsync-api/handler"
	"github.com/mkamrul9/spotsync-api/repository"
	"github.com/mkamrul9/spotsync-api/routes"
	"github.com/mkamrul9/spotsync-api/service"
	"github.com/mkamrul9/spotsync-api/utils"
)

func main() {
	// 1. Load Environment Variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found.")
	}

	// 2. Connect to Database & Run Auto-Migrations
	config.ConnectDB()

	// 3. Setup Echo framework & Custom Validator
	e := echo.New()

	// CORS Middleware (Crucial for frontend integration)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"}, // Change this to your frontend URL in production
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))

	e.Validator = &utils.CustomValidator{Validator: validator.New()}

	// ==========================================
	// 4. DEPENDENCY INJECTION (Strictly Wired)
	// ==========================================

	// Repositories
	userRepo := repository.NewUserRepository(config.DB)
	zoneRepo := repository.NewZoneRepository(config.DB)
	resRepo := repository.NewReservationRepository(config.DB)

	// Services
	authService := service.NewAuthService(userRepo)
	zoneService := service.NewZoneService(zoneRepo, resRepo)
	resService := service.NewReservationService(resRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	zoneHandler := handler.NewZoneHandler(zoneService)
	resHandler := handler.NewReservationHandler(resService)

	// ==========================================

	// 5. Setup Routes
	routes.SetupRoutes(e, authHandler, zoneHandler, resHandler)

	// 6. Health Check Route
	e.GET("/health", func(c echo.Context) error {
		return utils.SendSuccess(c, 200, "SpotSync API is running!", nil)
	})

	// 7. Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s\n", port)
	e.Logger.Fatal(e.Start(":" + port))
}

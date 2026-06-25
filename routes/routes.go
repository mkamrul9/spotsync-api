package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/mkamrul9/spotsync-api/handler"
	"github.com/mkamrul9/spotsync-api/middleware"
)

func SetupRoutes(
	e *echo.Echo,
	authHandler *handler.AuthHandler,
	zoneHandler *handler.ZoneHandler,
	resHandler *handler.ReservationHandler,
) {
	api := e.Group("/api/v1")

	// 🔹 Public Routes
	api.POST("/auth/register", authHandler.Register)
	api.POST("/auth/login", authHandler.Login)
	api.GET("/zones", zoneHandler.GetAllZones)
	api.GET("/zones/:id", zoneHandler.GetZoneByID)

	// 🔹 Authenticated Routes (Requires valid JWT)
	authGroup := api.Group("", middleware.JWTMiddleware)

	// Reservations
	authGroup.POST("/reservations", resHandler.CreateReservation)
	authGroup.GET("/reservations/my-reservations", resHandler.GetMyReservations)
	authGroup.DELETE("/reservations/:id", resHandler.CancelReservation)

	// 🔹 Admin-Only Routes (Requires valid JWT AND 'admin' role)
	adminGroup := authGroup.Group("", middleware.AdminOnlyMiddleware)

	// Zones
	adminGroup.POST("/zones", zoneHandler.CreateZone)
	adminGroup.GET("/reservations", resHandler.GetAllReservations) // Now wired!
}

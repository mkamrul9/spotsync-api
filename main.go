package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/yourusername/spotsync-api/config"
)

func main() {
	// 1. Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, relying on system environment variables.")
	}

	// 2. Connect to Database
	config.ConnectDB()

	// 3. Initialize Echo instance
	e := echo.New()

	// 4. Basic Health Check Route (just to test if server runs)
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "SpotSync API is up and running!",
		})
	})

	// 5. Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s\n", port)
	e.Logger.Fatal(e.Start(":" + port))
}

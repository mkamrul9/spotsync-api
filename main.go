package main

import (
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	"github.com/mkamrul9/spotsync-api/config"
	"github.com/mkamrul9/spotsync-api/utils"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, relying on system variables.")
	}

	config.ConnectDB()

	e := echo.New()

	// Register the custom validator
	e.Validator = &utils.CustomValidator{Validator: validator.New()}

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "SpotSync API is up and running!",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s\n", port)
	e.Logger.Fatal(e.Start(":" + port))
}

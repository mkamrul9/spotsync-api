package utils

import "github.com/labstack/echo/v4"

// SendSuccess formats a successful API response
func SendSuccess(c echo.Context, statusCode int, message string, data interface{}) error {
	return c.JSON(statusCode, map[string]interface{}{
		"success": true,
		"message": message,
		"data":    data,
	})
}

// SendError formats an error API response
func SendError(c echo.Context, statusCode int, message string, errDetails interface{}) error {
	return c.JSON(statusCode, map[string]interface{}{
		"success": false,
		"message": message,
		"errors":  errDetails,
	})
}

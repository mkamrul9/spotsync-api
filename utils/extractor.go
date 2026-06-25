package utils

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// GetUserDetails extracts the user ID and role from the JWT in the Echo context
func GetUserDetails(c echo.Context) (uint, string, error) {
	user := c.Get("user")
	if user == nil {
		return 0, "", errors.New("unauthorized: missing token in context")
	}

	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	// Safely parse ID (JSON numbers are parsed as float64)
	idFloat, ok := claims["id"].(float64)
	if !ok {
		return 0, "", errors.New("invalid token payload: missing id")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return 0, "", errors.New("invalid token payload: missing role")
	}

	return uint(idFloat), role, nil
}

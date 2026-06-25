package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/mkamrul9/spotsync-api/utils"
)

// JWTMiddleware validates the token and injects the user context
func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return utils.SendError(c, http.StatusUnauthorized, "Missing or invalid authorization header", nil)
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		secret := os.Getenv("JWT_SECRET")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return utils.SendError(c, http.StatusUnauthorized, "Invalid or expired token", nil)
		}

		// Inject token into Echo context so our utils.GetUserDetails can find it
		c.Set("user", token)
		return next(c)
	}
}

// AdminOnlyMiddleware enforces the "admin" role
func AdminOnlyMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, role, err := utils.GetUserDetails(c)
		if err != nil || role != "admin" {
			return utils.SendError(c, http.StatusForbidden, "Forbidden: Admin access required", nil)
		}
		return next(c)
	}
}

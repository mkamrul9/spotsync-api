package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mkamrul9/spotsync-api/dto"
	"github.com/mkamrul9/spotsync-api/service"
	"github.com/mkamrul9/spotsync-api/utils"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return utils.SendError(c, http.StatusBadRequest, "Invalid request payload", err.Error())
	}
	if err := c.Validate(&req); err != nil {
		return utils.SendError(c, http.StatusBadRequest, "Validation failed", err.Error())
	}

	res, err := h.authService.Register(req)
	if err != nil {
		return utils.SendError(c, http.StatusConflict, "Registration failed", err.Error())
	}

	return utils.SendSuccess(c, http.StatusCreated, "User registered successfully", res)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return utils.SendError(c, http.StatusBadRequest, "Invalid request payload", err.Error())
	}
	if err := c.Validate(&req); err != nil {
		return utils.SendError(c, http.StatusBadRequest, "Validation failed", err.Error())
	}

	res, err := h.authService.Login(req)
	if err != nil {
		return utils.SendError(c, http.StatusUnauthorized, "Login failed", err.Error())
	}

	return utils.SendSuccess(c, http.StatusOK, "Login successful", res)
}

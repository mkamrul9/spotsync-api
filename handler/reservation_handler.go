package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/mkamrul9/spotsync-api/dto"
	"github.com/mkamrul9/spotsync-api/service"
	"github.com/mkamrul9/spotsync-api/utils"
)

type ReservationHandler struct {
	resService service.ReservationService
}

func NewReservationHandler(resService service.ReservationService) *ReservationHandler {
	return &ReservationHandler{resService}
}

func (h *ReservationHandler) CreateReservation(c echo.Context) error {
	userID, _, err := utils.GetUserDetails(c)
	if err != nil {
		return utils.SendError(c, http.StatusUnauthorized, "Unauthorized", err.Error())
	}

	var req dto.CreateReservationRequest
	if err := c.Bind(&req); err != nil {
		return utils.SendError(c, http.StatusBadRequest, "Invalid request payload", err.Error())
	}
	if err := c.Validate(&req); err != nil {
		return utils.SendError(c, http.StatusBadRequest, "Validation failed", err.Error())
	}

	res, err := h.resService.CreateReservation(userID, req)
	if err != nil {
		// Differentiating between full zone vs server error
		if err.Error() == "zone is at maximum capacity" {
			return utils.SendError(c, http.StatusConflict, "Zone is full", err.Error())
		}
		return utils.SendError(c, http.StatusInternalServerError, "Failed to create reservation", err.Error())
	}

	return utils.SendSuccess(c, http.StatusCreated, "Reservation confirmed successfully", res)
}

func (h *ReservationHandler) GetMyReservations(c echo.Context) error {
	userID, _, err := utils.GetUserDetails(c)
	if err != nil {
		return utils.SendError(c, http.StatusUnauthorized, "Unauthorized", err.Error())
	}

	res, err := h.resService.GetMyReservations(userID)
	if err != nil {
		return utils.SendError(c, http.StatusInternalServerError, "Failed to fetch reservations", err.Error())
	}

	return utils.SendSuccess(c, http.StatusOK, "My reservations retrieved successfully", res)
}

func (h *ReservationHandler) CancelReservation(c echo.Context) error {
	userID, role, err := utils.GetUserDetails(c)
	if err != nil {
		return utils.SendError(c, http.StatusUnauthorized, "Unauthorized", err.Error())
	}

	resIDParam := c.Param("id")
	resID, err := strconv.ParseUint(resIDParam, 10, 32)
	if err != nil {
		return utils.SendError(c, http.StatusBadRequest, "Invalid reservation ID", nil)
	}

	if err := h.resService.CancelReservation(userID, uint(resID), role); err != nil {
		if err.Error() == "forbidden: you can only cancel your own reservations" {
			return utils.SendError(c, http.StatusForbidden, "Forbidden", err.Error())
		}
		return utils.SendError(c, http.StatusNotFound, "Failed to cancel reservation", err.Error())
	}

	return utils.SendSuccess(c, http.StatusOK, "Reservation cancelled successfully", nil)
}

func (h *ReservationHandler) GetAllReservations(c echo.Context) error {
	res, err := h.resService.GetAllSystemReservations()
	if err != nil {
		return utils.SendError(c, http.StatusInternalServerError, "Failed to fetch all reservations", err.Error())
	}
	return utils.SendSuccess(c, http.StatusOK, "All reservations retrieved", res)
}

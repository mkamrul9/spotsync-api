package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mkamrul9/spotsync-api/dto"
	"github.com/mkamrul9/spotsync-api/service"
	"github.com/mkamrul9/spotsync-api/utils"
)

type ZoneHandler struct {
	zoneService service.ZoneService
}

func NewZoneHandler(zoneService service.ZoneService) *ZoneHandler {
	return &ZoneHandler{zoneService}
}

func (h *ZoneHandler) CreateZone(c echo.Context) error {
	var req dto.CreateZoneRequest
	if err := c.Bind(&req); err != nil {
		return utils.SendError(c, http.StatusBadRequest, "Invalid request payload", err.Error())
	}
	if err := c.Validate(&req); err != nil {
		return utils.SendError(c, http.StatusBadRequest, "Validation failed", err.Error())
	}

	res, err := h.zoneService.CreateZone(req)
	if err != nil {
		return utils.SendError(c, http.StatusInternalServerError, "Failed to create parking zone", err.Error())
	}

	return utils.SendSuccess(c, http.StatusCreated, "Parking zone created successfully", res)
}

func (h *ZoneHandler) GetAllZones(c echo.Context) error {
	res, err := h.zoneService.GetAllZones()
	if err != nil {
		return utils.SendError(c, http.StatusInternalServerError, "Failed to fetch parking zones", err.Error())
	}

	return utils.SendSuccess(c, http.StatusOK, "Parking zones retrieved successfully", res)
}

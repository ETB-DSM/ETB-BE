package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Heiji57/ETB-BE/internal/domain"
	"github.com/Heiji57/ETB-BE/internal/service"
)

type DeviceStatusHandler struct {
	svc service.DeviceStatusService
}

func NewDeviceStatus(svc service.DeviceStatusService) *DeviceStatusHandler {
	return &DeviceStatusHandler{svc: svc}
}

func (h *DeviceStatusHandler) Update(c *gin.Context) {
	var req domain.UpdateDeviceStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		domain.Fail(c, domain.ErrInvalidRequest)
		return
	}
	res, err := h.svc.Update(c.Request.Context(), req)
	if err != nil {
		domain.Fail(c, toAppErr(err))
		return
	}
	domain.OK(c, http.StatusOK, res)
}

func (h *DeviceStatusHandler) Get(c *gin.Context) {
	deviceID := c.Param("deviceId")
	res, err := h.svc.Get(c.Request.Context(), deviceID)
	if err != nil {
		domain.Fail(c, toAppErr(err))
		return
	}
	domain.OK(c, http.StatusOK, res)
}

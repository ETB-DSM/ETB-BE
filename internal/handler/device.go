package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Heiji57/ETB-BE/internal/domain"
	"github.com/Heiji57/ETB-BE/internal/service"
)

type DeviceHandler struct {
	svc service.DeviceService
}

func NewDevice(svc service.DeviceService) *DeviceHandler {
	return &DeviceHandler{svc: svc}
}

func (h *DeviceHandler) Create(c *gin.Context) {
	var req domain.CreateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		domain.Fail(c, domain.ErrInvalidRequest)
		return
	}
	userID := c.GetString("userId")
	res, err := h.svc.Create(c.Request.Context(), userID, req.Name)
	if err != nil {
		domain.Fail(c, toAppErr(err))
		return
	}
	domain.OK(c, http.StatusCreated, res)
}

func (h *DeviceHandler) List(c *gin.Context) {
	userID := c.GetString("userId")
	res, err := h.svc.List(c.Request.Context(), userID)
	if err != nil {
		domain.Fail(c, toAppErr(err))
		return
	}
	domain.OK(c, http.StatusOK, res)
}

func (h *DeviceHandler) Delete(c *gin.Context) {
	userID := c.GetString("userId")
	deviceID := c.Param("deviceId")
	if err := h.svc.Delete(c.Request.Context(), userID, deviceID); err != nil {
		domain.Fail(c, toAppErr(err))
		return
	}
	domain.OK(c, http.StatusOK, gin.H{"message": "디바이스가 삭제되었습니다."})
}

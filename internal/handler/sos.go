package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Heiji57/ETB-BE/internal/domain"
	"github.com/Heiji57/ETB-BE/internal/service"
)

type SosHandler struct {
	svc service.SosService
}

func NewSos(svc service.SosService) *SosHandler {
	return &SosHandler{svc: svc}
}

func (h *SosHandler) Create(c *gin.Context) {
	var req domain.CreateSosRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		domain.Fail(c, domain.ErrInvalidRequest)
		return
	}
	userID := c.GetString("userId")
	res, err := h.svc.Create(c.Request.Context(), userID, req)
	if err != nil {
		domain.Fail(c, toAppErr(err))
		return
	}
	domain.OK(c, http.StatusCreated, res)
}

func (h *SosHandler) List(c *gin.Context) {
	userID := c.GetString("userId")
	res, err := h.svc.List(c.Request.Context(), userID)
	if err != nil {
		domain.Fail(c, toAppErr(err))
		return
	}
	domain.OK(c, http.StatusOK, res)
}

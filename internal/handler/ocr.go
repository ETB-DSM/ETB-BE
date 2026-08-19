package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Heiji57/ETB-BE/internal/domain"
	"github.com/Heiji57/ETB-BE/internal/service"
)

type OcrHandler struct {
	svc service.OcrService
}

func NewOcr(svc service.OcrService) *OcrHandler {
	return &OcrHandler{svc: svc}
}

func (h *OcrHandler) Create(c *gin.Context) {
	var req domain.CreateOcrLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		domain.Fail(c, domain.ErrInvalidRequest)
		return
	}
	id, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		domain.Fail(c, toAppErr(err))
		return
	}
	domain.OK(c, http.StatusCreated, gin.H{"logId": id})
}

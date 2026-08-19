package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Heiji57/ETB-BE/internal/domain"
	"github.com/Heiji57/ETB-BE/internal/service"
)

type GuardianHandler struct {
	svc service.GuardianService
}

func NewGuardian(svc service.GuardianService) *GuardianHandler {
	return &GuardianHandler{svc: svc}
}

func (h *GuardianHandler) Create(c *gin.Context) {
	var req domain.CreateGuardianRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		domain.Fail(c, domain.ErrInvalidRequest)
		return
	}
	userID := c.GetString("userId")
	res, err := h.svc.Create(c.Request.Context(), userID, req.Name, req.Phone)
	if err != nil {
		domain.Fail(c, toAppErr(err))
		return
	}
	domain.OK(c, http.StatusCreated, res)
}

func (h *GuardianHandler) List(c *gin.Context) {
	userID := c.GetString("userId")
	res, err := h.svc.List(c.Request.Context(), userID)
	if err != nil {
		domain.Fail(c, toAppErr(err))
		return
	}
	domain.OK(c, http.StatusOK, res)
}

func (h *GuardianHandler) Delete(c *gin.Context) {
	userID := c.GetString("userId")
	guardianID := c.Param("guardianId")
	if err := h.svc.Delete(c.Request.Context(), userID, guardianID); err != nil {
		domain.Fail(c, toAppErr(err))
		return
	}
	domain.OK(c, http.StatusOK, gin.H{"message": "보호자가 삭제되었습니다."})
}

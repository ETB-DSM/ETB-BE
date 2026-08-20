package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Heiji57/ETB-BE/internal/domain"
	"github.com/Heiji57/ETB-BE/internal/service"
)

type NavigationHandler struct {
	svc service.NavigationService
}

func NewNavigation(svc service.NavigationService) *NavigationHandler {
	return &NavigationHandler{svc: svc}
}

func (h *NavigationHandler) CreateSession(c *gin.Context) {
	var req domain.CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		domain.Fail(c, domain.ErrInvalidRequest)
		return
	}
	userID := c.GetString("userId")
	res, err := h.svc.CreateSession(c.Request.Context(), userID, req)
	if err != nil {
		domain.Fail(c, toAppErr(err))
		return
	}
	domain.OK(c, http.StatusCreated, res)
}

func (h *NavigationHandler) UpdateInstruction(c *gin.Context) {
	var req domain.UpdateInstructionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		domain.Fail(c, domain.ErrInvalidRequest)
		return
	}
	userID := c.GetString("userId")
	sessionID := c.Param("sessionId")
	res, err := h.svc.UpdateInstruction(c.Request.Context(), sessionID, userID, req)
	if err != nil {
		domain.Fail(c, toAppErr(err))
		return
	}
	domain.OK(c, http.StatusCreated, res)
}

func (h *NavigationHandler) GetLatestInstruction(c *gin.Context) {
	sessionID := c.Param("sessionId")
	res, err := h.svc.GetLatestInstruction(c.Request.Context(), sessionID)
	if err != nil {
		domain.Fail(c, toAppErr(err))
		return
	}
	domain.OK(c, http.StatusOK, res)
}

func (h *NavigationHandler) UpdateSessionStatus(c *gin.Context) {
	var req domain.UpdateSessionStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		domain.Fail(c, domain.ErrInvalidRequest)
		return
	}
	userID := c.GetString("userId")
	sessionID := c.Param("sessionId")
	res, err := h.svc.UpdateSessionStatus(c.Request.Context(), sessionID, userID, req)
	if err != nil {
		domain.Fail(c, toAppErr(err))
		return
	}
	domain.OK(c, http.StatusOK, res)
}

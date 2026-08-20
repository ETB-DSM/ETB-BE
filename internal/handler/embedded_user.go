package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Heiji57/ETB-BE/internal/domain"
	"github.com/Heiji57/ETB-BE/internal/service"
)

type EmbeddedUserHandler struct {
	svc service.EmbeddedUserService
}

func NewEmbeddedUser(svc service.EmbeddedUserService) *EmbeddedUserHandler {
	return &EmbeddedUserHandler{svc: svc}
}

func (h *EmbeddedUserHandler) Register(c *gin.Context) {
	var req domain.EmbeddedUserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		domain.Fail(c, domain.ErrInvalidRequest)
		return
	}
	res, err := h.svc.Register(c.Request.Context(), req)
	if err != nil {
		domain.Fail(c, toAppErr(err))
		return
	}
	domain.OK(c, http.StatusCreated, res)
}

func (h *EmbeddedUserHandler) Get(c *gin.Context) {
	userID := c.Param("userId")
	res, err := h.svc.Get(c.Request.Context(), userID)
	if err != nil {
		domain.Fail(c, toAppErr(err))
		return
	}
	domain.OK(c, http.StatusOK, res)
}

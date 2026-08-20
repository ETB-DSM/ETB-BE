package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Heiji57/ETB-BE/internal/domain"
	"github.com/Heiji57/ETB-BE/internal/service"
)

type DestinationHandler struct {
	svc service.DestinationService
}

func NewDestination(svc service.DestinationService) *DestinationHandler {
	return &DestinationHandler{svc: svc}
}

func (h *DestinationHandler) Create(c *gin.Context) {
	var req domain.CreateDestinationRequest
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

func (h *DestinationHandler) List(c *gin.Context) {
	userID := c.GetString("userId")
	res, err := h.svc.List(c.Request.Context(), userID)
	if err != nil {
		domain.Fail(c, toAppErr(err))
		return
	}
	domain.OK(c, http.StatusOK, res)
}

// CreateEmbedded handles POST /api/destinations (no JWT, userId in body)
func (h *DestinationHandler) CreateEmbedded(c *gin.Context) {
	var req domain.EmbeddedCreateDestinationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		domain.Fail(c, domain.ErrInvalidRequest)
		return
	}
	inner := domain.CreateDestinationRequest{
		Name:       req.Name,
		Latitude:   req.Latitude,
		Longitude:  req.Longitude,
		RadiusM:    req.RadiusM,
		TargetText: req.TargetText,
	}
	res, err := h.svc.Create(c.Request.Context(), req.UserID, inner)
	if err != nil {
		domain.Fail(c, toAppErr(err))
		return
	}
	domain.OK(c, http.StatusCreated, res)
}

// ListEmbedded handles GET /api/destinations?userId=xxx (no JWT)
func (h *DestinationHandler) ListEmbedded(c *gin.Context) {
	userID := c.Query("userId")
	if userID == "" {
		domain.Fail(c, domain.ErrInvalidRequest)
		return
	}
	res, err := h.svc.List(c.Request.Context(), userID)
	if err != nil {
		domain.Fail(c, toAppErr(err))
		return
	}
	domain.OK(c, http.StatusOK, res)
}

func (h *DestinationHandler) Delete(c *gin.Context) {
	userID := c.GetString("userId")
	destinationID := c.Param("destinationId")
	if err := h.svc.Delete(c.Request.Context(), userID, destinationID); err != nil {
		domain.Fail(c, toAppErr(err))
		return
	}
	domain.OK(c, http.StatusOK, gin.H{"message": "목적지가 삭제되었습니다."})
}

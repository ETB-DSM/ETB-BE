package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Heiji57/ETB-BE/internal/domain"
	"github.com/Heiji57/ETB-BE/internal/service"
)

type AuthHandler struct {
	svc service.AuthService
}

func NewAuth(svc service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Signup(c *gin.Context) {
	var req domain.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		domain.Fail(c, domain.ErrInvalidRequest)
		return
	}
	if err := h.svc.Signup(c.Request.Context(), req); err != nil {
		domain.Fail(c, toAppErr(err))
		return
	}
	domain.OK(c, http.StatusCreated, gin.H{"message": "인증 코드가 발송되었습니다."})
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req domain.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		domain.Fail(c, domain.ErrInvalidRequest)
		return
	}
	if err := h.svc.VerifyEmail(c.Request.Context(), req); err != nil {
		domain.Fail(c, toAppErr(err))
		return
	}
	domain.OK(c, http.StatusOK, gin.H{"message": "이메일 인증이 완료되었습니다."})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		domain.Fail(c, domain.ErrInvalidRequest)
		return
	}
	tokens, err := h.svc.Login(c.Request.Context(), req)
	if err != nil {
		domain.Fail(c, toAppErr(err))
		return
	}
	domain.OK(c, http.StatusOK, tokens)
}

func (h *AuthHandler) LoginGoogle(c *gin.Context) {
	var req domain.GoogleLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		domain.Fail(c, domain.ErrInvalidRequest)
		return
	}
	tokens, err := h.svc.LoginGoogle(c.Request.Context(), req.IDToken)
	if err != nil {
		domain.Fail(c, toAppErr(err))
		return
	}
	domain.OK(c, http.StatusOK, tokens)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req domain.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		domain.Fail(c, domain.ErrInvalidRequest)
		return
	}
	accessToken, err := h.svc.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		domain.Fail(c, toAppErr(err))
		return
	}
	domain.OK(c, http.StatusOK, gin.H{"accessToken": accessToken})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userID := c.GetString("userId")
	if err := h.svc.Logout(c.Request.Context(), userID); err != nil {
		domain.Fail(c, toAppErr(err))
		return
	}
	domain.OK(c, http.StatusOK, gin.H{"message": "로그아웃 되었습니다."})
}

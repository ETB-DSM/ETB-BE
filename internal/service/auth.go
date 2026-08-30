package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/idtoken"

	"github.com/Heiji57/ETB-BE/config"
	"github.com/Heiji57/ETB-BE/internal/domain"
	"github.com/Heiji57/ETB-BE/internal/infra"
	repository "github.com/Heiji57/ETB-BE/internal/repository/sqlc"
)

type AuthService interface {
	Signup(ctx context.Context, req domain.SignupRequest) error
	ResendCode(ctx context.Context, email string) error
	VerifyEmail(ctx context.Context, req domain.VerifyEmailRequest) error
	Login(ctx context.Context, req domain.LoginRequest) (domain.TokenResponse, error)
	LoginGoogle(ctx context.Context, idToken string) (domain.TokenResponse, error)
	Refresh(ctx context.Context, refreshToken string) (domain.TokenResponse, error)
	Logout(ctx context.Context, userID string) error
}

type authService struct {
	repo   *repository.Queries
	redis  *redis.Client
	cfg    *config.Config
	mailer *infra.Mailer
}

func NewAuth(repo *repository.Queries, redis *redis.Client, cfg *config.Config, mailer *infra.Mailer) AuthService {
	return &authService{repo: repo, redis: redis, cfg: cfg, mailer: mailer}
}

func (s *authService) Signup(ctx context.Context, req domain.SignupRequest) error {
	_, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return domain.ErrConflict
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return domain.ErrInternalError
	}
	hashStr := string(hash)

	if _, err := s.repo.CreateUser(ctx, repository.CreateUserParams{
		Email:         req.Email,
		Nickname:      req.Nickname,
		PasswordHash:  &hashStr,
		EmailVerified: false,
	}); err != nil {
		return domain.ErrInternalError
	}

	return s.sendVerificationCode(ctx, req.Email)
}

func (s *authService) ResendCode(ctx context.Context, email string) error {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return domain.ErrNotFound
	}
	if user.EmailVerified {
		return domain.ErrEmailAlreadyVerified
	}
	return s.sendVerificationCode(ctx, email)
}

// sendVerificationCode generates a 6-digit code, stores it in Redis with a
// TTL, and emails it. Shared by Signup and ResendCode.
func (s *authService) sendVerificationCode(ctx context.Context, email string) error {
	code, err := generateCode()
	if err != nil {
		return domain.ErrInternalError
	}

	key := fmt.Sprintf("verify:%s", email)
	ttl := time.Duration(s.cfg.Email.CodeExpireMin) * time.Minute
	if err := s.redis.Set(ctx, key, code, ttl).Err(); err != nil {
		return domain.ErrInternalError
	}

	body := fmt.Sprintf("<p>인증 코드: <strong>%s</strong> (5분 이내 입력)</p>", code)
	if err := s.mailer.Send(email, "[AI-Cane] 이메일 인증 코드", body); err != nil {
		return domain.ErrInternalError
	}
	return nil
}

func (s *authService) VerifyEmail(ctx context.Context, req domain.VerifyEmailRequest) error {
	key := fmt.Sprintf("verify:%s", req.Email)
	stored, err := s.redis.Get(ctx, key).Result()
	if err != nil || stored != req.Code {
		return domain.ErrInvalidVerifyCode
	}
	if err := s.repo.UpdateEmailVerified(ctx, repository.UpdateEmailVerifiedParams{EmailVerified: true, Email: req.Email}); err != nil {
		return domain.ErrInternalError
	}
	s.redis.Del(ctx, key)
	return nil
}

func (s *authService) Login(ctx context.Context, req domain.LoginRequest) (domain.TokenResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return domain.TokenResponse{}, domain.ErrUnauthorized
	}
	if !user.EmailVerified {
		return domain.TokenResponse{}, domain.ErrEmailNotVerified
	}
	if user.PasswordHash == nil {
		return domain.TokenResponse{}, domain.ErrUnauthorized
	}
	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(req.Password)); err != nil {
		return domain.TokenResponse{}, domain.ErrUnauthorized
	}
	return s.issueTokens(ctx, user.ID)
}

func (s *authService) LoginGoogle(ctx context.Context, idTokenStr string) (domain.TokenResponse, error) {
	payload, err := idtoken.Validate(ctx, idTokenStr, s.cfg.Google.ClientID)
	if err != nil {
		return domain.TokenResponse{}, domain.ErrUnauthorized
	}
	email, ok := payload.Claims["email"].(string)
	if !ok || email == "" {
		return domain.TokenResponse{}, domain.ErrUnauthorized
	}

	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		nickname, _ := payload.Claims["name"].(string)
		if nickname == "" {
			nickname = email
		}
		user, err = s.repo.CreateUser(ctx, repository.CreateUserParams{
			Email:         email,
			Nickname:      nickname,
			PasswordHash:  nil,
			EmailVerified: true,
		})
		if err != nil {
			return domain.TokenResponse{}, domain.ErrInternalError
		}
	}
	return s.issueTokens(ctx, user.ID)
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (domain.TokenResponse, error) {
	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrUnauthorized
		}
		return []byte(s.cfg.JWT.RefreshSecret), nil
	})
	if err != nil || !token.Valid {
		return domain.TokenResponse{}, domain.ErrUnauthorized
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return domain.TokenResponse{}, domain.ErrUnauthorized
	}
	userID, ok := claims["sub"].(string)
	if !ok {
		return domain.TokenResponse{}, domain.ErrUnauthorized
	}

	key := fmt.Sprintf("refresh:%s", userID)
	stored, err := s.redis.Get(ctx, key).Result()
	if err != nil || stored != refreshToken {
		return domain.TokenResponse{}, domain.ErrUnauthorized
	}

	return s.issueTokens(ctx, userID)
}

func (s *authService) Logout(ctx context.Context, userID string) error {
	key := fmt.Sprintf("refresh:%s", userID)
	s.redis.Del(ctx, key)
	return nil
}

func (s *authService) issueTokens(ctx context.Context, userID string) (domain.TokenResponse, error) {
	accessToken, err := s.generateAccessToken(userID)
	if err != nil {
		return domain.TokenResponse{}, domain.ErrInternalError
	}
	refreshToken, err := s.generateRefreshToken(userID)
	if err != nil {
		return domain.TokenResponse{}, domain.ErrInternalError
	}

	key := fmt.Sprintf("refresh:%s", userID)
	ttl := time.Duration(s.cfg.JWT.RefreshExpireDay) * 24 * time.Hour
	if err := s.redis.Set(ctx, key, refreshToken, ttl).Err(); err != nil {
		return domain.TokenResponse{}, domain.ErrInternalError
	}

	return domain.TokenResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *authService) generateAccessToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Duration(s.cfg.JWT.AccessExpireMin) * time.Minute).Unix(),
		"iat": time.Now().Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.cfg.JWT.AccessSecret))
}

func (s *authService) generateRefreshToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Duration(s.cfg.JWT.RefreshExpireDay) * 24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.cfg.JWT.RefreshSecret))
}

func generateCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n), nil
}

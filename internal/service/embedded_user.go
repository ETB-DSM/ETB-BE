package service

import (
	"context"

	"github.com/Heiji57/ETB-BE/internal/domain"
	repository "github.com/Heiji57/ETB-BE/internal/repository/sqlc"
)

type EmbeddedUserService interface {
	Register(ctx context.Context, req domain.EmbeddedUserRegisterRequest) (domain.EmbeddedUserResponse, error)
	Get(ctx context.Context, userID string) (domain.EmbeddedUserResponse, error)
}

type embeddedUserService struct {
	repo *repository.Queries
}

func NewEmbeddedUser(repo *repository.Queries) EmbeddedUserService {
	return &embeddedUserService{repo: repo}
}

func (s *embeddedUserService) Register(ctx context.Context, req domain.EmbeddedUserRegisterRequest) (domain.EmbeddedUserResponse, error) {
	if _, err := s.repo.GetEmbeddedUserByID(ctx, req.UserID); err == nil {
		return domain.EmbeddedUserResponse{}, domain.ErrConflict
	}
	if _, err := s.repo.GetEmbeddedUserByDeviceID(ctx, req.DeviceID); err == nil {
		return domain.EmbeddedUserResponse{}, domain.ErrConflict
	}

	u, err := s.repo.CreateEmbeddedUser(ctx, repository.CreateEmbeddedUserParams{
		UserID:        req.UserID,
		Name:          req.Name,
		GuardianName:  req.GuardianName,
		GuardianPhone: req.GuardianPhone,
		DeviceID:      req.DeviceID,
	})
	if err != nil {
		return domain.EmbeddedUserResponse{}, domain.ErrInternalError
	}
	return domain.EmbeddedUserResponse{UserID: u.UserID, DeviceID: u.DeviceID}, nil
}

func (s *embeddedUserService) Get(ctx context.Context, userID string) (domain.EmbeddedUserResponse, error) {
	u, err := s.repo.GetEmbeddedUserByID(ctx, userID)
	if err != nil {
		return domain.EmbeddedUserResponse{}, domain.ErrNotFound
	}
	return domain.EmbeddedUserResponse{UserID: u.UserID, DeviceID: u.DeviceID}, nil
}

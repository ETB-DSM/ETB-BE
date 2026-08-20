package service

import (
	"context"
	"time"

	"github.com/Heiji57/ETB-BE/internal/domain"
	repository "github.com/Heiji57/ETB-BE/internal/repository/sqlc"
)

type LocationService interface {
	Save(ctx context.Context, req domain.CreateLocationRequest) error
}

type locationService struct {
	repo *repository.Queries
}

func NewLocation(repo *repository.Queries) LocationService {
	return &locationService{repo: repo}
}

func (s *locationService) Save(ctx context.Context, req domain.CreateLocationRequest) error {
	if _, err := s.repo.GetEmbeddedUserByID(ctx, req.UserID); err != nil {
		return domain.ErrNotFound
	}

	recordedAt, err := time.Parse(time.RFC3339, req.Timestamp)
	if err != nil {
		recordedAt = time.Now()
	}

	_, err = s.repo.CreateLocation(ctx, repository.CreateLocationParams{
		UserID:     req.UserID,
		DeviceID:   req.DeviceID,
		Latitude:   req.Latitude,
		Longitude:  req.Longitude,
		RecordedAt: recordedAt,
	})
	if err != nil {
		return domain.ErrInternalError
	}
	return nil
}

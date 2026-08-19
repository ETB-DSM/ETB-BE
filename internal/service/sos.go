package service

import (
	"context"

	"github.com/Heiji57/ETB-BE/internal/domain"
	repository "github.com/Heiji57/ETB-BE/internal/repository/sqlc"
)

type SosService interface {
	Create(ctx context.Context, userID string, req domain.CreateSosRequest) (domain.SosResponse, error)
	List(ctx context.Context, userID string) ([]domain.SosResponse, error)
}

type sosService struct {
	repo *repository.Queries
}

func NewSos(repo *repository.Queries) SosService {
	return &sosService{repo: repo}
}

func (s *sosService) Create(ctx context.Context, userID string, req domain.CreateSosRequest) (domain.SosResponse, error) {
	if _, err := s.repo.GetDevice(ctx, req.DeviceID, userID); err != nil {
		return domain.SosResponse{}, domain.ErrNotFound
	}

	e, err := s.repo.CreateSosEvent(ctx, repository.CreateSosEventParams{
		UserID:    userID,
		DeviceID:  req.DeviceID,
		EventType: req.EventType,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	})
	if err != nil {
		return domain.SosResponse{}, domain.ErrInternalError
	}
	return domain.SosResponse{
		SosID:     e.ID,
		DeviceID:  e.DeviceID,
		EventType: e.EventType,
		Latitude:  e.Latitude,
		Longitude: e.Longitude,
		CreatedAt: e.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (s *sosService) List(ctx context.Context, userID string) ([]domain.SosResponse, error) {
	events, err := s.repo.ListSosEventsByUser(ctx, userID)
	if err != nil {
		return nil, domain.ErrInternalError
	}
	res := make([]domain.SosResponse, len(events))
	for i, e := range events {
		res[i] = domain.SosResponse{
			SosID:     e.ID,
			DeviceID:  e.DeviceID,
			EventType: e.EventType,
			Latitude:  e.Latitude,
			Longitude: e.Longitude,
			CreatedAt: e.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}
	return res, nil
}

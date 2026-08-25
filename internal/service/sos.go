package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/Heiji57/ETB-BE/internal/domain"
	repository "github.com/Heiji57/ETB-BE/internal/repository/sqlc"
)

type SosService interface {
	Create(ctx context.Context, userID string, req domain.CreateSosRequest) (domain.SosResponse, error)
	CreateEmbedded(ctx context.Context, userID, deviceID string, req domain.EmbeddedCreateSosRequest) (domain.SosResponse, error)
	List(ctx context.Context, userID string) ([]domain.SosResponse, error)
}

type sosService struct {
	repo *repository.Queries
}

func NewSos(repo *repository.Queries) SosService {
	return &sosService{repo: repo}
}

func (s *sosService) Create(ctx context.Context, userID string, req domain.CreateSosRequest) (domain.SosResponse, error) {
	if _, err := s.repo.GetDevice(ctx, repository.GetDeviceParams{ID: req.DeviceID, UserID: userID}); err != nil {
		return domain.SosResponse{}, domain.ErrNotFound
	}
	return s.createEvent(ctx, userID, req.DeviceID, req.EventType, req.Latitude, req.Longitude, req.Battery)
}

func (s *sosService) CreateEmbedded(ctx context.Context, userID, deviceID string, req domain.EmbeddedCreateSosRequest) (domain.SosResponse, error) {
	return s.createEvent(ctx, userID, deviceID, req.EventType, req.Latitude, req.Longitude, req.Battery)
}

func (s *sosService) createEvent(ctx context.Context, userID, deviceID, eventType string, lat, lon float64, battery *int32) (domain.SosResponse, error) {
	var bat pgtype.Int4
	if battery != nil {
		bat = pgtype.Int4{Int32: *battery, Valid: true}
	}
	e, err := s.repo.CreateSosEvent(ctx, repository.CreateSosEventParams{
		UserID:    userID,
		DeviceID:  deviceID,
		EventType: eventType,
		Latitude:  lat,
		Longitude: lon,
		Battery:   bat,
	})
	if err != nil {
		return domain.SosResponse{}, domain.ErrInternalError
	}
	return domain.SosResponse{
		SosID:          e.ID,
		DeviceID:       e.DeviceID,
		EventType:      e.EventType,
		Latitude:       e.Latitude,
		Longitude:      e.Longitude,
		SentToGuardian: e.SentToGuardian,
		CreatedAt:      e.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
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
			SosID:          e.ID,
			DeviceID:       e.DeviceID,
			EventType:      e.EventType,
			Latitude:       e.Latitude,
			Longitude:      e.Longitude,
			SentToGuardian: e.SentToGuardian,
			CreatedAt:      e.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}
	return res, nil
}

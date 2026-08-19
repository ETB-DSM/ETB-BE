package service

import (
	"context"

	"github.com/Heiji57/ETB-BE/internal/domain"
	repository "github.com/Heiji57/ETB-BE/internal/repository/sqlc"
)

type DestinationService interface {
	Create(ctx context.Context, userID string, req domain.CreateDestinationRequest) (domain.DestinationResponse, error)
	List(ctx context.Context, userID string) ([]domain.DestinationResponse, error)
	Delete(ctx context.Context, userID, destinationID string) error
}

type destinationService struct {
	repo *repository.Queries
}

func NewDestination(repo *repository.Queries) DestinationService {
	return &destinationService{repo: repo}
}

func (s *destinationService) Create(ctx context.Context, userID string, req domain.CreateDestinationRequest) (domain.DestinationResponse, error) {
	d, err := s.repo.CreateDestination(ctx, repository.CreateDestinationParams{
		UserID:     userID,
		Name:       req.Name,
		Latitude:   req.Latitude,
		Longitude:  req.Longitude,
		RadiusM:    req.RadiusM,
		TargetText: req.TargetText,
	})
	if err != nil {
		return domain.DestinationResponse{}, domain.ErrInternalError
	}
	return toDestinationResponse(d), nil
}

func (s *destinationService) List(ctx context.Context, userID string) ([]domain.DestinationResponse, error) {
	dests, err := s.repo.ListDestinationsByUser(ctx, userID)
	if err != nil {
		return nil, domain.ErrInternalError
	}
	res := make([]domain.DestinationResponse, len(dests))
	for i, d := range dests {
		res[i] = toDestinationResponse(d)
	}
	return res, nil
}

func (s *destinationService) Delete(ctx context.Context, userID, destinationID string) error {
	if _, err := s.repo.GetDestination(ctx, destinationID, userID); err != nil {
		return domain.ErrNotFound
	}
	if err := s.repo.DeleteDestination(ctx, destinationID, userID); err != nil {
		return domain.ErrInternalError
	}
	return nil
}

func toDestinationResponse(d repository.Destination) domain.DestinationResponse {
	return domain.DestinationResponse{
		DestinationID: d.ID,
		Name:          d.Name,
		Latitude:      d.Latitude,
		Longitude:     d.Longitude,
		RadiusM:       d.RadiusM,
		TargetText:    d.TargetText,
		CreatedAt:     d.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

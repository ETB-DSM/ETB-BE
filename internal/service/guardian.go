package service

import (
	"context"

	"github.com/Heiji57/ETB-BE/internal/domain"
	repository "github.com/Heiji57/ETB-BE/internal/repository/sqlc"
)

type GuardianService interface {
	Create(ctx context.Context, userID, name, phone string) (domain.GuardianResponse, error)
	List(ctx context.Context, userID string) ([]domain.GuardianResponse, error)
	Delete(ctx context.Context, userID, guardianID string) error
}

type guardianService struct {
	repo *repository.Queries
}

func NewGuardian(repo *repository.Queries) GuardianService {
	return &guardianService{repo: repo}
}

func (s *guardianService) Create(ctx context.Context, userID, name, phone string) (domain.GuardianResponse, error) {
	count, err := s.repo.CountGuardiansByUser(ctx, userID)
	if err != nil {
		return domain.GuardianResponse{}, domain.ErrInternalError
	}
	if count >= 5 {
		return domain.GuardianResponse{}, domain.ErrGuardianLimitExceeded
	}

	g, err := s.repo.CreateGuardian(ctx, repository.CreateGuardianParams{
		UserID: userID,
		Name:   name,
		Phone:  phone,
	})
	if err != nil {
		return domain.GuardianResponse{}, domain.ErrInternalError
	}
	return domain.GuardianResponse{GuardianID: g.ID, Name: g.Name, Phone: g.Phone}, nil
}

func (s *guardianService) List(ctx context.Context, userID string) ([]domain.GuardianResponse, error) {
	guardians, err := s.repo.ListGuardiansByUser(ctx, userID)
	if err != nil {
		return nil, domain.ErrInternalError
	}
	res := make([]domain.GuardianResponse, len(guardians))
	for i, g := range guardians {
		res[i] = domain.GuardianResponse{GuardianID: g.ID, Name: g.Name, Phone: g.Phone}
	}
	return res, nil
}

func (s *guardianService) Delete(ctx context.Context, userID, guardianID string) error {
	if _, err := s.repo.GetGuardian(ctx, repository.GetGuardianParams{ID: guardianID, UserID: userID}); err != nil {
		return domain.ErrNotFound
	}
	if err := s.repo.DeleteGuardian(ctx, repository.DeleteGuardianParams{ID: guardianID, UserID: userID}); err != nil {
		return domain.ErrInternalError
	}
	return nil
}

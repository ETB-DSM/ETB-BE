package service

import (
	"context"

	"github.com/Heiji57/ETB-BE/internal/domain"
	repository "github.com/Heiji57/ETB-BE/internal/repository/sqlc"
)

type OcrService interface {
	Create(ctx context.Context, req domain.CreateOcrLogRequest) (string, error)
}

type ocrService struct {
	repo *repository.Queries
}

func NewOcr(repo *repository.Queries) OcrService {
	return &ocrService{repo: repo}
}

func (s *ocrService) Create(ctx context.Context, req domain.CreateOcrLogRequest) (string, error) {
	log, err := s.repo.CreateOcrLog(ctx, repository.CreateOcrLogParams{
		DestinationID:  req.DestinationID,
		RecognizedText: req.RecognizedText,
		TargetText:     req.TargetText,
		Matched:        req.Matched,
		Confidence:     req.Confidence,
	})
	if err != nil {
		return "", domain.ErrInternalError
	}
	return log.ID, nil
}

package service

import (
	"context"

	"github.com/Heiji57/ETB-BE/internal/domain"
	repository "github.com/Heiji57/ETB-BE/internal/repository/sqlc"
)

type NavigationService interface {
	CreateSession(ctx context.Context, userID string, req domain.CreateSessionRequest) (domain.SessionResponse, error)
	UpdateInstruction(ctx context.Context, sessionID, userID string, req domain.UpdateInstructionRequest) (domain.InstructionResponse, error)
	GetLatestInstruction(ctx context.Context, sessionID string) (domain.InstructionResponse, error)
	UpdateSessionStatus(ctx context.Context, sessionID, userID string, req domain.UpdateSessionStatusRequest) (domain.SessionResponse, error)
}

type navigationService struct {
	repo *repository.Queries
}

func NewNavigation(repo *repository.Queries) NavigationService {
	return &navigationService{repo: repo}
}

func (s *navigationService) CreateSession(ctx context.Context, userID string, req domain.CreateSessionRequest) (domain.SessionResponse, error) {
	sess, err := s.repo.CreateNavigationSession(ctx, repository.CreateNavigationSessionParams{
		UserID:         userID,
		DeviceID:       req.DeviceID,
		DestinationID:  req.DestinationID,
		StartLatitude:  req.StartLatitude,
		StartLongitude: req.StartLongitude,
	})
	if err != nil {
		return domain.SessionResponse{}, domain.ErrInternalError
	}
	return toSessionResponse(sess), nil
}

func (s *navigationService) UpdateInstruction(ctx context.Context, sessionID, userID string, req domain.UpdateInstructionRequest) (domain.InstructionResponse, error) {
	sess, err := s.repo.GetNavigationSession(ctx, sessionID)
	if err != nil {
		return domain.InstructionResponse{}, domain.ErrNotFound
	}
	if sess.UserID != userID {
		return domain.InstructionResponse{}, domain.ErrForbidden
	}

	instr, err := s.repo.CreateNavigationInstruction(ctx, repository.CreateNavigationInstructionParams{
		SessionID:      sessionID,
		Action:         req.Action,
		DistanceMeters: req.DistanceMeters,
		Message:        req.Message,
	})
	if err != nil {
		return domain.InstructionResponse{}, domain.ErrInternalError
	}
	return toInstructionResponse(instr), nil
}

func (s *navigationService) GetLatestInstruction(ctx context.Context, sessionID string) (domain.InstructionResponse, error) {
	instr, err := s.repo.GetLatestNavigationInstruction(ctx, sessionID)
	if err != nil {
		return domain.InstructionResponse{}, domain.ErrNotFound
	}
	return toInstructionResponse(instr), nil
}

func (s *navigationService) UpdateSessionStatus(ctx context.Context, sessionID, userID string, req domain.UpdateSessionStatusRequest) (domain.SessionResponse, error) {
	sess, err := s.repo.GetNavigationSession(ctx, sessionID)
	if err != nil {
		return domain.SessionResponse{}, domain.ErrNotFound
	}
	if sess.UserID != userID {
		return domain.SessionResponse{}, domain.ErrForbidden
	}

	updated, err := s.repo.UpdateNavigationSessionStatus(ctx, sessionID, req.Status)
	if err != nil {
		return domain.SessionResponse{}, domain.ErrInternalError
	}
	return toSessionResponse(updated), nil
}

func toSessionResponse(s repository.NavigationSession) domain.SessionResponse {
	return domain.SessionResponse{
		SessionID:      s.ID,
		Status:         s.Status,
		StartLatitude:  s.StartLatitude,
		StartLongitude: s.StartLongitude,
		CreatedAt:      s.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toInstructionResponse(i repository.NavigationInstruction) domain.InstructionResponse {
	return domain.InstructionResponse{
		InstructionID:  i.ID,
		SessionID:      i.SessionID,
		Action:         i.Action,
		DistanceMeters: i.DistanceMeters,
		Message:        i.Message,
		CreatedAt:      i.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

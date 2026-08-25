package service

import (
	"context"
	"crypto/rand"
	"fmt"

	"github.com/Heiji57/ETB-BE/internal/domain"
	repository "github.com/Heiji57/ETB-BE/internal/repository/sqlc"
)

type DeviceService interface {
	Create(ctx context.Context, userID, name string) (domain.DeviceResponse, error)
	List(ctx context.Context, userID string) ([]domain.DeviceResponse, error)
	Delete(ctx context.Context, userID, deviceID string) error
}

type deviceService struct {
	repo *repository.Queries
}

func NewDevice(repo *repository.Queries) DeviceService {
	return &deviceService{repo: repo}
}

func (s *deviceService) Create(ctx context.Context, userID, name string) (domain.DeviceResponse, error) {
	count, err := s.repo.CountDevicesByUser(ctx, userID)
	if err != nil {
		return domain.DeviceResponse{}, domain.ErrInternalError
	}
	if count >= 5 {
		return domain.DeviceResponse{}, domain.ErrDeviceLimitExceeded
	}

	apiKey, err := generateAPIKey()
	if err != nil {
		return domain.DeviceResponse{}, domain.ErrInternalError
	}

	d, err := s.repo.CreateDevice(ctx, repository.CreateDeviceParams{UserID: userID, Name: name, APIKey: apiKey})
	if err != nil {
		return domain.DeviceResponse{}, domain.ErrInternalError
	}
	res := toDeviceResponse(d)
	res.APIKey = d.APIKey
	return res, nil
}

func generateAPIKey() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}

func (s *deviceService) List(ctx context.Context, userID string) ([]domain.DeviceResponse, error) {
	devices, err := s.repo.ListDevicesByUser(ctx, userID)
	if err != nil {
		return nil, domain.ErrInternalError
	}
	res := make([]domain.DeviceResponse, len(devices))
	for i, d := range devices {
		res[i] = toDeviceResponse(d)
	}
	return res, nil
}

func (s *deviceService) Delete(ctx context.Context, userID, deviceID string) error {
	if _, err := s.repo.GetDevice(ctx, repository.GetDeviceParams{ID: deviceID, UserID: userID}); err != nil {
		return domain.ErrNotFound
	}
	if err := s.repo.DeleteDevice(ctx, repository.DeleteDeviceParams{ID: deviceID, UserID: userID}); err != nil {
		return domain.ErrInternalError
	}
	return nil
}

func toDeviceResponse(d repository.Device) domain.DeviceResponse {
	return domain.DeviceResponse{
		DeviceID:  d.ID,
		Name:      d.Name,
		IsActive:  d.IsActive,
		CreatedAt: d.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}


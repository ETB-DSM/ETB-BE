package service

import (
	"context"

	"github.com/Heiji57/ETB-BE/internal/domain"
	repository "github.com/Heiji57/ETB-BE/internal/repository/sqlc"
)

type DeviceStatusService interface {
	Update(ctx context.Context, req domain.UpdateDeviceStatusRequest) (domain.DeviceStatusResponse, error)
	Get(ctx context.Context, deviceID string) (domain.DeviceStatusResponse, error)
	GetForUser(ctx context.Context, userID, deviceID string) (domain.DeviceStatusResponse, error)
}

type deviceStatusService struct {
	repo *repository.Queries
}

func NewDeviceStatus(repo *repository.Queries) DeviceStatusService {
	return &deviceStatusService{repo: repo}
}

func (s *deviceStatusService) Update(ctx context.Context, req domain.UpdateDeviceStatusRequest) (domain.DeviceStatusResponse, error) {
	ds, err := s.repo.UpsertDeviceStatus(ctx, repository.UpsertDeviceStatusParams{
		DeviceID:  req.DeviceID,
		Battery:   req.Battery,
		LidarOk:   req.LidarOk,
		CameraOk:  req.CameraOk,
		GpsOk:     req.GpsOk,
		NetworkOk: req.NetworkOk,
	})
	if err != nil {
		return domain.DeviceStatusResponse{}, domain.ErrInternalError
	}
	return toDeviceStatusResponse(ds), nil
}

func (s *deviceStatusService) Get(ctx context.Context, deviceID string) (domain.DeviceStatusResponse, error) {
	ds, err := s.repo.GetLatestDeviceStatus(ctx, deviceID)
	if err != nil {
		return domain.DeviceStatusResponse{}, domain.ErrNotFound
	}
	return toDeviceStatusResponse(ds), nil
}

// GetForUser는 deviceID가 userID 소유인지 먼저 검증한 뒤 상태를 반환한다.
// (Get은 Embedded API 전용 — JWT 없이 호출되므로 소유권 검증이 없다)
func (s *deviceStatusService) GetForUser(ctx context.Context, userID, deviceID string) (domain.DeviceStatusResponse, error) {
	if _, err := s.repo.GetDevice(ctx, repository.GetDeviceParams{ID: deviceID, UserID: userID}); err != nil {
		return domain.DeviceStatusResponse{}, domain.ErrNotFound
	}
	return s.Get(ctx, deviceID)
}

func toDeviceStatusResponse(ds repository.DeviceStatus) domain.DeviceStatusResponse {
	return domain.DeviceStatusResponse{
		DeviceID:  ds.DeviceID,
		Battery:   ds.Battery,
		LidarOk:   ds.LidarOk,
		CameraOk:  ds.CameraOk,
		GpsOk:     ds.GpsOk,
		NetworkOk: ds.NetworkOk,
		CreatedAt: ds.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

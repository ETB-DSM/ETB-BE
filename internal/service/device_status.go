package service

import (
	"context"

	"github.com/Heiji57/ETB-BE/internal/domain"
	repository "github.com/Heiji57/ETB-BE/internal/repository/sqlc"
)

type DeviceStatusService interface {
	Update(ctx context.Context, req domain.UpdateDeviceStatusRequest) (domain.DeviceStatusResponse, error)
	Get(ctx context.Context, deviceID string) (domain.DeviceStatusResponse, error)
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

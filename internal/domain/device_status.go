package domain

type UpdateDeviceStatusRequest struct {
	DeviceID  string `json:"deviceId"  binding:"required"`
	Battery   int32  `json:"battery"   binding:"required,min=0,max=100"`
	LidarOk   bool   `json:"lidarOk"`
	CameraOk  bool   `json:"cameraOk"`
	GpsOk     bool   `json:"gpsOk"`
	NetworkOk bool   `json:"networkOk"`
}

type DeviceStatusResponse struct {
	DeviceID  string `json:"deviceId"`
	Battery   int32  `json:"battery"`
	LidarOk   bool   `json:"lidarOk"`
	CameraOk  bool   `json:"cameraOk"`
	GpsOk     bool   `json:"gpsOk"`
	NetworkOk bool   `json:"networkOk"`
	CreatedAt string `json:"createdAt"`
}

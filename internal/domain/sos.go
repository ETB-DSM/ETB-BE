package domain

type CreateSosRequest struct {
	DeviceID  string  `json:"deviceId"   binding:"required"`
	EventType string  `json:"eventType"  binding:"required,oneof=fall manual_sos"`
	Latitude  float64 `json:"latitude"   binding:"required"`
	Longitude float64 `json:"longitude"  binding:"required"`
}

type SosResponse struct {
	SosID     string  `json:"sosId"`
	DeviceID  string  `json:"deviceId"`
	EventType string  `json:"eventType"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	CreatedAt string  `json:"createdAt"`
}

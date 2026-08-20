package domain

type CreateLocationRequest struct {
	UserID    string  `json:"userId"    binding:"required"`
	DeviceID  string  `json:"deviceId"  binding:"required"`
	Latitude  float64 `json:"latitude"  binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	Timestamp string  `json:"timestamp" binding:"required"`
}

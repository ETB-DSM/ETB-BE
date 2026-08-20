package domain

type CreateSosRequest struct {
	DeviceID  string  `json:"deviceId"   binding:"required"`
	EventType string  `json:"eventType"  binding:"required,oneof=fall manual_sos fall_detected emergency_button"`
	Latitude  float64 `json:"latitude"   binding:"required"`
	Longitude float64 `json:"longitude"  binding:"required"`
	Battery   *int32  `json:"battery"`
	Timestamp string  `json:"timestamp"`
}

// EmbeddedCreateSosRequest is used by the Embedded API (no JWT, userId in body)
type EmbeddedCreateSosRequest struct {
	UserID    string  `json:"userId"     binding:"required"`
	DeviceID  string  `json:"deviceId"   binding:"required"`
	EventType string  `json:"eventType"  binding:"required,oneof=fall manual_sos fall_detected emergency_button"`
	Latitude  float64 `json:"latitude"   binding:"required"`
	Longitude float64 `json:"longitude"  binding:"required"`
	Battery   *int32  `json:"battery"`
	Timestamp string  `json:"timestamp"`
}

type SosResponse struct {
	SosID          string  `json:"sosId"`
	DeviceID       string  `json:"deviceId"`
	EventType      string  `json:"eventType"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	SentToGuardian bool    `json:"sentToGuardian"`
	CreatedAt      string  `json:"createdAt"`
}

package domain

type CreateSessionRequest struct {
	DeviceID       string  `json:"deviceId"       binding:"required"`
	DestinationID  string  `json:"destinationId"  binding:"required"`
	StartLatitude  float64 `json:"startLatitude"  binding:"required"`
	StartLongitude float64 `json:"startLongitude" binding:"required"`
}

type SessionResponse struct {
	SessionID      string  `json:"sessionId"`
	Status         string  `json:"status"`
	StartLatitude  float64 `json:"startLatitude"`
	StartLongitude float64 `json:"startLongitude"`
	CreatedAt      string  `json:"createdAt"`
}

type UpdateInstructionRequest struct {
	Action         string  `json:"action"         binding:"required,oneof=straight prepare_left left prepare_right right arrived reroute"`
	DistanceMeters *int32  `json:"distanceMeters"`
	Message        *string `json:"message"`
}

type InstructionResponse struct {
	InstructionID  string  `json:"instructionId"`
	SessionID      string  `json:"sessionId"`
	Action         string  `json:"action"`
	DistanceMeters *int32  `json:"distanceMeters,omitempty"`
	Message        *string `json:"message,omitempty"`
	CreatedAt      string  `json:"createdAt"`
}

type UpdateSessionStatusRequest struct {
	Status string  `json:"status" binding:"required,oneof=active paused arrived canceled error"`
	Reason *string `json:"reason"`
}

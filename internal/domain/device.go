package domain

type CreateDeviceRequest struct {
	Name string `json:"name" binding:"required"`
}

type DeviceResponse struct {
	DeviceID  string `json:"deviceId"`
	Name      string `json:"name"`
	IsActive  bool   `json:"isActive"`
	APIKey    string `json:"apiKey,omitempty"`
	CreatedAt string `json:"createdAt"`
}

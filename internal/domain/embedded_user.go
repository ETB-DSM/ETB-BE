package domain

type EmbeddedUserRegisterRequest struct {
	UserID        string `json:"userId"        binding:"required"`
	Name          string `json:"name"          binding:"required"`
	GuardianName  string `json:"guardianName"  binding:"required"`
	GuardianPhone string `json:"guardianPhone" binding:"required"`
	DeviceID      string `json:"deviceId"      binding:"required"`
}

type EmbeddedUserResponse struct {
	UserID   string `json:"userId"`
	DeviceID string `json:"deviceId"`
}

package domain

type CreateGuardianRequest struct {
	Name  string `json:"name"  binding:"required"`
	Phone string `json:"phone" binding:"required"`
}

type GuardianResponse struct {
	GuardianID string `json:"guardianId"`
	Name       string `json:"name"`
	Phone      string `json:"phone"`
}

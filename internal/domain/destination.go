package domain

type CreateDestinationRequest struct {
	Name       string  `json:"name"       binding:"required"`
	Latitude   float64 `json:"latitude"   binding:"required"`
	Longitude  float64 `json:"longitude"  binding:"required"`
	RadiusM    int32   `json:"radiusM"    binding:"required,min=1"`
	TargetText string  `json:"targetText" binding:"required"`
}

type DestinationResponse struct {
	DestinationID string  `json:"destinationId"`
	Name          string  `json:"name"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	RadiusM       int32   `json:"radiusM"`
	TargetText    string  `json:"targetText"`
	CreatedAt     string  `json:"createdAt"`
}

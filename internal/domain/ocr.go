package domain

type CreateOcrLogRequest struct {
	DestinationID  string  `json:"destinationId"  binding:"required"`
	RecognizedText string  `json:"recognizedText" binding:"required"`
	TargetText     string  `json:"targetText"     binding:"required"`
	Matched        bool    `json:"matched"`
	Confidence     float32 `json:"confidence"     binding:"required,min=0,max=1"`
}

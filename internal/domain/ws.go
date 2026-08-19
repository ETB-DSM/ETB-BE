package domain

type WSMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type LocationPayload struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp string  `json:"timestamp"`
}

type DeviceStatusPayload struct {
	Battery       int    `json:"battery"`
	SensorStatus  string `json:"sensorStatus"`
	NetworkStatus string `json:"networkStatus"`
}

package dto

type Response struct {
	Message string            `json:"message"`
	Data    interface{}       `json:"data"`
	Errors  map[string]string `json:"errors"`
	ErrorID string            `json:"errorID,omitempty"`
}

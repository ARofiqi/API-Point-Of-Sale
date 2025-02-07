package models

type Response struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Errors  []string    `json:"errors"`
	ErrorID string      `json:"errorID,omitempty"`
}

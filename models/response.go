package models

type Response struct {
	Data    interface{}         `json:"data"`
	Message string              `json:"message"`
	Errors  []map[string]string `json:"errors"`
	ErrorID string             `json:"errorID,omitempty"`
}

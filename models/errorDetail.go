package models

type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

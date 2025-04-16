package dto

type ErrorDetails map[string]string

type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

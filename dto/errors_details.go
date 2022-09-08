package dto

type ErrorDetails struct {
	Msg         string `json:"error"`
	Description string `json:"error_description,omitempty"`
}

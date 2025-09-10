package models

type ApiResponse struct {
	Success      bool    `json:"success"`
	ErrorMessage *string `json:"errorMessage"`
	Data         any     `json:"data"`
}

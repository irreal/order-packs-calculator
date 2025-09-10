package utils

import (
	"encoding/json"
	"net/http"

	"github.com/irreal/order-packs/models"
)

// TODO!: Handle serialization errors
func WriteAPIResponse(w http.ResponseWriter, status int, response models.ApiResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

func WriteAPISuccessResponse(w http.ResponseWriter, data any) {
	WriteAPIResponse(w, http.StatusOK, models.ApiResponse{
		Success: true,
		Data:    data,
	})
}

func WriteAPIErrorResponse(w http.ResponseWriter, status int, errorMessage string) {
	WriteAPIResponse(w, status, models.ApiResponse{
		Success:      false,
		ErrorMessage: &errorMessage,
	})
}

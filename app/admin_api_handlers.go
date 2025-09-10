package app

import (
	"encoding/json"
	"net/http"

	"github.com/irreal/order-packs/models"
	"github.com/irreal/order-packs/utils"
)

func (a *App) handleGetPacks(w http.ResponseWriter, r *http.Request) {
	packs, err := a.packsService.GetPacks()
	if err != nil {
		utils.WriteAPIErrorResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	utils.WriteAPISuccessResponse(w, packs)
}

func (a *App) handleSetPacks(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Packs []int `json:"packs"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.WriteAPIErrorResponse(w, http.StatusBadRequest, "invalid JSON format")
		return
	}

	if len(request.Packs) == 0 {
		utils.WriteAPIErrorResponse(w, http.StatusBadRequest, "at least one pack is required")
		return
	}

	var newPacks models.Packs
	for _, packSize := range request.Packs {
		if packSize <= 0 {
			utils.WriteAPIErrorResponse(w, http.StatusBadRequest, "pack size must be positive")
			return
		}
		newPacks = append(newPacks, models.Pack(packSize))
	}

	if err := a.packsService.SavePacks(newPacks); err != nil {
		utils.WriteAPIErrorResponse(w, http.StatusInternalServerError, "failed to save packs")
		return
	}
	utils.WriteAPISuccessResponse(w, "packs saved successfully")
}

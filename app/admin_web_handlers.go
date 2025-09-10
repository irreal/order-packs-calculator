package app

import (
	"net/http"
	"strconv"

	"github.com/irreal/order-packs/app/pages"
	"github.com/irreal/order-packs/models"
	"github.com/irreal/order-packs/utils"
)

func (a *App) handleAdminPageGet(w http.ResponseWriter, r *http.Request) {

	packs, err := a.packsService.GetPacks()
	if err != nil {
		http.Error(w, "Failed to get packs", http.StatusInternalServerError)
		return
	}

	success := r.URL.Query().Get("success") == "1"

	utils.Render(w, r, pages.AdminPage(packs, success))
}

func (a *App) handleAdminPageSetPacks(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	packValues := r.Form["packs"]
	if len(packValues) == 0 {
		http.Error(w, "At least one pack is required", http.StatusBadRequest)
		return
	}

	// convert to slice
	var newPacks models.Packs
	for _, packStr := range packValues {
		packSize, err := strconv.Atoi(packStr)
		if err != nil || packSize <= 0 {
			http.Error(w, "Invalid pack size: "+packStr, http.StatusBadRequest)
			return
		}
		newPacks = append(newPacks, models.Pack(packSize))
	}

	// persist to repo
	if err := a.packsService.SavePacks(newPacks); err != nil {
		http.Error(w, "Failed to save packs", http.StatusInternalServerError)
		return
	}

	// redirect to admin page
	http.Redirect(w, r, "/admin?success=1", http.StatusSeeOther)
}

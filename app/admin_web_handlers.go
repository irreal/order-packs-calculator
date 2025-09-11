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
		utils.Render(w, r, pages.ErrorPage(err.Error()))
		return
	}

	success := r.URL.Query().Get("success") == "1"

	utils.Render(w, r, pages.AdminPage(packs, success))
}

func (a *App) handleAdminPageSetPacks(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		utils.Render(w, r, pages.ErrorPage(err.Error()))
		return
	}

	packValues := r.Form["packs"]
	if len(packValues) == 0 {
		utils.Render(w, r, pages.ErrorPage("At least one pack is required"))
		return
	}

	// convert to slice
	var newPacks models.Packs
	for _, packStr := range packValues {
		packSize, err := strconv.Atoi(packStr)
		if err != nil || packSize <= 0 {
			utils.Render(w, r, pages.ErrorPage("Invalid pack size: "+packStr))
			return
		}
		newPacks = append(newPacks, models.Pack(packSize))
	}

	// persist to repo
	if err := a.packsService.SavePacks(newPacks); err != nil {
		utils.Render(w, r, pages.ErrorPage(err.Error()))
		return
	}

	// redirect to admin page
	http.Redirect(w, r, "/admin?success=1", http.StatusSeeOther)
}

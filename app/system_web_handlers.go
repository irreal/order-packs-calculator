package app

import (
	"net/http"

	"github.com/irreal/order-packs/app/pages"
	"github.com/irreal/order-packs/utils"
)

func (a *App) handleHomePage(w http.ResponseWriter, r *http.Request) {
	utils.Render(w, r, pages.HomePage())
}

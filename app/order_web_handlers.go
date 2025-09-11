package app

import (
	"net/http"
	"strconv"

	"github.com/irreal/order-packs/app/pages"
	"github.com/irreal/order-packs/models"
	"github.com/irreal/order-packs/utils"
)

func (a *App) handleOrderPage(w http.ResponseWriter, r *http.Request) {
	orders, err := a.orderService.GetLast10Orders()
	if err != nil {
		utils.Render(w, r, pages.ErrorPage(err.Error()))
		return
	}

	packs, err := a.packsService.GetPacks()
	if err != nil {
		utils.Render(w, r, pages.ErrorPage(err.Error()))
		return
	}

	maxCount := int32(a.orderService.MaxOrderItemCount)
	success := r.URL.Query().Get("success") == "1"

	utils.Render(w, r, pages.OrderPage(orders, packs, maxCount, success))
}

func (a *App) handleCreateOrderWeb(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		utils.Render(w, r, pages.ErrorPage(err.Error()))
		return
	}

	amountStr := r.Form.Get("amount")
	if amountStr == "" {
		utils.Render(w, r, pages.ErrorPage("Amount is required"))
		return
	}

	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		utils.Render(w, r, pages.ErrorPage("Invalid amount format"))
		return
	}

	orderRequest := models.OrderRequest{
		ItemCount: amount,
	}

	packs, err := a.packsService.GetPacks()
	if err != nil {
		utils.Render(w, r, pages.ErrorPage(err.Error()))
		return
	}

	_, err = a.orderService.CreateOrder(orderRequest, packs)
	if err != nil {
		utils.Render(w, r, pages.ErrorPage(err.Error()))
		return
	}

	// Redirect back to the order page with success
	http.Redirect(w, r, "/order?success=1", http.StatusSeeOther)
}

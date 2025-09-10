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
		//render error page
	}

	packs, err := a.packsService.GetPacks()
	if err != nil {
		//render error page
	}

	maxCount := int32(a.orderService.MaxOrderItemCount)
	success := r.URL.Query().Get("success") == "1"

	utils.Render(w, r, pages.OrderPage(orders, packs, maxCount, success))
}

func (a *App) handleCreateOrderWeb(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	amountStr := r.Form.Get("amount")
	if amountStr == "" {
		http.Error(w, "Amount is required", http.StatusBadRequest)
		return
	}

	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		http.Error(w, "Invalid amount format", http.StatusBadRequest)
		return
	}

	orderRequest := models.OrderRequest{
		ItemCount: amount,
	}

	packs, err := a.packsService.GetPacks()
	if err != nil {
		http.Error(w, "Failed to get available packs", http.StatusInternalServerError)
		return
	}

	_, err = a.orderService.CreateOrder(orderRequest, packs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Redirect back to the order page with success
	http.Redirect(w, r, "/order?success=1", http.StatusSeeOther)
}

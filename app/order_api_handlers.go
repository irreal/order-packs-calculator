package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/irreal/order-packs/models"
	"github.com/irreal/order-packs/orders"
	"github.com/irreal/order-packs/utils"
)

func (a *App) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	var orderRequest models.OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&orderRequest); err != nil {
		utils.WriteAPIErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Invalid JSON request: %v", err))
		return
	}

	order, err := a.orderService.CreateOrder(orderRequest, []models.Pack{250, 500, 1000, 2000, 5000})
	if err != nil {
		fmt.Fprintf(a.stderr, "error creating order: %v\n", err)

		// customize response code based on error type
		if errors.Is(err, orders.InvalidOrderItemCountError) {
			utils.WriteAPIErrorResponse(w, http.StatusBadRequest, err.Error())
		} else {
			utils.WriteAPIErrorResponse(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	utils.WriteAPISuccessResponse(w, order)
}

package app

import (
	"net/http"

	"github.com/irreal/order-packs/utils"
)

func (a *App) handleHealth(w http.ResponseWriter, r *http.Request) {
	// would check operational stuff such as db online, etc.
	utils.WriteAPISuccessResponse(w, map[string]string{"api_status": "ok"})
}

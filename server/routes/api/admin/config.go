package adminApi

import (
	"encoding/json"
	"fmt"
	"localapps/utils"
	"net/http"
)

func getConfig(w http.ResponseWriter, r *http.Request) {
	config := utils.CachedConfig

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(config); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON: %s", err), http.StatusInternalServerError)
		return
	}
}

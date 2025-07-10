package api

import (
	"encoding/json"
	"localapps/constants"
	"localapps/types"
	"localapps/utils"
	"net/http"
)

func Link(w http.ResponseWriter, r *http.Request) {
	apiKey := utils.ServerConfig

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(apiKey); err != nil {
		response := types.ApiError{
			Code:    constants.ErrorEncode,
			Message: "Error while encoding JSON",
			Error:   err,
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
}

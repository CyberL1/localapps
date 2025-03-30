package middlewares

import (
	"encoding/json"
	"localapps/types"
	"localapps/utils"
	"net/http"
)

func ApiAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			noAccessError := types.ApiError{
				Code:    "ACCESS_DENIED",
				Message: "Missing API Key",
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(noAccessError)
			return
		}

		if r.Header.Get("Authorization") != utils.CachedConfig.ApiKey {
			noAccessError := types.ApiError{
				Code:    "ACCESS_DENIED",
				Message: "Invalid API Key",
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(noAccessError)
			return
		}

		next.ServeHTTP(w, r)
	})
}

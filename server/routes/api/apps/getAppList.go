package appsApi

import (
	"context"
	"encoding/json"
	"localapps/constants"
	dbClient "localapps/db/client"
	"localapps/types"
	"net/http"
)

func getAppList(w http.ResponseWriter, r *http.Request) {
	db, _ := dbClient.GetClient()
	apps, err := db.ListApps(context.Background())
	if err != nil {
		response := types.ApiError{
			Code:    constants.ErrorDb,
			Message: "Error while fetching apps",
			Error:   err,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	var list []types.ApiAppResponse
	for _, appData := range apps {
		list = append(list, types.ApiAppResponse{
			Id:          appData.ID,
			AppId:       appData.AppID,
			Name:        appData.Name,
			Icon:        appData.Icon,
			InstalledAt: appData.InstalledAt.String(),
			Parts: func() map[string]string {
				var parts map[string]string
				if err := json.Unmarshal([]byte(appData.Parts), &parts); err != nil {
					parts = make(map[string]string) // default to empty map on error
				}
				return parts
			}(),
		})
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(list); err != nil {
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

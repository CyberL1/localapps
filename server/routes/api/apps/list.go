package appsApi

import (
	"context"
	"encoding/json"
	"fmt"
	dbClient "localapps/db/client"
	"localapps/types"
	"localapps/utils"
	"net/http"
)

func getAppList(w http.ResponseWriter, r *http.Request) {
	db, _ := dbClient.GetClient()
	apps, err := db.ListApps(context.Background())
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching apps: %s", err), http.StatusInternalServerError)
		return
	}

	var list []types.ApiAppListResponse

	for _, appData := range apps {
		app, err := utils.GetAppData(appData.ID)

		if err != nil {
			continue
		}

		list = append(list, types.ApiAppListResponse{
			Id:          appData.ID,
			Name:        app.Name,
			Icon:        app.Icon,
			InstalledAt: appData.InstalledAt.Time.String(),
			Parts: func() map[string]string {
				var parts map[string]string
				if err := json.Unmarshal(appData.Parts, &parts); err != nil {
					parts = make(map[string]string) // default to empty map on error
				}
				return parts
			}(),
		})
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(list); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON: %s", err), http.StatusInternalServerError)
		return
	}
}

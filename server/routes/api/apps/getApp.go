package appsApi

import (
	"context"
	"encoding/json"
	"fmt"
	dbClient "localapps/db/client"
	"localapps/types"
	"net/http"
)

func getApp(w http.ResponseWriter, r *http.Request) {
	db, _ := dbClient.GetClient()

	app, err := db.GetApp(context.Background(), r.PathValue("appId"))
	if err != nil {
		http.Error(w, fmt.Sprintf("App \"%s\" not found", r.PathValue("appId")), http.StatusInternalServerError)
		return
	}

	response := types.ApiAppResponse{
		Id:          app.ID,
		Name:        app.Name,
		Icon:        app.Icon,
		InstalledAt: app.InstalledAt.Time.String(),
		Parts: func() map[string]string {
			var parts map[string]string
			if err := json.Unmarshal(app.Parts, &parts); err != nil {
				parts = make(map[string]string) // default to empty map on error
			}
			return parts
		}(),
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON: %s", err), http.StatusInternalServerError)
		return
	}
}

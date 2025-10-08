package appsApi

import (
	"context"
	"encoding/json"
	"fmt"
	"localapps/constants"
	dbClient "localapps/db/client"
	"localapps/types"
	"net/http"
)

func getApp(w http.ResponseWriter, r *http.Request) {
	db, _ := dbClient.GetClient()

	app, err := db.GetAppByAppId(context.Background(), r.PathValue("appId"))
	if err != nil {
		response := types.ApiError{
			Code:    constants.ErrorNotFound,
			Message: fmt.Sprintf("App \"%s\" not found", r.PathValue("appId")),
			Error:   err,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := types.ApiAppResponse{
		Id:          app.ID,
		AppId:       app.AppID,
		Name:        app.Name,
		Icon:        app.Icon,
		InstalledAt: app.InstalledAt.String(),
		Parts: func() map[string]string {
			var parts map[string]string
			if err := json.Unmarshal([]byte(app.Parts), &parts); err != nil {
				parts = make(map[string]string) // default to empty map on error
			}
			return parts
		}(),
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		response := types.ApiError{
			Code:    constants.ErrorEncode,
			Message: "Error while encoding JSON",
			Error:   err,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
}

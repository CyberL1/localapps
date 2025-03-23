package appsApi

import (
	"encoding/json"
	"fmt"
	dbClient "localapps/db/client"
	"localapps/types"
	"net/http"
	"strings"
)

func uninstallApp(w http.ResponseWriter, r *http.Request) {
	client, _ := dbClient.GetClient()

	var body types.ApiAppUninstallRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %s", err), http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(body.Id) == "" {
		http.Error(w, "App ID is required", http.StatusBadRequest)
		return
	}

	appId := body.Id

	if strings.Contains(appId, " ") {
		http.Error(w, "App ID cannot contain spaces", http.StatusBadRequest)
		return
	}

	_, err := client.GetApp(dbClient.Ctx, appId)
	if err != nil {
		http.Error(w, "App not installed", http.StatusInternalServerError)
		return
	}

	err = client.DeleteApp(dbClient.Ctx, appId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting app record from db: %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

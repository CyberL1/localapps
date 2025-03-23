package appsApi

import (
	"encoding/json"
	"fmt"
	dbClient "localapps/db/client"
	db "localapps/db/generated"
	"localapps/types"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func installApp(w http.ResponseWriter, r *http.Request) {
	client, _ := dbClient.GetClient()

	var body types.ApiAppInstallRequestBody
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

	appWithTheSameId, _ := client.GetApp(dbClient.Ctx, appId)

	if appWithTheSameId.ID == appId {
		http.Error(w, "App already installed", http.StatusInternalServerError)
		return
	}

	_, err := client.CreateApp(dbClient.Ctx, db.CreateAppParams{
		ID:          appId,
		InstalledAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
	})

	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating app record in db: %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

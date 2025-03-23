package appsApi

import (
	"encoding/json"
	"fmt"
	"io"
	dbClient "localapps/db/client"
	db "localapps/db/generated"
	"localapps/types"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"gopkg.in/yaml.v2"
)

func installApp(w http.ResponseWriter, r *http.Request) {
	appInfoFile, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading uploaded file: %s", err), http.StatusBadRequest)
		return
	}
	defer appInfoFile.Close()

	fileContent, err := io.ReadAll(appInfoFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading file content: %s", err), http.StatusInternalServerError)
		return
	}

	var appInfo types.App
	if err := yaml.Unmarshal(fileContent, &appInfo); err != nil {
		http.Error(w, fmt.Sprintf("Error unmarshaling YAML content: %s", err), http.StatusBadRequest)
		return
	}

	var appId string

	if appInfo.Id != "" {
		appId = appInfo.Id
	} else {
		appId = strings.ToLower(strings.ReplaceAll(appInfo.Name, " ", "-"))
	}

	if strings.Contains(appId, " ") {
		http.Error(w, "App ID cannot contain spaces", http.StatusBadRequest)
		return
	}

	client, _ := dbClient.GetClient()
	appWithTheSameId, _ := client.GetApp(dbClient.Ctx, appId)

	if appWithTheSameId.ID == appId {
		http.Error(w, "App already installed", http.StatusInternalServerError)
		return
	}

	var appParts = make(map[string]string)
	for name, part := range appInfo.Parts {
		appParts[name] = part.Path
	}

	partsJson, err := json.Marshal(appParts)
	if err != nil {
		fmt.Println("Error marshaling map to JSON:", err)
	}

	_, err = client.CreateApp(dbClient.Ctx, db.CreateAppParams{
		ID:          appId,
		InstalledAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
		Name:        appInfo.Name,
		Parts:       partsJson,
	})

	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating app record in db: %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

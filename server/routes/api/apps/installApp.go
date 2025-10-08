package appsApi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"localapps/constants"
	dbClient "localapps/db/client"
	db "localapps/db/generated"
	"localapps/types"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
)

func installApp(w http.ResponseWriter, r *http.Request) {
	appInfoFile, _, err := r.FormFile("file")
	if err != nil {
		response := types.ApiError{
			Code:    constants.ErrorNotFound,
			Message: "Field 'file' not found",
			Error:   err,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	defer appInfoFile.Close()

	fileContent, err := io.ReadAll(appInfoFile)
	if err != nil {
		response := types.ApiError{
			Code:    constants.ErrorRead,
			Message: "Error while reading file content",
			Error:   err,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	var appInfo types.App
	if err := yaml.Unmarshal(fileContent, &appInfo); err != nil {
		response := types.ApiError{
			Code:    constants.ErrorParse,
			Message: "Error while parsing YAML content",
			Error:   err,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if appInfo.Name == "" {
		response := types.ApiError{
			Code:    constants.ErrorFieldRequired,
			Message: "Field 'name' required",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if appInfo.Id == "" {
		appInfo.Id = strings.ToLower(strings.ReplaceAll(appInfo.Name, " ", "-"))
	}

	if appInfo.Parts == nil {
		respose := types.ApiError{
			Code:    constants.ErrorFieldRequired,
			Message: "Field 'parts' required",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(respose)
		return
	}

	if strings.Contains(appInfo.Id, " ") {
		respose := types.ApiError{
			Code:    constants.ErrorParse,
			Message: "App ID cannot contain spaces",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respose)
		return
	}

	client, _ := dbClient.GetClient()
	appWithTheSameId, _ := client.GetAppByAppId(context.Background(), appInfo.Id)

	if r.FormValue("update") != "true" && appWithTheSameId.AppID == appInfo.Id {
		response := types.ApiError{
			Code:    constants.ErrorAppInstalled,
			Message: "App already installed",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	var appParts = make(map[string]string)
	for name, part := range appInfo.Parts {
		appParts[name] = part.Path
	}

	partsJson, err := json.Marshal(appParts)
	if err != nil {
		response := types.ApiError{
			Code:    constants.ErrorParse,
			Message: "Error while marshaling map to JSON",
			Error:   err,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if appInfo.Icon != "" {
		appInfo.Icon = strings.ReplaceAll(uuid.NewString(), "-", "")

		if _, err := os.Stat(constants.LocalappsAppIconsDir); errors.Is(err, fs.ErrNotExist) {
			os.MkdirAll(constants.LocalappsAppIconsDir, 0755)
		}

		os.Create(filepath.Join(constants.LocalappsAppIconsDir, appInfo.Icon))

		iconFile, _, err := r.FormFile("icon")
		if err != nil {
			response := types.ApiError{
				Code:    constants.ErrorRead,
				Message: "Error while reading app icon",
				Error:   err,
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
		defer iconFile.Close()

		iconData, err := io.ReadAll(iconFile)
		if err != nil {
			response := types.ApiError{
				Code:    constants.ErrorRead,
				Message: "Error while reading app icon",
				Error:   err,
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		if http.DetectContentType(iconData[:512]) != "image/png" {
			response := types.ApiError{
				Code:    constants.ErrorInvalidIcon,
				Message: "App icon must be a png file",
			}

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		os.WriteFile(filepath.Join(constants.LocalappsAppIconsDir, appInfo.Icon), iconData, 0644)
	}

	if r.FormValue("update") == "true" {
		if appWithTheSameId.Icon != "" {
			os.Remove(filepath.Join(constants.LocalappsAppIconsDir, appWithTheSameId.Icon))
		}

		_, err = client.UpdateApp(context.Background(), db.UpdateAppParams{
			AppID: appInfo.Id,
			Name:  appInfo.Name,
			Parts: string(partsJson),
			Icon:  appInfo.Icon,
		})
	} else {
		_, err = client.CreateApp(context.Background(), db.CreateAppParams{
			AppID:       appInfo.Id,
			InstalledAt: time.Now(),
			Name:        appInfo.Name,
			Parts:       string(partsJson),
			Icon:        appInfo.Icon,
		})
	}

	if err != nil {
		response := types.ApiError{
			Code:    constants.ErrorEncode,
			Message: "Error while creating DB record",
			Error:   err,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

package api

import (
	"encoding/json"
	"fmt"
	"localapps/constants"
	"localapps/types"
	"localapps/utils"
	"net/http"
	"os"
	"path/filepath"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes() *http.ServeMux {
	r := http.NewServeMux()

	r.HandleFunc("/apps", appList)

	return r
}

func appList(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(filepath.Join(constants.LocalappsDir, "apps"))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading directory: %s", err), http.StatusInternalServerError)
		return
	}

	var list []types.AppNameWithSubdomain

	for _, file := range files {
		app, err := utils.GetApp(file.Name())

		if err != nil {
			continue
		}

		list = append(list, types.AppNameWithSubdomain{
			Name:      app.Name,
			Subdomain: file.Name(),
		})
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(list); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON: %s", err), http.StatusInternalServerError)
		return
	}
}

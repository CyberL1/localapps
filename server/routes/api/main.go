package api

import (
	"encoding/json"
	"fmt"
	dbClient "localapps/db/client"
	adminApi "localapps/server/routes/api/admin"
	"localapps/types"
	"localapps/utils"
	"net/http"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes() *http.ServeMux {
	r := http.NewServeMux()

	r.Handle("/admin/", http.StripPrefix("/admin", adminApi.NewHandler().RegisterRoutes()))

	r.HandleFunc("/apps", appList)

	return r
}

func appList(w http.ResponseWriter, r *http.Request) {
	db, _ := dbClient.GetClient()
	apps, err := db.ListApps(dbClient.Ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching apps: %s", err), http.StatusInternalServerError)
		return
	}

	var list []types.ApiAppResponse

	for _, appData := range apps {
		app, err := utils.GetApp(appData.ID)

		if err != nil {
			continue
		}

		list = append(list, types.ApiAppResponse{
			Id:          appData.ID,
			Name:        app.Name,
			InstalledAt: appData.InstalledAt.Time.String(),
		})
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(list); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON: %s", err), http.StatusInternalServerError)
		return
	}
}

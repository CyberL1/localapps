package routes

import (
	"encoding/json"
	"fmt"
	"localapps/resources"
	"localapps/types"
	"localapps/utils"
	"net/http"
	"path/filepath"
	"text/template"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes() *http.ServeMux {
	r := http.NewServeMux()

	r.HandleFunc("/", homePage)

	return r
}

func homePage(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://localhost:8080/api/apps")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching apps: %s", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var list []*types.ApiAppResponse
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding response: %s", err), http.StatusInternalServerError)
		return
	}

	fileContent, err := resources.Resources.ReadFile(filepath.Join("pages", "home.html"))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading file: %s", err), http.StatusInternalServerError)
		return
	}

	templ, err := template.New("home.html").Parse(string(fileContent))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	data := struct {
		Config types.Config
		Apps   []*types.ApiAppResponse
	}{
		Config: utils.CachedConfig,
		Apps:   list,
	}

	if err = template.Must(templ.Clone()).Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

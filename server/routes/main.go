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
	req, err := http.NewRequest("GET", "http://localhost:8080/api/apps/", nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating request: %s", err), http.StatusInternalServerError)
		return
	}
	req.Header.Add("Authorization", utils.ServerConfig.ApiKey)

	resp, err := http.DefaultClient.Do(req)
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
		Config types.ServerConfig
		Apps   []*types.ApiAppResponse
	}{
		Config: utils.ServerConfig,
		Apps:   list,
	}

	if err = template.Must(templ.Clone()).Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

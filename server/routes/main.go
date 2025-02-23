package routes

import (
	"encoding/json"
	"fmt"
	"localapps/constants"
	"localapps/types"
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

	var list []*types.AppNameWithSubdomain
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding response: %s", err), http.StatusInternalServerError)
		return
	}

	templ, err := template.New("home.html").ParseFiles(filepath.Join(constants.LocalappsDir, "pages", "home.html"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if err = template.Must(templ.Clone()).Execute(w, list); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

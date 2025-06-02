package routes

import (
	"localapps/resources"
	"mime"
	"net/http"
	"path/filepath"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes() *http.ServeMux {
	r := http.NewServeMux()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[1:]
		if path == "" {
			path = "index"
		}

		ext := filepath.Ext(r.URL.Path)
		if ext == "" {
			ext = ".html"
			path += ext
		}

		var fileContent []byte
		fileContent, err := resources.Resources.ReadFile(filepath.Join("pages", path))
		if err != nil {
			fileContent, _ = resources.Resources.ReadFile(filepath.Join("pages", "fallback.html"))
		}

		w.Header().Set("Content-Type", mime.TypeByExtension(ext))
		w.Write(fileContent)
	})

	return r
}

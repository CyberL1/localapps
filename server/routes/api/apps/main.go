package appsApi

import "net/http"

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes() *http.ServeMux {
	r := http.NewServeMux()

	r.HandleFunc("GET /", getAppList)
	r.HandleFunc("GET /{appId}", getApp)
	r.HandleFunc("POST /", installApp)
	r.HandleFunc("DELETE /{appId}", uninstallApp)

	return r
}

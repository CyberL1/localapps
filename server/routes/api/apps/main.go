package appsApi

import (
	"net/http"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes() *http.ServeMux {
	r := http.NewServeMux()

	r.HandleFunc("GET /list", getAppList)
	r.HandleFunc("PUT /install", installApp)
	r.HandleFunc("PUT /uninstall", uninstallApp)

	return r
}

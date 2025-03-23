package api

import (
	adminApi "localapps/server/routes/api/admin"
	appsApi "localapps/server/routes/api/apps"
	"net/http"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes() *http.ServeMux {
	r := http.NewServeMux()

	r.Handle("/admin/", http.StripPrefix("/admin", adminApi.NewHandler().RegisterRoutes()))
	r.Handle("/apps/", http.StripPrefix("/apps", appsApi.NewHandler().RegisterRoutes()))

	return r
}

package iconsApi

import "net/http"

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes() *http.ServeMux {
	r := http.NewServeMux()
	
	r.HandleFunc("GET /apps/{icon}", getAppIcon)
	
	return r
}

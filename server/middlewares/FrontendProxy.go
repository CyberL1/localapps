package middlewares

import (
	"localapps/utils"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func FrontendProxy(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip if the request is for an app
		if len(strings.Split(r.Host, ".")) == len(strings.Split(utils.ServerConfig.Domain, "."))+1 {
			next.ServeHTTP(w, r)
			return
		}

		if !strings.HasPrefix(r.URL.Path, "/api") {
			frontendUrl, _ := url.Parse("http://localhost:8081")
			httputil.NewSingleHostReverseProxy(frontendUrl).ServeHTTP(w, r)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

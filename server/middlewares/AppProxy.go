package middlewares

import (
	"fmt"
	"localapps/types"
	"localapps/utils"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func AppProxy(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.Host, "localhost:8080") {
			if strings.HasPrefix(r.URL.Path, "/api") {
				next.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusForbidden)
			}
			return
		}

		if len(strings.Split(r.Host, ".")) == 3 {

			appName := strings.Split(r.Host, ".")[0]
			app, err := utils.GetApp(appName)

			if err != nil {
				w.Write([]byte(fmt.Sprintf("App \"%s\" not found", appName)))
				return
			}

			// Check if Docker is installed
			if _, err := exec.LookPath("docker"); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Docker is not installed"))
				return
			}

			// Check if Docker daemon is running
			if err := exec.Command("docker", "info").Run(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Docker daemon is not running"))
				return
			}

			var dockerAppName string
			var dockerImageName string

			var currentPart types.Part

			var fallbackPart types.Part
			var fallbackPartName string

			for partName, part := range app.Parts {
				if part.Path == "" {
					fallbackPart = part
					fallbackPartName = partName
				}

				if strings.Split(r.URL.Path, "/")[1] == part.Path {
					currentPart = part
					dockerAppName = "localapps-app-" + strings.ToLower(app.Name) + "-" + partName
					dockerImageName = "localapps/" + strings.ToLower(app.Name) + "/" + partName
					break
				} else {
					currentPart = fallbackPart
					dockerAppName = "localapps-app-" + strings.ToLower(app.Name) + "-" + fallbackPartName
					dockerImageName = "localapps/" + strings.ToLower(app.Name) + "/" + fallbackPartName
				}
			}

			findProcess := exec.Command("docker", "ps", "--format", "{{.Names}}", "-f", "name="+dockerAppName+"$")
			out, _ := findProcess.Output()

			var freePort int
			if strings.TrimSpace(string(out)) == dockerAppName {
				portCmd := exec.Command("docker", "port", dockerAppName, "80")
				out, _ := portCmd.Output()
				freePort, _ = strconv.Atoi(strings.TrimSpace(string(out)[8:]))
			} else {
				freePort, _ = utils.GetFreePort()

				runCmd := exec.Command("docker", "run", "--rm", "--name", dockerAppName, "-p", strconv.Itoa(freePort)+":80", "-e", "PORT=80", dockerImageName)
				runCmd.Dir = filepath.Join(utils.GetAppDirectory(appName), currentPart.Src)

				fmt.Println("[app:"+app.Name+"]", "Got a http request while stopped - starting")

				if err := runCmd.Start(); err != nil {
					w.Write([]byte(fmt.Sprintf("Failed to start app \"%s\": %s", appName, err)))
					return
				}

				go func() {
					time.Sleep(30 * time.Second)
					stopCmd := exec.Command("docker", "stop", dockerAppName)
					fmt.Println("[app:"+app.Name+"]", "Exceeded timeout (30s) - stopping")
					stopCmd.Start()
				}()
			}

			// Wait for the app to be ready
			for {
				_, err := http.Get(fmt.Sprintf("http://localhost:%d", freePort))
				if err == nil {
					break
				}
				time.Sleep(500 * time.Millisecond)
			}

			appUrl, _ := url.Parse(fmt.Sprintf("http://localhost:%d", freePort))
			httputil.NewSingleHostReverseProxy(appUrl).ServeHTTP(w, r)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

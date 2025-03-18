package middlewares

import (
	"context"
	"fmt"
	"localapps/utils"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
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

			cli, err := client.NewClientWithOpts(client.FromEnv)
			if err != nil {
				w.Write([]byte(fmt.Sprintf("Failed to connect to Docker daemon: %s", err)))
				return
			}

			var dockerAppName string
			var dockerImageName string
			var fallbackPartName string

			appId := strings.Split(r.Host, ".")[0]

			for partName, part := range app.Parts {
				if part.Path == "" {
					fallbackPartName = partName
				}

				if strings.Split(r.URL.Path, "/")[1] == part.Path {
					dockerAppName = "localapps-app-" + appId + "-" + partName
					dockerImageName = "localapps/apps/" + appId + "/" + partName
					break
				} else {
					dockerAppName = "localapps-app-" + appId + "-" + fallbackPartName
					dockerImageName = "localapps/apps/" + appId + "/" + fallbackPartName
				}
			}

			containersByName, _ := cli.ContainerList(context.Background(), container.ListOptions{
				Filters: filters.NewArgs(
					filters.Arg("name", dockerAppName),
				),
			})

			var freePort int
			if len(containersByName) > 0 {
				appContainer := containersByName[0]

				containerPort := appContainer.Ports[0].PublicPort
				freePort = int(containerPort)
			} else {
				freePort, _ = utils.GetFreePort()

				config := container.Config{
					Image: dockerImageName,
					Env:   []string{"PORT=80"},
					ExposedPorts: nat.PortSet{
						"80": struct{}{},
					},
				}

				hostConfig := container.HostConfig{
					PortBindings: nat.PortMap{
						"80": {
							{
								HostIP:   "0.0.0.0",
								HostPort: strconv.Itoa(freePort),
							},
						},
					}, Binds: []string{fmt.Sprintf("%s:/storage", filepath.Join(utils.GetAppDirectory(appId), "storage"))},
				}

				createdContainer, _ := cli.ContainerCreate(context.Background(), &config, &hostConfig, nil, nil, dockerAppName)
				fmt.Println("[app:"+appId+"]", "Got a http request while stopped - creating container")

				if err := cli.ContainerStart(context.Background(), createdContainer.ID, container.StartOptions{}); err != nil {
					w.Write([]byte(fmt.Sprintf("Failed to start app \"%s\": %s", appName, err)))
					return
				}

				go func() {
					time.Sleep(30 * time.Second)

					fmt.Println("[app:"+appId+"]", "Exceeded timeout (30s) - removing container")
					cli.ContainerRemove(context.Background(), createdContainer.ID, container.RemoveOptions{
						Force: true,
					})
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

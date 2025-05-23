package middlewares

import (
	"context"
	"fmt"
	"localapps/utils"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func AppProxy(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.Host, utils.ServerConfig.Domain) {
			if strings.HasPrefix(r.URL.Path, "/api") {
				next.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusForbidden)
			}
			return
		}

		if len(strings.Split(r.Host, ".")) == len(strings.Split(utils.ServerConfig.Domain, "."))+1 {
			appId := strings.Split(r.Host, ".")[0]

			appData, err := utils.GetAppData(appId)
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}

			cli, _ := client.NewClientWithOpts(client.FromEnv)

			_, err = cli.Ping(context.Background())
			if err != nil {
				w.Write([]byte(fmt.Sprintf("Failed to connect to Docker engine: %s", err)))
				return
			}

			var dockerAppName string
			var dockerImageName string
			var fallbackPartName string

			for name, path := range appData.Parts {
				if path == "" {
					fallbackPartName = name
				}

				if strings.Split(r.URL.Path, "/")[1] == path {
					dockerAppName = "localapps-app-" + appId + "-" + name
					dockerImageName = "localapps/apps/" + appId + "/" + name
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
					AutoRemove: true,
					Mounts:     []mount.Mount{{Type: mount.TypeVolume, Source: "localapps-storage-" + appId, Target: "/storage"}},
					PortBindings: nat.PortMap{
						"80": {
							{
								HostIP:   "0.0.0.0",
								HostPort: strconv.Itoa(freePort),
							},
						},
					},
				}

				appNameWithPart := strings.Replace(strings.TrimPrefix(dockerAppName, "localapps-app-"), "-", ":", 1)

				createdContainer, _ := cli.ContainerCreate(context.Background(), &config, &hostConfig, nil, nil, dockerAppName)
				fmt.Println("[app:"+appNameWithPart+"]", "Got a http request while stopped - creating container")

				if err := cli.ContainerStart(context.Background(), createdContainer.ID, container.StartOptions{}); err != nil {
					w.Write([]byte(fmt.Sprintf("Failed to start app \"%s\": %s", appId, err)))
					return
				}

				go func() {
					time.Sleep(20 * time.Second)

					fmt.Println("[app:"+appNameWithPart+"]", "Exceeded timeout (20s) - stopping container")
					cli.ContainerStop(context.Background(), createdContainer.ID, container.StopOptions{})
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

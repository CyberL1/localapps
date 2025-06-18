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
		var appId string
		var appEnvironmentVars []string

		if len(strings.Split(r.Host, ".")) == len(strings.Split(utils.ServerConfig.AccessUrl, "."))+1 {
			appId = strings.Split(r.Host, ".")[0]
		} else {
			appId = "home"
			appEnvironmentVars = append(appEnvironmentVars, "ORIGIN="+utils.ServerConfig.AccessUrl)

			if strings.HasPrefix(r.URL.Path, "/api") {
				next.ServeHTTP(w, r)
				return
			}
		}

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

		var currentPartName string
		var fallbackPartName string

		for name, path := range appData.Parts {
			if path == "" {
				fallbackPartName = name
			}

			if strings.Split(r.URL.Path, "/")[1] == path {
				currentPartName = name
				break
			} else {
				currentPartName = fallbackPartName
			}
		}

		containersByLabels, _ := cli.ContainerList(context.Background(), container.ListOptions{
			Filters: filters.NewArgs(
				filters.Arg("label", "LOCALAPPS_APP_ID="+appId),
				filters.Arg("label", "LOCALAPPS_APP_PART="+currentPartName),
			),
		})

		var freePort int
		if len(containersByLabels) > 0 {
			appContainer := containersByLabels[0]

			containerPort := appContainer.Ports[0].PublicPort
			freePort = int(containerPort)
		} else {
			freePort, _ = utils.GetFreePort()

			config := container.Config{
				Image: "localapps/apps/" + appId + "/" + currentPartName,
				Env:   append(appEnvironmentVars, "PORT=80"),
				ExposedPorts: nat.PortSet{
					"80": struct{}{},
				},
				Labels: map[string]string{
					"LOCALAPPS_APP_ID":   appId,
					"LOCALAPPS_APP_PART": currentPartName,
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
				NetworkMode: "localapps-network",
			}

			appNameWithPart := appId + "/" + currentPartName

			createdContainer, _ := cli.ContainerCreate(context.Background(), &config, &hostConfig, nil, nil, "")
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

		next.ServeHTTP(w, r)
	})
}

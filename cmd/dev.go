package cmd

import (
	"context"
	"fmt"
	"io"
	"localapps/types"
	"localapps/utils"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func init() {
	rootCmd.AddCommand(devCmd)
}

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Start localapps server",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		currentDir, _ := os.Getwd()

		appFilePath := filepath.Join(currentDir, "app.yml")
		file, err := os.Open(appFilePath)

		if err != nil {
			cmd.PrintErrln("No app.yml file detected")
			return
		}
		defer file.Close()

		appFileContents, err := io.ReadAll(file)
		if err != nil {
			cmd.PrintErrf("failed to read app file: %v\n", err)
		}

		var app types.App
		err = yaml.Unmarshal(appFileContents, &app)
		if err != nil {
			cmd.PrintErrf("failed to parse app file: %v\n", err)
		}

		cmd.Println("Running app in development mode")
		cli, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			fmt.Println([]byte(fmt.Sprintf("Failed to connect to Docker daemon: %s", err)))
			return
		}

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			var appId string
			if app.Id != "" {
				appId = app.Id
			} else {
				appId = strings.ToLower(strings.ReplaceAll(app.Name, " ", "-"))
			}

			var dockerAppName string
			var dockerImageName string

			var currentPartSource string
			var fallbackPartSource string

			var fallbackPartName string

			for name, part := range app.Parts {
				if part.Path == "" {
					fallbackPartName = name
					fallbackPartSource = part.Src
				}

				if strings.Split(r.URL.Path, "/")[1] == part.Path {
					dockerAppName = "localapps-app-" + appId + "-" + name
					dockerImageName = "localapps/apps/" + appId + "/" + name
					currentPartSource = part.Src
					break
				} else {
					dockerAppName = "localapps-app-" + appId + "-" + fallbackPartName
					dockerImageName = "localapps/apps/" + appId + "/" + fallbackPartName
					currentPartSource = fallbackPartSource
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
					Mounts: []mount.Mount{{Type: mount.TypeVolume, Source: "localapps-storage-" + appId, Target: "/storage"},
						{Type: mount.TypeBind, Source: filepath.Join(currentDir, currentPartSource), Target: "/app"}},
					PortBindings: nat.PortMap{
						"80": {
							{
								HostIP:   "0.0.0.0",
								HostPort: strconv.Itoa(freePort),
							},
						},
					},
				}

				createdContainer, _ := cli.ContainerCreate(context.Background(), &config, &hostConfig, nil, nil, dockerAppName)

				if err := cli.ContainerStart(context.Background(), createdContainer.ID, container.StartOptions{}); err != nil {
					w.Write([]byte(fmt.Sprintf("Failed to start app \"%s\": %s", appId, err)))
					return
				}
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
		})

		// Exit handler
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-stop

			cmd.Println("Stopping development containers")

			containers, err := cli.ContainerList(context.Background(), container.ListOptions{})
			if err != nil {
				cmd.Printf("Failed to list containers: %s\n", err)
				return
			}

			for _, c := range containers {
				if strings.HasPrefix(c.Names[0], "/localapps-app-") {
					if err := cli.ContainerRemove(context.Background(), c.ID, container.RemoveOptions{Force: true}); err != nil {
						cmd.Printf("Failed to stop development container: %s\n", err)
						break
					}
				}

			}
			os.Exit(0)
		}()

		cmd.Println("Your app is ready on http://localhost:8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Printf("Failed to bind to port 8080: %s\n", err)
			os.Exit(1)
		}
	},
}

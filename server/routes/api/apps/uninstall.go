package appsApi

import (
	"context"
	"encoding/json"
	"fmt"
	dbClient "localapps/db/client"
	"localapps/types"
	"net/http"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/volume"
	dockerClient "github.com/docker/docker/client"
)

func uninstallApp(w http.ResponseWriter, r *http.Request) {
	client, _ := dbClient.GetClient()

	var body types.ApiAppUninstallRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %s", err), http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(body.Id) == "" {
		http.Error(w, "App ID is required", http.StatusBadRequest)
		return
	}

	appId := body.Id

	if strings.Contains(appId, " ") {
		http.Error(w, "App ID cannot contain spaces", http.StatusBadRequest)
		return
	}

	_, err := client.GetApp(context.Background(), appId)
	if err != nil {
		http.Error(w, "App not installed", http.StatusInternalServerError)
		return
	}

	cli, err := dockerClient.NewClientWithOpts(dockerClient.FromEnv)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to connect to docker engine: %s", err), http.StatusInternalServerError)
		return
	}

	appContainers, _ := cli.ContainerList(context.Background(), container.ListOptions{Filters: filters.NewArgs(filters.Arg("name", "localapps-app-"+appId))})
	if len(appContainers) > 0 {
		for _, c := range appContainers {
			err := cli.ContainerRemove(context.Background(), c.ID, container.RemoveOptions{Force: true})
			if err != nil {
				http.Error(w, fmt.Sprintf("Error removing app conianer: %s", err), http.StatusInternalServerError)
				return
			}
		}
	}

	storageVolumes, _ := cli.VolumeList(context.Background(), volume.ListOptions{Filters: filters.NewArgs(filters.Arg("name", "localapps-storage-"+appId))})
	if len(storageVolumes.Volumes) > 0 {
		for _, volume := range storageVolumes.Volumes {
			err = cli.VolumeRemove(context.Background(), volume.Name, false)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error removing app storage: %s", err), http.StatusInternalServerError)
				return
			}
		}
	}

	appImages, _ := cli.ImageList(context.Background(), image.ListOptions{Filters: filters.NewArgs(filters.Arg("reference", "localapps/apps/"+appId+"/*"))})
	if len(appImages) > 0 {
		for _, im := range appImages {
			_, err = cli.ImageRemove(context.Background(), im.ID, image.RemoveOptions{})
			if err != nil {
				http.Error(w, fmt.Sprintf("Error removing app image: %s", err), http.StatusInternalServerError)
				return
			}
		}
	}

	err = client.DeleteApp(context.Background(), appId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting app record from db: %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

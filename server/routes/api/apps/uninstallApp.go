package appsApi

import (
	"context"
	"encoding/json"
	"fmt"
	"localapps/constants"
	dbClient "localapps/db/client"
	"localapps/types"
	"net/http"

	"strconv"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/volume"
	dockerClient "github.com/docker/docker/client"
)

func uninstallApp(w http.ResponseWriter, r *http.Request) {
	client, _ := dbClient.GetClient()
	appId := r.PathValue("appId")

	appIdInt64, err := strconv.ParseInt(appId, 10, 64)
	if err != nil {
		response := types.ApiError{
			Code:    constants.ErrorParse,
			Message: fmt.Sprintf("Could not parse app ID \"%s\"", appId),
			Error:   err,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	_, err = client.GetAppById(context.Background(), appIdInt64)
	if err != nil {
		response := types.ApiError{
			Code:    constants.ErrorNotFound,
			Message: fmt.Sprintf("App \"%s\" not installed", appId),
			Error:   err,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	cli, _ := dockerClient.NewClientWithOpts(dockerClient.FromEnv)

	_, err = cli.Ping(context.Background())
	if err != nil {
		response := types.ApiError{
			Code:    constants.ErrorDockerEngine,
			Message: "Failed to connect to Docker engine",
			Error:   err,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	appContainers, _ := cli.ContainerList(context.Background(), container.ListOptions{Filters: filters.NewArgs(filters.Arg("label", "LOCALAPPS_APP_ID="+appId))})
	if len(appContainers) > 0 {
		for _, c := range appContainers {
			err := cli.ContainerRemove(context.Background(), c.ID, container.RemoveOptions{Force: true})
			if err != nil {
				response := types.ApiError{
					Code:    constants.ErrorUninstall,
					Message: "Error while removing app container",
					Error:   err,
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response)
				return
			}
		}
	}

	storageVolumes, _ := cli.VolumeList(context.Background(), volume.ListOptions{Filters: filters.NewArgs(filters.Arg("name", "localapps-storage-"+appId))})
	if len(storageVolumes.Volumes) > 0 {
		for _, volume := range storageVolumes.Volumes {
			err = cli.VolumeRemove(context.Background(), volume.Name, false)
			if err != nil {
				response := types.ApiError{
					Code:    constants.ErrorUninstall,
					Message: "Error while removing app storage",
					Error:   err,
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response)
				return
			}
		}
	}

	appImages, _ := cli.ImageList(context.Background(), image.ListOptions{Filters: filters.NewArgs(filters.Arg("reference", "localapps/apps/"+appId+"/*"))})
	if len(appImages) > 0 {
		for _, im := range appImages {
			_, err = cli.ImageRemove(context.Background(), im.ID, image.RemoveOptions{Force: true})
			if err != nil {
				response := types.ApiError{
					Code:    constants.ErrorUninstall,
					Message: "Error while removing app image",
					Error:   err,
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response)
				return
			}
		}
	}

	err = client.DeleteApp(context.Background(), appIdInt64)
	if err != nil {
		response := types.ApiError{
			Code:    constants.ErrorDb,
			Message: "Error while deleting DB record",
			Error:   err,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

package iconsApi

import (
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"localapps/constants"
	"localapps/types"
	"net/http"
	"os"
	"path/filepath"
)

func getAppIcon(w http.ResponseWriter, r *http.Request) {
	if _, err := os.Stat(constants.LocalappsAppIconsDir); errors.Is(err, fs.ErrNotExist) {
		if err := os.MkdirAll(constants.LocalappsAppIconsDir, 0755); err != nil {
			response := types.ApiError{
				Code:    constants.ErrorFsCreateDir,
				Message: "Failed to create directory",
				Error:   err,
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	_, err := os.Stat(filepath.Join(constants.LocalappsAppIconsDir, r.PathValue("icon")))
	if errors.Is(err, fs.ErrNotExist) {
		response := types.ApiError{
			Code:    constants.ErrorNotFound,
			Message: "Icon not found",
			Error:   err,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	file, err := os.Open(filepath.Join(constants.LocalappsAppIconsDir, r.PathValue("icon")))
	if err != nil {
		response := types.ApiError{
			Code:    constants.ErrorAccessDenied,
			Message: "Failed to open icon file",
			Error:   err,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "image/png")
	io.Copy(w, file)
}

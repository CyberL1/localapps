package utils

import (
	"fmt"
	"io"
	"localapps/constants"
	"localapps/types"
	"os"
	"path/filepath"

	"github.com/go-yaml/yaml"
)

func GetApp(appName string) (*types.App, error) {
	appFilePath := filepath.Join(constants.LocalappsDir, appName, "app.yml")

	file, err := os.Open(appFilePath)
	if err != nil {
		return nil, fmt.Errorf("app \"%s\" not found", appName)
	}

	defer file.Close()
	appFileContents, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read app file: %v", err)
	}

	var app types.App
	err = yaml.Unmarshal(appFileContents, &app)
	if err != nil {
		return nil, fmt.Errorf("failed to parse app file: %v", err)
	}

	return &app, nil
}

func GetAppDirectory(appName string) string {
	return filepath.Join(constants.LocalappsDir, appName)
}

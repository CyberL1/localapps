package main

import (
	"fmt"
	"localapps/cmd"
	"localapps/constants"
	"os"
	"path/filepath"
)

func main() {
	// Check for all required resources
	if _, err := os.Stat(constants.LocalappsDir); os.IsNotExist(err) {
		if err := os.Mkdir(constants.LocalappsDir, os.ModePerm); err != nil {
			fmt.Println("Failed to create ~/.config/localapps directory:", err)
			return
		}
	}

	if _, err := os.Stat(filepath.Join(constants.LocalappsDir, "apps")); os.IsNotExist(err) {
		if err := os.Mkdir(filepath.Join(constants.LocalappsDir, "apps"), os.ModePerm); err != nil {
			fmt.Println("Failed to create ~/.config/localapps/apps directory:", err)
			return
		}
	}

	cmd.Execute()
}

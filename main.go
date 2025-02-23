package main

import (
	"embed"
	"fmt"
	"localapps/cmd"
	"localapps/constants"
	"os"
	"path/filepath"
)

//go:embed resources
var resources embed.FS

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

	if _, err := os.Stat(filepath.Join(constants.LocalappsDir, "pages")); os.IsNotExist(err) {
		os.Mkdir(filepath.Join(constants.LocalappsDir, "pages"), 0775)

		pagesDir, err := resources.ReadDir("resources/pages")
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("pagesDir", pagesDir)
		for _, file := range pagesDir {
			fmt.Println("file", file.Name())
			contents, _ := resources.ReadFile(fmt.Sprintf("resources/pages/%v", file.Name()))
			os.WriteFile(filepath.Join(constants.LocalappsDir, "pages", file.Name()), contents, 0664)
		}
	}

	cmd.Execute()
}

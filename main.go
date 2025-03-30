package main

import (
	"fmt"
	"localapps/cmd"
	"localapps/constants"
	"os"
)

func main() {
	// Check for all required resources
	if _, err := os.Stat(constants.LocalappsDir); os.IsNotExist(err) {
		if err := os.Mkdir(constants.LocalappsDir, os.ModePerm); err != nil {
			fmt.Println("Failed to create ~/.config/localapps directory:", err)
			return
		}
	}

	cmd.Execute()
}

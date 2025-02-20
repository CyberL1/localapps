package cmd

import (
	"localapps/constants"
	"localapps/utils"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed localapps",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		files, err := os.ReadDir(constants.LocalappsDir)
		if err != nil {
			cmd.PrintErrln("Error reading directory:", err)
			return
		}

		for _, file := range files {
			app, _ := utils.GetApp(file.Name())
			cmd.Println(app.Name, "->", filepath.Join(constants.LocalappsDir, file.Name()))
		}
	},
}

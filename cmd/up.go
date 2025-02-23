package cmd

import (
	"fmt"
	"localapps/server/middlewares"
	"localapps/server/routes"
	"localapps/server/routes/api"
	"net/http"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(upCmd)
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Start localapps server",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("Starting server")

		router := http.NewServeMux()

		homeHandler := routes.NewHandler().RegisterRoutes()
		router.Handle("/", homeHandler)

		apiHandler := api.NewHandler().RegisterRoutes()
		router.Handle("/api/", http.StripPrefix("/api", apiHandler))

		err := http.ListenAndServe(":8080", middlewares.AppProxy(router))
		if err != nil {
			fmt.Println(err)
		}
	},
}

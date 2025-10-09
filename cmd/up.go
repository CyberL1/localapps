package cmd

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"localapps/constants"
	dbClient "localapps/db/client"
	db "localapps/db/generated"
	"localapps/server/middlewares"
	"localapps/server/routes"
	"localapps/server/routes/api"
	"localapps/utils"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func init() {
	rootCmd.AddCommand(upCmd)
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Start localapps server",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// Check for all required resources
		if _, err := os.Stat(constants.LocalappsDir); os.IsNotExist(err) {
			if err := os.Mkdir(constants.LocalappsDir, 0755); err != nil {
				fmt.Println("Failed to create ~/.config/localapps directory:", err)
				return
			}
		}

		if !constants.IsRunningInContainer() {
			fmt.Println("----- 🚨 Running on host 🚨 -----\nApp ports will be exposed to a random port on host.\nIt is recommended to run on docker in production.\n----- 🚨 Running on host 🚨 -----")
		}

		cli, _ := client.NewClientWithOpts(client.FromEnv)

		_, err := cli.Ping(context.Background())
		if err != nil {
			fmt.Printf("Failed to connect to Docker engine: %s\n", err)
			return
		}

		if staleContainers, _ := cli.ContainerList(context.Background(), container.ListOptions{Filters: filters.NewArgs(filters.Arg("label", "LOCALAPPS_APP_ID")), All: true}); len(staleContainers) > 0 {
			fmt.Printf("Found %d stale containers, removing\n", len(staleContainers))

			for _, c := range staleContainers {
				cli.ContainerRemove(context.Background(), c.ID, container.RemoveOptions{Force: true})
			}
		}

		if networks, _ := cli.NetworkList(context.Background(), network.ListOptions{Filters: filters.NewArgs(filters.Arg("name", "localapps-network"))}); len(networks) == 0 {
			cmd.Println("Creating localapps network")
			cli.NetworkCreate(context.Background(), "localapps-network", network.CreateOptions{})
		}

		fmt.Println("Running database migrations")
		dbClient.Migrate()

		fmt.Println("Fetching server configuration")
		err = utils.UpdateServerConfigCache()
		if err != nil {
			fmt.Printf("Error updating config cache: %s\n", err)
			return
		}

		accessUrlFilePath := filepath.Join(constants.LocalappsDir, "access-url.txt")
		if _, err := os.Stat(accessUrlFilePath); err == nil {
			fmt.Println("Found access-url.txt file, updating server configuration")
			client, _ := dbClient.GetClient()

			file, err := os.Open(accessUrlFilePath)
			if err != nil {
				fmt.Printf("Error opening file: %s\n", err)
				return
			}
			defer file.Close()

			accessUrlFileContents, err := io.ReadAll(file)
			if err != nil {
				fmt.Printf("Error reading file: %s\n", err)
				return
			}

			accessUrlRaw := strings.Split(string(accessUrlFileContents), "\n")[0]
			accessUrlParsed, err := json.Marshal(string(accessUrlRaw))
			if err != nil {
				fmt.Printf("Error parsing file: %s\n", err)
				return
			}

			_, err = client.UpdateConfigKey(context.Background(), db.UpdateConfigKeyParams{
				Key:   "AccessUrl",
				Value: sql.NullString{String: string(accessUrlParsed), Valid: true},
			})
			if err != nil {
				fmt.Printf("Error updating access URL: %s\n", err)
			}

			err = utils.UpdateServerConfigCache()
			if err != nil {
				fmt.Printf("Error updating config cache: %s\n", err)
				return
			}

			fmt.Println("Success, removing the file")
			err = os.Remove(accessUrlFilePath)
			if err != nil {
				fmt.Printf("Error removing file: %s\n", err)
				return
			}
		}

		if utils.ServerConfig.ApiKey == "" {
			fmt.Println("Server API Key is empty, using a random value")
			client, _ := dbClient.GetClient()

			apiKeyParsed, err := json.Marshal(strings.ReplaceAll(uuid.NewString(), "-", ""))
			if err != nil {
				fmt.Printf("Error parsing api key: %s\n", err)
				return
			}

			_, err = client.UpdateConfigKey(context.Background(), db.UpdateConfigKeyParams{
				Key:   "ApiKey",
				Value: sql.NullString{String: string(apiKeyParsed), Valid: true},
			})
			if err != nil {
				fmt.Printf("Error updating domain: %s\n", err)
			}

			err = utils.UpdateServerConfigCache()
			if err != nil {
				fmt.Printf("Error updating config cache: %s\n", err)
				return
			}
		}

		cmd.Println("Starting HTTP server")

		router := http.NewServeMux()

		router.Handle("/", routes.NewHandler().RegisterRoutes())
		router.Handle("/api/", http.StripPrefix("/api", api.NewHandler().RegisterRoutes()))

		if err := http.ListenAndServe(":8080", middlewares.AppProxy(router)); err != nil {
			fmt.Printf("Failed to bind to port 8080: %s\n", err)
			os.Exit(1)
		}
	},
}

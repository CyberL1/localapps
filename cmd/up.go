package cmd

import (
	"context"
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
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/spf13/cobra"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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
		cmd.Println("Creating database server container")
		freePort, _ := utils.GetFreePort()

		cli, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			fmt.Printf("Failed to connect to Docker daemon: %s\n", err)
			return
		}

		var databasePassword string
		pgPassFilePath := filepath.Join(constants.LocalappsDir, ".pgpasswd")
		if _, err := os.Stat(pgPassFilePath); err == nil {
			file, err := os.Open(pgPassFilePath)
			if err != nil {
				fmt.Printf("Error opening .pgpasswd file: %s\n", err)
				return
			}
			defer file.Close()

			content, err := io.ReadAll(file)
			if err != nil {
				fmt.Printf("Error reading .pgpasswd file: %s\n", err)
				return
			}
			databasePassword = string(content)
		} else {
			databasePassword = strings.ReplaceAll(uuid.NewString(), "-", "")

			file, err := os.Create(pgPassFilePath)
			if err != nil {
				fmt.Printf("Error creating .pgpasswd file: %s\n", err)
				return
			}
			defer file.Close()

			_, err = file.WriteString(databasePassword)
			if err != nil {
				fmt.Printf("Error writing to .pgpasswd file: %s\n", err)
				return
			}
		}

		if c, _ := cli.ContainerList(context.Background(), container.ListOptions{Filters: filters.NewArgs(filters.Arg("name", "localapps-database")), All: true}); len(c) > 0 {
			cli.ContainerRemove(context.Background(), c[0].ID, container.RemoveOptions{Force: true})
		}

		config := container.Config{
			Image: "postgres:17-alpine",
			Env: []string{
				"POSTGRES_USER=localapps",
				fmt.Sprintf("POSTGRES_PASSWORD=%s", databasePassword),
			},
			ExposedPorts: nat.PortSet{"5432": struct{}{}},
		}

		hostConfig := container.HostConfig{
			PortBindings: nat.PortMap{
				"5432": {
					{
						HostIP:   "0.0.0.0",
						HostPort: strconv.Itoa(freePort),
					},
				},
			},
			Binds: []string{fmt.Sprintf("%s:/var/lib/postgresql/data", filepath.Join(constants.LocalappsDir, "data"))},
		}

		databaseContainer, err := cli.ContainerCreate(context.Background(), &config, &hostConfig, nil, nil, "localapps-database")
		if err != nil {
			fmt.Printf("Failed to create database container: %s\n", err)
			return
		}

		if err := cli.ContainerStart(context.Background(), databaseContainer.ID, container.StartOptions{}); err != nil {
			cmd.Printf("Failed to start database container: %s\n", err)
			return
		}

		os.Setenv("LOCALAPPS_DB", fmt.Sprintf("postgres://localapps:%s@localhost:%d/localapps?sslmode=disable", databasePassword, freePort))

		fmt.Println("Waiting for database client to connect")
		for {
			_, err := dbClient.GetClient()
			if err == nil {
				break
			}
			fmt.Println("Database server not ready, retrying in 1 second")
			time.Sleep(1 * time.Second)
		}

		fmt.Println("Running database migrations")
		dbClient.Migrate()

		fmt.Println("Fetching server configuration")
		err = utils.UpdateConfigCache()
		if err != nil {
			fmt.Printf("Error updating config cache: %s\n", err)
			return
		}

		domainFilePath := filepath.Join(constants.LocalappsDir, "domain.txt")
		if _, err := os.Stat(domainFilePath); err == nil {
			fmt.Println("Found domain.txt file, updating server configuration")
			client, _ := dbClient.GetClient()

			file, err := os.Open(domainFilePath)
			if err != nil {
				fmt.Printf("Error opening file: %s\n", err)
				return
			}
			defer file.Close()

			domainFileContents, err := io.ReadAll(file)
			if err != nil {
				fmt.Printf("Error reading file: %s\n", err)
				return
			}

			domainRaw := strings.Split(string(domainFileContents), "\n")[0]
			domainParsed, err := json.Marshal(string(domainRaw))
			if err != nil {
				fmt.Printf("Error parsing file: %s\n", err)
				return
			}

			_, err = client.UpdateConfigKey(dbClient.Ctx, db.UpdateConfigKeyParams{
				Key:   "Domain",
				Value: pgtype.Text{String: string(domainParsed), Valid: true},
			})
			if err != nil {
				fmt.Printf("Error updating domain: %s\n", err)
			}

			err = utils.UpdateConfigCache()
			if err != nil {
				fmt.Printf("Error updating config cache: %s\n", err)
				return
			}

			fmt.Println("Success, removing the file")
			err = os.Remove(domainFilePath)
			if err != nil {
				fmt.Printf("Error removing file: %s\n", err)
				return
			}
		}

		cmd.Println("Starting HTTP server")

		router := http.NewServeMux()

		homeHandler := routes.NewHandler().RegisterRoutes()
		router.Handle("/", homeHandler)

		apiHandler := api.NewHandler().RegisterRoutes()
		router.Handle("/api/", http.StripPrefix("/api", apiHandler))

		// Exit handler
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-stop

			cmd.Println("Removing database container")
			if err := cli.ContainerRemove(context.Background(), databaseContainer.ID, container.RemoveOptions{Force: true}); err != nil {
				cmd.Printf("Failed to remove database container: %s\n", err)
			}
			os.Exit(0)
		}()

		if err := http.ListenAndServe(":8080", middlewares.AppProxy(router)); err != nil {
			fmt.Printf("Failed to bind to port 8080: %s\n", err)
			os.Exit(1)
		}
	},
}

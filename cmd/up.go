package cmd

import (
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
	"strings"
	"syscall"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/spf13/cobra"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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
		var isExternalDb bool
		var database *embeddedpostgres.EmbeddedPostgres

		if _, found := os.LookupEnv("LOCALAPPS_DB"); found {
			isExternalDb = true
		} else {
			cmd.Println("Starting built-in database server")

			freePort, _ := utils.GetFreePort()
			database = embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
				Username("localapps").
				Password("localapps").
				Database("localapps").
				RuntimePath(filepath.Join(constants.LocalappsDir, "postgres")).
				DataPath(filepath.Join(constants.LocalappsDir, "data")).
				Port(uint32(freePort)))

			if err := database.Start(); err != nil {
				cmd.Println(err)
				return
			}

			os.Setenv("LOCALAPPS_DB", fmt.Sprintf("postgres://localapps:localapps@localhost:%d/localapps?sslmode=disable", freePort))
		}

		fmt.Println("Running database migrations")
		dbClient.Migrate()

		fmt.Println("Fetching server configuration")
		err := utils.UpdateConfigCache()
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

			if !isExternalDb {
				cmd.Println("Stopping built-in database server")
				if err := database.Stop(); err != nil {
					cmd.PrintErrln(err)
				}
			}

			os.Exit(0)
		}()

		if err := http.ListenAndServe(":8080", middlewares.AppProxy(router)); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

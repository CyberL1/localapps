package cmd

import (
	"fmt"
	"io"
	"localapps/types"
	"localapps/utils"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func init() {
	rootCmd.AddCommand(devCmd)
}

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Start localapps server",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		currentDir, _ := os.Getwd()

		appFilePath := filepath.Join(currentDir, "app.yml")
		file, err := os.Open(appFilePath)

		if err != nil {
			cmd.PrintErrln("No app.yml file detected")
			return
		}
		defer file.Close()

		appFileContents, err := io.ReadAll(file)
		if err != nil {
			cmd.PrintErrf("failed to read app file: %v\n", err)
		}

		var app types.App
		err = yaml.Unmarshal(appFileContents, &app)
		if err != nil {
			cmd.PrintErrf("failed to parse app file: %v\n", err)
		}

		for partName, part := range app.Parts {
			if strings.TrimSpace(part.Dev) == "" {
				println("Part", partName, "has no dev command")
				os.Exit(1)
			}
		}

		cmd.Println("Running app in dev mode")

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			var currentPart types.Part
			var fallbackPart types.Part

			for _, part := range app.Parts {
				if part.Path == "" {
					fallbackPart = part
				}

				if strings.Split(r.URL.Path, "/")[1] == part.Path {
					currentPart = part
					break
				} else {
					currentPart = fallbackPart
				}
			}

			findProcess := exec.Command("pgrep", "-f", currentPart.Dev)

			out, _ := findProcess.Output()
			splitted := strings.Split(string(out), "\n")

			var freePort int
			var isRunning bool

			for _, pid := range splitted {
				getDirectory := exec.Command("readlink", fmt.Sprintf("/proc/%s/cwd", string(pid)))
				out, _ = getDirectory.Output()

				if strings.TrimSpace(string(out)) == filepath.Join(currentDir, currentPart.Src) {
					isRunning = true
					findPort := exec.Command("cat", fmt.Sprintf("/proc/%s/environ", string(pid)))
					out, _ = findPort.Output()

					splitted := strings.Split(string(out), "\x00")

					for _, env := range splitted {
						if strings.HasPrefix(env, "PORT=") {
							freePort, _ = strconv.Atoi(strings.Split(env, "=")[1])
							break
						}
					}
				}
			}

			if !isRunning {
				freePort, _ = utils.GetFreePort()

				cmd := exec.Command("/bin/sh", "-c", currentPart.Dev)
				cmd.Dir = currentPart.Src

				cmd.Env = append(cmd.Env, "PATH="+os.Getenv("PATH"))
				cmd.Env = append(cmd.Env, fmt.Sprintf("PORT=%s", strconv.Itoa(freePort)))

				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				if err := cmd.Start(); err != nil {
					w.Write([]byte(fmt.Sprintf("Failed to start app \"%s\": %s", app.Name, err)))
					return
				}

				// Wait for the app to be ready
				for {
					_, err := http.Get(fmt.Sprintf("http://localhost:%s", strconv.Itoa(freePort)))
					if err == nil {
						break
					}
					time.Sleep(500 * time.Millisecond)
				}
			}

			appUrl, _ := url.Parse(fmt.Sprintf("http://localhost:%s", strconv.Itoa(freePort)))
			httputil.NewSingleHostReverseProxy(appUrl).ServeHTTP(w, r)
		})

		cmd.Println("Your app is ready on", fmt.Sprintf("https://%s.apps.localhost", strings.Split(currentDir, "/")[4]))
		http.ListenAndServe(":8080", nil)
	},
}

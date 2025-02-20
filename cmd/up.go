package cmd

import (
	"fmt"
	"localapps/types"
	"localapps/utils"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.Host, "localhost:8080") {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			// Check if Docker is installed
			if _, err := exec.LookPath("docker"); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Docker is not installed"))
				return
			}

			// Check if Docker daemon is running
			if err := exec.Command("docker", "info").Run(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Docker daemon is not running"))
				return
			}

			if len(strings.Split(r.Host, ".")) == 2 {
				http.Redirect(w, r, "//www."+r.Host, http.StatusTemporaryRedirect)
				return
			}

			if len(strings.Split(r.Host, ".")) == 3 {
				appName := strings.Split(r.Host, ".")[0]
				app, err := utils.GetApp(appName)

				if err != nil {
					w.Write([]byte(fmt.Sprintf("App \"%s\" not found", appName)))
					return
				}

				var dockerAppName string
				var dockerImageName string

				var currentPart types.Part

				var fallbackPart types.Part
				var fallbackPartName string

				for partName, part := range app.Parts {
					if part.Path == "" {
						fallbackPart = part
						fallbackPartName = partName
					}

					if strings.Split(r.URL.Path, "/")[1] == part.Path {
						currentPart = part
						dockerAppName = "localapps-app-" + strings.ToLower(app.Name) + "-" + partName
						dockerImageName = "localapps/" + strings.ToLower(app.Name) + "/" + partName
						break
					} else {
						currentPart = fallbackPart
						dockerAppName = "localapps-app-" + strings.ToLower(app.Name) + "-" + fallbackPartName
						dockerImageName = "localapps/" + strings.ToLower(app.Name) + "/" + fallbackPartName
					}
				}

				findProcess := exec.Command("docker", "ps", "--format", "{{.Names}}", "-f", "name="+dockerAppName+"$")
				out, _ := findProcess.Output()

				var freePort int
				if strings.TrimSpace(string(out)) == dockerAppName {
					portCmd := exec.Command("docker", "port", dockerAppName, "80")
					out, _ := portCmd.Output()
					freePort, _ = strconv.Atoi(strings.TrimSpace(string(out)[8:]))
				} else {
					freePort, _ = GetFreePort()

					runCmd := exec.Command("docker", "run", "--rm", "--name", dockerAppName, "-p", strconv.Itoa(freePort)+":80", "-e", "PORT=80", dockerImageName)
					runCmd.Dir = filepath.Join(utils.GetAppDirectory(appName), currentPart.Src)

					cmd.Println("[app:"+app.Name+"]", "Got a http request while stopped - starting")

					if err := runCmd.Start(); err != nil {
						w.Write([]byte(fmt.Sprintf("Failed to start app \"%s\": %s", appName, err)))
						return
					}

					go func() {
						time.Sleep(30 * time.Second)
						stopCmd := exec.Command("docker", "stop", dockerAppName)
						cmd.Println("[app:"+app.Name+"]", "Exceeded timeout (30s) - stopping")
						stopCmd.Start()
					}()
				}

				// Wait for the app to be ready
				for {
					_, err := http.Get(fmt.Sprintf("http://localhost:%s", strconv.Itoa(freePort)))
					if err == nil {
						break
					}
					time.Sleep(500 * time.Millisecond)
				}

				appUrl, _ := url.Parse(fmt.Sprintf("http://localhost:%s", strconv.Itoa(freePort)))
				httputil.NewSingleHostReverseProxy(appUrl).ServeHTTP(w, r)
			}
		})

		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func GetFreePort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return
}

package cmd

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/CyberL1/localapps/types"
	"github.com/CyberL1/localapps/utils"
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

				var currentPart types.Part
				for _, part := range app.Parts {
					if strings.Split(r.URL.Path, "/")[1] == part.Path {
						currentPart = part
						break
					}
				}

				freePort, _ := GetFreePort()

				cmd := exec.Command(strings.Split(currentPart.Run, " ")[0], strings.Split(currentPart.Run, " ")[1:]...)
				cmd.Dir = filepath.Join(utils.GetAppDirectory(appName), currentPart.Src)
				cmd.Env = append(cmd.Env, fmt.Sprintf("PORT=%s", strconv.Itoa(freePort)))

				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				if err := cmd.Start(); err != nil {
					w.Write([]byte(fmt.Sprintf("Failed to start app \"%s\": %s", appName, err)))
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

				appUrl, _ := url.Parse(fmt.Sprintf("http://localhost:%s", strconv.Itoa(freePort)))
				httputil.NewSingleHostReverseProxy(appUrl).ServeHTTP(w, r)

				go func() {
					time.Sleep(30 * time.Second)
					cmd.Process.Kill()
				}()

				return
			}
		})

		http.ListenAndServe(":8080", nil)
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

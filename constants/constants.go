package constants

import (
	"os"
	"path/filepath"
)

var configDir, _ = os.UserConfigDir()
var runningInContainer string

var (
	LocalappsDir         = filepath.Join(configDir, "localapps")
	LocalappsFilesDir    = filepath.Join(LocalappsDir, "files")
	LocalappsIconsDir    = filepath.Join(LocalappsFilesDir, "icons")
	LocalappsAppIconsDir = filepath.Join(LocalappsIconsDir, "apps")

	Version          string
	GithubReleaseUrl = "https://api.github.com/repos/CyberL1/localapps/releases/latest"
)

func IsRunningInContainer() bool {
	return runningInContainer == "true"
}

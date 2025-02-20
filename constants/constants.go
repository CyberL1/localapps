package constants

import (
	"os"
	"path/filepath"
)

var homeDir = os.Getenv("HOME")

var (
	LocalappsDir     = filepath.Join(homeDir, ".config", "localapps")
	Version          string
	GithubReleaseUrl = "https://api.github.com/repos/CyberL1/localapps/releases/latest"
)

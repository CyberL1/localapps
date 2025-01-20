package constants

import (
	"os"
	"path/filepath"
)

var homeDir = os.Getenv("HOME")

var (
	LocalappsDir = filepath.Join(homeDir, "localapps")
)

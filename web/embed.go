package web

import (
	"embed"
	"io/fs"
)

//go:embed all:build
var BuildDir embed.FS
var BuildDirFS, _ = fs.Sub(BuildDir, "build")

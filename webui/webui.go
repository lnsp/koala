package webui

import (
	"embed"
	"io/fs"
)

//go:embed dist
var DistFS embed.FS

var FS, _ = fs.Sub(DistFS, "dist")

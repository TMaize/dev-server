package ui

import (
	"embed"
)

//go:embed dist
var FS embed.FS

var Prefix = "dist"
package static

import (
	"embed"
	"io/fs"
)

//go:embed root/*
var root embed.FS

func Root() (fs.FS, error) {
	return fs.Sub(root, "root")
}

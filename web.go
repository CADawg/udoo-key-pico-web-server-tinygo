package main

import (
	"embed"
	"io/fs"
	"path/filepath"
)

//go:embed web/include/*
var webFilesNonRoot embed.FS

var webFiles = MustSubFS(webFilesNonRoot, "web/include/")

func subFS(currentFs fs.FS, root string) (fs.FS, error) {
	root = filepath.ToSlash(filepath.Clean(root)) // note: fs.FS operates only with slashes. `ToSlash` is necessary for Windows
	return fs.Sub(currentFs, root)
}

func MustSubFS(currentFs fs.FS, root string) fs.FS {
	fsOut, err := subFS(currentFs, root)
	if err != nil {
		panic(err)
	}
	return fsOut
}

func HttpGetFile(filename string) ([]byte, error) {
	if filename == "" {
		filename = "index.html"
	}

	// flash led to show file access
	TurnOnLed()

	return fs.ReadFile(webFiles, filename)
}

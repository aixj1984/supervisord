//go:build release
// +build release

package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
)

//go:embed webgui/*
var content embed.FS

var HTTP http.FileSystem

func init() {
	webgui, err := fs.Sub(content, "webgui")
	if err != nil {
		panic(err)
	}

	HTTP = http.FS(webgui)
}

func ReadStaticFile(path string) ([]byte, error) {
	data, err := content.ReadFile(path)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	return data, nil
}

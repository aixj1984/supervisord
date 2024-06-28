//go:build !release
// +build !release

package main

import (
	"io"
	"net/http"
	"os"
)

// HTTP auto generated
var HTTP http.FileSystem = http.Dir("./webgui")

func ReadStaticFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return b, nil
}

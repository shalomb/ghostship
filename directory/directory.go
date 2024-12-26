package directory

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	name = "directory"
)

type directoryRenderer struct{}

// Renderer ...
func Renderer() *directoryRenderer {
	return &directoryRenderer{}
}

func (i *directoryRenderer) Name() string {
	return name
}

// Render ...
func (r *directoryRenderer) Render() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// return pwd, filepath.Base(pwd)
	return filepath.Base(pwd), err
}

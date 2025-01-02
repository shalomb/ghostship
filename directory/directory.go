package directory

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/shalomb/ghostship/config"
)

const (
	name = "directory"
)

// DirectoryRenderer ...
type DirectoryRenderer struct{}

// Renderer ...
func Renderer() *DirectoryRenderer {
	return &DirectoryRenderer{}
}

// Name ...
func (i *DirectoryRenderer) Name() string {
	return name
}

// Render ...
func (r *DirectoryRenderer) Render(c config.TomlConfig, e interface{}) (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// return pwd, filepath.Base(pwd)
	return filepath.Base(pwd), err
}

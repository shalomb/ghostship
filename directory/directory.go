package directory

import (
	"fmt"
	"os"
	"path/filepath"

	colors "github.com/shalomb/ghostship/colors"
	"github.com/shalomb/ghostship/config"
	log "github.com/sirupsen/logrus"
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
func (r *DirectoryRenderer) Name() string {
	return name
}

// Render ...
func (r *DirectoryRenderer) Render(c config.AppConfig, e config.EnvironmentConfig) (string, error) {
	cfg := c.DirectoryConfig

	// TODO: Parse and use format
	log.Debugf("DirectoryConfig:\nc:%+v\ns:%+v\nv:%+v", cfg, e, cfg.Format)

	pwd, err := os.Getwd()
	if err != nil {
		return pwd, err
	}
	// return pwd, filepath.Base(pwd)

	style := colors.ByExpression(cfg.Style)
	symbol := filepath.Base(pwd)

	ret := fmt.Sprintf(
		"\\[%s\\]%s\\[%s\\]",
		style,
		symbol,
		colors.Reset,
	)

	return ret, nil
}

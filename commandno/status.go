package commandno

import (
	"fmt"

	colors "github.com/shalomb/ghostship/colors"
	config "github.com/shalomb/ghostship/config"
	log "github.com/sirupsen/logrus"
)

const (
	name = ""
)

type CommandNumberRenderer struct{}

// Renderer ...
func Renderer() *CommandNumberRenderer {
	return &CommandNumberRenderer{}
}

// Name ...
func (r *CommandNumberRenderer) Name() string {
	return name
}

// Render ...
func (r *CommandNumberRenderer) Render(c config.AppConfig, e config.EnvironmentConfig) (string, error) {
	cfg := c.CommandNumberConfig
	log.Debugf("CommandNumberConfig: %+v, s: %+v, v: %+v", cfg, e, cfg.Format)

	symbol := cfg.Format
	style := colors.ByExpression(cfg.Style)
	log.Debugf("symbol:%+v", symbol)

	ret := fmt.Sprintf(
		"\\[%s\\]%s\\[%s\\]",
		style,
		symbol,
		colors.Reset,
	)

	return ret, nil
}

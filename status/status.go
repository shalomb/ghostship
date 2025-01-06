package status

import (
	"fmt"

	colors "github.com/shalomb/ghostship/colors"
	config "github.com/shalomb/ghostship/config"
	log "github.com/sirupsen/logrus"
)

const (
	name = ""
)

type StatusRenderer struct{}

// Renderer ...
func Renderer() *StatusRenderer {
	return &StatusRenderer{}
}

// Name ...
func (r *StatusRenderer) Name() string {
	return name
}

// Render ...
func (r *StatusRenderer) Render(c config.AppConfig, e config.EnvironmentConfig) (string, error) {
	cfg := c.StatusConfig
	log.Debugf("StatusConfig: %+v, s: %+v, v: %+v", cfg, e, cfg.Format)

	var symbol string
	if e["status"] == uint16(0) && e["pipestatus"] == uint16(0) {
		symbol = ""
	} else {
        if e["status"] != e["pipestatus"] {
            symbol = fmt.Sprintf("%d|%d", e["status"], e["pipestatus"])
        } else {
            symbol = fmt.Sprintf("%d", e["status"])
        }
	}
	style := colors.ByExpression(cfg.Style)

	ret := fmt.Sprintf(
		"\\[%s\\]%s\\[%s\\]",
		style,
		symbol,
		colors.Reset,
	)

	return ret, nil
}

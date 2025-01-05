package duration

import (
	"fmt"

	colors "github.com/shalomb/ghostship/colors"
	config "github.com/shalomb/ghostship/config"
	log "github.com/sirupsen/logrus"
)

const (
	name = ""
)

type DurationRenderer struct{}

// Renderer ...
func Renderer() *DurationRenderer {
	return &DurationRenderer{}
}

// Name ...
func (r *DurationRenderer) Name() string {
	return name
}

// Render ...
func (r *DurationRenderer) Render(c config.AppConfig, e config.EnvironmentConfig) (string, error) {
	cfg := c.DurationConfig
	log.Debugf("DurationConfig: %+v, s: %+v, v: %+v", cfg, e, cfg.Format)

    if e["cmd-duration"] != nil && e["cmd-duration"].(uint16) < cfg.MinTime {
        return "", nil
    }

    symbol := fmt.Sprintf(cfg.Format, e["cmd-duration"])
	style := colors.ByExpression(cfg.Style)
    log.Debugf("Sty:%+v Sty:%+v", symbol, style)

	ret := fmt.Sprintf(
		"\\[%s\\]%s\\[%s\\]",
		style,
		symbol,
		colors.Reset,
	)

	return ret, nil
}

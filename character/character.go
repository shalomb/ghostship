package character

import (
	"fmt"

	colors "github.com/shalomb/ghostship/colors"
	config "github.com/shalomb/ghostship/config"
	log "github.com/sirupsen/logrus"
)

const (
	name = ""
)

// CharacterRenderer ...
type CharacterRenderer struct{}

// Renderer ...
func Renderer() *CharacterRenderer {
	return &CharacterRenderer{}
}

// Name ...
func (r *CharacterRenderer) Name() string {
	return name
}

// Render ...
func (r *CharacterRenderer) Render(c config.AppConfig, e config.EnvironmentConfig) (string, error) {
	cfg := c.CharacterConfig
	log.Debugf("CharacterConfig: e: %+v\n, s: %+v\n, v: %+v\n", e, cfg.SuccessColour, cfg.Format)

	symbol := " " + e["prompt-character"].(string)
	stl := cfg.SuccessColour
	if e["status"] != uint16(0) && e["pipestatus"] != uint16(0) {
		stl = cfg.ErrorColour
	}

	style := colors.ByExpression(stl)

	ret := fmt.Sprintf(
		"\\[%s\\]%s\\[%s\\] ",
		style,
		symbol,
		colors.Reset,
	)

	return ret, nil
}

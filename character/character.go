package character

import (
    "fmt"
    log "github.com/sirupsen/logrus"
    config "github.com/shalomb/ghostship/config"
	colors "github.com/shalomb/ghostship/colors"
)

const (
	name = ""
)

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
func (r *CharacterRenderer) Render(c config.TomlConfig, e interface{}) (string, error) {
	cfg := c.CharacterConfig
	log.Debugf("CharacterConfig: %+v, s: %+v, v: %+v", c, e, cfg.Format)

	style := cfg.Style
	ret := fmt.Sprintf(
		"%s%s%s",
		colors.ByName(style),
		cfg.Format,
		colors.Reset,
	)

	return ret, nil
}

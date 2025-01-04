package character

import (
	config "github.com/shalomb/ghostship/config"
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
func (r *CharacterRenderer) Render(c config.AppConfig, e config.EnvironmentConfig) (string, error) {
	return "\n", nil
}

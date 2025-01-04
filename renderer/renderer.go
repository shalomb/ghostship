// Package renderer ...
package renderer

import (
	config "github.com/shalomb/ghostship/config"
	log "github.com/sirupsen/logrus"
)

// Renderer ...
type Renderer interface {
	Name() string
	Render(config.AppConfig, config.EnvironmentConfig) (string, error)
}

// ComponentRenderer ...
type ComponentRenderer struct {
	Renderer Renderer
}

// New ...
func New() *ComponentRenderer {
	ip := &ComponentRenderer{}
	log.Debugf("init Renderer: %+v", ip)
	return ip
}

// SetRenderer ...
func (i *ComponentRenderer) SetRenderer(renderer Renderer) {
	i.Renderer = renderer
}

// Render ...
func (i *ComponentRenderer) Render(c config.AppConfig, e config.EnvironmentConfig) (string, error) {
	return i.Renderer.Render(c, e)
}

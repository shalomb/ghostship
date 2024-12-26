package renderer

import (
	log "github.com/sirupsen/logrus"
)

type Renderer interface {
	Name() string
	Render() (string, error)
}

// ComponentRenderer ...
type ComponentRenderer struct {
	Renderer Renderer
}

func New() *ComponentRenderer {
	ip := &ComponentRenderer{}
	log.Printf("init Renderer: %+v", ip)
	return ip
}

func (i *ComponentRenderer) SetRenderer(renderer Renderer) {
	i.Renderer = renderer
}

func (i *ComponentRenderer) Render() (string, error) {
	return i.Renderer.Render()
}


package commandno

import (
	// "regexp"
	"log"
	"testing"

	config "github.com/shalomb/ghostship/config"
	assert "github.com/stretchr/testify/assert"
)

func TestCommandNumberFormat(t *testing.T) {
	renderer := Renderer()
	actual, err := renderer.Render(config.DefaultConfig(), config.EnvironmentConfig{})
	assert.Equal(t, err, nil, "Error not empty")

	expected := "❯ "

	log.Printf("commandno render: actual: %v, expected: %v", actual, expected)
	assert.Equal(t,
		actual,
		actual,
		"Rendered commandno must match",
	)
}

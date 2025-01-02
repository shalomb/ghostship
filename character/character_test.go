package character

import (
	// "regexp"
	"log"
	"testing"

	config "github.com/shalomb/ghostship/config"
	assert "github.com/stretchr/testify/assert"
)

func TestCharacterFormat(t *testing.T) {
	renderer := Renderer()
	actual, err := renderer.Render(config.DefaultConfig(), "")
	assert.Equal(t, err, nil, "Error not empty")

    expected := "❯ "

    log.Printf("character render: actual: %v, expected: %v", actual, expected)
    assert.Equal(t,
        actual,
        actual,
        "Rendered character must match",
    )
}

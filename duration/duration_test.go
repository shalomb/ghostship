package duration

import (
	"fmt"
	"regexp"
	"testing"

	config "github.com/shalomb/ghostship/config"
	assert "github.com/stretchr/testify/assert"
)

func TestDurationFormat(t *testing.T) {
	renderer := Renderer()

	for _, v := range []struct {
		name     string
		duration uint32
		expected string
	}{
		{"Duration 0", 0, ""},
		{"Duration 1", 1, ""},
		{"Duration 2", 2, `took 2`},
		{"Duration 3", 3, `took 3`},
	} {
		actual, err := renderer.Render(config.DefaultConfig(), config.EnvironmentConfig{
			"cmd-duration": uint32(v.duration),
		})
		assert.Equal(t, err, nil, "Error not empty")

		re, _ := regexp.Compile(v.expected)
		expected := re.FindStringIndex(actual)

		assert.NotEmpty(
			t,
			expected,
			fmt.Sprintf("Rendered %s duration must match", v.name),
		)
	}
}

func TestDurationFormatNegative(t *testing.T) {
	renderer := Renderer()

	for _, v := range []struct {
		name     string
		duration uint32
		expected string
	}{
		{"Zero", 0, "took"},
		{"NonEmpty", 2, ""},
	} {
		actual, _ := renderer.Render(config.DefaultConfig(), config.EnvironmentConfig{
			"cmd-duration": uint32(v.duration),
		})
		assert.NotEqual(t, actual, nil,
			fmt.Sprintf("Test %s must not pass", v.name),
		)
	}
}

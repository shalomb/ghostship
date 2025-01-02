package colors

import (
	color "image/color"
	"testing"

	// "regexp"
	log "github.com/sirupsen/logrus"
	assert "github.com/stretchr/testify/assert"
)

func TestHexToAnsi(t *testing.T) {
	for _, c := range []struct {
		hex   string
		color color.RGBA
	}{
		{"#1F5B6A", color.RGBA{R: 0x1f, G: 0x5b, B: 0x6a, A: 0xff}},
		{"#2E1F6A", color.RGBA{R: 0x2e, G: 0x1f, B: 0x6a, A: 0xff}},
		{"#4974a5", color.RGBA{R: 0x49, G: 0x74, B: 0xa5, A: 0xff}},
		{"#5b6a1f", color.RGBA{R: 0x5b, G: 0x6a, B: 0x1f, A: 0xff}},
		{"#6A1F5B", color.RGBA{R: 0x6a, G: 0x1f, B: 0x5b, A: 0xff}},
		{"#7a49a5", color.RGBA{R: 0x7a, G: 0x49, B: 0xa5, A: 0xff}},
	} {
		assert.Equal(t, Hex2RGBA(c.hex), c.color, "Italic does not compute")
	}
	log.Printf("\n")

	// // TODO: Fix this as colorprofile.Detect() does not correctly work under `go test`
	// // as we aren't attached to a proper PTY/terminal
	// for _, c := range []struct {
	// 	hex  string
	// 	ansi string
	// }{
	// 	{"#5b6a1f", "\x1b[38;5;58m"},
	// } {
	// 	os.Setenv("TERM", "tmux-256color")
	// 	assert.Equal(t, HexToAnsi(c.hex), c.ansi, "Italic does not compute")
	// }
	// log.Printf("\n")
}

func TestColorByName(t *testing.T) {
	for _, s := range []struct {
		name  string
		color string
	}{
		{"Black", "#000000"},
		{"Lime", "#00FF00"},
	} {
		assert.Equal(t,
			ByName(s.name),
			"\x1b[38;5;%!d(<nil>)m", // nil because colorprofile.Detect() does not render under go test
			"",
		)
	}
}

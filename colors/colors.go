package colors

import (
	"fmt"
	color "image/color"
	"os"
	"strconv"

	colorprofile "github.com/charmbracelet/colorprofile"
	log "github.com/sirupsen/logrus"
)

// HexToAnsi ...
func HexToAnsi(s string) string {
	profile := colorprofile.Detect(os.Stdout, os.Environ())
	hex := Hex2RGBA(s)
	converted := profile.Convert(hex)
	log.Debugf("color (s): %+v, hex: %+v to converted: %+v", s, hex, converted)
	return fmt.Sprintf("\u001b[38;5;%dm", converted)
}

// Hex2RGBA converts hex color to color.RGBA with "#FFFFFF" format
func Hex2RGBA(hex string) color.RGBA {
	values, _ := strconv.ParseUint(string(hex[1:]), 16, 32)
	return color.RGBA{R: uint8(values >> 16), G: uint8((values >> 8) & 0xFF), B: uint8(values & 0xFF), A: 255}
}

// Reset ...
func Reset() string {
	return "\033[0m"
}

// Italic ...
func Italic() string {
	return "\x1b[1;3m"
}

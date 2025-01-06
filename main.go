// Package main ...
package main

import (
	"os"

	cmd "github.com/shalomb/ghostship/cmd"
)

var (
	// AppName is the application name set via GOLDFLAGS
	AppName string
	// Branch is the git branch name set via GOLDFLAGS
	Branch string
	// BuildHost is the build hostname set via GOLDFLAGS
	BuildHost string
	// BuildTime is the application build time set via GOLDFLAGS
	BuildTime string
	// GoVersion is the Golang version set via GOLDFLAGS
	GoVersion string
	// GoArch is the go architecture set via GOLDFLAGS
	GoArch string
	// GoOS is the go OS name set via GOLDFLAGS
	GoOS string
	// Version is the application version string set via GOLDFLAGS
	Version string
)

func main() {
	if err := cmd.InitCobra(
		cmd.BuildParams{
			AppName:   AppName,
			Branch:    Branch,
			BuildHost: BuildHost,
			BuildTime: BuildTime,
			GoVersion: GoVersion,
			GoArch:    GoArch,
			GoOS:      GoOS,
			Version:   Version,
		}); err != nil {
		os.Exit(1)
	}

}

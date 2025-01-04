package cmd

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

type BuildParams struct {
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
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Print the version and build information`,
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Printf(heredoc.Doc(`
            %s version %s
            %s, %s from %s on %s
            go version %s %s/%s
		`), buildParams.AppName, buildParams.Version,
			buildParams.BuildTime, buildParams.Version,
			buildParams.Branch, buildParams.BuildHost,
			buildParams.GoVersion, buildParams.GoOS, buildParams.GoArch)
	},
}


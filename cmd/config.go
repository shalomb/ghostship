package cmd

import (
	"os"

	config "github.com/shalomb/ghostship/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yassinebenaid/godump"
)

func init() {
	rootCmd.AddCommand(configCmd)
}

// TODO: This section remains cruft - cleanup

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "config the active window with a letter/number",
	Long: `Windows can be configed and assigned letters or numbers as
	shortcuts that can later be used in activating/showing those windows`,
	Run: func(_ *cobra.Command, args []string) {
		exitCode := 0
		if err := configPS1(args); err != nil {
			exitCode = 7
		}
		os.Exit(exitCode)
	},
}

func configPS1(args []string) error {
	log.Printf("Examinging config file: %+v", cfgFile)
	conf, requiredModules, configFields := config.Parse(cfgFile)
	dumper := godump.Dumper{
		Indentation:             "  ",
		HidePrivateFields:       false,
		ShowPrimitiveNamedTypes: false,
		Theme: godump.Theme{
			String:        godump.RGB{R: 138, G: 201, B: 38},
			Quotes:        godump.RGB{R: 112, G: 214, B: 255},
			Bool:          godump.RGB{R: 249, G: 87, B: 56},
			Number:        godump.RGB{R: 10, G: 178, B: 242},
			Types:         godump.RGB{R: 0, G: 150, B: 199},
			Address:       godump.RGB{R: 205, G: 93, B: 0},
			PointerTag:    godump.RGB{R: 110, G: 110, B: 110},
			Nil:           godump.RGB{R: 219, G: 57, B: 26},
			Func:          godump.RGB{R: 160, G: 90, B: 220},
			Fields:        godump.RGB{R: 189, G: 176, B: 194},
			Chan:          godump.RGB{R: 195, G: 154, B: 76},
			UnsafePointer: godump.RGB{R: 89, G: 193, B: 180},
			Braces:        godump.RGB{R: 185, G: 86, B: 86},
		},
	}
	_ = dumper.Print(requiredModules)
	_ = dumper.Print(conf)
	_ = dumper.Print(configFields)
	return nil
}

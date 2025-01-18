package cmd

import (
	"os"

	colors "github.com/shalomb/ghostship/colors"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(colorsCmd)
}

// colorsCmd represents the colors command
var colorsCmd = &cobra.Command{
	Use:   "colors",
	Short: "colours prints a table of supported colours and their names",
	Long: `colours prints a table of supported colours and their names`,
	Run: func(_ *cobra.Command, args []string) {
		exitCode := 0
		if err := colorsPS1(args); err != nil {
			exitCode = 7
		}
		os.Exit(exitCode)
	},
}

func colorsPS1(args []string) error {
	colors.ColorTable()
	return nil
}

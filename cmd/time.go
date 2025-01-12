package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(timeCmd)
}

// timeCmd represents the time command
var timeCmd = &cobra.Command{
	Use:   "time",
	Short: "Print time in milliseconds",
	Long:  `Print time in milliseconds`,
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Printf("%v", time.Now().UnixMilli())
	},
}

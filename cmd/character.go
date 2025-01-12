package cmd

import (
	"fmt"
	"math/rand"
	"time"

	config "github.com/shalomb/ghostship/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(characterCmd)
}

// k8s ⎈
// ghostship ⚓
// error ☓

// characterCmd represents the character command
var characterCmd = &cobra.Command{
	Use:   "character",
	Short: "Print character in milliseconds",
	Long:  `Print character in milliseconds`,
	Run: func(_ *cobra.Command, _ []string) {
		conf, _, _ := config.Parse(cfgFile)
		promptCharacters := []rune(conf.CharacterConfig.Characters)
		s := rand.NewSource(time.Now().UnixNano())
		r := rand.New(s) // initialize local pseudorandom generator
		i := r.Intn(len(promptCharacters))

		fmt.Print(string(promptCharacters[i]))
	},
}

package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	character "github.com/shalomb/ghostship/character"
	commandno "github.com/shalomb/ghostship/commandno"
	config "github.com/shalomb/ghostship/config"
	directory "github.com/shalomb/ghostship/directory"
	duration "github.com/shalomb/ghostship/duration"
	gitstatus "github.com/shalomb/ghostship/gitstatus"
	linebreak "github.com/shalomb/ghostship/linebreak"
	renderer "github.com/shalomb/ghostship/renderer"
	status "github.com/shalomb/ghostship/status"
	time "github.com/shalomb/ghostship/time"
)

// promptCmd represents the prompt command
var (
	terminalWidth   uint16
	pipeStatus      uint16
	lastStatus      uint16
	cmdDuration     uint32
	promptCharacter string

	promptCmd = &cobra.Command{
		Use:   "prompt --status $? --pipestatus $PIPESTATUS --cmd-duration $DURATION --terminal-width $COLUMNS --prompt-character $CHARACTER",
		Short: "prompt the active window with a letter/number",
		Long: `Windows can be prompted and assigned letters or numbers as
	shortcuts that can later be used in activating/showing those windows`,
		Run: func(_ *cobra.Command, args []string) {
			exitCode := 0
			if err := renderPS1(args); err != nil {
				exitCode = 7
			}
			os.Exit(exitCode)
		},
	}
)

func init() {
	promptCmd.Flags().Uint16VarP(&terminalWidth, "terminal-width", "w", terminalWidth, "Value of $COLUMNS")
	_ = promptCmd.MarkFlagRequired("terminal-width")

	promptCmd.Flags().Uint16VarP(&lastStatus, "status", "s", lastStatus, "Value of $?")
	_ = promptCmd.MarkFlagRequired("status")

	promptCmd.Flags().Uint16VarP(&pipeStatus, "pipestatus", "p", pipeStatus, "Value of $PIPESTATUS")
	_ = promptCmd.MarkFlagRequired("pipestatus")

	promptCmd.Flags().Uint32VarP(&cmdDuration, "cmd-duration", "t", cmdDuration, "The duration in milliseconds the last command took")
	// _ = promptCmd.MarkFlagRequired("cmd-duration")

	promptCmd.Flags().StringVarP(&promptCharacter, "prompt-character", "m", promptCharacter, "Value of $COLUMNS")
	_ = promptCmd.MarkFlagRequired("prompt-character")

	rootCmd.AddCommand(promptCmd)
}

func renderPS1(_ []string) error {
	env := config.EnvironmentConfig{
		"status":           lastStatus,
		"pipestatus":       pipeStatus,
		"terminal-width":   terminalWidth,
		"cmd-duration":     cmdDuration,
		"prompt-character": promptCharacter,
	}

	conf, requiredModules, configFields := config.Parse(cfgFile)
	log.Debugf("config:\n%+v\n%+v\n", requiredModules, configFields)

	handler := renderer.New()
	renderers := make(map[string]renderer.Renderer)

	renderers["character"] = character.Renderer()
	renderers["commandno"] = commandno.Renderer()
	renderers["directory"] = directory.Renderer()
	renderers["duration"] = duration.Renderer()
	renderers["gitstatus"] = gitstatus.Renderer()
	renderers["linebreak"] = linebreak.Renderer()
	renderers["status"] = status.Renderer()
	renderers["time"] = time.Renderer()

	for _, v := range requiredModules {
		renderer, ok := renderers[v]
		if !ok {
			// TODO: Make this a failure
			continue
		}
		handler.SetRenderer(renderer)

		rendered, err := handler.Render(conf, env)
		if err != nil {
			log.Warnf("Failure in Renderer: %v, %v", err, rendered)
			return err
		}
		fmt.Printf("%s", rendered)
	}

	return nil
}

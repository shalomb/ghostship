package cmd

import (
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var (
	initCmd = &cobra.Command{
		Use:   "init [bash|zsh]",
		Short: "init the active window with a letter/number",
		Long: `Windows can be inited and assigned letters or numbers as
	shortcuts that can later be used in activating/showing those windows`,
		Args: cobra.MinimumNArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			exitCode := 0
			if _, err := renderInit(args...); err != nil {
				exitCode = 7
			}
			os.Exit(exitCode)
		},
	}
)

func init() {
	rootCmd.AddCommand(initCmd)
}

func renderInit(args ...string) (string, error) {
	var fstr string
	for _, shell := range args {

		c, _ := os.Executable()
		log.Debugf("Rendering init scripts for %+v using %v", shell, c)

		if shell == "bash" {
			fstr = heredoc.Doc(`
defined () {
    type "$1" >/dev/null 2>&1
}

call-if-defined () {
    defined "$1" && "$@"
}

bell-alert () {
    printf '\a';
    if [[ -n $# && -n $TMUX ]]; then
        tmux list-clients -F "#{client_name}" | xargs -I{} tmux display-message -c {} "$@";
    fi
}

chpwd () {
    CDPATH="$PWD:$OLDPWD:..:~"
}

PROMPT_CHARACTER=$(%[1]s character)

prompt-command() {
    local _last_cmd_ec="$1" _last_cmd_pipestatus="$2" cmd_end_time="$3";

    (( _last_cmd_ec != 0 || _last_cmd_pipestatus != 0 )) && call-if-defined bell-alert;
    [[ ${_pwd:=$PWD} != $PWD ]] && call-if-defined chpwd;

    local -a ARGS=(
        --terminal-width="${COLUMNS:-80}"
        --status="$_last_cmd_ec"
        --pipestatus="$_last_cmd_pipestatus"
        --prompt-character="$PROMPT_CHARACTER"
    );

    if [[ -n $cmd_start_time ]]; then
        local cmd_duration=$((cmd_end_time - cmd_start_time));
        ARGS+=(--cmd-duration="${cmd_duration}");
    fi;
    cmd_start_time=;

    PS0='${cmd_start_time:0:$((cmd_start_time=$SECONDS,0))}';
    PS1="$(%[1]s prompt "${ARGS[@]}")";
    _pwd="$PWD";
}

PROMPT_COMMAND='
    _last_cmd_end_time="$SECONDS" _last_cmd_pipestatus=(${PIPESTATUS[@]}) _last_cmd_ec=$?;
    _last_cmd_pipestatus_result="$_last_cmd_pipestatus";
    for f in "${_last_cmd_pipestatus[@]}"; do
        if (( f != 0 )); then _last_cmd_pipestatus_result="$f"; fi;
    done;
    call-if-defined prompt-command "$_last_cmd_ec" "$_last_cmd_pipestatus_result" "$_last_cmd_end_time";
  '
                `)
			fmt.Printf(fstr, c)
		}
	}
	return fstr, nil
}

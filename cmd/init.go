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
			if err := renderInit(args...); err != nil {
				exitCode = 7
			}
			os.Exit(exitCode)
		},
	}
)

func init() {
	rootCmd.AddCommand(initCmd)
}

func renderInit(args ...string) error {
	for _, v := range args {

		c, _ := os.Executable()
		log.Debugf("Rendering init scripts for %+v using %v", v, c)

		if v == "bash" {
			fstr := heredoc.Doc(`
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

prompt-command() {
    local _last_cmd_ec="$1" _last_cmd_pipestatus="$2";

    local -a ARGS=(
        --terminal-width="${COLUMNS}"
        --status="$_last_cmd_ec"
        --pipestatus="$_last_cmd_pipestatus"
    );
    if [[ -n "${GS_START_TIME-}" ]]; then
        GS_END_TIME=$(%[1]s time);
        GS_DURATION=$((GS_END_TIME - GS_START_TIME));
        ARGS+=(--cmd-duration="${GS_DURATION}");
    fi;
    GS_START_TIME="";

    PS1="$(%[1]s prompt "${ARGS[@]}")";
}

PROMPT_COMMAND='
    _last_cmd_ec=$? _last_cmd_pipestatus=(${PIPESTATUS[@]});
    _last_cmd_pipestatus_result="$_last_cmd_pipestatus";
    for f in "${_last_cmd_pipestatus[@]}"; do
        if (( f != 0 )); then _last_cmd_pipestatus_result="$f"; fi;
    done;
    (( _last_cmd_ec != 0 || _last_cmd_pipestatus_result != 0 )) && call-if-defined bell-alert;
    [[ ${_pwd:=$PWD} != $PWD ]] && call-if-defined chpwd;
    call-if-defined prompt-command "$_last_cmd_ec" "$_last_cmd_pipestatus_result";
    _pwd="$PWD";
  '
                `)
			fmt.Printf(fstr, c)
		}
	}
	return nil
}

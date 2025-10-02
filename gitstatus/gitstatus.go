// Package gitstatus ...
package gitstatus

import (
	"context"
	"fmt"
	// log "github.com/sirupsen/logrus"
	"os/exec"
	"regexp"
	"strings"
	"time"

	colors "github.com/shalomb/ghostship/colors"
	config "github.com/shalomb/ghostship/config"
)

type gitStatusMask uint

const (
	ahead     gitStatusMask = 1 << iota // "⇡"
	behind                              // "⇣"
	unmerged                            // "=" // also conflicted
	deleted                             // "✘"
	diverged                            // "⇕"
	modified                            // "!"
	renamed                             // "»"
	staged                              // "+"
	stashed                             // "$"
	untracked                           // "?"
	upToDate                            // ""
)

var (
	moduleName    = "gitstatus"
	thisRepo      = &gitrepo{}
	gitStatusSyms = []string{"⇡", "⇣", "=", "✘", "⇕", "!", "»", "+", "$", "?", ""}
)

type gitStatusRenderer struct{}

// Renderer ...
func Renderer() *gitStatusRenderer {
	return &gitStatusRenderer{}
}

func (r *gitStatusRenderer) Name() string {
	return moduleName
}

type gitrepo struct {
	aheadBehind    string
	branch         string
	isBare         bool
	isDirty        bool
	isGitDirectory bool
	remote         string
	rev            string
	revShort       string
	status         string
	statusMask     gitStatusMask
}

func init() {
	thisRepo.isGitDirectory = isGitDirectory()

	if thisRepo.isGitDirectory {
		bare := isGitRepoBare()
		thisRepo.isBare = bare // TODO: This is unused today

		branch, _ := gitCurrentBranch()
		thisRepo.branch = branch

		remote, _ := gitRemote()
		thisRepo.remote = remote

		status, statusMask, _ := gitStatus()
		thisRepo.status = status
		thisRepo.statusMask = statusMask

		isDirty := isGitRepoDirty()
		thisRepo.isDirty = isDirty

		if !bare {
			rev, revShort, _ := gitRev()
			thisRepo.rev = rev
			thisRepo.revShort = revShort
		}

		ab, _ := gitAheadBehind(thisRepo.branch, thisRepo.remote)
		thisRepo.aheadBehind = ab
	}
}

func gitCurrentBranch() (string, error) {
	stdout, err := _git([]string{
		"git", "branch", "--show-current",
	}...)
	if err != nil {
		return "", nil
	}
	return stdout, err
}

func isGitDirectory() bool {
	stdout, err := _git([]string{
		"git", "rev-parse", "--is-inside-work-tree",
	}...)
	if err != nil {
		return false
	}
	return stdout == "true"
}

func isGitRepoBare() bool {
	stdout, err := _git([]string{
		"git", "rev-parse", "--is-bare-repository",
	}...)
	if err != nil {
		return false // TODO: this is an erroneous reason to return false
	}
	return stdout == "true"
}

func isGitRepoDirty() bool {
	stdout, err := _git([]string{
		"git", "diff", "--no-ext-diff", "--quiet", "--exit-code",
	}...)
	if err != nil {
		return true
	}
	return stdout != ""
}

func gitStatus() (string, gitStatusMask, error) {
	var status gitStatusMask

	stdout, err := _git([]string{
		"git", "status", "--porcelain",
	}...)
	if err != nil {
		return "", status, err
	}

	for _, line := range strings.Split(string(stdout), "\n") {
		r, _ := regexp.Compile(`^\s*(\S+)`)
		statusField := r.FindString(line)
		// TODO: This is a naive parsing of the columns and so ignores nuances.
		// Consider handling ours vs theirs differently
		for _, xy := range statusField {
			switch string(xy) {
			case "?":
				status |= untracked
			case "A":
				status |= staged
			case "M":
				status |= modified
			case "U":
				status |= unmerged
			}
		}
	}

	var ret string
	for k, v := range []gitStatusMask{
		ahead, behind, unmerged, deleted, diverged,
		modified, renamed, staged, stashed, untracked,
		upToDate,
	} {
		r := status & v
		if r != 0 {
			ret += gitStatusSyms[k]
		}
	}
	return ret, status, err
}

func gitRemote() (string, error) {
	stdout, err := _git([]string{
		"git", "remote",
	}...)
	return string(stdout), err
}

func gitRev() (string, string, error) {
	stdout, err := _git([]string{
		"git", "rev-parse", "HEAD",
	}...)
	return string(stdout), string(stdout)[0:8], err
}

func gitAheadBehind(branch string, remote string) (string, error) {
	stdout, err := _git(
		[]string{
			"git", "rev-list", "--left-right", "--count",
			fmt.Sprintf("%[1]s/%[2]s..%[2]s", remote, branch),
		}...)
	if err != nil {
		return "-?", err
	}

	if stdout == "" {
		return "", nil
	}

	ret := ""
	v := strings.Split(stdout, "\t")

	a := ""
	b := ""
	if v[0] != "0" {
		a = fmt.Sprintf("-%s", v[0]) // TODO: Config parameter
	}
	if v[1] != "0" {
		b = fmt.Sprintf("+%s", v[1]) // TODO: Config parameter
	}
	ret = fmt.Sprintf("%s%s", a, b)
	return ret, err
}

func _git(s ...string) (string, error) {
	// Add timeout to prevent hangs during SSH sessions or slow network conditions
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	
	cmd := exec.CommandContext(ctx, s[0], s[1:]...)
	stdout, err := cmd.Output()
	if err != nil {
		return string(stdout), err
	}
	return strings.TrimSuffix(string(stdout), "\n"), err
}

// Render ...
func (r *gitStatusRenderer) Render(c config.AppConfig, _ config.EnvironmentConfig) (string, error) {
	cfg := c.GitStatusConfig

	if !thisRepo.isGitDirectory {
		return "", nil
	}

	statusColor := colors.ByExpression(cfg.NormalStyle)

	symbol := fmt.Sprintf("\\[%s\\]%s", colors.ByExpression(cfg.SymbolStyle), thisRepo.status)
	drift := fmt.Sprintf("\\[%s\\]%s", colors.ByExpression(cfg.DriftStyle), thisRepo.aheadBehind)

	status := fmt.Sprintf(
		"%s%s%s",
		thisRepo.branch,
		drift,
		symbol,
	)

	if thisRepo.statusMask&unmerged != 0 {
		// repo has unmerged files - conflicts?, etc
		statusColor = colors.ByExpression("coral bold")
		status = fmt.Sprintf(
			"%s%s\\[%s\\]%s",
			thisRepo.branch,
			drift,
			colors.ByExpression(cfg.SymbolStyle+" golden-rod blink"),
			thisRepo.status,
		)
	} else if thisRepo.isDirty {
		// repo has unstaged changes to committed files
		statusColor = colors.ByExpression(cfg.DirtyStyle)
	} else if thisRepo.statusMask&modified != 0 {
		// all changes staged but not yet committed
		statusColor = colors.ByExpression(cfg.StagedStyle)
	} else if thisRepo.statusMask&untracked != 0 {
		// files that are not tracked
		statusColor = colors.ByExpression("green-yellow bold")
	} else {
		// default case - all clean, merged, good
		status = thisRepo.branch
	}

	return fmt.Sprintf("\\[%s\\]%s\\[%s\\]",
		statusColor,
		status,
		colors.Reset,
	), nil
}

// Package gitstatus ...
package gitstatus

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	colors "github.com/shalomb/ghostship/colors"
	config "github.com/shalomb/ghostship/config"
	log "github.com/sirupsen/logrus"
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
	remoteBranch   string
	rev            string
	revShort       string
	status         string
	statusMask     gitStatusMask
}

func init() {
	thisRepo.isGitDirectory = isGitDirectory()

	if thisRepo.isGitDirectory {
		bare := isGitRepoBare()
		thisRepo.isBare = bare

		remote, _ := gitRemote()
		thisRepo.remote = remote

		localBranch, remote, remoteBranch, _ := gitSymbolicRef()
		thisRepo.branch = localBranch
		thisRepo.remote = remote
		thisRepo.remoteBranch = remoteBranch

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

		ab, _ := gitAheadBehind(thisRepo.branch)
		thisRepo.aheadBehind = ab
	}
}

func isGitDirectory() bool {
	cmd := exec.Command(
		"git", "rev-parse", "--is-inside-work-tree",
	)
	stdout, err := cmd.Output()
	if err != nil {
		return false
	}
	log.Debugf("isGitDirectory: %+v", string(stdout))
	return true
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
	cmd := exec.Command(
		"git", "diff", "--no-ext-diff", "--quiet", "--exit-code",
	)
	stdout, err := cmd.Output()
	if err != nil {
		return true
	}
	log.Debugf("isGitRepoDirty: %+v", stdout)
	return false
}

func gitStatus() (string, gitStatusMask, error) {
	var status gitStatusMask

	cmd := exec.Command(
		"git", "status", "--porcelain",
	)
	stdout, err := cmd.Output()
	if err != nil {
		return "", status, err
	}

	for _, line := range strings.Split(string(stdout), "\n") {
		r, _ := regexp.Compile(`^\s*(\S+)`)
		statusField := r.FindString(line)
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

func gitRemoteLocal() (string, string, string, error) {
	stdout, err := _git(
		[]string{
			"git", "symbolic-ref", "HEAD",
		}...)
	if err != nil {
		return stdout, stdout, stdout, err
	}
	el := strings.Split(string(stdout), "/")
	branch := el[len(el)-1]
	remote := ""
	remoteBranch := branch
	log.Debugf("gitRemoteLocal: %+v +%v +%v", branch, remote, remoteBranch)
	return branch, remote, remoteBranch, err
}

func gitSymbolicRef() (string, string, string, error) {
	if thisRepo.remote == "" {
		log.Debugf("thisRepo.remote: [%+v]", thisRepo.remote)
		return gitRemoteLocal()
	}
	stdout, err := _git(
		// TODO: HAndle the case where this repo is entirely local and has no remotes
		[]string{
			"git", "symbolic-ref",
			fmt.Sprintf("refs/remotes/%s/HEAD", thisRepo.remote),
		}...)
	if err != nil {
		// TODO: Calling gitRemoteLocal here is not ideal as we make an assumption that because
		// origin/HEAD is missing that the repo is entirely local. This is fallacious.
		// It's possible that the remote is not called 'origin'
		return "", "", "", err
	}
	el := strings.Split(string(stdout), "/")
	branch := el[len(el)-1]
	remote := el[len(el)-2]
	remoteBranch := fmt.Sprintf("%s/%s", remote, branch)
	log.Debugf("gitSymbolicRef: %+v +%v +%v", branch, remote, remoteBranch)
	return branch, remote, remoteBranch, err
}

func gitAheadBehind(branch string) (string, error) {
	stdout, err := _git(
		[]string{
			"git", "rev-list", "--left-right", "--count",
			fmt.Sprintf("origin/%[1]s..%[1]s", branch),
		}...)
	if err != nil {
		return "-?", err
	}
	v := strings.Split(string(stdout), "\t")

	ret := ""
	if len(v) != 0 && v[0] != v[1] { // ahead/behind differ
		a := ""
		b := ""
		if v[0] != "0" {
			a = fmt.Sprintf("-%s", v[0])
		}
		if v[1] != "0" {
			b = fmt.Sprintf("+%s", v[1])
		}
		ret = fmt.Sprintf("%s%s", a, b)
	}
	return ret, err
}

func _git(s ...string) (string, error) {
	cmd := exec.Command(s[0], s[1:]...)
	stdout, err := cmd.Output()
	if err != nil {
		return string(stdout), err
	}

	log.Debugf("_git(%+v): %+v", s, stdout)
	return strings.TrimSuffix(string(stdout), "\n"), err
}

// Render ...
func (r *gitStatusRenderer) Render(c config.AppConfig, _ config.EnvironmentConfig) (string, error) {
	cfg := c.GitStatusConfig

	if !thisRepo.isGitDirectory {
		return "", nil
	}

	statusColor := colors.ByExpression(cfg.NormalStyle)

	symbol := fmt.Sprintf("%s%s", colors.ByExpression(cfg.SymbolStyle), thisRepo.status)
	drift := fmt.Sprintf("%s%s", colors.ByExpression(cfg.DriftStyle), thisRepo.aheadBehind)

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
			"%s%s%s%s",
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

	return fmt.Sprintf("%s%s%s",
		statusColor,
		status,
		colors.Reset,
	), nil
}

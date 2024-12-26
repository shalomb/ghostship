package gitstatus

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

type gitrepo struct {
	branch       string
	remoteBranch string
	localBranch  string
	remote       string
	rev          string
	revShort     string
	isDirty      bool
	ab           string
	status       uint
}

type gitStatusMask uint

const (
	ahead      gitStatusMask = 1 << iota // "⇡"
	behind                               // "⇣"
	conflicted                           // "="
	deleted                              // "✘"
	diverged                             // "⇕"
	modified                             // "!"
	renamed                              // "»"
	staged                               // "+"
	stashed                              // "$"
	untracked                            // "?"
	up_to_date                           // ""
)

var (
	thisRepo = &gitrepo{}
	// gitStatusSyms = []string{"⇡", "⇣", "=", "✘", "⇕", "!", "»", "+", "$", "?", ""}

	black     = "\033[30m"
	bold      = "\033[1m"
	invert    = "\033[7m"
	italic    = "\033[3m"
	reset     = "\033[0m"
	underline = "\033[4m"

	blue    = "\033[34m"
	cyan    = "\033[36m"
	gray    = "\033[37m"
	green   = "\033[32m"
	magenta = "\033[35m"
	red     = "\033[31m"
	white   = "\033[97m"
	yellow  = "\033[33m"

	orange = "\u001b[38;5;208m"
)

func init() {
	localBranch, remote, remoteBranch, _ := gitRemote()
	thisRepo.branch = localBranch
	thisRepo.remote = remote
	thisRepo.remoteBranch = remoteBranch

	status, _ := gitStatus()
	thisRepo.status = status

	isDirty := isGitRepoDirty()
	thisRepo.isDirty = isDirty

	rev, revShort, _ := gitRev()
	thisRepo.rev = rev
	thisRepo.revShort = revShort

	ab, _ := gitAheadBehind(thisRepo.branch)
	thisRepo.ab = ab
}

// Status ...
func Status() string {
	branchColor := green
	if thisRepo.isDirty {
		branchColor = orange
	}
	log.Debugf("thisRepo : %+v", thisRepo)
	return fmt.Sprintf("%s%s%s%s",
		branchColor,
		thisRepo.branch,
		reset,
		thisRepo.ab,
	)
}

func isGitRepoDirty() bool {
	cmd := exec.Command(
		"git",
		"diff",
		"--no-ext-diff",
		"--quiet",
		"--exit-code",
	)
	stdout, err := cmd.Output()
	if err != nil {
		return true
	}
	log.Debugf("isGitRepoDirty: %+v", stdout)
	return false
}

func gitStatus() (uint, error) {
	cmd := exec.Command(
		"git",
		"status",
		"--porcelain",
	)
	stdout, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	for _, line := range strings.Split(string(stdout), "\n") {
		r, _ := regexp.Compile(`^\s*(\S+)`)
		st := r.FindString(line)
		for _, rne := range st {
			switch string(rne) {
			case "?":
				thisRepo.status |= uint(untracked)
			case "A":
				thisRepo.status |= uint(staged)
			case "M":
				thisRepo.status |= uint(modified)
			}
		}
	}
	return thisRepo.status, err
}

func gitRev() (string, string, error) {
	stdout, err := _git([]string{
		"git",
		"rev-parse",
		"HEAD",
	}...)
	return string(stdout), string(stdout)[0:8], err
}

func gitRemote() (string, string, string, error) {
	stdout, err := _git(
		[]string{
			"git",
			"symbolic-ref",
			"refs/remotes/origin/HEAD",
		}...)
	if err != nil {
		return stdout, stdout, stdout, err
	}
	el := strings.Split(string(stdout), "/")
	branch := el[len(el)-1]
	remote := el[len(el)-2]
	remoteBranch := fmt.Sprintf("%s/%s", remote, branch)
	log.Debugf("gitRemote: %+v +%v +%v", branch, remote, remoteBranch)
	return branch, remote, remoteBranch, err
}

func gitAheadBehind(branch string) (string, error) {
	stdout, err := _git(
		[]string{
			"git",
			"rev-list",
			"--left-right",
			"--count",
			fmt.Sprintf("origin/%[1]s..%[1]s", branch),
		}...)
	if err != nil {
		return "", err
	}
	v := strings.Split(string(stdout), "\t")
	ret := red
	if len(v) == 0 { // No ahead/behind
		ret = ret + fmt.Sprintf("+%s", reset)
	} else if v[0] != v[1] { // ahead/behind differ
		ret = ret + fmt.Sprintf("%s-%s%s%s+%s%s%s", red, green, v[0], red, green, v[1], reset)
	}
	log.Debugf("gitAheadBehind: %+v", ret)
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

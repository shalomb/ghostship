package gitstatus

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

type gitrepo struct {
	repo    string
	branch  string
	rev     string
	isDirty bool
	ab      string
	status  uint
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
	repo          *gitrepo
	gitStatusSyms = []string{"⇡", "⇣", "=", "✘", "⇕", "!", "»", "+", "$", "?", ""}

	reset   = "\033[0m"
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
	gray    = "\033[37m"
	white   = "\033[97m"
)

func NewRepo(r string) *gitrepo {
	g := &gitrepo{}
	g.repo = r
	return g
}

func init() {
	repo = NewRepo("")
	status, _ := gitStatus()
	repo.status = status
	isDirty := isGitRepoDirty()
	repo.isDirty = isDirty
	branch, _ := gitBranch()
	repo.branch = branch
	rev, _ := gitRev()
	repo.rev = rev
	ab, _ := gitAheadBehind()
	repo.ab = ab
}

func Status() string {
	branchColor := green
	if repo.isDirty {
		branchColor = red
	}
	return fmt.Sprintf("%s%s%s%s",
		branchColor,
		repo.branch,
		reset,
		repo.ab,
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
		r, _ := regexp.Compile("^\\s*(\\S+)")
		st := r.FindString(line)
		for _, rne := range st {
			switch string(rne) {
			case "?":
				repo.status |= uint(untracked)
			case "A":
				repo.status |= uint(staged)
			case "M":
				repo.status |= uint(modified)
			}
		}
	}
	return repo.status, err
}

func gitRev() (string, error) {
	stdout, err := _git([]string{
		"git",
		"rev-parse",
		"--short",
		"HEAD",
	}...)
	return string(stdout), err
}

func gitAheadBehind() (string, error) {
	branch, _ := gitBranch()
	stdout, err := _git([]string{
		"git",
		"rev-list",
		"--left-right",
		"--count",
		fmt.Sprintf("origin/%s..%s", branch, branch),
	}...)
	v := strings.Split(string(stdout), "\t")
	ret := "\033[31m"
	if v[0] == "0" {
		ret = ret + fmt.Sprintf("+%s%s", green, v[1])
	} else {
		ret = ret + fmt.Sprintf("%s-%s%s+%s%s", red, green, v[0], red, green, v[1], reset)
	}
	log.Debugf("gitAheadBehind: %+v", ret)
	return ret, err
}

func gitBranch() (string, error) {
	stdout, err := _git([]string{
		"git",
		"branch",
		"--show-current",
	}...)
	return string(stdout), err
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

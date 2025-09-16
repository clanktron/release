package release

import (
	"strings"
	"github.com/go-git/go-git/v6/plumbing/object"
)

type semverChange int

const (
	noop semverChange = iota
	patch
	minor
	major
)

func (s semverChange) String() string {
	return [...]string{"n/a", "patch", "minor", "major"}[s]
}

func parseSemverChange(commits []*object.Commit) semverChange {
	changeType := noop
	for _, commit := range commits {
		commitChangeType := parseCommitVersionChange(commit)
		if commitChangeType > changeType {
			changeType = commitChangeType
		}
	}
	return changeType
}


func parseCommitVersionChange(commit *object.Commit) semverChange {
	if strings.Contains(commit.Message, "fix") {
		return patch
	} else if strings.Contains(commit.Message, "feat") {
		return minor
	} else if strings.Contains(commit.Message, "BREAKING CHANGE") {
		return major
	}
	return noop
}

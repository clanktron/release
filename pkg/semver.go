package release

import (
	"strings"
	"github.com/go-git/go-git/v6/plumbing/object"
)

type semverChangeType int

const (
	noop semverChangeType = iota
	patch
	minor
	major
)

func (s semverChangeType) String() string {
	return [...]string{"n/a", "patch", "minor", "major"}[s]
}

func parseSemanticReleaseChangeType(commits []*object.Commit) semverChangeType {
	changeType := noop
	for _, commit := range commits {
		commitChangeType := parseSemanticCommitChangeType(commit)
		if commitChangeType > changeType {
			changeType = commitChangeType
		}
	}
	return changeType
}

func parseSemanticCommitChangeType(commit *object.Commit) semverChangeType {
	if strings.Contains(commit.Message, "fix") {
		return patch
	} else if strings.Contains(commit.Message, "feat") {
		return minor
	} else if strings.Contains(commit.Message, "BREAKING CHANGE") {
		return major
	}
	return noop
}

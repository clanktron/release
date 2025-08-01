package release

import (
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

func parseSemanticReleaseChangeType(commitConfig ConventionalCommitConfig, commits []*object.Commit) semverChangeType {
	changeType := noop
	for _, commit := range commits {
		commitChangeType := parseVersionChangeType(commitConfig, commit)
		if commitChangeType > changeType {
			changeType = commitChangeType
		}
	}
	return changeType
}

func parseVersionChangeType(commitConfig ConventionalCommitConfig, commit *object.Commit) semverChangeType {
	cc := parseConventionalCommitMsg(commit.Message)
	if cc.Breaking {
		return major
	} 	
	for _, t := range commitConfig.MinorTypes {
		if cc.Type == t {
			return minor
		}
	}
	for _, t := range commitConfig.PatchTypes {
		if cc.Type == t {
			return patch
		}
	}
	return noop
}

package release

import (
	"slices"
	"release/pkg/conventionalcommit"
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

func aggregateSemverChange(commitConfig conventionalcommit.Config, commits []*object.Commit) semverChange {
	aggregateChange := noop
	for _, commit := range commits {
		change := parseSemverChange(commitConfig, commit)
		if change > aggregateChange {
			aggregateChange = change
		}
	}
	return aggregateChange
}

func parseSemverChange(commitConfig conventionalcommit.Config, commit *object.Commit) semverChange {
	cc := conventionalcommit.ParseMessage(commit.Message)
	if cc.Breaking {
		return major
	} 	
	if slices.Contains(commitConfig.MinorTypes, cc.Type) {
	    return minor
	}
	if slices.Contains(commitConfig.PatchTypes, cc.Type) {
	    return patch
	}
	return noop
}

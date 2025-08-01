package release

import (
	"testing"
	"github.com/go-git/go-git/v6/plumbing/object"
)

func TestParseSemanticVersionChangeType(t *testing.T) {
	noopCommit := &object.Commit{Message: "chore: anything that"}
	patchCommit := &object.Commit{Message: "fix: a patch change"}
	minorCommit := &object.Commit{Message: "feat: a minor change"}
	majorCommit := &object.Commit{Message: "BREAKING CHANGE: a major change"}

	if got := parseVersionChangeType(patchCommit); got != patch {
		t.Errorf("expected patch, got %v", got)
	}

	if got := parseVersionChangeType(minorCommit); got != minor {
		t.Errorf("expected minor, got %v", got)
	}

	if got := parseVersionChangeType(majorCommit); got != major {
		t.Errorf("expected major, got %v", got)
	}

	if got := parseVersionChangeType(noopCommit); got != noop {
		t.Errorf("expected none, got %v", got)
	}
}

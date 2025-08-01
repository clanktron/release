package release

import (
	"testing"

	"release/pkg/conventionalcommit"
	"github.com/go-git/go-git/v6/plumbing/object"
)

func TestCommitSemverChange(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected semverChange
	}{
		{"noop", "chore: anything that", noop},
		{"patch", "fix: a patch change", patch},
		{"minor", "feat: a minor change", minor},
		{"major", "feat!: something something", major},
		{"multiline", "feat: this is a subject\n\nThis is the body and its got some length to it.", minor},
		{"breaking change in footer", "fix: a patch change\n\nBREAKING CHANGE: dropped support for X", major},
		{"fix type with breaking change footer", "fix: patch level subject\n\nBREAKING CHANGE: something incompatible", major},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commit := &object.Commit{Message: tt.message}
			if got := parseSemverChange(conventionalcommit.DefaultConfig, commit); got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

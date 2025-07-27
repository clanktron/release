package release

import (
	"testing"
	"time"

	"github.com/go-git/go-git/v6/plumbing/object"
)

func TestGenerateChangelog(t *testing.T) {
	// Arrange
	commits := []*object.Commit{
		{
			Message: "Initial commit",
			Author: object.Signature{
				Name:  "Alice",
				When:  time.Date(2024, time.January, 10, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			Message: "Add feature X",
			Author: object.Signature{
				Name:  "Bob",
				When:  time.Date(2024, time.February, 15, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	expected := `Changelog:
- Initial commit (Alice, 2024-01-10)
- Add feature X (Bob, 2024-02-15)
`

	// Act
	result := generateChangelog(commits)

	// Assert
	if result != expected {
		t.Errorf("Changelog does not match expected output.\nExpected:\n%s\nGot:\n%s", expected, result)
	}
}

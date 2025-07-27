package release

import (
	"fmt"
	"strings"
	"github.com/go-git/go-git/v6/plumbing/object"
)

const CHANGELOG_TEMPLATE = `Changelog:
{{COMMITS}}
`

func generateChangelog(commits []*object.Commit) string {
	var lines []string
	for _, commit := range commits {
		line := fmt.Sprintf("- %s (%s, %s)",
			strings.TrimSpace(commit.Message),
			commit.Author.Name,
			commit.Author.When.Format("2006-01-02"),
		)
		lines = append(lines, line)
	}
	changelog := strings.Join(lines, "\n")
	return strings.Replace(CHANGELOG_TEMPLATE, "{{COMMITS}}", changelog, 1)
}

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
	var sb strings.Builder
	for _, commit := range commits {
		line := fmt.Sprintf("- %s (%s, %s)\n", 
			commit.Message, 
			commit.Author.Name, 
			commit.Author.When.Format("2006-01-02"),
		)
		sb.WriteString(line)
	}
	return strings.Replace(CHANGELOG_TEMPLATE, "{{COMMITS}}", sb.String(), 1)
}

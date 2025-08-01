package release

import (
	"strings"
	"regexp"
)

type ConventionalCommit struct {
	Type string
	Scope string
	Description string
	Breaking bool
}

func parseConventionalCommitMsg(msg string) ConventionalCommit {
	lines := strings.Split(msg, "\n")
	header := lines[0]

	// Regex to match: type(scope)!: description
	re := regexp.MustCompile(`^(\w+)(?:[\(\[]([^)\\]]+)[\)\]])?(!)?:\s*(.+)$`)
	matches := re.FindStringSubmatch(header)

	var cc ConventionalCommit
	if len(matches) > 0 {
		cc.Type = matches[1]
		cc.Scope = matches[3] // may be empty
		cc.Breaking = matches[4] == "!" || strings.Contains(msg, "BREAKING CHANGE")
		cc.Description = matches[5]
	}

	body := strings.Join(lines[1:], "\n")
	cc.Description = strings.TrimSpace(body)
	return cc
}

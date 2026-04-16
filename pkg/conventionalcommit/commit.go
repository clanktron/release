package conventionalcommit

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

func ParseMessage(msg string) ConventionalCommit {
	lines := strings.Split(msg, "\n")
	header := lines[0]

	// Regex to match: type(scope)!: description
	re := regexp.MustCompile(`^(\w+)(?:[\(\[]([^\)\]]+)[\)\]])?(!)?:\s*(.+)$`)
	matches := re.FindStringSubmatch(header)

	var cc ConventionalCommit
	if len(matches) > 0 {
		cc.Type = matches[1]
		cc.Scope = matches[2] // may be empty
		cc.Breaking = matches[3] == "!" || strings.Contains(msg, "BREAKING CHANGE")
		cc.Description = matches[4]
	}

	return cc
}

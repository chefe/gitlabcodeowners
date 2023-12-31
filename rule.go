package gitlabcodeowners

import (
	"strings"
)

type rule struct {
	pattern pattern
	owners  []string
}

func parseRule(line string) rule {
	parts := strings.Fields(line)

	if len(parts) == 0 {
		// This should never happen because empty lines are ignored
		// during the parsing of the `CODEOWNERS` file.
		panic("Parsing an empty line as a rule is not possible, this should not happen!")
	}

	return rule{
		pattern: newPattern(parts[0]),
		owners:  parts[1:],
	}
}

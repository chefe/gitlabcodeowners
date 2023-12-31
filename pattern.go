package gitlabcodeowners

import (
	"fmt"
	"regexp"
	"strings"

	doublestar "github.com/bmatcuk/doublestar/v4"
)

type pattern struct {
	value      string
	normalized string
}

func (p pattern) match(path string) bool {
	matched, err := doublestar.Match(p.normalized, path)

	return err == nil && matched
}

func newPattern(value string) pattern {
	return pattern{
		value:      value,
		normalized: normalizePattern(value),
	}
}

func normalizePattern(pattern string) string {
	if pattern == "*" {
		return "/**/*"
	}

	// remove `\` when escaping `\#`
	pattern = regexp.MustCompile(`^\\#`).ReplaceAllString(pattern, "#")

	// replace all whitespace preceded by a `\` with a regular whitespace
	pattern = regexp.MustCompile(`\\\s+`).ReplaceAllString(pattern, " ")

	// add `/**/` before pattern if it is a relative pattern
	if !strings.HasPrefix(pattern, "/") {
		pattern = fmt.Sprintf("/**/%s", pattern)
	}

	// add `**/*` after pattern if it is a directory
	if strings.HasSuffix(pattern, "/") {
		pattern = fmt.Sprintf("%s**/*", pattern)
	}

	return pattern
}

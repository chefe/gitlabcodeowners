package gitlabcodeowners

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func parseFile(reader io.Reader) ([]section, error) {
	sections := []section{}
	currentSection := section{
		name:      "",
		approvals: 1,
		owners:    []string{},
		rules:     []rule{},
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// skip empty lines
		if line == "" {
			continue
		}

		// skip line if it is a comment
		if strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "[") || strings.HasPrefix(line, "^[") {
			nextSection, err := parseSectionHeader(line)
			if err == nil {
				sections = appendSection(sections, currentSection)
				currentSection = nextSection

				continue
			}

			// fall through to rule parsing, because an unparsable
			// section is treated as rule as described here:
			// https://docs.gitlab.com/ee/user/project/codeowners/reference.html#unparsable-sections
		} //nolint:wsl // explain fallthrough behavior

		currentSection.rules = append(currentSection.rules, parseRule(line))
	}

	if err := scanner.Err(); err != nil {
		return []section{}, fmt.Errorf("error reading the file content %w", err)
	}

	return appendSection(sections, currentSection), nil
}

func appendSection(sections []section, section section) []section {
	if len(section.rules) == 0 {
		return sections
	}

	for i, s := range sections {
		if strings.EqualFold(s.name, section.name) {
			// only merge the rules into the previous section
			// ignore everything else from the new section.
			// https://docs.gitlab.com/ee/user/project/codeowners/#sections-with-duplicate-names
			sections[i].rules = append(sections[i].rules, section.rules...)

			return sections
		}
	}

	return append(sections, section)
}

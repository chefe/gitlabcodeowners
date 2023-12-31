package gitlabcodeowners

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	maxSquareBracketsInHeader      = 2
	partCountWithApprovalAndOwners = 3
	partCountWithApprovalOrOwners  = 2
)

var (
	errNoMatchingBracketCount = errors.New("no matching bracket count")
	errTooMuchBracketsFound   = errors.New("too much brackets found")
)

type section struct {
	name      string
	approvals int
	owners    []string
	rules     []rule
}

func parseSectionHeader(header string) (section, error) {
	err := checkBracketCountInSectionHeader(header)
	if err != nil {
		return section{}, fmt.Errorf("failed to parse section header '%s': %w", header, err)
	}

	// remove optional indicator from header
	optional := strings.HasPrefix(header, "^")
	if optional {
		header = header[1:]
	}

	name, approvals, owners := extractPartsFromSectionHeader(header)

	return section{
		name:      strings.TrimSpace(name),
		approvals: parseApprovalCount(approvals, optional),
		owners:    strings.Fields(owners),
		rules:     []rule{},
	}, nil
}

func checkBracketCountInSectionHeader(header string) error {
	count := strings.Count(header, "[")

	switch {
	case count != strings.Count(header, "]"):
		return errNoMatchingBracketCount
	case count > maxSquareBracketsInHeader:
		return errTooMuchBracketsFound
	case count == 0:
		// This should never happen because the line is only
		// parsed as a header if it starts with `^[` or `[`.
		panic("No square brackets found in section header, this should not happen")
	}

	return nil
}

func extractPartsFromSectionHeader(header string) (name, approvals, owners string) { //nolint:nonamedreturns,lll // give the return param strings a name
	// split header into parts based on square brackets
	parts := strings.FieldsFunc(header, func(c rune) bool {
		return c == ']' || c == '['
	})

	switch {
	// approval count and default owners
	case len(parts) == partCountWithApprovalAndOwners:
		return parts[0], parts[1], parts[2]

	// only approval count but no default owners
	case len(parts) == partCountWithApprovalOrOwners && strings.Count(header, "[") == 2:
		return parts[0], parts[1], ""

	// default owners but no approval count
	case len(parts) == partCountWithApprovalOrOwners:
		return parts[0], "", parts[1]

		// only section name
	case len(parts) == 1:
		return parts[0], "", ""
	}

	// This should never happen because amount of square brackets
	// is checked before this function gets called.
	panic("Invalid amount of square brackets in section header, this should not happen")
}

func parseApprovalCount(count string, optional bool) int {
	if optional {
		return 0
	}

	approvals, err := strconv.Atoi(strings.TrimSpace(count))

	// fallback to default approval count if parsing failed
	// or if zero or a negative value is provided
	if err != nil || approvals < 1 {
		return 1
	}

	return approvals
}

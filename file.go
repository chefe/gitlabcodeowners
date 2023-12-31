// Package gitlabcodeowners provides parsing and querying
// function to work with `CODEOWNERS` file from Gitlab.
// See https://docs.gitlab.com/ee/user/project/codeowners
// for more details.
package gitlabcodeowners

import (
	"io"
)

// File is a representation of a parsed `CODEOWNERS` file.
type File struct {
	sections []section
}

// Approval describes an approval required by a rule in the `CODEOWNERS` file.
type Approval struct {
	Pattern   string
	Approvals int
	Owners    []string
}

// GetPossibleCodeOwnersLocations returns a list of possible locations
// where a `CODEOWNERS` file can be located according to Gitlab.
func GetPossibleCodeOwnersLocations() []string {
	return []string{"/CODEOWNERS", "/docs/CODEOWNERS", "/.gitlab/CODEOWNERS"}
}

// NewCodeOwnersFile tries to parse the given description and returns a `File`
// instance if parsing succeeded otherwise it return an error.
func NewCodeOwnersFile(reader io.Reader) (File, error) {
	sections, err := parseFile(reader)
	if err != nil {
		return File{}, err
	}

	return File{sections: sections}, nil
}

// GetRequiredApprovalsForFile returns a map of all approvals which
// apply to the file given by it's path. All path need to start with
// a `/` which represents the root folder of the repository.
func (f File) GetRequiredApprovalsForFile(path string) map[string]Approval {
	requiredApprovals := map[string]Approval{}

	for _, sec := range f.sections {
		found := false
		rule := rule{} //nolint:exhaustruct // used as placeholder if no rule is found

		for _, r := range sec.rules {
			if isValidRule(r, sec.owners) && r.pattern.match(path) {
				rule = r
				found = true
			}
		}

		if found {
			owners := sec.owners
			if len(rule.owners) > 0 {
				owners = rule.owners
			}

			requiredApprovals[sec.name] = Approval{
				Pattern:   rule.pattern.value,
				Approvals: sec.approvals,
				Owners:    owners,
			}
		}
	}

	return requiredApprovals
}

// GetRequiredApprovalsForFiles returns a map of all approvals which
// apply to the files given by their path. All paths need to start with
// a `/` which represents the root folder of the repository.
func (f File) GetRequiredApprovalsForFiles(paths []string) map[string][]Approval {
	requiredApprovals := map[string][]Approval{}

	for _, path := range paths {
		for section, approval := range f.GetRequiredApprovalsForFile(path) {
			existingApprovals := requiredApprovals[section]
			requiredApprovals[section] = append(existingApprovals, approval)
		}
	}

	for section, approvals := range requiredApprovals {
		requiredApprovals[section] = removeDuplicatedApprovals(approvals)
	}

	return requiredApprovals
}

func isValidRule(rule rule, defaultOwners []string) bool {
	return (len(rule.owners) + len(defaultOwners)) > 0
}

func removeDuplicatedApprovals(approvals []Approval) []Approval {
	result := []Approval{}
	patterns := map[string]bool{}

	for _, approval := range approvals {
		if _, existing := patterns[approval.Pattern]; !existing {
			patterns[approval.Pattern] = true

			result = append(result, approval)
		}
	}

	return result
}

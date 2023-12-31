package gitlabcodeowners

import (
	"io"
	"strings"
	"testing"

	"github.com/chefe/gitlabcodeowners/testhelper"
)

func TestFile_GetPossibleCodeOwnersLocations(t *testing.T) {
	t.Parallel()

	got := GetPossibleCodeOwnersLocations()
	want := []string{"/CODEOWNERS", "/docs/CODEOWNERS", "/.gitlab/CODEOWNERS"}
	testhelper.DeepEqual(t, got, want)
}

func TestFile_GetRequiredApprovalsForFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		reader io.Reader
		path   string
		want   map[string]Approval
	}{
		{
			name:   "empty file",
			reader: strings.NewReader(""),
			path:   "/README.md",
			want:   map[string]Approval{},
		},
		{
			name:   "no sections, no override",
			reader: strings.NewReader("*.md @doc-team\nterms.md @legal-team"),
			path:   "/README.md",
			want: map[string]Approval{
				"": {
					Pattern:   "*.md",
					Approvals: 1,
					Owners:    []string{"@doc-team"},
				},
			},
		},
		{
			name:   "no sections, with override",
			reader: strings.NewReader("*.md @doc-team\nterms.md @legal-team"),
			path:   "/terms.md",
			want: map[string]Approval{
				"": {
					Pattern:   "terms.md",
					Approvals: 1,
					Owners:    []string{"@legal-team"},
				},
			},
		},
		{
			name:   "section with default owner",
			reader: strings.NewReader("[Database] @database-team\nmodel/db/\nconfig/db/database-setup.md @docs-team"),
			path:   "/model/db/backup.sql",
			want: map[string]Approval{
				"Database": {
					Pattern:   "model/db/",
					Approvals: 1,
					Owners:    []string{"@database-team"},
				},
			},
		},
		{
			name:   "section with default owner and override",
			reader: strings.NewReader("[Database] @database-team\nmodel/db/\nconfig/db/database-setup.md @docs-team"),
			path:   "/config/db/database-setup.md",
			want: map[string]Approval{
				"Database": {
					Pattern:   "config/db/database-setup.md",
					Approvals: 1,
					Owners:    []string{"@docs-team"},
				},
			},
		},
		{
			name:   "section with default owner and approval count",
			reader: strings.NewReader("[Documentation][2] @docs-team\ndocs/\nREADME.md"),
			path:   "/README.md",
			want: map[string]Approval{
				"Documentation": {
					Pattern:   "README.md",
					Approvals: 2,
					Owners:    []string{"@docs-team"},
				},
			},
		},
		{
			name:   "optional section with default owner",
			reader: strings.NewReader("^[Database] @database-team\nmodel/db/"),
			path:   "/model/db/backup.sql",
			want: map[string]Approval{
				"Database": {
					Pattern:   "model/db/",
					Approvals: 0,
					Owners:    []string{"@database-team"},
				},
			},
		},
		{
			name:   "ignore approval count on optional section with",
			reader: strings.NewReader("^[Documentation][2]\ndocs/ @docs-team"),
			path:   "/docs/intro.md",
			want: map[string]Approval{
				"Documentation": {
					Pattern:   "docs/",
					Approvals: 0,
					Owners:    []string{"@docs-team"},
				},
			},
		},
		{
			name: "multiple section with overrides",
			reader: strings.NewReader(`
# Required for all files
* @general-approvers

[Documentation] @docs-team
docs/
README.md
*.txt

[Database] @database-team
model/db/
config/db/database-setup.md @docs-team
        `),
			path: "/model/db/CHANGELOG.txt",
			want: map[string]Approval{
				"": {
					Approvals: 1,
					Pattern:   "*",
					Owners:    []string{"@general-approvers"},
				},
				"Documentation": {
					Pattern:   "*.txt",
					Approvals: 1,
					Owners:    []string{"@docs-team"},
				},
				"Database": {
					Pattern:   "model/db/",
					Approvals: 1,
					Owners:    []string{"@database-team"},
				},
			},
		},
		{
			name:   "rule with zero owner and no default owners for section",
			reader: strings.NewReader("[Documentation]\nREADME.md"),
			path:   "/README.md",
			want:   map[string]Approval{},
		},
		{
			name:   "rule with zero owner and no default owners for section",
			reader: strings.NewReader("[Documentation]\n*.md @username\nREADME.md"),
			path:   "/README.md",
			want: map[string]Approval{
				"Documentation": {
					Pattern:   "*.md",
					Approvals: 1,
					Owners:    []string{"@username"},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			file, err := NewCodeOwnersFile(tt.reader)
			if err != nil {
				t.Errorf("Failed to create code owners file: %v", err)
			}

			got := file.GetRequiredApprovalsForFile(tt.path)
			testhelper.DeepEqual(t, got, tt.want)
		})
	}
}

func TestFile_GetRequiredApprovalsForFiles(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		reader io.Reader
		paths  []string
		want   map[string][]Approval
	}{
		{
			name: "files matching multiple sections and rules",
			reader: strings.NewReader(`
# Required for all files
* @general-approvers

[Documentation] @docs-team
docs/
README.md
*.txt

[Database] @database-team
model/db/
config/db/database-setup.md @docs-team
        `),
			paths: []string{
				"/README.md",
				"/model/db/CHANGELOG.txt",
				"/integration/run-integration-tests.sh",
			},
			want: map[string][]Approval{
				"": {
					{
						Approvals: 1,
						Pattern:   "*",
						Owners:    []string{"@general-approvers"},
					},
				},
				"Documentation": {
					{
						Pattern:   "README.md",
						Approvals: 1,
						Owners:    []string{"@docs-team"},
					},
					{
						Pattern:   "*.txt",
						Approvals: 1,
						Owners:    []string{"@docs-team"},
					},
				},
				"Database": {
					{
						Pattern:   "model/db/",
						Approvals: 1,
						Owners:    []string{"@database-team"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			file, err := NewCodeOwnersFile(tt.reader)
			if err != nil {
				t.Errorf("Failed to create code owners file: %v", err)
			}

			got := file.GetRequiredApprovalsForFiles(tt.paths)
			testhelper.DeepEqual(t, got, tt.want)
		})
	}
}

func TestFile_isValidRule(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		rule          rule
		defaultOwners []string
		want          bool
	}{
		{
			name:          "no owners",
			rule:          rule{pattern: newPattern(""), owners: []string{}},
			defaultOwners: []string{},
			want:          false,
		},
		{
			name:          "only default owners",
			rule:          rule{pattern: newPattern(""), owners: []string{}},
			defaultOwners: []string{"@foo"},
			want:          true,
		},
		{
			name:          "only rule owners",
			rule:          rule{pattern: newPattern(""), owners: []string{"@bar"}},
			defaultOwners: []string{},
			want:          true,
		},
		{
			name:          "rule and default owners",
			rule:          rule{pattern: newPattern(""), owners: []string{"@bar"}},
			defaultOwners: []string{"@foo"},
			want:          true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := isValidRule(tt.rule, tt.defaultOwners)

			if got != tt.want {
				t.Errorf("got %t, wanted %t", got, tt.want)
			}
		})
	}
}

func TestFile_removeDuplicatedApprovals(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		approvals []Approval
		want      []Approval
	}{
		{
			name:      "empty list",
			approvals: []Approval{},
			want:      []Approval{},
		},
		{
			name: "no duplicated approvals",
			approvals: []Approval{
				{Pattern: "*.md", Approvals: 1, Owners: []string{"@foo"}},
				{Pattern: "README.md", Approvals: 1, Owners: []string{"@bar"}},
			},
			want: []Approval{
				{Pattern: "*.md", Approvals: 1, Owners: []string{"@foo"}},
				{Pattern: "README.md", Approvals: 1, Owners: []string{"@bar"}},
			},
		},
		{
			name: "duplicated approvals",
			approvals: []Approval{
				{Pattern: "*.md", Approvals: 1, Owners: []string{"@foo"}},
				{Pattern: "*.md", Approvals: 1, Owners: []string{"@foo"}},
			},
			want: []Approval{
				{Pattern: "*.md", Approvals: 1, Owners: []string{"@foo"}},
			},
		},
		{
			name: "duplicated approvals and others",
			approvals: []Approval{
				{Pattern: "*.md", Approvals: 1, Owners: []string{"@foo"}},
				{Pattern: "README.md", Approvals: 1, Owners: []string{"@bar"}},
				{Pattern: "*.md", Approvals: 1, Owners: []string{"@foo"}},
			},
			want: []Approval{
				{Pattern: "*.md", Approvals: 1, Owners: []string{"@foo"}},
				{Pattern: "README.md", Approvals: 1, Owners: []string{"@bar"}},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := removeDuplicatedApprovals(tt.approvals)
			testhelper.DeepEqual(t, got, tt.want)
		})
	}
}

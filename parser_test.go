package gitlabcodeowners

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/chefe/gitlabcodeowners/testhelper"
)

var errSomethingFailed = errors.New("something failed")

type errorReader struct{}

func (errorReader) Read(_ []byte) (int, error) {
	return 0, errSomethingFailed
}

func TestParser_parseFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		reader  io.Reader
		want    []section
		wantErr bool
	}{
		{
			name:    "empty file",
			reader:  strings.NewReader(""),
			want:    []section{},
			wantErr: false,
		},
		{
			name:   "without sections",
			reader: strings.NewReader("# A comment\n*.md @doc-team\n\nterms.md @legal-team"),
			want: []section{
				{
					name:      "",
					approvals: 1,
					owners:    []string{},
					rules: []rule{
						{
							pattern: pattern{value: "*.md", normalized: "/**/*.md"},
							owners:  []string{"@doc-team"},
						},
						{
							pattern: pattern{value: "terms.md", normalized: "/**/terms.md"},
							owners:  []string{"@legal-team"},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:   "with sections",
			reader: strings.NewReader("[README Owners]\nREADME.md @user1 @user2\ninternal/README.md @user4\n\n[README other owners]\nREADME.md @user3"),
			want: []section{
				{
					name:      "README Owners",
					approvals: 1,
					owners:    []string{},
					rules: []rule{
						{
							pattern: pattern{value: "README.md", normalized: "/**/README.md"},
							owners:  []string{"@user1", "@user2"},
						},
						{
							pattern: pattern{value: "internal/README.md", normalized: "/**/internal/README.md"},
							owners:  []string{"@user4"},
						},
					},
				},
				{
					name:      "README other owners",
					approvals: 1,
					owners:    []string{},
					rules: []rule{
						{
							pattern: pattern{value: "README.md", normalized: "/**/README.md"},
							owners:  []string{"@user3"},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:   "with optional sections, approval count and default owners",
			reader: strings.NewReader("[Documentation][2] @docs-team\ndocs/\nREADME.md\n\n^[Database] @database-team\nmodel/db/\nconfig/db/database-setup.md @docs-team"),
			want: []section{
				{
					name:      "Documentation",
					approvals: 2,
					owners:    []string{"@docs-team"},
					rules: []rule{
						{
							pattern: pattern{value: "docs/", normalized: "/**/docs/**/*"},
							owners:  []string{},
						},
						{
							pattern: pattern{value: "README.md", normalized: "/**/README.md"},
							owners:  []string{},
						},
					},
				},
				{
					name:      "Database",
					approvals: 0,
					owners:    []string{"@database-team"},
					rules: []rule{
						{
							pattern: pattern{value: "model/db/", normalized: "/**/model/db/**/*"},
							owners:  []string{},
						},
						{
							pattern: pattern{value: "config/db/database-setup.md", normalized: "/**/config/db/database-setup.md"},
							owners:  []string{"@docs-team"},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:   "duplicate section names",
			reader: strings.NewReader("[Documentation]\nee/docs/ @docs\ndocs/ @docs\n\n[Database]\nREADME.md @database\nmodel/db/ @database\n\n[DOCUMENTATION]\nREADME.md  @docs"),
			want: []section{
				{
					name:      "Documentation",
					approvals: 1,
					owners:    []string{},
					rules: []rule{
						{
							pattern: pattern{value: "ee/docs/", normalized: "/**/ee/docs/**/*"},
							owners:  []string{"@docs"},
						},
						{
							pattern: pattern{value: "docs/", normalized: "/**/docs/**/*"},
							owners:  []string{"@docs"},
						},
						{
							pattern: pattern{value: "README.md", normalized: "/**/README.md"},
							owners:  []string{"@docs"},
						},
					},
				},
				{
					name:      "Database",
					approvals: 1,
					owners:    []string{},
					rules: []rule{
						{
							pattern: pattern{value: "README.md", normalized: "/**/README.md"},
							owners:  []string{"@database"},
						},
						{
							pattern: pattern{value: "model/db/", normalized: "/**/model/db/**/*"},
							owners:  []string{"@database"},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			// https://docs.gitlab.com/ee/user/project/codeowners/reference.html#unparsable-sections
			name:   "unparsable section example",
			reader: strings.NewReader("* @group\n\n[Section name\ndocs/ @docs_group"),
			want: []section{
				{
					name:      "",
					approvals: 1,
					owners:    []string{},
					rules: []rule{
						{
							pattern: pattern{value: "*", normalized: "/**/*"},
							owners:  []string{"@group"},
						},
						{
							pattern: pattern{value: "[Section", normalized: "/**/[Section"},
							owners:  []string{"name"},
						},
						{
							pattern: pattern{value: "docs/", normalized: "/**/docs/**/*"},
							owners:  []string{"@docs_group"},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "error while reading",
			reader:  errorReader{},
			want:    []section{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := parseFile(tt.reader)

			testhelper.DeepEqual(t, got, tt.want)

			if (err != nil) != tt.wantErr {
				t.Errorf("wantErr=%t but got error %v", tt.wantErr, err)
			}
		})
	}
}

func TestParser_appendSection(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		sections []section
		section  section
		want     []section
	}{
		{
			name:     "add valid section to empty list",
			sections: []section{},
			section: section{
				name:      "Section",
				approvals: 2,
				owners:    []string{},
				rules: []rule{
					{
						pattern: newPattern("/foo"),
						owners:  []string{"@bar"},
					},
				},
			},
			want: []section{
				{
					name:      "Section",
					approvals: 2,
					owners:    []string{},
					rules: []rule{
						{
							pattern: newPattern("/foo"),
							owners:  []string{"@bar"},
						},
					},
				},
			},
		},
		{
			name:     "ignore section with empty rule set",
			sections: []section{},
			section: section{
				name:      "Empty",
				approvals: 0,
				owners:    []string{},
				rules:     []rule{},
			},
			want: []section{},
		},
		{
			name: "merge sections with the same name",
			sections: []section{
				{
					name:      "Documentation",
					approvals: 2,
					owners:    []string{},
					rules: []rule{
						{
							pattern: newPattern("/foo"),
							owners:  []string{"@foo"},
						},
					},
				},
			},
			section: section{
				name:      "DOCUMENTATION",
				approvals: 0,
				owners:    []string{},
				rules: []rule{
					{
						pattern: newPattern("/bar"),
						owners:  []string{"@bar"},
					},
				},
			},
			want: []section{
				{
					name:      "Documentation",
					approvals: 2,
					owners:    []string{},
					rules: []rule{
						{
							pattern: newPattern("/foo"),
							owners:  []string{"@foo"},
						},
						{
							pattern: newPattern("/bar"),
							owners:  []string{"@bar"},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := appendSection(tt.sections, tt.section)
			testhelper.DeepEqual(t, got, tt.want)
		})
	}
}

package gitlabcodeowners

import (
	"testing"

	"github.com/chefe/gitlabcodeowners/testhelper"
)

func TestRule_parseRule(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		rule string
		want rule
	}{
		{
			name: "single owner",
			rule: "/*.md @username",
			want: rule{
				pattern: pattern{value: "/*.md", normalized: "/*.md"},
				owners:  []string{"@username"},
			},
		},
		{
			name: "multiple owners",
			rule: "/path/to/entry.txt @group @group/subgroup @user",
			want: rule{
				pattern: pattern{value: "/path/to/entry.txt", normalized: "/path/to/entry.txt"},
				owners:  []string{"@group", "@group/subgroup", "@user"},
			},
		},
		{
			name: "multiple owners with tabs",
			rule: "/path/to/entry.txt\t@username\tjanedoe@gitlab.com",
			want: rule{
				pattern: pattern{value: "/path/to/entry.txt", normalized: "/path/to/entry.txt"},
				owners:  []string{"@username", "janedoe@gitlab.com"},
			},
		},
		{
			name: "entries with spaces",
			rule: "folder with spaces/*.md @group",
			want: rule{
				pattern: pattern{value: "folder", normalized: "/**/folder"},
				owners:  []string{"with", "spaces/*.md", "@group"},
			},
		},
		{
			name: "no owner",
			rule: "/file.md",
			want: rule{
				pattern: pattern{value: "/file.md", normalized: "/file.md"},
				owners:  []string{},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := parseRule(tt.rule)
			testhelper.DeepEqual(t, got, tt.want)
		})
	}
}

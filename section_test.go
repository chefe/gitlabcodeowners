package gitlabcodeowners

import (
	"testing"

	"github.com/chefe/gitlabcodeowners/testhelper"
)

func TestSection_parseSectionHeader(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		header  string
		want    section
		wantErr bool
	}{
		{
			name:   "required section with no approval count and no default owners",
			header: "[Section name]",
			want: section{
				name:      "Section name",
				approvals: 1,
				owners:    []string{},
				rules:     []rule{},
			},
			wantErr: false,
		},
		{
			name:   "optional section with no approval count and no default owners",
			header: "^[Section name]",
			want: section{
				name:      "Section name",
				approvals: 0,
				owners:    []string{},
				rules:     []rule{},
			},
			wantErr: false,
		},
		{
			name:   "required section with approval count and no default owners",
			header: "[Section name][5]",
			want: section{
				name:      "Section name",
				approvals: 5,
				owners:    []string{},
				rules:     []rule{},
			},
			wantErr: false,
		},
		{
			name:   "optional section with approval count and no default owners",
			header: "^[Section name][5]",
			want: section{
				name:      "Section name",
				approvals: 0,
				owners:    []string{},
				rules:     []rule{},
			},
			wantErr: false,
		},
		{
			name:   "required section with no approval count and one default owner",
			header: "[Section name] @username",
			want: section{
				name:      "Section name",
				approvals: 1,
				owners:    []string{"@username"},
				rules:     []rule{},
			},
			wantErr: false,
		},
		{
			name:   "optional section with no approval count and one default owner",
			header: "^[Section name] @username",
			want: section{
				name:      "Section name",
				approvals: 0,
				owners:    []string{"@username"},
				rules:     []rule{},
			},
			wantErr: false,
		},
		{
			name:   "required section with approval count and multiple default owners",
			header: "[Docs][2] @group @subgroup",
			want: section{
				name:      "Docs",
				approvals: 2,
				owners:    []string{"@group", "@subgroup"},
				rules:     []rule{},
			},
			wantErr: false,
		},
		{
			name:   "optional section with approval count and multiple default owners",
			header: "^[Docs][2] @group @subgroup",
			want: section{
				name:      "Docs",
				approvals: 0,
				owners:    []string{"@group", "@subgroup"},
				rules:     []rule{},
			},
			wantErr: false,
		},
		{
			name:   "required section with zero as approval count",
			header: "[Testing][0]",
			want: section{
				name:      "Testing",
				approvals: 1,
				owners:    []string{},
				rules:     []rule{},
			},
			wantErr: false,
		},
		{
			name:   "required section with negative approval count",
			header: "[Testing][-42]",
			want: section{
				name:      "Testing",
				approvals: 1,
				owners:    []string{},
				rules:     []rule{},
			},
			wantErr: false,
		},
		{
			name:   "required section where approval count is not a number",
			header: "[Legal][abc]",
			want: section{
				name:      "Legal",
				approvals: 1,
				owners:    []string{},
				rules:     []rule{},
			},
			wantErr: false,
		},
		{
			name:    "missing square closing bracket for section",
			header:  "[Section name",
			want:    section{}, //nolint:exhaustruct // default is returned on error
			wantErr: true,
		},
		{
			name:    "missing square closing bracket for approval count",
			header:  "[Section name][1 @username",
			want:    section{}, //nolint:exhaustruct // default is returned on error
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := parseSectionHeader(tt.header)

			testhelper.DeepEqual(t, got, tt.want)

			if (err != nil) != tt.wantErr {
				t.Errorf("wantErr=%t but got error %v", tt.wantErr, err)
			}
		})
	}
}

func TestSection_checkBracketCountInSectionHeader(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		header  string
		wantErr bool
	}{
		{
			name:    "33 - too much brackets",
			header:  "[One][Two][Three]",
			wantErr: true,
		},
		{
			name:    "22 - name with approval count",
			header:  "[Documentation][4] @docs-team",
			wantErr: false,
		},
		{
			name:    "11 - only name without approval count",
			header:  "[Documentation] @docs-team",
			wantErr: false,
		},
		{
			name:    "21 - missing closing brackets for approval count",
			header:  "[Documentation][4 @docs-team",
			wantErr: true,
		},
		{
			name:    "10 - missing closing brackets for name",
			header:  "[Documentation @docs-team",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := checkBracketCountInSectionHeader(tt.header)

			if (err != nil) != tt.wantErr {
				t.Errorf("wantErr=%t but got error %v", tt.wantErr, err)
			}
		})
	}
}

func TestSection_extractPartsFromSectionHeader(t *testing.T) {
	t.Parallel()

	type parts struct {
		name      string
		approvals string
		owners    string
	}

	tests := []struct {
		name   string
		header string
		want   parts
	}{
		{
			name:   "name with approval count and default owners",
			header: "[Documentation][4] @docs-team @specialuser",
			want: parts{
				name:      "Documentation",
				approvals: "4",
				owners:    " @docs-team @specialuser",
			},
		},
		{
			name:   "name with approval count but no default owners",
			header: "[Documentation][4]",
			want: parts{
				name:      "Documentation",
				approvals: "4",
				owners:    "",
			},
		},
		{
			name:   "name with default owners but no approval count",
			header: "[Documentation] @docs-team @specialuser",
			want: parts{
				name:      "Documentation",
				approvals: "",
				owners:    " @docs-team @specialuser",
			},
		},
		{
			name:   "only name",
			header: "[Documentation]",
			want: parts{
				name:      "Documentation",
				approvals: "",
				owners:    "",
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			name, approvals, owners := extractPartsFromSectionHeader(tt.header)
			got := parts{name: name, approvals: approvals, owners: owners}
			testhelper.DeepEqual(t, got, tt.want)
		})
	}
}

func TestSection_parseApprovalCount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		count    string
		optional bool
		want     int
	}{
		{
			name:     "required - valid number",
			count:    "5",
			optional: false,
			want:     5,
		},
		{
			name:     "required - number with whitespace",
			count:    " 42\t",
			optional: false,
			want:     42,
		},
		{
			name:     "required - zero",
			count:    "0",
			optional: false,
			want:     1,
		},
		{
			name:     "required - negative number",
			count:    "-10",
			optional: false,
			want:     1,
		},
		{
			name:     "required - empty string",
			count:    "",
			optional: false,
			want:     1,
		},
		{
			name:     "required - not a number",
			count:    "abc",
			optional: false,
			want:     1,
		},
		{
			name:     "optional - valid number",
			count:    "5",
			optional: true,
			want:     0,
		},
		{
			name:     "optional - number with whitespace",
			count:    " 42\t",
			optional: true,
			want:     0,
		},
		{
			name:     "optional - zero",
			count:    "0",
			optional: true,
			want:     0,
		},
		{
			name:     "optional - negative number",
			count:    "-10",
			optional: true,
			want:     0,
		},
		{
			name:     "optional - empty string",
			count:    "",
			optional: true,
			want:     0,
		},
		{
			name:     "optional - not a number",
			count:    "abc",
			optional: true,
			want:     0,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := parseApprovalCount(tt.count, tt.optional)

			if got != tt.want {
				t.Errorf("got %d, wanted %d", got, tt.want)
			}
		})
	}
}

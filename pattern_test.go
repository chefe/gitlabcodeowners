package gitlabcodeowners

import (
	"testing"
)

func TestPattern_match(t *testing.T) {
	t.Parallel()

	type example struct {
		path string
		want bool
	}

	tests := []struct {
		name     string
		pattern  string
		examples []example
	}{
		{
			name:    "relative directory with no star",
			pattern: "docs/",
			examples: []example{
				{path: "/todo.txt", want: false},
				{path: "/README.md", want: false},
				{path: "/CONTRIBUTING.md", want: false},
				{path: "/docs/README.md", want: true},
				{path: "/docs/usage.md", want: true},
				{path: "/docs/internal/notes.txt", want: true},
				{path: "/internal/README.md", want: false},
				{path: "/internal/testing/README.md", want: false},
				{path: "/internal/testing/run.sh", want: false},
				{path: "/internal/testing/docs/content/integration/intro.md", want: true},
			},
		},
		{
			name:    "relative directory with single star",
			pattern: "docs/*/",
			examples: []example{
				{path: "/todo.txt", want: false},
				{path: "/README.md", want: false},
				{path: "/CONTRIBUTING.md", want: false},
				{path: "/docs/README.md", want: false},
				{path: "/docs/usage.md", want: false},
				{path: "/docs/internal/notes.txt", want: true},
				{path: "/internal/README.md", want: false},
				{path: "/internal/testing/README.md", want: false},
				{path: "/internal/testing/run.sh", want: false},
				{path: "/internal/testing/docs/content/integration/intro.md", want: true},
			},
		},
		{
			name:    "relative directory with two stars",
			pattern: "testing/**/integration/",
			examples: []example{
				{path: "/todo.txt", want: false},
				{path: "/README.md", want: false},
				{path: "/CONTRIBUTING.md", want: false},
				{path: "/docs/README.md", want: false},
				{path: "/docs/usage.md", want: false},
				{path: "/docs/internal/notes.txt", want: false},
				{path: "/internal/README.md", want: false},
				{path: "/internal/testing/README.md", want: false},
				{path: "/internal/testing/run.sh", want: false},
				{path: "/internal/testing/docs/content/integration/intro.md", want: true},
			},
		},
		{
			name:    "relative file with no star",
			pattern: "README.md",
			examples: []example{
				{path: "/todo.txt", want: false},
				{path: "/README.md", want: true},
				{path: "/CONTRIBUTING.md", want: false},
				{path: "/docs/README.md", want: true},
				{path: "/docs/usage.md", want: false},
				{path: "/docs/internal/notes.txt", want: false},
				{path: "/internal/README.md", want: true},
				{path: "/internal/testing/README.md", want: true},
				{path: "/internal/testing/run.sh", want: false},
				{path: "/internal/testing/docs/content/integration/intro.md", want: false},
			},
		},
		{
			name:    "relative file with single star",
			pattern: "*.md",
			examples: []example{
				{path: "/todo.txt", want: false},
				{path: "/README.md", want: true},
				{path: "/CONTRIBUTING.md", want: true},
				{path: "/docs/README.md", want: true},
				{path: "/docs/usage.md", want: true},
				{path: "/docs/internal/notes.txt", want: false},
				{path: "/internal/README.md", want: true},
				{path: "/internal/testing/README.md", want: true},
				{path: "/internal/testing/run.sh", want: false},
				{path: "/internal/testing/docs/content/integration/intro.md", want: true},
			},
		},
		{
			name:    "absolute directory with no star",
			pattern: "/docs/",
			examples: []example{
				{path: "/todo.txt", want: false},
				{path: "/README.md", want: false},
				{path: "/CONTRIBUTING.md", want: false},
				{path: "/docs/README.md", want: true},
				{path: "/docs/usage.md", want: true},
				{path: "/docs/internal/notes.txt", want: true},
				{path: "/internal/README.md", want: false},
				{path: "/internal/testing/README.md", want: false},
				{path: "/internal/testing/run.sh", want: false},
				{path: "/internal/testing/docs/content/integration/intro.md", want: false},
			},
		},
		{
			name:    "absolute directory with single star",
			pattern: "/internal/*/",
			examples: []example{
				{path: "/todo.txt", want: false},
				{path: "/README.md", want: false},
				{path: "/CONTRIBUTING.md", want: false},
				{path: "/docs/README.md", want: false},
				{path: "/docs/usage.md", want: false},
				{path: "/docs/internal/notes.txt", want: false},
				{path: "/internal/README.md", want: false},
				{path: "/internal/testing/README.md", want: true},
				{path: "/internal/testing/run.sh", want: true},
				{path: "/internal/testing/docs/content/integration/intro.md", want: true},
			},
		},
		{
			name:    "absolute directory with two stars",
			pattern: "/internal/**/integration/",
			examples: []example{
				{path: "/todo.txt", want: false},
				{path: "/README.md", want: false},
				{path: "/CONTRIBUTING.md", want: false},
				{path: "/docs/README.md", want: false},
				{path: "/docs/usage.md", want: false},
				{path: "/docs/internal/notes.txt", want: false},
				{path: "/internal/README.md", want: false},
				{path: "/internal/testing/README.md", want: false},
				{path: "/internal/testing/run.sh", want: false},
				{path: "/internal/testing/docs/content/integration/intro.md", want: true},
			},
		},
		{
			name:    "absolute file with no star",
			pattern: "/README.md",
			examples: []example{
				{path: "/todo.txt", want: false},
				{path: "/README.md", want: true},
				{path: "/CONTRIBUTING.md", want: false},
				{path: "/docs/README.md", want: false},
				{path: "/docs/usage.md", want: false},
				{path: "/docs/internal/notes.txt", want: false},
				{path: "/internal/README.md", want: false},
				{path: "/internal/testing/README.md", want: false},
				{path: "/internal/testing/run.sh", want: false},
				{path: "/internal/testing/docs/content/integration/intro.md", want: false},
			},
		},
		{
			name:    "absolute file with single star",
			pattern: "/docs/*",
			examples: []example{
				{path: "/todo.txt", want: false},
				{path: "/README.md", want: false},
				{path: "/CONTRIBUTING.md", want: false},
				{path: "/docs/README.md", want: true},
				{path: "/docs/usage.md", want: true},
				{path: "/docs/internal/notes.txt", want: false},
				{path: "/internal/README.md", want: false},
				{path: "/internal/testing/README.md", want: false},
				{path: "/internal/testing/run.sh", want: false},
				{path: "/internal/testing/docs/content/integration/intro.md", want: false},
			},
		},
		{
			name:    "absolute file with single star",
			pattern: "/internal/**/*.md",
			examples: []example{
				{path: "/todo.txt", want: false},
				{path: "/README.md", want: false},
				{path: "/CONTRIBUTING.md", want: false},
				{path: "/docs/README.md", want: false},
				{path: "/docs/usage.md", want: false},
				{path: "/docs/internal/notes.txt", want: false},
				{path: "/internal/README.md", want: true},
				{path: "/internal/testing/README.md", want: true},
				{path: "/internal/testing/run.sh", want: false},
				{path: "/internal/testing/docs/content/integration/intro.md", want: true},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			for _, ex := range tt.examples {
				got := newPattern(tt.pattern).match(ex.path)

				if got != ex.want {
					t.Errorf("path %s -> got %t, wanted %t", ex.path, got, ex.want)
				}
			}
		})
	}
}

func TestPattern_normalizePattern(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		pattern string
		want    string
	}{
		{
			name:    "catch all",
			pattern: "*",
			want:    "/**/*",
		},
		{
			name:    "relative file pattern",
			pattern: "*.md",
			want:    "/**/*.md",
		},
		{
			name:    "absolute file pattern",
			pattern: "/scripts/*.sh",
			want:    "/scripts/*.sh",
		},
		{
			name:    "relative directory pattern",
			pattern: "build/",
			want:    "/**/build/**/*",
		},
		{
			name:    "absolute directory pattern",
			pattern: "/tmp/",
			want:    "/tmp/**/*",
		},
		{
			name:    "special case files starting with a pound",
			pattern: "\\#file\\#with\\#pound.txt",
			want:    "/**/#file\\#with\\#pound.txt",
		},
		{
			name:    "special case whitespaces",
			pattern: "file\\ with\\ spaces.txt",
			want:    "/**/file with spaces.txt",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := normalizePattern(tt.pattern)

			if got != tt.want {
				t.Errorf("got %s, wanted %s", got, tt.want)
			}
		})
	}
}

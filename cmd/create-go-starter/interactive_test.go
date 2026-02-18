package main

import (
	"bufio"
	"os"
	"strings"
	"testing"
)

// emptyDefaults returns an InteractiveDefaults with zero values (all defaults).
func emptyDefaults() InteractiveDefaults {
	return InteractiveDefaults{}
}

// TestPromptStringFromReader tests the promptStringFromReader helper
func TestPromptStringFromReader(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		defaultVal string
		want       string
		wantErr    bool
	}{
		{
			name:       "user provides a value",
			input:      "my-project\n",
			defaultVal: "default-project",
			want:       "my-project",
			wantErr:    false,
		},
		{
			name:       "user presses enter to accept default",
			input:      "\n",
			defaultVal: "default-project",
			want:       "default-project",
			wantErr:    false,
		},
		{
			name:       "user provides value with spaces trimmed",
			input:      "  my-app  \n",
			defaultVal: "",
			want:       "my-app",
			wantErr:    false,
		},
		{
			name:       "EOF returns error (Ctrl+C simulation)",
			input:      "",
			defaultVal: "default",
			want:       "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.input))
			got, err := promptStringFromReader(reader, "Test prompt", tt.defaultVal)
			if (err != nil) != tt.wantErr {
				t.Errorf("promptStringFromReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("promptStringFromReader() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestPromptSelectFromReader tests the promptSelectFromReader helper
func TestPromptSelectFromReader(t *testing.T) {
	options := []string{"full", "minimal", "graphql"}
	descriptions := []string{"Full API with JWT", "Basic REST API", "GraphQL API"}

	tests := []struct {
		name       string
		input      string
		defaultIdx int
		want       string
		wantErr    bool
	}{
		{
			name:       "user selects option 1",
			input:      "1\n",
			defaultIdx: 0,
			want:       "full",
			wantErr:    false,
		},
		{
			name:       "user selects option 2",
			input:      "2\n",
			defaultIdx: 0,
			want:       "minimal",
			wantErr:    false,
		},
		{
			name:       "user selects option 3",
			input:      "3\n",
			defaultIdx: 0,
			want:       "graphql",
			wantErr:    false,
		},
		{
			name:       "user presses enter to accept default (index 0)",
			input:      "\n",
			defaultIdx: 0,
			want:       "full",
			wantErr:    false,
		},
		{
			name:       "user presses enter to accept default (index 1)",
			input:      "\n",
			defaultIdx: 1,
			want:       "minimal",
			wantErr:    false,
		},
		{
			name:       "invalid then valid selection",
			input:      "99\n1\n",
			defaultIdx: 0,
			want:       "full",
			wantErr:    false,
		},
		{
			name:       "zero (invalid) then valid selection",
			input:      "0\n2\n",
			defaultIdx: 0,
			want:       "minimal",
			wantErr:    false,
		},
		{
			name:       "EOF returns error (Ctrl+C simulation)",
			input:      "",
			defaultIdx: 0,
			want:       "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.input))
			got, err := promptSelectFromReader(reader, "Select template", options, descriptions, tt.defaultIdx)
			if (err != nil) != tt.wantErr {
				t.Errorf("promptSelectFromReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("promptSelectFromReader() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestRunInteractiveModeCancel tests cancellation at confirmation step
func TestRunInteractiveModeCancel(t *testing.T) {
	// Input: project name, template (1=full), database (1=postgres), observability (1=none), confirm NO
	input := "my-project\n1\n1\n1\nn\n"
	err := runInteractiveModeWithReader(strings.NewReader(input), emptyDefaults())
	if err == nil {
		t.Error("Expected error when user cancels at confirmation, got nil")
	}
}

// TestRunInteractiveModeEOFAtProjectName tests EOF (Ctrl+C) handling at project name prompt
func TestRunInteractiveModeEOFAtProjectName(t *testing.T) {
	// No input = EOF immediately at project name prompt
	err := runInteractiveModeWithReader(strings.NewReader(""), emptyDefaults())
	if err == nil {
		t.Error("Expected error on EOF, got nil")
	}
}

// TestRunInteractiveModeSuccess tests the complete happy path using a temp directory.
// This test changes the working directory, so it's skipped in short mode to avoid
// interfering with other tests that rely on the process working directory.
func TestRunInteractiveModeSuccess(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode (uses os.Chdir)")
	}

	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Cannot get working dir: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Cannot change to temp dir: %v", err)
	}
	defer os.Chdir(origDir) //nolint:errcheck

	// Input: project name, template (2=minimal), database (3=sqlite), observability (1=none), confirm YES
	// We use minimal+sqlite to avoid Docker requirements in tests
	input := "test-interactive-project\n2\n3\n1\ny\n"
	err = runInteractiveModeWithReader(strings.NewReader(input), emptyDefaults())
	if err != nil {
		t.Errorf("runInteractiveModeWithReader() unexpected error = %v", err)
	}
}

// TestRunInteractiveModeInvalidProjectName tests re-prompting on invalid project name
func TestRunInteractiveModeInvalidProjectName(t *testing.T) {
	// Invalid project name "@invalid!" (contains invalid characters) then valid name, then cancel
	input := "@invalid!\nmy-project\n1\n1\n1\nn\n"
	err := runInteractiveModeWithReader(strings.NewReader(input), emptyDefaults())
	// Expecting error because we cancel at confirmation
	if err == nil {
		t.Error("Expected cancellation error, got nil")
	}
}

// TestRunInteractiveModeInvalidSelection tests validation re-prompting
func TestRunInteractiveModeInvalidSelection(t *testing.T) {
	// Invalid template selection "99", then valid "1" (full), database "1", observability "1", cancel
	input := "my-project\n99\n1\n1\n1\nn\n"
	err := runInteractiveModeWithReader(strings.NewReader(input), emptyDefaults())
	// Expecting cancellation error (not invalid selection error)
	if err == nil {
		t.Error("Expected cancellation error, got nil")
	}
}

// TestRunInteractiveModeWithDefaults tests that CLI flags are pre-populated as defaults
func TestRunInteractiveModeWithDefaults(t *testing.T) {
	// Press enter for all prompts (accept defaults) then cancel
	// With defaults: template=minimal (index 1), database=sqlite (index 2)
	input := "my-project\n\n\n\nn\n"
	defaults := InteractiveDefaults{
		Template: TemplateMinimal,
		Database: DatabaseSQLite,
	}
	err := runInteractiveModeWithReader(strings.NewReader(input), defaults)
	// Expecting cancellation error
	if err == nil {
		t.Error("Expected cancellation error, got nil")
	}
}

// TestRunInteractiveModeObservabilityValidation tests that advanced observability
// with non-full template is rejected
func TestRunInteractiveModeObservabilityValidation(t *testing.T) {
	// Select project, template=minimal (2), database=postgres (1), observability=advanced (3)
	input := "my-project\n2\n1\n3\n"
	err := runInteractiveModeWithReader(strings.NewReader(input), emptyDefaults())
	if err == nil {
		t.Error("Expected error for invalid observability+template combination, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "observability") {
		t.Errorf("Expected observability validation error, got: %v", err)
	}
}

// TestPromptConfirmFromReader tests the promptConfirmFromReader helper
func TestPromptConfirmFromReader(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		defaultYes bool
		want       bool
		wantErr    bool
	}{
		{
			name:       "explicit yes",
			input:      "y\n",
			defaultYes: true,
			want:       true,
			wantErr:    false,
		},
		{
			name:       "explicit YES",
			input:      "YES\n",
			defaultYes: false,
			want:       true,
			wantErr:    false,
		},
		{
			name:       "explicit no",
			input:      "n\n",
			defaultYes: true,
			want:       false,
			wantErr:    false,
		},
		{
			name:       "empty with defaultYes=true",
			input:      "\n",
			defaultYes: true,
			want:       true,
			wantErr:    false,
		},
		{
			name:       "empty with defaultYes=false",
			input:      "\n",
			defaultYes: false,
			want:       false,
			wantErr:    false,
		},
		{
			name:       "EOF returns error",
			input:      "",
			defaultYes: true,
			want:       false,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.input))
			got, err := promptConfirmFromReader(reader, "Confirm?", tt.defaultYes)
			if (err != nil) != tt.wantErr {
				t.Errorf("promptConfirmFromReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("promptConfirmFromReader() = %v, want %v", got, tt.want)
			}
		})
	}
}

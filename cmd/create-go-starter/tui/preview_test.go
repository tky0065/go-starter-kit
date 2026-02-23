package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewPreviewModel(t *testing.T) {
	files := []PreviewFile{
		{Path: "cmd/main.go", Content: "package main", Size: 100},
		{Path: "internal/domain/user.go", Content: "package user", Size: 200},
		{Path: "go.mod", Content: "module test", Size: 50},
	}

	m := NewPreviewModel("test-project", files, 80, 24)

	if m.projectName != "test-project" {
		t.Errorf("Expected project name 'test-project', got '%s'", m.projectName)
	}

	if len(m.files) != 3 {
		t.Errorf("Expected 3 files, got %d", len(m.files))
	}

	expectedSize := int64(350)
	if m.totalSize != expectedSize {
		t.Errorf("Expected total size %d, got %d", expectedSize, m.totalSize)
	}

	if !m.ready {
		t.Error("Expected model to be ready")
	}
}

func TestPreviewModel_Update_Navigation(t *testing.T) {
	files := []PreviewFile{
		{Path: "file1.go", Content: "content1", Size: 100},
		{Path: "file2.go", Content: "content2", Size: 200},
	}

	m := NewPreviewModel("test-project", files, 80, 24)
	initialY := m.viewport.YOffset

	// Test down navigation
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	if m.viewport.YOffset < initialY {
		t.Error("Expected viewport to scroll down")
	}

	// Test up navigation
	m.viewport.YOffset = 10 // Set some scroll position
	prevY := m.viewport.YOffset
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	if m.viewport.YOffset >= prevY {
		t.Error("Expected viewport to scroll up")
	}
}

func TestPreviewModel_Update_PageNavigation(t *testing.T) {
	files := []PreviewFile{
		{Path: "file1.go", Content: "content1", Size: 100},
	}

	m := NewPreviewModel("test-project", files, 80, 24)

	// Test page down
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyPgDown})
	// Viewport should handle page down

	// Test page up
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyPgUp})
	// Viewport should handle page up
}

func TestPreviewModel_Update_WindowResize(t *testing.T) {
	files := []PreviewFile{
		{Path: "file1.go", Content: "content1", Size: 100},
	}

	m := NewPreviewModel("test-project", files, 80, 24)

	// Resize window
	m, _ = m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	if m.width != 120 {
		t.Errorf("Expected width 120, got %d", m.width)
	}

	if m.height != 40 {
		t.Errorf("Expected height 40, got %d", m.height)
	}

	expectedViewportWidth := 120 - 4
	if m.viewport.Width != expectedViewportWidth {
		t.Errorf("Expected viewport width %d, got %d", expectedViewportWidth, m.viewport.Width)
	}

	expectedViewportHeight := 40 - 10
	if m.viewport.Height != expectedViewportHeight {
		t.Errorf("Expected viewport height %d, got %d", expectedViewportHeight, m.viewport.Height)
	}
}

func TestPreviewModel_View(t *testing.T) {
	files := []PreviewFile{
		{Path: "cmd/main.go", Content: "package main\n\nfunc main() {}", Size: 100},
		{Path: "go.mod", Content: "module test", Size: 50},
	}

	m := NewPreviewModel("test-project", files, 80, 24)
	view := m.View()

	// Check that view contains expected elements
	if !strings.Contains(view, "test-project") {
		t.Error("View should contain project name")
	}

	if !strings.Contains(view, "2 files") {
		t.Error("View should show file count")
	}

	if !strings.Contains(view, "KB") {
		t.Error("View should show total size in KB")
	}

	// Check for navigation hints
	if !strings.Contains(view, "↑↓") || !strings.Contains(view, "Navigate") {
		t.Error("View should contain navigation hints")
	}

	if !strings.Contains(view, "Esc: Back") {
		t.Error("View should contain back hint")
	}
}

func TestPreviewModel_View_NotReady(t *testing.T) {
	files := []PreviewFile{}
	m := NewPreviewModel("test", files, 80, 24)
	m.ready = false

	view := m.View()
	if !strings.Contains(view, "Initializing") {
		t.Error("View should show initializing message when not ready")
	}
}

func TestBuildFileTree(t *testing.T) {
	files := []PreviewFile{
		{Path: "cmd/main.go", Content: "", Size: 100},
		{Path: "internal/domain/user.go", Content: "", Size: 200},
		{Path: "internal/adapters/handler.go", Content: "", Size: 150},
		{Path: "go.mod", Content: "", Size: 50},
		{Path: "README.md", Content: "", Size: 80},
	}

	m := NewPreviewModel("test-project", files, 80, 24)
	tree := m.buildFileTree()

	// Check that tree contains expected directories
	if !strings.Contains(tree, "cmd") {
		t.Error("Tree should contain 'cmd' directory")
	}

	if !strings.Contains(tree, "internal") {
		t.Error("Tree should contain 'internal' directory")
	}

	// Check that tree contains files
	if !strings.Contains(tree, "main.go") {
		t.Error("Tree should contain 'main.go'")
	}

	if !strings.Contains(tree, "go.mod") {
		t.Error("Tree should contain 'go.mod'")
	}

	// Check for diff-style markers
	if !strings.Contains(tree, "+") {
		t.Error("Tree should contain diff-style '+' markers")
	}

	// Check for file size information
	if !strings.Contains(tree, "B") && !strings.Contains(tree, "KB") {
		t.Error("Tree should contain file size information")
	}
}

func TestHighlightGoCode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string // Substrings that should be present
	}{
		{
			name:     "package keyword",
			input:    "package main",
			expected: []string{"package", "main"},
		},
		{
			name:     "function definition",
			input:    "func main() {",
			expected: []string{"func", "main"},
		},
		{
			name:     "string literal",
			input:    `fmt.Println("Hello, World!")`,
			expected: []string{"Hello", "World"},
		},
		{
			name:     "comment",
			input:    "// This is a comment",
			expected: []string{"//", "comment"},
		},
		{
			name:     "type definition",
			input:    "type User struct {",
			expected: []string{"type", "struct"},
		},
		{
			name:     "import statement",
			input:    `import "fmt"`,
			expected: []string{"import", "fmt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := highlightGoCode(tt.input)

			// Result should not be empty
			if result == "" {
				t.Error("highlightGoCode returned empty string")
			}

			// Check that expected substrings are present (they may be styled)
			for _, expected := range tt.expected {
				if !strings.Contains(result, expected) {
					t.Errorf("Expected output to contain '%s', got '%s'", expected, result)
				}
			}
		})
	}
}

func TestShouldShowPreview(t *testing.T) {
	tests := []struct {
		filename string
		expected bool
	}{
		{"main.go", true},
		{"cmd/main.go", true},
		{"config.go", true},
		{"internal/config.go", true},
		{"server.go", true},
		{"user.go", true},
		{"service.go", true},
		{".env", true},
		{"Dockerfile", true},
		{"README.md", true},
		{"random.go", false},
		{"test.txt", false},
		{"helper.go", false},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			result := shouldShowPreview(tt.filename)
			if result != tt.expected {
				t.Errorf("shouldShowPreview(%s) = %v, expected %v", tt.filename, result, tt.expected)
			}
		})
	}
}

func TestFormatSize(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{100, "100 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1024 * 1024, "1.0 MB"},
		{1536 * 1024, "1.5 MB"},
		{2048 * 1024, "2.0 MB"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := formatSize(tt.bytes)
			if result != tt.expected {
				t.Errorf("formatSize(%d) = %s, expected %s", tt.bytes, result, tt.expected)
			}
		})
	}
}

func TestSortStrings(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "already sorted",
			input:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "reverse order",
			input:    []string{"c", "b", "a"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "mixed order",
			input:    []string{"zebra", "apple", "banana"},
			expected: []string{"apple", "banana", "zebra"},
		},
		{
			name:     "single element",
			input:    []string{"single"},
			expected: []string{"single"},
		},
		{
			name:     "empty slice",
			input:    []string{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arr := make([]string, len(tt.input))
			copy(arr, tt.input)

			sortStrings(arr)

			if len(arr) != len(tt.expected) {
				t.Errorf("Length mismatch: got %d, expected %d", len(arr), len(tt.expected))
				return
			}

			for i := range arr {
				if arr[i] != tt.expected[i] {
					t.Errorf("Element %d: got '%s', expected '%s'", i, arr[i], tt.expected[i])
				}
			}
		})
	}
}

func TestRenderCodePreview(t *testing.T) {
	m := PreviewModel{}

	code := `package main

import "fmt"

func main() {
	fmt.Println("Hello")
	fmt.Println("World")
	// More lines here
	x := 42
	y := "test"
}
// Extra line 1
// Extra line 2`

	preview := m.renderCodePreview(code)

	// Should be truncated to max lines
	lines := strings.Split(preview, "\n")
	if len(lines) > 10 {
		t.Errorf("Preview should be truncated, got %d lines", len(lines))
	}

	// Should contain ellipsis for truncated content
	if !strings.Contains(preview, "...") {
		t.Error("Preview should contain '...' for truncated content")
	}

	// Should contain some of the original code
	if !strings.Contains(preview, "package") {
		t.Error("Preview should contain code content")
	}
}

func TestRenderCodePreview_ShortContent(t *testing.T) {
	m := PreviewModel{}

	code := `package main

func main() {
	fmt.Println("Hi")
}`

	preview := m.renderCodePreview(code)

	// Should not contain ellipsis for short content
	if strings.Contains(preview, "...") {
		t.Error("Short preview should not contain '...'")
	}
}

package tui

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// PreviewModel represents the dry-run preview screen state.
// It displays a scrollable preview of files that will be created with syntax highlighting.
type PreviewModel struct {
	viewport       viewport.Model
	projectName    string
	files          []PreviewFile
	totalSize      int64
	ready          bool
	width          int
	height         int
}

// PreviewFile represents a file that will be created, with content preview.
type PreviewFile struct {
	Path    string
	Content string
	Size    int64
}

// dirEntry represents a node in the file tree (directory or file).
type dirEntry struct {
	isDir    bool
	name     string
	size     int64
	content  string
	children map[string]*dirEntry
}

// NewPreviewModel creates a new PreviewModel with the given files.
func NewPreviewModel(projectName string, files []PreviewFile, width, height int) PreviewModel {
	// Calculate total size
	var totalSize int64
	for _, f := range files {
		totalSize += f.Size
	}

	// Create viewport for scrollable content
	vp := viewport.New(width-4, height-10) // Reserve space for header/footer
	vp.YPosition = 0
	vp.HighPerformanceRendering = false

	m := PreviewModel{
		viewport:    vp,
		projectName: projectName,
		files:       files,
		totalSize:   totalSize,
		width:       width,
		height:      height,
		ready:       true,
	}

	// Set viewport content
	m.viewport.SetContent(m.renderContent())

	return m
}

// Update handles messages for the preview screen.
func (m PreviewModel) Update(msg tea.Msg) (PreviewModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Update viewport size on window resize
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width - 4
		m.viewport.Height = msg.Height - 10
		m.viewport.SetContent(m.renderContent())
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.viewport.LineUp(1)
		case "down", "j":
			m.viewport.LineDown(1)
		case "pgup", "b":
			m.viewport.ViewUp()
		case "pgdown", "f", " ":
			m.viewport.ViewDown()
		case "home", "g":
			m.viewport.GotoTop()
		case "end", "G":
			m.viewport.GotoBottom()
		}
	}

	// Update viewport
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

// View renders the preview screen.
func (m PreviewModel) View() string {
	if !m.ready {
		return "\nInitializing preview..."
	}

	var b strings.Builder

	// Header
	b.WriteString("\n")
	title := fmt.Sprintf("Dry-Run Preview: %s", m.projectName)
	b.WriteString(RenderTitleBox(title, m.width-4))
	b.WriteString("\n\n")

	// Summary info
	summary := fmt.Sprintf("%s %d files will be created",
		RenderInfo("📁"),
		len(m.files))
	summary += fmt.Sprintf("  •  %s ~%.1f KB",
		RenderMuted("📦"),
		float64(m.totalSize)/1024.0)
	b.WriteString(summary)
	b.WriteString("\n\n")

	// Viewport with file list
	b.WriteString(m.viewport.View())
	b.WriteString("\n")

	// Scroll indicator - Story 10.7 Task 9: "↓ More..." when scrollable
	if !m.viewport.AtBottom() {
		moreIndicator := InfoStyle.Render("↓ More... (scroll down to see more files)")
		b.WriteString(moreIndicator)
		b.WriteString("\n")
	}

	// Footer with navigation hints
	scrollInfo := fmt.Sprintf("%.0f%% ", m.viewport.ScrollPercent()*100)
	if m.viewport.AtTop() {
		scrollInfo = "Top "
	} else if m.viewport.AtBottom() {
		scrollInfo = "Bottom "
	}

	footer := scrollInfo + "• ↑↓/k/j: Navigate • PgUp/PgDn: Page • Home/End: Jump • Enter: Confirm • Esc: Back"
	b.WriteString(RenderFooter(footer))
	b.WriteString("\n")

	return b.String()
}

// renderContent generates the full content for the viewport (file tree + previews).
func (m PreviewModel) renderContent() string {
	var b strings.Builder

	// Build file tree structure
	tree := m.buildFileTree()

	b.WriteString(RenderHeader("📂 Project Structure:"))
	b.WriteString("\n\n")
	b.WriteString(tree)

	return b.String()
}

// buildFileTree creates a tree-style representation of the files.
func (m PreviewModel) buildFileTree() string {
	var b strings.Builder

	// Create root node
	root := &dirEntry{
		isDir:    true,
		name:     m.projectName,
		children: make(map[string]*dirEntry),
	}

	// Build tree structure
	for _, file := range m.files {
		parts := strings.Split(file.Path, "/")
		current := root

		// Navigate/create directory structure
		for i := 0; i < len(parts)-1; i++ {
			dirName := parts[i]
			if _, exists := current.children[dirName]; !exists {
				current.children[dirName] = &dirEntry{
					isDir:    true,
					name:     dirName,
					children: make(map[string]*dirEntry),
				}
			}
			current = current.children[dirName]
		}

		// Add file
		fileName := parts[len(parts)-1]
		current.children[fileName] = &dirEntry{
			isDir:   false,
			name:    fileName,
			size:    file.Size,
			content: file.Content,
		}
	}

	// Render tree
	m.renderTreeNode(&b, root, "", true)

	return b.String()
}

// renderTreeNode recursively renders a tree node with proper indentation.
func (m PreviewModel) renderTreeNode(b *strings.Builder, node *dirEntry, prefix string, isLast bool) {
	if node.name != m.projectName { // Skip root project name, we already show it in title
		// Create diff-style prefix
		diffStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorSuccess))

		icon := "📄"
		if node.isDir {
			icon = "📁"
		}

		// Tree structure characters
		branch := "├── "
		if isLast {
			branch = "└── "
		}

		nameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorText))
		if node.isDir {
			nameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorInfo)).Bold(true)
		}

		line := fmt.Sprintf("%s%s%s %s",
			prefix,
			diffStyle.Render("+"),
			branch,
			icon)

		// File name
		line += " " + nameStyle.Render(node.name)

		// File size for files
		if !node.isDir && node.size > 0 {
			sizeStr := formatSize(node.size)
			line += " " + RenderMuted(fmt.Sprintf("(%s)", sizeStr))
		}

		b.WriteString(line)
		b.WriteString("\n")

		// Show code preview for key files (main.go, config files, etc.)
		if !node.isDir && shouldShowPreview(node.name) && node.content != "" {
			preview := m.renderCodePreview(node.content)
			// Indent preview
			previewLines := strings.Split(preview, "\n")
			newPrefix := prefix
			if isLast {
				newPrefix += "    "
			} else {
				newPrefix += "│   "
			}
			for _, line := range previewLines {
				if line != "" {
					b.WriteString(newPrefix + "  " + line + "\n")
				}
			}
			b.WriteString("\n")
		}
	}

	// Render children (directories first, then files)
	if node.isDir && len(node.children) > 0 {
		// Collect and sort children
		var dirs, files []string
		for name, child := range node.children {
			if child.isDir {
				dirs = append(dirs, name)
			} else {
				files = append(files, name)
			}
		}

		// Sort alphabetically (simple bubble sort for small lists)
		sortStrings(dirs)
		sortStrings(files)

		// Combine dirs + files
		children := append(dirs, files...)

		// Render each child
		for i, childName := range children {
			child := node.children[childName]
			isLastChild := i == len(children)-1

			newPrefix := prefix
			if node.name != m.projectName {
				if isLast {
					newPrefix += "    "
				} else {
					newPrefix += "│   "
				}
			}

			m.renderTreeNode(b, child, newPrefix, isLastChild)
		}
	}
}

// renderCodePreview renders a syntax-highlighted preview of Go code.
func (m PreviewModel) renderCodePreview(code string) string {
	lines := strings.Split(code, "\n")

	// Show first 5-10 lines for preview
	maxLines := 8
	if len(lines) > maxLines {
		lines = lines[:maxLines]
		lines = append(lines, RenderMuted("  ..."))
	}

	var highlighted []string
	for _, line := range lines {
		highlighted = append(highlighted, highlightGoCode(line))
	}

	return strings.Join(highlighted, "\n")
}

// highlightGoCode applies basic syntax highlighting to Go code.
func highlightGoCode(line string) string {
	// Color styles for syntax highlighting
	keywordStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00b0ff")) // Blue
	stringStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00c853"))  // Green
	commentStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")) // Gray
	functionStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff6d00")) // Orange

	// Handle comments first (everything after // is a comment)
	if strings.Contains(line, "//") {
		parts := strings.SplitN(line, "//", 2)
		if len(parts) == 2 {
			return parts[0] + commentStyle.Render("//"+parts[1])
		}
	}

	// Highlight strings (basic - handles simple cases)
	stringRegex := regexp.MustCompile(`"[^"]*"`)
	line = stringRegex.ReplaceAllStringFunc(line, func(s string) string {
		return stringStyle.Render(s)
	})

	// Highlight keywords
	keywords := []string{
		"package", "import", "func", "type", "struct", "interface",
		"return", "if", "else", "for", "range", "switch", "case",
		"var", "const", "map", "chan", "go", "defer", "select",
		"break", "continue", "goto", "fallthrough", "default",
	}

	for _, kw := range keywords {
		// Use word boundaries to avoid partial matches
		pattern := regexp.MustCompile(`\b` + kw + `\b`)
		line = pattern.ReplaceAllStringFunc(line, func(s string) string {
			return keywordStyle.Render(s)
		})
	}

	// Highlight function calls (word followed by '(')
	funcRegex := regexp.MustCompile(`\b([a-zA-Z_][a-zA-Z0-9_]*)\s*\(`)
	line = funcRegex.ReplaceAllStringFunc(line, func(s string) string {
		funcName := strings.TrimSpace(strings.TrimSuffix(s, "("))
		return functionStyle.Render(funcName) + "("
	})

	return line
}

// shouldShowPreview determines if a file should have its content previewed.
func shouldShowPreview(filename string) bool {
	// Show preview for key files
	previewFiles := []string{
		"main.go",
		"config.go",
		"server.go",
		"user.go",
		"service.go",
		".env",
		"Dockerfile",
		"README.md",
	}

	for _, pf := range previewFiles {
		if strings.HasSuffix(filename, pf) {
			return true
		}
	}

	return false
}

// formatSize formats a file size in bytes to a human-readable string.
func formatSize(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	} else if bytes < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(bytes)/1024.0)
	} else {
		return fmt.Sprintf("%.1f MB", float64(bytes)/(1024.0*1024.0))
	}
}

// sortStrings sorts a slice of strings in-place (simple bubble sort).
func sortStrings(arr []string) {
	n := len(arr)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
}

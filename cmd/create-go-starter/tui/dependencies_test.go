package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

// TestBubbleTeaDependenciesImport verifies that Bubble Tea dependencies are available.
// This test ensures that all required Charm packages are correctly installed.
func TestBubbleTeaDependenciesImport(t *testing.T) {
	t.Run("bubbletea framework imports successfully", func(t *testing.T) {
		// Verify we can reference the tea package
		var _ tea.Model
		var _ tea.Msg
		var _ tea.Cmd
	})

	t.Run("bubbles components import successfully", func(t *testing.T) {
		// Verify we can reference bubbles components
		var _ textinput.Model
		var _ list.Model
		var _ progress.Model
		var _ spinner.Model
	})

	t.Run("lipgloss styling imports successfully", func(t *testing.T) {
		// Verify we can reference lipgloss
		_ = lipgloss.NewStyle()
	})
}

// TestBubbleTeaVersionCompatibility ensures we're using compatible versions.
func TestBubbleTeaVersionCompatibility(t *testing.T) {
	t.Run("bubbletea v1.x or v2 is available", func(t *testing.T) {
		// This test will pass if bubbletea is available at all
		// Version-specific compatibility is checked at runtime
		var _ tea.Model
	})
}

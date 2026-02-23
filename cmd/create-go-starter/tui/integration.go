package tui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// IsTTY checks if stdout is connected to a terminal (TTY).
// Returns false in CI/CD environments, piped output, or redirected output.
func IsTTY() bool {
	fileInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// ShouldUseTUI determines if the Bubble Tea TUI should be used.
// Requires: TTY available AND NO_COLOR not set
// Falls back to text mode for CI/CD, piped output, or when NO_COLOR is set.
func ShouldUseTUI() bool {
	return IsTTY() && !IsNoColorMode()
}

// RunInteractiveTUI launches the Bubble Tea interactive mode.
// This is the modern TUI experience with rich components.
// The generatorFunc is called to actually generate the project files.
// Returns an error if the TUI fails to run.
func RunInteractiveTUI(defaults InteractiveDefaults, generatorFunc GeneratorFunc) error {
	model := NewInteractiveModel(defaults, generatorFunc)

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),       // Use alternate screen buffer (clean restore)
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("failed to run interactive TUI: %w", err)
	}

	// Extract final model to check for errors
	m, ok := finalModel.(Model)
	if !ok {
		return fmt.Errorf("unexpected model type returned from TUI")
	}

	// Check if user quit before completing
	if m.state != StateDone {
		return fmt.Errorf("user cancelled interactive mode")
	}

	// Check if generation failed
	if m.err != nil {
		return fmt.Errorf("project generation failed: %w", m.err)
	}

	// Success - the generation already happened during StateGenerating
	return nil
}

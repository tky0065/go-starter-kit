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

// TUIResult contains the results of a successful TUI generation.
// This is used by the caller to display post-generation messages
// after the alternate screen buffer has been restored.
type TUIResult struct {
	ProjectName string
	Database    string
}

// RunInteractiveTUI launches the Bubble Tea interactive mode.
// This is the modern TUI experience with rich components.
// The generatorFunc is called to actually generate the project files.
// On success, returns a TUIResult with the project info so the caller can
// display post-generation messages on the normal terminal (after alt screen restore).
// Returns an error if the TUI fails to run or the user cancels.
func RunInteractiveTUI(defaults InteractiveDefaults, generatorFunc GeneratorFunc) (*TUIResult, error) {
	model := NewInteractiveModel(defaults, generatorFunc)

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),       // Use alternate screen buffer (clean restore)
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	finalModel, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to run interactive TUI: %w", err)
	}

	// Extract final model to check for errors
	m, ok := finalModel.(Model)
	if !ok {
		return nil, fmt.Errorf("unexpected model type returned from TUI")
	}

	// Check if user quit before completing
	if m.state != StateDone {
		return nil, fmt.Errorf("user cancelled interactive mode")
	}

	// Check if generation failed
	if m.err != nil {
		return nil, fmt.Errorf("project generation failed: %w", m.err)
	}

	// Success - return project info for post-generation display
	return &TUIResult{
		ProjectName: m.projectName,
		Database:    m.database,
	}, nil
}

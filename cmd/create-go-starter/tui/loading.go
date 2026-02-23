package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// LoadingModel represents a loading spinner with contextual text.
// It provides elegant visual feedback for async operations with colored spinners.
// Story 10.7 Task 4 (AC #4): Elegant spinner with contextual text during long operations.
type LoadingModel struct {
	spinner     spinner.Model
	text        string // Contextual text (e.g., "Initializing Git...")
	state       LoadingState
	spinnerType spinner.Spinner // Configurable spinner type
}

// LoadingState represents the current state of the loading operation.
type LoadingState int

const (
	// LoadingInProgress indicates the operation is in progress (blue spinner)
	LoadingInProgress LoadingState = iota
	// LoadingSuccess indicates the operation completed successfully (green checkmark)
	LoadingSuccess
	// LoadingError indicates the operation failed (red cross)
	LoadingError
)

// NewLoadingModel creates a new loading model with the specified text.
// The spinner starts with a blue color (info) to indicate in-progress state.
func NewLoadingModel(text string) LoadingModel {
	s := spinner.New()
	s.Spinner = spinner.Dot // Default to Dot spinner (elegant and minimal)
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorInfo)) // Blue during operation

	return LoadingModel{
		spinner:     s,
		text:        text,
		state:       LoadingInProgress,
		spinnerType: spinner.Dot,
	}
}

// NewLoadingModelWithSpinner creates a loading model with a custom spinner type.
// Available spinner types: spinner.Line, spinner.Dot, spinner.MiniDot, spinner.Jump, etc.
func NewLoadingModelWithSpinner(text string, spinnerType spinner.Spinner) LoadingModel {
	s := spinner.New()
	s.Spinner = spinnerType
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorInfo))

	return LoadingModel{
		spinner:     s,
		text:        text,
		state:       LoadingInProgress,
		spinnerType: spinnerType,
	}
}

// SetText updates the contextual text of the loading spinner.
// Useful for showing different stages of the same operation.
func (m *LoadingModel) SetText(text string) {
	m.text = text
}

// SetSuccess marks the operation as successful.
// The spinner is replaced with a green checkmark (✓).
func (m *LoadingModel) SetSuccess() {
	m.state = LoadingSuccess
}

// SetError marks the operation as failed.
// The spinner is replaced with a red cross (✗).
func (m *LoadingModel) SetError() {
	m.state = LoadingError
}

// IsInProgress returns true if the loading is still in progress.
func (m LoadingModel) IsInProgress() bool {
	return m.state == LoadingInProgress
}

// Init initializes the loading model and starts the spinner animation.
func (m LoadingModel) Init() tea.Cmd {
	return m.spinner.Tick
}

// Update handles messages for the loading model.
// It updates the spinner animation when in progress.
func (m LoadingModel) Update(msg tea.Msg) (LoadingModel, tea.Cmd) {
	if m.state != LoadingInProgress {
		// Don't update spinner if operation is complete
		return m, nil
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

// View renders the loading model.
// Format: "[spinner/icon] Contextual text..." with appropriate color.
func (m LoadingModel) View() string {
	switch m.state {
	case LoadingInProgress:
		// Show animated spinner with blue color
		spinnerView := m.spinner.View()
		return fmt.Sprintf("%s %s", spinnerView, InfoStyle.Render(m.text))

	case LoadingSuccess:
		// Show green checkmark
		checkmark := SuccessStyle.Render("✓")
		return fmt.Sprintf("%s %s", checkmark, SuccessStyle.Render(m.text))

	case LoadingError:
		// Show red cross
		cross := ErrorStyle.Render("✗")
		return fmt.Sprintf("%s %s", cross, ErrorStyle.Render(m.text))

	default:
		return m.text
	}
}

// Common loading messages for typical operations.
// These can be used as predefined text for common tasks.
const (
	LoadingInitializingGit      = "Initializing Git repository..."
	LoadingInstallingDeps       = "Installing Go dependencies..."
	LoadingRunningGoModTidy     = "Running go mod tidy..."
	LoadingGeneratingFiles      = "Generating project files..."
	LoadingCreatingDirectories  = "Creating directories..."
	LoadingWritingConfig        = "Writing configuration files..."
	LoadingSettingUpDatabase    = "Setting up database schema..."
	LoadingConfiguringDocker    = "Configuring Docker containers..."
	LoadingRunningTests         = "Running tests..."
	LoadingBuildingProject      = "Building project..."
	LoadingVerifyingInstall     = "Verifying installation..."
)

// LoadingOperation represents a single loading operation with its text and state.
// Useful for tracking multiple sequential operations.
type LoadingOperation struct {
	Text  string
	State LoadingState
}

// MultiLoadingModel manages multiple loading operations sequentially.
// Story 10.7 Task 4: Display contextual text for different operations.
type MultiLoadingModel struct {
	operations     []LoadingOperation
	currentIndex   int
	spinner        spinner.Model
	showCompleted  bool // Whether to show completed operations
}

// NewMultiLoadingModel creates a model for managing multiple sequential operations.
func NewMultiLoadingModel(operations []string, showCompleted bool) MultiLoadingModel {
	ops := make([]LoadingOperation, len(operations))
	for i, text := range operations {
		ops[i] = LoadingOperation{
			Text:  text,
			State: LoadingInProgress,
		}
	}

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorInfo))

	return MultiLoadingModel{
		operations:    ops,
		currentIndex:  0,
		spinner:       s,
		showCompleted: showCompleted,
	}
}

// SetCurrentSuccess marks the current operation as successful and advances.
func (m *MultiLoadingModel) SetCurrentSuccess() {
	if m.currentIndex < len(m.operations) {
		m.operations[m.currentIndex].State = LoadingSuccess
		m.currentIndex++
	}
}

// SetCurrentError marks the current operation as failed.
func (m *MultiLoadingModel) SetCurrentError() {
	if m.currentIndex < len(m.operations) {
		m.operations[m.currentIndex].State = LoadingError
	}
}

// IsComplete returns true if all operations are complete.
func (m MultiLoadingModel) IsComplete() bool {
	return m.currentIndex >= len(m.operations)
}

// CurrentOperation returns the current operation being processed.
func (m MultiLoadingModel) CurrentOperation() string {
	if m.currentIndex < len(m.operations) {
		return m.operations[m.currentIndex].Text
	}
	return ""
}

// Init initializes the multi-loading model.
func (m MultiLoadingModel) Init() tea.Cmd {
	return m.spinner.Tick
}

// Update handles messages for the multi-loading model.
func (m MultiLoadingModel) Update(msg tea.Msg) (MultiLoadingModel, tea.Cmd) {
	if m.IsComplete() {
		return m, nil
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

// View renders all operations with their current state.
func (m MultiLoadingModel) View() string {
	if len(m.operations) == 0 {
		return ""
	}

	var views []string

	for i, op := range m.operations {
		var line string

		switch op.State {
		case LoadingInProgress:
			if i == m.currentIndex {
				// Current operation - show animated spinner
				line = fmt.Sprintf("%s %s", m.spinner.View(), InfoStyle.Render(op.Text))
			} else if m.showCompleted && i < m.currentIndex {
				// Completed (shouldn't happen, but handle gracefully)
				line = fmt.Sprintf("%s %s", SuccessStyle.Render("✓"), MutedStyle.Render(op.Text))
			} else {
				// Future operation - don't show yet
				continue
			}

		case LoadingSuccess:
			if m.showCompleted {
				// Show completed operation with green checkmark
				line = fmt.Sprintf("%s %s", SuccessStyle.Render("✓"), MutedStyle.Render(op.Text))
			} else {
				// Don't show completed operations
				continue
			}

		case LoadingError:
			// Always show errors
			line = fmt.Sprintf("%s %s", ErrorStyle.Render("✗"), ErrorStyle.Render(op.Text))
		}

		views = append(views, line)
	}

	return lipgloss.JoinVertical(lipgloss.Left, views...)
}

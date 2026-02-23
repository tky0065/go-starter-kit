package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// FormModel manages multi-step form navigation and state persistence.
// Story 10.7 Task 2: Multi-step forms with navigation (AC #2, #10)
type FormModel struct {
	// Form state
	currentStep int   // Current step (1-indexed)
	totalSteps  int   // Total number of steps
	history     []int // Navigation history stack for back button

	// Form values (persisted during navigation)
	projectName   string
	template      string
	database      string
	observability string

	// Validation
	validationError string
}

// NewFormModel creates a new FormModel with initial state.
func NewFormModel() FormModel {
	return FormModel{
		currentStep: 1,
		totalSteps:  4, // Project Name, Template, Database, Observability
		history:     []int{},
	}
}

// Init initializes the FormModel (required by tea.Model interface).
func (f FormModel) Init() tea.Cmd {
	return nil
}

// Update handles form navigation and input messages.
// Implements the Bubble Tea Update pattern for FormModel.
func (f FormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return f.handleKeyMsg(msg)
	}
	return f, nil
}

// handleKeyMsg processes keyboard input for form navigation.
func (f FormModel) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		// Forward navigation
		return f.goNext(), nil

	case tea.KeyLeft:
		// Back navigation with arrow key
		return f.goBack(), nil

	case tea.KeyRight:
		// Forward navigation with arrow key
		return f.goNext(), nil

	case tea.KeyRunes:
		// Handle 'b' for back, 'n' for next
		if len(msg.Runes) > 0 {
			switch msg.Runes[0] {
			case 'b', 'B':
				return f.goBack(), nil
			case 'n', 'N':
				return f.goNext(), nil
			}
		}
	}

	return f, nil
}

// goNext advances to the next step and updates history.
// Story 10.7 Task 2: Forward navigation (→/n/Enter)
func (f FormModel) goNext() tea.Model {
	if f.currentStep >= f.totalSteps {
		// Already at last step, don't advance
		return f
	}

	// Push current step to history
	f.history = append(f.history, f.currentStep)
	// Advance to next step
	f.currentStep++

	return f
}

// goBack returns to the previous step using history stack.
// Story 10.7 Task 2: Back navigation (←/b)
func (f FormModel) goBack() tea.Model {
	if len(f.history) == 0 {
		// No history, stay at current step
		return f
	}

	// Pop from history
	f.currentStep = f.history[len(f.history)-1]
	f.history = f.history[:len(f.history)-1]

	return f
}

// RenderProgressIndicator renders the visual progress indicator.
// Story 10.7 Task 2: Progress indicators "Step 1/5: Project Name"
func (f FormModel) RenderProgressIndicator() string {
	stepName := f.GetStepName(f.currentStep)
	return fmt.Sprintf("Step %d/%d: %s", f.currentStep, f.totalSteps, stepName)
}

// RenderBreadcrumb renders the breadcrumb indicator at the top.
// Story 10.7 Task 2: Breadcrumb visual en haut de l'écran
//
// Example output:
// [1/4] Project Name → [2/4] Template → [3/4] Database → [4/4] Observability
//   ●                     ●                  ○                  ○
func (f FormModel) RenderBreadcrumb() string {
	var b strings.Builder

	// Step names row
	stepNames := []string{}
	indicators := []string{}

	for step := 1; step <= f.totalSteps; step++ {
		stepName := f.GetStepName(step)
		label := fmt.Sprintf("[%d/%d] %s", step, f.totalSteps, stepName)

		if step == f.currentStep {
			// Current step - highlighted
			stepNames = append(stepNames, HighlightStyle.Render(label))
			indicators = append(indicators, SuccessStyle.Render("●"))
		} else if step < f.currentStep {
			// Completed step - success color
			stepNames = append(stepNames, SuccessStyle.Render(label))
			indicators = append(indicators, SuccessStyle.Render("●"))
		} else {
			// Future step - muted
			stepNames = append(stepNames, MutedStyle.Render(label))
			indicators = append(indicators, MutedStyle.Render("○"))
		}
	}

	// Join with arrows
	b.WriteString(strings.Join(stepNames, "  →  "))
	b.WriteString("\n")
	// Indicators row
	b.WriteString(strings.Join(indicators, "     "))

	return b.String()
}

// GetStepName returns the name of a given step.
func (f FormModel) GetStepName(step int) string {
	switch step {
	case 1:
		return "Project Name"
	case 2:
		return "Template"
	case 3:
		return "Database"
	case 4:
		return "Observability"
	default:
		return "Unknown"
	}
}

// View renders the form based on current step.
// This is called from the main Model.View() when in form state.
func (f FormModel) View() string {
	var b strings.Builder

	// Render breadcrumb at top
	b.WriteString(f.RenderBreadcrumb())
	b.WriteString("\n\n")

	// Render progress indicator
	b.WriteString(RenderHeader(f.RenderProgressIndicator()))
	b.WriteString("\n\n")

	// Render current step content
	switch f.currentStep {
	case 1:
		b.WriteString(f.renderProjectNameStep())
	case 2:
		b.WriteString(f.renderTemplateStep())
	case 3:
		b.WriteString(f.renderDatabaseStep())
	case 4:
		b.WriteString(f.renderObservabilityStep())
	}

	b.WriteString("\n\n")

	// Render navigation help
	b.WriteString(f.renderNavigationHelp())

	// Render validation error if present
	if f.validationError != "" {
		b.WriteString("\n\n")
		b.WriteString(ErrorStyle.Render("⚠ " + f.validationError))
	}

	return b.String()
}

// renderProjectNameStep renders step 1: project name input.
func (f FormModel) renderProjectNameStep() string {
	var b strings.Builder

	b.WriteString("Enter the project name:\n\n")

	if f.projectName != "" {
		b.WriteString(FocusedStyle.Render(f.projectName))
	} else {
		b.WriteString(BlurredStyle.Render("(will use placeholder: my-project)"))
	}

	b.WriteString("\n\n")
	b.WriteString(MutedStyle.Render("Project name will be used for the directory and module name."))

	return b.String()
}

// renderTemplateStep renders step 2: template selection.
func (f FormModel) renderTemplateStep() string {
	var b strings.Builder

	b.WriteString("Select a template:\n\n")

	if f.template != "" {
		b.WriteString(FocusedStyle.Render("Selected: " + f.template))
	} else {
		b.WriteString(BlurredStyle.Render("(use template selection from main flow)"))
	}

	return b.String()
}

// renderDatabaseStep renders step 3: database selection.
func (f FormModel) renderDatabaseStep() string {
	var b strings.Builder

	b.WriteString("Select a database:\n\n")

	if f.database != "" {
		b.WriteString(FocusedStyle.Render("Selected: " + f.database))
	} else {
		b.WriteString(BlurredStyle.Render("(use database selection from main flow)"))
	}

	return b.String()
}

// renderObservabilityStep renders step 4: observability selection.
func (f FormModel) renderObservabilityStep() string {
	var b strings.Builder

	b.WriteString("Select observability level:\n\n")

	if f.observability != "" {
		b.WriteString(FocusedStyle.Render("Selected: " + f.observability))
	} else {
		b.WriteString(BlurredStyle.Render("(use observability selection from main flow)"))
	}

	return b.String()
}

// renderNavigationHelp renders navigation keyboard shortcuts.
// Story 10.7 Task 2: Keyboard shortcuts help
func (f FormModel) renderNavigationHelp() string {
	var parts []string

	// Back navigation (only show if we can go back)
	if len(f.history) > 0 {
		parts = append(parts, "← or b: Back")
	}

	// Forward navigation (only show if we can go forward)
	if f.currentStep < f.totalSteps {
		parts = append(parts, "→ or n: Next")
	}

	// Always show Enter to confirm
	parts = append(parts, "Enter: Confirm")

	// Always show help and quit
	parts = append(parts, "?: Help")
	parts = append(parts, "Ctrl+C: Quit")

	return RenderFooter(strings.Join(parts, " • "))
}

// RenderFormStep renders a complete form step (used by main Model).
// This is a convenience function that wraps View() for integration.
func RenderFormStep(form FormModel, width, height int) string {
	// Apply width constraints if needed
	style := lipgloss.NewStyle().Width(width - 4)
	return style.Render(form.View())
}

package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View renders the current state of the model to a string.
// This is a pure function with ZERO side-effects (critical for Elm Architecture).
func (m Model) View() string {
	// Handle different states
	switch m.state {
	case StateWelcome:
		return m.viewWelcome()
	case StateProjectName:
		return m.viewProjectName()
	case StateTemplateSelect:
		return m.viewTemplateSelect()
	case StateDatabaseSelect:
		return m.viewDatabaseSelect()
	case StateFrameworkSelect:
		return m.viewFrameworkSelect()
	case StateObservabilitySelect:
		return m.viewObservabilitySelect()
	case StateSummary:
		return m.viewSummary()
	case StatePreview:
		return m.viewPreview()
	case StateGenerating:
		return m.viewGenerating()
	case StateDone:
		return m.viewDone()
	case StateHelp:
		return m.viewHelp()
	default:
		return "Unknown state"
	}
}

// viewProjectName renders the project name input screen.
func (m Model) viewProjectName() string {
	var b strings.Builder

	b.WriteString("\n")
	titleWidth := 64
	if m.width-4 < titleWidth {
		titleWidth = m.width - 4
	}

	title := GradientHeaderStyle.Width(titleWidth).Render("CREATE-GO-STARTER")
	b.WriteString(CenterHorizontal(title, m.width))
	b.WriteString("\n\n")

	b.WriteString(CenterHorizontal(RenderHeader("Project Name"), m.width))
	b.WriteString("\n")

	// Center the input field
	inputView := m.projectInput.View()
	b.WriteString(CenterHorizontal(inputView, m.width))
	b.WriteString("\n")

	// Show validation error if present
	if m.validationError != "" {
		b.WriteString("\n")
		errBox := lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorError)).
			Bold(true).
			Padding(0, 1)
		b.WriteString(CenterHorizontal(errBox.Render("✗ "+m.validationError), m.width))
		b.WriteString("\n")
	} else {
		b.WriteString("\n")
		info := MutedStyle.Render("Standard Go module name (e.g., github.com/user/project)")
		b.WriteString(CenterHorizontal(info, m.width))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	footer := renderStyledFooter("Enter", "Continue", "?", "Help", "Ctrl+C", "Quit")
	b.WriteString(CenterHorizontal(footer, m.width))
	b.WriteString("\n")

	return b.String()
}

// viewTemplateSelect renders the template selection screen.
func (m Model) viewTemplateSelect() string {
	var b strings.Builder

	b.WriteString("\n")
	breadcrumb := RenderBreadcrumb(m.projectName, false) + MutedStyle.Render(" / ") + RenderBreadcrumb("Template", true)
	b.WriteString(CenterHorizontal(breadcrumb, m.width))
	b.WriteString("\n\n")

	// List.View() already handles its own layout, but we need to ensure it's centered
	b.WriteString(m.templateList.View())
	b.WriteString("\n")

	footer := renderStyledFooter("↑↓", "Navigate", "Enter", "Select", "?", "Help", "Esc", "Back")
	b.WriteString(CenterHorizontal(footer, m.width))
	b.WriteString("\n")

	return b.String()
}

// viewDatabaseSelect renders the database selection screen.
func (m Model) viewDatabaseSelect() string {
	var b strings.Builder

	b.WriteString("\n")
	breadcrumb := RenderBreadcrumb(m.projectName, false) + MutedStyle.Render(" / ") +
		RenderBreadcrumb(m.template, false) + MutedStyle.Render(" / ") +
		RenderBreadcrumb("Database", true)
	b.WriteString(CenterHorizontal(breadcrumb, m.width))
	b.WriteString("\n\n")

	b.WriteString(m.databaseList.View())
	b.WriteString("\n")

	footer := renderStyledFooter("↑↓", "Navigate", "Enter", "Select", "?", "Help", "Esc", "Back")
	b.WriteString(CenterHorizontal(footer, m.width))
	b.WriteString("\n")

	return b.String()
}

// viewFrameworkSelect renders the framework selection screen.
func (m Model) viewFrameworkSelect() string {
	var b strings.Builder

	b.WriteString("\n")
	breadcrumb := RenderBreadcrumb(m.projectName, false) + MutedStyle.Render(" / ") +
		RenderBreadcrumb(m.template, false) + MutedStyle.Render(" / ") +
		RenderBreadcrumb(m.database, false) + MutedStyle.Render(" / ") +
		RenderBreadcrumb("Framework", true)
	b.WriteString(CenterHorizontal(breadcrumb, m.width))
	b.WriteString("\n\n")

	b.WriteString(m.frameworkList.View())
	b.WriteString("\n")

	footer := renderStyledFooter("↑↓", "Navigate", "Enter", "Select", "?", "Help", "Esc", "Back")
	b.WriteString(CenterHorizontal(footer, m.width))
	b.WriteString("\n")

	return b.String()
}

// viewObservabilitySelect renders the observability selection screen.
func (m Model) viewObservabilitySelect() string {
	var b strings.Builder

	b.WriteString("\n")
	breadcrumb := RenderBreadcrumb(m.projectName, false) + MutedStyle.Render(" / ") +
		RenderBreadcrumb(m.template, false) + MutedStyle.Render(" / ") +
		RenderBreadcrumb(m.database, false) + MutedStyle.Render(" / ") +
		RenderBreadcrumb(m.framework, false) + MutedStyle.Render(" / ") +
		RenderBreadcrumb("Observability", true)
	b.WriteString(CenterHorizontal(breadcrumb, m.width))
	b.WriteString("\n\n")

	b.WriteString(m.obsList.View())
	b.WriteString("\n")

	footer := renderStyledFooter("↑↓", "Navigate", "Enter", "Select", "?", "Help", "Esc", "Back")
	b.WriteString(CenterHorizontal(footer, m.width))
	b.WriteString("\n")

	return b.String()
}

// viewSummary renders the configuration summary screen.
func (m Model) viewSummary() string {
	var b strings.Builder

	b.WriteString("\n")

	// Title with capped width, centered
	titleWidth := 64
	if m.width-4 < titleWidth {
		titleWidth = m.width - 4
	}

	title := GradientHeaderStyle.Width(titleWidth).Render("CONFIGURATION SUMMARY")
	b.WriteString(CenterHorizontal(title, m.width))
	b.WriteString("\n\n")

	// Summary content with aligned labels
	labelStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(ColorMuted)).
		Width(18).
		Align(lipgloss.Right)

	valueStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(ColorSuccess))

	rows := []struct {
		label string
		value string
	}{
		{"Project Name", m.projectName},
		{"Template", m.template},
		{"Database", m.database},
		{"Framework", m.framework},
		{"Observability", m.observability},
	}

	var summary string
	for i, row := range rows {
		line := fmt.Sprintf("%s  %s",
			labelStyle.Render(row.label),
			valueStyle.Render(row.value))

		summary += line
		if i < len(rows)-1 {
			summary += "\n\n"
		}
	}

	// Content box with a more sophisticated border
	contentBox := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color(ColorSuccess)).
		Padding(2, 4).
		Width(titleWidth)

	b.WriteString(CenterHorizontal(contentBox.Render(summary), m.width))
	b.WriteString("\n\n")

	// Styled footer
	instruction := MutedStyle.Render("Press Enter to confirm and generate your project")
	b.WriteString(CenterHorizontal(instruction, m.width))
	b.WriteString("\n\n")

	footer := renderStyledFooter("Enter", "Generate", "?", "Help", "Esc", "Back", "Ctrl+C", "Quit")
	b.WriteString(CenterHorizontal(footer, m.width))
	b.WriteString("\n")

	return b.String()
}

// viewGenerating renders the file generation progress screen.
// Story 10.7 Task 3: Enhanced with real-time statistics panel
func (m Model) viewGenerating() string {
	var b strings.Builder

	b.WriteString("\n")

	// Title with spinner, centered
	titleWidth := 64
	if m.width-4 < titleWidth {
		titleWidth = m.width - 4
	}

	title := fmt.Sprintf("%s GENERATING PROJECT", m.loadSpinner.View())
	b.WriteString(CenterHorizontal(GradientHeaderStyle.Width(titleWidth).Render(title), m.width))
	b.WriteString("\n\n")

	if m.projectName != "" {
		projectLine := fmt.Sprintf("%s %s", RenderMuted("Project:"), RenderHighlight(m.projectName))
		b.WriteString(CenterHorizontal(projectLine, m.width))
		b.WriteString("\n\n")
	}

	if m.totalFiles > 0 {
		percent := float64(m.filesGenerated) / float64(m.totalFiles)

		// Progress bar with percentage, centered
		progressBarWidth := 50
		if m.width-10 < progressBarWidth {
			progressBarWidth = m.width - 10
		}
		m.progressBar.Width = progressBarWidth

		progressInfo := fmt.Sprintf("%s  %s",
			RenderHighlight(fmt.Sprintf("%d/%d files", m.filesGenerated, m.totalFiles)),
			RenderMuted(fmt.Sprintf("%.0f%%", percent*100)))

		b.WriteString(CenterHorizontal(progressInfo, m.width))
		b.WriteString("\n")

		// Render progress bar with pulsing gradient animation (Task 10)
		pb := renderProgressBarWithPulse(m.progressBar, percent, m.progressPulse)
		b.WriteString(CenterHorizontal(pb, m.width))
		b.WriteString("\n\n")

		// Statistics panel (Story 10.7 Task 3)
		statsPanel := renderGenerationStats(m.genStats, m.totalFiles)
		b.WriteString(CenterHorizontal(statsPanel, m.width))
		b.WriteString("\n\n")

		// Current step/file indicator (Story 10.7 Task 3)
		stepMessage := renderCurrentStep(m.genStats.CurrentStep, m.genStats.CurrentFile)
		if stepMessage != "" {
			b.WriteString(CenterHorizontal(stepMessage, m.width))
			b.WriteString("\n")
		}
	} else {
		b.WriteString(CenterHorizontal(RenderSimpleInfo("Initializing..."), m.width))
		b.WriteString("\n")
	}

	b.WriteString("\n")

	return b.String()
}

// viewDone renders the completion screen.
func (m Model) viewDone() string {
	var b strings.Builder

	b.WriteString("\n")

	// Check if there was an error
	if m.err != nil {
		errorBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorError)).
			Padding(1, 4).
			Width(64)

		header := ErrorStyle.Render("GENERATION FAILED")
		errMsg := RenderMuted("The following error occurred:") + "\n\n" + ErrorStyle.Render(m.err.Error())

		b.WriteString(CenterHorizontal(errorBox.Render(header+"\n\n"+errMsg), m.width))
		b.WriteString("\n\n")
		footer := renderStyledFooter("any key", "Exit")
		b.WriteString(CenterHorizontal(footer, m.width))
		b.WriteString("\n")
		return b.String()
	}

	// Success case
	successBox := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color(ColorSuccess)).
		Padding(1, 4).
		Width(64)

	labelStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(ColorMuted)).
		Width(18).
		Align(lipgloss.Right)

	header := SuccessStyle.Render("PROJECT READY")

	details := fmt.Sprintf("\n\n%s  %s\n%s  %s\n%s  %s",
		labelStyle.Render("Project:"),
		RenderHighlight(m.projectName),
		labelStyle.Render("Status:"),
		SuccessStyle.Render("Successfully Generated"),
		labelStyle.Render("Files:"),
		RenderHighlight(fmt.Sprintf("%d files created", m.filesGenerated)))

	nextSteps := "\n\n" + RenderHeader("Next Steps:") + "\n" +
		MutedStyle.Render(fmt.Sprintf("  cd %s\n  ./setup.sh\n  make run", m.projectName))

	b.WriteString(CenterHorizontal(successBox.Render(header+details+nextSteps), m.width))
	b.WriteString("\n\n")

	instruction := SuccessStyle.Render("Your project is ready for development.")
	b.WriteString(CenterHorizontal(instruction, m.width))
	b.WriteString("\n\n")

	footer := renderStyledFooter("any key", "Exit")
	b.WriteString(CenterHorizontal(footer, m.width))
	b.WriteString("\n")

	return b.String()
}

// viewPreview renders the dry-run preview screen (Story 10.7 Task 6).
func (m Model) viewPreview() string {
	return m.previewModel.View()
}

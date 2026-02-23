package tui

import (
	"fmt"
	"strings"
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
	title := "create-go-starter — Interactive Mode"
	b.WriteString(RenderTitleBox(title, m.width-4))
	b.WriteString("\n\n")
	b.WriteString(RenderHeader("📝 Project Name:"))
	b.WriteString("\n")
	b.WriteString(m.projectInput.View())
	b.WriteString("\n")

	// Show validation error if present (MEDIUM-2 fix)
	if m.validationError != "" {
		b.WriteString("\n")
		b.WriteString(RenderSimpleError(fmt.Sprintf("❌ %s", m.validationError)))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(RenderFooter("Press Enter to continue • ? for help • Ctrl+C to quit"))
	b.WriteString("\n")

	return b.String()
}

// viewTemplateSelect renders the template selection screen.
func (m Model) viewTemplateSelect() string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("📝 Project: %s\n", m.projectName))
	b.WriteString("\n")
	b.WriteString(m.templateList.View())
	b.WriteString("\n")
	b.WriteString("↑↓ Navigate • Enter Select • ? Help • Esc Back • Ctrl+C Quit\n")

	return b.String()
}

// viewDatabaseSelect renders the database selection screen.
func (m Model) viewDatabaseSelect() string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("📝 Project: %s | Template: %s\n", m.projectName, m.template))
	b.WriteString("\n")
	b.WriteString(m.databaseList.View())
	b.WriteString("\n")
	b.WriteString("↑↓ Navigate • Enter Select • ? Help • Esc Back • Ctrl+C Quit\n")

	return b.String()
}

// viewObservabilitySelect renders the observability selection screen.
func (m Model) viewObservabilitySelect() string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("📝 Project: %s | Template: %s | Database: %s\n",
		m.projectName, m.template, m.database))
	b.WriteString("\n")
	b.WriteString(m.obsList.View())
	b.WriteString("\n")
	b.WriteString("↑↓ Navigate • Enter Select • ? Help • Esc Back • Ctrl+C Quit\n")

	return b.String()
}

// viewSummary renders the configuration summary screen.
func (m Model) viewSummary() string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(RenderTitleBox("Configuration Summary", m.width-4))
	b.WriteString("\n\n")

	// Create summary content
	summary := fmt.Sprintf("%s Project Name:   %s\n", RenderMuted("📝"), RenderHighlight(m.projectName))
	summary += fmt.Sprintf("%s Template:       %s\n", RenderMuted("📦"), RenderHighlight(m.template))
	summary += fmt.Sprintf("%s Database:       %s\n", RenderMuted("🗄️ "), RenderHighlight(m.database))
	summary += fmt.Sprintf("%s Observability:  %s", RenderMuted("📊"), RenderHighlight(m.observability))

	b.WriteString(RenderBox(summary))
	b.WriteString("\n")
	b.WriteString(RenderFooter("Press Enter to generate • ? for help • Esc to go back • Ctrl+C to quit"))
	b.WriteString("\n")

	return b.String()
}

// viewGenerating renders the file generation progress screen.
// Story 10.7 Task 3: Enhanced with real-time statistics panel
func (m Model) viewGenerating() string {
	var b strings.Builder

	b.WriteString("\n")

	// Title with spinner
	title := fmt.Sprintf("%s Generating Project: %s",
		m.loadSpinner.View(),
		RenderHighlight(m.projectName))
	b.WriteString(RenderTitleBox(title, m.width-4))
	b.WriteString("\n\n")

	if m.totalFiles > 0 {
		percent := float64(m.filesGenerated) / float64(m.totalFiles)

		// Progress bar with percentage and pulsing gradient (Story 10.7 Task 3 & Task 10)
		b.WriteString(fmt.Sprintf("Progress: %s (%s)\n",
			RenderHighlight(fmt.Sprintf("%d/%d files", m.filesGenerated, m.totalFiles)),
			RenderMuted(fmt.Sprintf("%.0f%%", percent*100))))
		b.WriteString("\n")
		// Render progress bar with pulsing gradient animation (Task 10)
		b.WriteString(renderProgressBarWithPulse(m.progressBar, percent, m.progressPulse))
		b.WriteString("\n\n")

		// Statistics panel (Story 10.7 Task 3)
		statsPanel := renderGenerationStats(m.genStats, m.totalFiles)
		b.WriteString(statsPanel)
		b.WriteString("\n\n")

		// Current step/file indicator (Story 10.7 Task 3)
		stepMessage := renderCurrentStep(m.genStats.CurrentStep, m.genStats.CurrentFile)
		if stepMessage != "" {
			b.WriteString(stepMessage)
			b.WriteString("\n")
		}
	} else {
		b.WriteString(RenderSimpleInfo("Initializing..."))
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
		b.WriteString(RenderSimpleError(fmt.Sprintf("❌ Generation failed: %v", m.err)))
		b.WriteString("\n\n")
		b.WriteString(RenderFooter("Press any key to exit"))
		b.WriteString("\n")
		return b.String()
	}

	// Success case
	b.WriteString(RenderSuccess("✅ Project generated successfully!"))
	b.WriteString("\n\n")

	info := fmt.Sprintf("📁 Project:       %s\n", RenderHighlight(m.projectName))
	info += fmt.Sprintf("📊 Files created: %s", RenderHighlight(fmt.Sprintf("%d", m.filesGenerated)))

	b.WriteString(RenderBox(info))
	b.WriteString("\n")
	b.WriteString(RenderFooter("Press any key to exit"))
	b.WriteString("\n")

	return b.String()
}

// viewPreview renders the dry-run preview screen (Story 10.7 Task 6).
func (m Model) viewPreview() string {
	return m.previewModel.View()
}

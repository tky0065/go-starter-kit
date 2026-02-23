package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
)

// KeyMap defines keyboard shortcuts for the application.
// This is used for generating contextual help screens.
type KeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Back   key.Binding
	Next   key.Binding
	Help   key.Binding
	Quit   key.Binding
	Enter  key.Binding
	Space  key.Binding
	PgUp   key.Binding
	PgDown key.Binding
	Home   key.Binding
	End    key.Binding
}

// DefaultKeyMap returns the default keyboard shortcuts for the application.
var DefaultKeyMap = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "go back"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "go forward"),
	),
	Back: key.NewBinding(
		key.WithKeys("b"),
		key.WithHelp("b", "back"),
	),
	Next: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "next"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("Ctrl+C", "quit"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("Enter", "confirm/select"),
	),
	Space: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("Space", "toggle"),
	),
	PgUp: key.NewBinding(
		key.WithKeys("pgup"),
		key.WithHelp("PgUp", "scroll up"),
	),
	PgDown: key.NewBinding(
		key.WithKeys("pgdown"),
		key.WithHelp("PgDn", "scroll down"),
	),
	Home: key.NewBinding(
		key.WithKeys("home"),
		key.WithHelp("Home", "go to top"),
	),
	End: key.NewBinding(
		key.WithKeys("end"),
		key.WithHelp("End", "go to bottom"),
	),
}

// viewHelp renders the contextual help screen (Story 10.7 Task 8 - AC#9).
// This screen is displayed when the user presses '?' in any state.
// The help adapts to show relevant shortcuts based on the current screen.
// Story 10.7 Task 9: Enhanced with scroll indicators for long content
func (m Model) viewHelp() string {
	// Build help content
	content := m.buildHelpContent()

	// Calculate if content is scrollable
	contentLines := strings.Count(content, "\n")
	viewportHeight := m.height - 6 // Reserve space for footer
	isScrollable := contentLines > viewportHeight

	var b strings.Builder

	// Header
	b.WriteString("\n")
	b.WriteString(RenderTitleBox("Keyboard Shortcuts", m.width-ListWidthMargin))
	b.WriteString("\n\n")

	// Main content
	b.WriteString(content)

	// Scroll indicator - Story 10.7 Task 9: "↓ More..." when scrollable
	if isScrollable {
		b.WriteString("\n")
		moreIndicator := InfoStyle.Render("↓ More help content below (this screen is scrollable)")
		b.WriteString(moreIndicator)
		b.WriteString("\n")
	}

	// Footer
	footer := "Press any key to return"
	if isScrollable {
		footer = "↑↓: Scroll • " + footer
	}
	b.WriteString(RenderFooter(footer))
	b.WriteString("\n")

	return b.String()
}

// buildHelpContent builds the help content body (separated for clarity).
func (m Model) buildHelpContent() string {
	var b strings.Builder

	// Global shortcuts (always available)
	b.WriteString(RenderHeader("Global Shortcuts"))
	b.WriteString("\n")
	b.WriteString(renderShortcut("?", "Show this help screen"))
	b.WriteString(renderShortcut("Ctrl+C", "Quit application"))
	b.WriteString(renderShortcut("Esc", "Go back / Cancel"))
	b.WriteString("\n")

	// Context-specific shortcuts based on current screen
	b.WriteString(RenderHeader(getScreenTitle(m.previousState)))
	b.WriteString("\n")
	b.WriteString(getScreenShortcuts(m.previousState))
	b.WriteString("\n")

	// Tips section
	b.WriteString(RenderHeader("Tips"))
	b.WriteString("\n")
	b.WriteString(renderTip("Use arrow keys (↑↓) to navigate lists"))
	b.WriteString(renderTip("Press Esc at any time to go back"))
	b.WriteString(renderTip("Your configuration is saved as you progress"))
	b.WriteString(renderTip("Press ? at any time to see context-specific shortcuts"))

	return b.String()
}

// getScreenTitle returns the title for the help section based on screen state.
func getScreenTitle(state AppState) string {
	switch state {
	case StateWelcome:
		return "Welcome Screen"
	case StateProjectName:
		return "Project Name Entry"
	case StateTemplateSelect:
		return "Template Selection"
	case StateDatabaseSelect:
		return "Database Selection"
	case StateObservabilitySelect:
		return "Observability Selection"
	case StateSummary:
		return "Configuration Summary"
	case StateGenerating:
		return "Generation Progress"
	case StateDone:
		return "Completion Screen"
	default:
		return "Current Screen"
	}
}

// getScreenShortcuts returns contextual shortcuts for the given screen state.
func getScreenShortcuts(state AppState) string {
	switch state {
	case StateWelcome:
		return renderShortcut("↑ / ↓", "Navigate menu items") +
			renderShortcut("Enter", "Select menu option") +
			renderShortcut("Esc", "Quit application")

	case StateProjectName:
		return renderShortcut("Enter", "Submit project name") +
			renderShortcut("Esc", "Go back to welcome screen") +
			renderShortcut("Type", "Enter your project name")

	case StateTemplateSelect, StateDatabaseSelect, StateObservabilitySelect:
		return renderShortcut("↑ / ↓", "Navigate list items") +
			renderShortcut("Enter", "Select highlighted item") +
			renderShortcut("Esc", "Go back to previous screen") +
			renderShortcut("j / k", "Navigate (vim-style)")

	case StateSummary:
		return renderShortcut("Enter", "Confirm and start generation") +
			renderShortcut("Esc", "Go back to edit configuration") +
			renderShortcut("← / →", "Navigate back/forward") +
			renderShortcut("b / n", "Back / Next (alternative)")

	case StateGenerating:
		return RenderMuted("  No shortcuts available during generation\n") +
			RenderMuted("  Please wait for the process to complete...\n")

	case StateDone:
		return renderShortcut("Any key", "Exit application") +
			renderShortcut("Enter", "Exit application")

	default:
		return RenderMuted("  No specific shortcuts for this screen\n")
	}
}

// renderShortcut renders a keyboard shortcut with consistent formatting.
func renderShortcut(key, description string) string {
	return "  " + KeyStyle.Render(key) + "  " + DescStyle.Render(description) + "\n"
}

// renderTip renders a tip with bullet point formatting.
func renderTip(tip string) string {
	bullet := InfoStyle.Render("•")
	return "  " + bullet + " " + tip + "\n"
}

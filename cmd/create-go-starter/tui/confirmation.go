package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ConfirmationModel represents a confirmation dialog with visual icons.
// Story 10.7 Task 5: Composants de confirmation et feedback visuels (AC #5)
type ConfirmationModel struct {
	// Message to display to the user
	Message string

	// ConfirmText is the text for the confirm button (default: "Yes")
	ConfirmText string

	// CancelText is the text for the cancel button (default: "No")
	CancelText string

	// Focused indicates which button is currently selected (true = confirm, false = cancel)
	Focused bool

	// Confirmed is set to true when the user confirms, false when they cancel
	Confirmed bool

	// Done is set to true when the user has made a selection
	Done bool
}

// NewConfirmationModel creates a new confirmation dialog.
func NewConfirmationModel(message string) ConfirmationModel {
	return ConfirmationModel{
		Message:     message,
		ConfirmText: "Yes",
		CancelText:  "No",
		Focused:     true, // Default to confirm button
		Done:        false,
	}
}

// Update handles input for the confirmation dialog.
func (m ConfirmationModel) Update(msg tea.Msg) (ConfirmationModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h", "shift+tab":
			// Move to cancel button
			m.Focused = false
			return m, nil

		case "right", "l", "tab":
			// Move to confirm button
			m.Focused = true
			return m, nil

		case "enter", " ":
			// Confirm or cancel based on focused button
			m.Done = true
			m.Confirmed = m.Focused
			return m, nil

		case "y", "Y":
			// Shortcut for "yes"
			m.Done = true
			m.Confirmed = true
			return m, nil

		case "n", "N", "esc":
			// Shortcut for "no"
			m.Done = true
			m.Confirmed = false
			return m, nil
		}
	}

	return m, nil
}

// View renders the confirmation dialog.
func (m ConfirmationModel) View() string {
	var b strings.Builder

	// Question icon and message
	questionIcon := "?"
	iconStyle := InfoStyle.Bold(true)
	b.WriteString(iconStyle.Render(questionIcon))
	b.WriteString(" ")
	b.WriteString(m.Message)
	b.WriteString("\n\n")

	// Render buttons with icons
	confirmButton := m.renderButton(m.ConfirmText, "✓", m.Focused)
	cancelButton := m.renderButton(m.CancelText, "✗", !m.Focused)

	buttons := lipgloss.JoinHorizontal(
		lipgloss.Left,
		confirmButton,
		"  ", // Space between buttons
		cancelButton,
	)

	b.WriteString(buttons)
	b.WriteString("\n\n")

	// Help text
	helpText := RenderMuted("Tab/←→ Navigate • Enter/Space Select • Y/N Shortcuts")
	b.WriteString(helpText)

	// Wrap in box
	return BoxStyle.Render(b.String())
}

// renderButton renders a button with an icon.
func (m ConfirmationModel) renderButton(text, icon string, focused bool) string {
	buttonText := icon + " " + text

	var style lipgloss.Style
	if focused {
		// Green for confirm, red for cancel
		if icon == "✓" {
			style = lipgloss.NewStyle().
				Foreground(lipgloss.Color(ColorSuccess)).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color(ColorSuccess)).
				Padding(0, 2).
				Bold(true)
		} else {
			style = lipgloss.NewStyle().
				Foreground(lipgloss.Color(ColorError)).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color(ColorError)).
				Padding(0, 2).
				Bold(true)
		}
	} else {
		// Blurred style
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorMuted)).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorMuted)).
			Padding(0, 2)
	}

	return style.Render(buttonText)
}

// RenderConfirmation is a helper function to render a standalone confirmation dialog.
// This can be used for quick confirmations without managing a full Model.
func RenderConfirmation(message string, confirmed bool) string {
	m := NewConfirmationModel(message)
	m.Focused = confirmed
	return m.View()
}

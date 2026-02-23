package tui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// welcomeMenuItem represents an item in the welcome menu.
// Implements list.Item interface for bubbles/list.
// Story 10.7 Task 1.3: Menu interactif avec bubbles/list
type welcomeMenuItem struct {
	title       string
	description string
	action      string // "create", "help", or "exit"
}

func (i welcomeMenuItem) Title() string       { return i.title }
func (i welcomeMenuItem) Description() string { return i.description }
func (i welcomeMenuItem) FilterValue() string { return i.title }

// welcomeItemDelegate is a custom delegate for rendering welcome menu items.
type welcomeItemDelegate struct{}

func (d welcomeItemDelegate) Height() int                               { return 2 }
func (d welcomeItemDelegate) Spacing() int                              { return 1 }
func (d welcomeItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d welcomeItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(welcomeMenuItem)
	if !ok {
		return
	}

	// Apply focused/blurred style based on selection
	var titleStyle, descStyle lipgloss.Style
	if index == m.Index() {
		titleStyle = FocusedStyle.Bold(true) // Copy() is deprecated, assignment copies automatically
		descStyle = InfoStyle
	} else {
		titleStyle = BlurredStyle
		descStyle = MutedStyle
	}

	// Render with icon
	icon := "  "
	if index == m.Index() {
		icon = "▸ "
	}

	title := titleStyle.Render(icon + i.title)
	desc := descStyle.Render("  " + i.description)

	fmt.Fprintf(w, "%s\n%s", title, desc)
}

// initializeWelcomeList creates and initializes the welcome menu list.
// Story 10.7 Task 1.3: Menu interactif avec bubbles/list
func initializeWelcomeList() list.Model {
	items := []list.Item{
		welcomeMenuItem{
			title:       "Create New Project",
			description: "Generate a new go-starter-kit project",
			action:      "create",
		},
		welcomeMenuItem{
			title:       "Help",
			description: "View keyboard shortcuts and usage guide",
			action:      "help",
		},
		welcomeMenuItem{
			title:       "Exit",
			description: "Quit the application",
			action:      "exit",
		},
	}

	delegate := welcomeItemDelegate{}
	l := list.New(items, delegate, 0, 0)
	l.Title = ""
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false) // We'll show custom help

	return l
}

// WelcomeLogo returns the ASCII art logo for go-starter-kit.
// Uses Option 2 (minimaliste) as recommended in Story 10.7.
func WelcomeLogo() string {
	logo := `
╔═══════════════════════════════════╗
║                                   ║
║     🚀  go-starter-kit  🚀       ║
║                                   ║
║   Production-Ready Go API         ║
║   in < 5 minutes                  ║
║                                   ║
╚═══════════════════════════════════╝
`
	return logo
}

// RenderWelcomeScreen renders the welcome screen with logo and interactive menu.
// Story 10.7 Task 1.3: Menu interactif avec bubbles/list
// Story 10.7 Task 1.5: Fade-in animation for logo
func RenderWelcomeScreen(width, height int, menuList list.Model, logoOpacity float64) string {
	var b strings.Builder

	// Render logo with fade-in animation (Story 10.7 Task 1.5)
	logo := WelcomeLogo()
	logoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorSuccess)). // Green for logo
		Bold(true).
		Align(lipgloss.Center)

	// Apply opacity to logo (simple approach: render only if opacity > 0)
	// For a smoother fade, we'd use color blending, but this is sufficient
	if logoOpacity > 0.0 {
		b.WriteString("\n")
		// Note: Lipgloss doesn't directly support opacity/alpha
		// A production implementation would blend colors or use extended ANSI codes
		// For now, we show/hide based on a threshold
		b.WriteString(logoStyle.Render(logo))
		b.WriteString("\n\n")
	} else {
		// Logo is invisible, add spacing to maintain layout
		b.WriteString("\n\n\n\n\n\n\n\n\n") // Roughly logo height in lines
	}

	// Render menu header
	b.WriteString(RenderHeader("What would you like to do?"))
	b.WriteString("\n\n")

	// Render interactive menu list (Story 10.7 Task 1.3)
	b.WriteString(menuList.View())
	b.WriteString("\n\n")

	// Footer with help text
	b.WriteString(RenderFooter("↑↓ Navigate • Enter Select • ? Help • Ctrl+C Quit"))
	b.WriteString("\n")

	return b.String()
}

// viewWelcome renders the welcome screen.
// This is called from Model.View() when in StateWelcome.
// Story 10.7 Task 1.5: Logo fade-in animation
func (m Model) viewWelcome() string {
	// Get current logo opacity from animation state
	logoOpacity := 1.0 // Default to fully visible
	if m.logoAnimation != nil {
		logoOpacity = m.logoAnimation.Opacity()
	}
	return RenderWelcomeScreen(m.width, m.height, m.welcomeList, logoOpacity)
}

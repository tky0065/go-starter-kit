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

func (d welcomeItemDelegate) Height() int                               { return 3 }
func (d welcomeItemDelegate) Spacing() int                              { return 0 }
func (d welcomeItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d welcomeItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(welcomeMenuItem)
	if !ok {
		return
	}

	isFocused := index == m.Index()
	width := m.Width()
	if width == 0 {
		width = 80 // Fallback
	}

	if isFocused {
		titleStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(ColorSuccess))
		descStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorMuted)).
			Italic(false)

		boxStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorSuccess)).
			Padding(0, 2).
			Width(44)

		content := fmt.Sprintf("%s\n%s",
			titleStyle.Render(i.title),
			descStyle.Render(i.description))

		centeredBox := CenterHorizontal(boxStyle.Render(content), width)
		fmt.Fprint(w, centeredBox)
	} else {
		titleStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorText))
		descStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorMuted))

		content := fmt.Sprintf("%s\n%s",
			titleStyle.Render(i.title),
			descStyle.Render(i.description))

		centeredContent := CenterHorizontal(content, width)
		fmt.Fprintf(w, "\n%s\n", centeredContent)
	}
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

// WelcomeLogo returns the styled logo for go-starter-kit.
// Renders the logo using lipgloss for consistent styling and alignment.
func WelcomeLogo() string {
	titleLine := "go-starter-kit"
	subtitle := "create-go-starter"
	tagline := "Scaffold production-ready Go APIs in minutes"

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(ColorSuccess)).
		Background(lipgloss.Color("#003d1f")). // Dark green background for impact
		Padding(0, 2).
		MarginBottom(1)

	taglineStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorInfo)).
		Italic(true).
		MarginBottom(1)

	subtitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorText)).
		Bold(true)

	metaStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorMuted)).
		Faint(true)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		titleStyle.Render(titleLine),
		subtitleStyle.Render(subtitle),
		taglineStyle.Render(tagline),
		metaStyle.Render("Interactive project generator"),
	)

	box := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color(ColorSuccess)).
		Padding(1, 4).
		Align(lipgloss.Center).
		Width(54)

	return box.Render(content)
}

// RenderWelcomeScreen renders the welcome screen with logo and interactive menu.
func RenderWelcomeScreen(width, height int, menuList list.Model, logoOpacity float64) string {
	var b strings.Builder

	if width <= 0 {
		width = 80
	}
	if height <= 0 {
		height = 24
	}
	if menuList.Items() == nil {
		menuList = initializeWelcomeList()
		menuList.SetSize(width-ListWidthMargin, 10)
	}

	// Render logo with fade-in animation
	if logoOpacity > 0.0 {
		logo := WelcomeLogo()
		b.WriteString("\n")
		b.WriteString(CenterHorizontal(logo, width))
		b.WriteString("\n")
	} else {
		b.WriteString("\n\n\n\n\n\n\n")
	}

	// Menu header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(ColorSuccess)).
		MarginTop(1).
		MarginBottom(1)

	headerText := "What would you like to do?"
	b.WriteString(CenterHorizontal(headerStyle.Render(headerText), width))

	// Separator below header
	sep := RenderSeparator(len(headerText) + 4)
	b.WriteString(CenterHorizontal(sep, width))
	b.WriteString("\n\n")

	// Menu list
	// Adjust list width for centering
	menuList.SetWidth(width)
	b.WriteString(menuList.View())
	b.WriteString("\n")

	// Footer
	footer := renderStyledFooter("↑↓", "Navigate", "Enter", "Select", "?", "Help", "Ctrl+C", "Quit")
	b.WriteString("\n")
	b.WriteString(CenterHorizontal(footer, width))
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

package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// FeedbackType represents the type of feedback message.
type FeedbackType int

const (
	FeedbackSuccess FeedbackType = iota
	FeedbackError
	FeedbackWarning
	FeedbackInfo
)

// ErrorModel represents an error message with visual styling and resolution suggestions.
// Story 10.7 Task 5: Composants de confirmation et feedback visuels (AC #6)
type ErrorModel struct {
	// Title is the main error message
	Title string

	// Details provides additional context about the error
	Details string

	// Suggestions are resolution suggestions for the user
	Suggestions []string

	// ShowIcon controls whether to display the error icon
	ShowIcon bool
}

// NewErrorModel creates a new error feedback model.
func NewErrorModel(title string) ErrorModel {
	return ErrorModel{
		Title:       title,
		ShowIcon:    true,
		Suggestions: []string{},
	}
}

// WithDetails adds details to the error model.
func (e ErrorModel) WithDetails(details string) ErrorModel {
	e.Details = details
	return e
}

// WithSuggestions adds resolution suggestions to the error model.
func (e ErrorModel) WithSuggestions(suggestions ...string) ErrorModel {
	e.Suggestions = suggestions
	return e
}

// View renders the error message.
func (e ErrorModel) View() string {
	var b strings.Builder

	// Error icon and title
	if e.ShowIcon {
		icon := "✗"
		iconStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorError)).
			Bold(true)
		b.WriteString(iconStyle.Render(icon))
		b.WriteString(" ")
	}

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorError)).
		Bold(true)
	b.WriteString(titleStyle.Render(e.Title))
	b.WriteString("\n")

	// Details
	if e.Details != "" {
		b.WriteString("\n")
		detailStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorText))
		b.WriteString(detailStyle.Render(e.Details))
		b.WriteString("\n")
	}

	// Resolution suggestions
	if len(e.Suggestions) > 0 {
		b.WriteString("\n")
		suggestionHeader := lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorWarning)).
			Bold(true).
			Render("💡 Suggestions:")
		b.WriteString(suggestionHeader)
		b.WriteString("\n")

		for _, suggestion := range e.Suggestions {
			bullet := lipgloss.NewStyle().
				Foreground(lipgloss.Color(ColorWarning)).
				Render("  •")
			suggestionText := lipgloss.NewStyle().
				Foreground(lipgloss.Color(ColorText)).
				Render(" " + suggestion)
			b.WriteString(bullet)
			b.WriteString(suggestionText)
			b.WriteString("\n")
		}
	}

	// Wrap in error box
	errorBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(ColorError)).
		Padding(1, 2).
		MarginTop(1).
		MarginBottom(1)

	return errorBox.Render(b.String())
}

// SuccessModel represents a success message with visual styling.
// Story 10.7 Task 5: Composants de confirmation et feedback visuels (AC #5)
type SuccessModel struct {
	// Title is the main success message
	Title string

	// Details provides additional context
	Details string

	// Items are additional success items to display
	Items []string

	// ShowIcon controls whether to display the success icon
	ShowIcon bool
}

// NewSuccessModel creates a new success feedback model.
func NewSuccessModel(title string) SuccessModel {
	return SuccessModel{
		Title:    title,
		ShowIcon: true,
		Items:    []string{},
	}
}

// WithDetails adds details to the success model.
func (s SuccessModel) WithDetails(details string) SuccessModel {
	s.Details = details
	return s
}

// WithItems adds success items to display.
func (s SuccessModel) WithItems(items ...string) SuccessModel {
	s.Items = items
	return s
}

// View renders the success message.
func (s SuccessModel) View() string {
	var b strings.Builder

	// Success icon and title
	if s.ShowIcon {
		icon := "✓"
		iconStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorSuccess)).
			Bold(true)
		b.WriteString(iconStyle.Render(icon))
		b.WriteString(" ")
	}

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorSuccess)).
		Bold(true)
	b.WriteString(titleStyle.Render(s.Title))
	b.WriteString("\n")

	// Details
	if s.Details != "" {
		b.WriteString("\n")
		detailStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorText))
		b.WriteString(detailStyle.Render(s.Details))
		b.WriteString("\n")
	}

	// Success items
	if len(s.Items) > 0 {
		b.WriteString("\n")
		for _, item := range s.Items {
			checkmark := lipgloss.NewStyle().
				Foreground(lipgloss.Color(ColorSuccess)).
				Render("  ✓")
			itemText := lipgloss.NewStyle().
				Foreground(lipgloss.Color(ColorText)).
				Render(" " + item)
			b.WriteString(checkmark)
			b.WriteString(itemText)
			b.WriteString("\n")
		}
	}

	// Wrap in success box
	successBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(ColorSuccess)).
		Padding(1, 2).
		MarginTop(1).
		MarginBottom(1)

	return successBox.Render(b.String())
}

// WarningModel represents a warning message with visual styling.
type WarningModel struct {
	// Title is the main warning message
	Title string

	// Details provides additional context
	Details string

	// ShowIcon controls whether to display the warning icon
	ShowIcon bool
}

// NewWarningModel creates a new warning feedback model.
func NewWarningModel(title string) WarningModel {
	return WarningModel{
		Title:    title,
		ShowIcon: true,
	}
}

// WithDetails adds details to the warning model.
func (w WarningModel) WithDetails(details string) WarningModel {
	w.Details = details
	return w
}

// View renders the warning message.
func (w WarningModel) View() string {
	var b strings.Builder

	// Warning icon and title
	if w.ShowIcon {
		icon := "⚠"
		iconStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorWarning)).
			Bold(true)
		b.WriteString(iconStyle.Render(icon))
		b.WriteString(" ")
	}

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorWarning)).
		Bold(true)
	b.WriteString(titleStyle.Render(w.Title))
	b.WriteString("\n")

	// Details
	if w.Details != "" {
		b.WriteString("\n")
		detailStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorText))
		b.WriteString(detailStyle.Render(w.Details))
		b.WriteString("\n")
	}

	// Wrap in warning box
	warningBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(ColorWarning)).
		Padding(1, 2).
		MarginTop(1).
		MarginBottom(1)

	return warningBox.Render(b.String())
}

// InfoModel represents an informational message with visual styling.
type InfoModel struct {
	// Title is the main info message
	Title string

	// Details provides additional context
	Details string

	// ShowIcon controls whether to display the info icon
	ShowIcon bool
}

// NewInfoModel creates a new info feedback model.
func NewInfoModel(title string) InfoModel {
	return InfoModel{
		Title:    title,
		ShowIcon: true,
	}
}

// WithDetails adds details to the info model.
func (i InfoModel) WithDetails(details string) InfoModel {
	i.Details = details
	return i
}

// View renders the info message.
func (i InfoModel) View() string {
	var b strings.Builder

	// Info icon and title
	if i.ShowIcon {
		icon := "ℹ"
		iconStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorInfo)).
			Bold(true)
		b.WriteString(iconStyle.Render(icon))
		b.WriteString(" ")
	}

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorInfo)).
		Bold(true)
	b.WriteString(titleStyle.Render(i.Title))
	b.WriteString("\n")

	// Details
	if i.Details != "" {
		b.WriteString("\n")
		detailStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorText))
		b.WriteString(detailStyle.Render(i.Details))
		b.WriteString("\n")
	}

	// Wrap in info box
	infoBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(ColorInfo)).
		Padding(1, 2).
		MarginTop(1).
		MarginBottom(1)

	return infoBox.Render(b.String())
}

// Common error suggestions for typical scenarios
var (
	// ProjectNameSuggestions provides suggestions for invalid project names
	ProjectNameSuggestions = []string{
		"Use lowercase letters, numbers, and hyphens only",
		"Start with a letter",
		"Avoid spaces and special characters",
		"Example: my-awesome-api",
	}

	// DirectoryExistsSuggestions provides suggestions when a directory already exists
	DirectoryExistsSuggestions = []string{
		"Choose a different project name",
		"Remove or rename the existing directory",
		"Use --force flag to overwrite (if available)",
	}

	// GitInitSuggestions provides suggestions when git initialization fails
	GitInitSuggestions = []string{
		"Ensure Git is installed: git --version",
		"Check if you have write permissions in the current directory",
		"Try running the command again",
	}

	// DependencySuggestions provides suggestions for dependency installation errors
	DependencySuggestions = []string{
		"Ensure Go is installed: go version",
		"Check your internet connection",
		"Verify GOPATH and GOROOT are set correctly",
		"Try running: go mod tidy",
	}

	// DatabaseSuggestions provides suggestions for database connection errors
	DatabaseSuggestions = []string{
		"Ensure PostgreSQL is running",
		"Check DATABASE_URL in .env file",
		"Verify database credentials",
		"Try: docker ps (if using Docker)",
	}
)

// Note: Helper functions for rendering complete feedback boxes
// are not provided here to avoid conflicts with existing simple text renderers
// in styles.go. Instead, use the models directly:
//
// Example:
//   err := NewErrorModel("Invalid input").
//           WithDetails("Project name contains invalid characters").
//           WithSuggestions(ProjectNameSuggestions...)
//   fmt.Println(err.View())

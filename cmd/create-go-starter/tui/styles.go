package tui

import (
	"os"

	"github.com/charmbracelet/lipgloss"
)

// Color palette for go-starter-kit (cohérente avec le branding actuel)
const (
	ColorSuccess = "#00c853" // Green - Success/Primary (used in Green())
	ColorInfo    = "#00b0ff" // Blue - Info/Focus
	ColorWarning = "#ff6d00" // Orange - Warning (used in Yellow())
	ColorError   = "#ff1744" // Red - Error (used in Red())
	ColorText    = "#ffffff" // White - Normal text
	ColorMuted   = "#666666" // Gray - Muted/secondary text
	ColorBorder  = "#444444" // Dark gray - Borders
)

var (
	// Header style for titles and section headers
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(ColorSuccess)).
			MarginTop(1).
			MarginBottom(1).
			PaddingLeft(1)

	// Success style for positive messages and checkmarks
	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorSuccess)).
			Bold(true)

	// Info style for informational messages
	InfoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorInfo))

	// Warning style for warnings
	WarningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorWarning)).
			Bold(true)

	// Error style for errors
	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorError)).
			Bold(true)

	// Muted style for secondary text
	MutedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorMuted))

	// Focused style for active/selected items
	FocusedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorInfo)).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorInfo)).
			Padding(0, 1)

	// Blurred style for inactive items
	BlurredStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorMuted))

	// Box style for containers
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorBorder)).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1)

	// Title box style for prominent titles
	TitleBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color(ColorSuccess)).
			Padding(0, 2).
			Bold(true).
			Align(lipgloss.Center)

	// Footer style for help text and instructions
	FooterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorMuted)).
			MarginTop(1).
			Italic(true)

	// Highlight style for emphasized text
	HighlightStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorInfo)).
			Bold(true)

	// SubHeader style for section subtitles
	SubHeaderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorInfo)).
			Bold(false).
			MarginBottom(1)

	// SuccessBox style for success messages with border
	SuccessBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorSuccess)).
			Foreground(lipgloss.Color(ColorSuccess)).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1)

	// ErrorBox style for error messages with border
	ErrorBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorError)).
			Foreground(lipgloss.Color(ColorError)).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1)

	// WarningBox style for warning messages with border
	WarningBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorWarning)).
			Foreground(lipgloss.Color(ColorWarning)).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1)

	// InfoBox style for info messages with border
	InfoBox = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(ColorInfo)).
		Foreground(lipgloss.Color(ColorInfo)).
		Padding(1, 2).
		MarginTop(1).
		MarginBottom(1)

	// GradientHeaderStyle for headers with visual emphasis
	GradientHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color(ColorSuccess)).
				Background(lipgloss.Color("#003d1f")). // Dark green background
				Padding(0, 2).
				MarginTop(1).
				MarginBottom(1)

	// ProgressBarFilledStyle for filled portion of progress bars
	ProgressBarFilledStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(ColorSuccess)).
				Background(lipgloss.Color(ColorSuccess))

	// ProgressBarEmptyStyle for empty portion of progress bars
	ProgressBarEmptyStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(ColorMuted)).
				Background(lipgloss.Color("#222222"))

	// DimmedStyle for very subtle text
	DimmedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#444444"))

	// KeyStyle for keyboard shortcuts (e.g., "↑↓")
	KeyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorInfo)).
			Bold(true).
			Background(lipgloss.Color("#1a1a1a")).
			Padding(0, 1)

	// DescStyle for descriptions next to keys
	DescStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorText))

	// SeparatorStyle for visual separators
	SeparatorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorBorder))

	// ListItemStyle for list items (unselected)
	ListItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorText)).
			PaddingLeft(2)

	// SelectedListItemStyle for selected list items
	SelectedListItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(ColorSuccess)).
				Bold(true).
				PaddingLeft(1).
				BorderLeft(true).
				BorderStyle(lipgloss.ThickBorder()).
				BorderForeground(lipgloss.Color(ColorSuccess))

	// ProgressTextStyle for text beside progress bars
	ProgressTextStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(ColorInfo))

	// SpinnerStyle for spinner text
	SpinnerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorInfo))

	// CodeStyle for inline code snippets
	CodeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorInfo)).
			Background(lipgloss.Color("#1a1a1a")).
			Padding(0, 1)

	// KeywordStyle for syntax highlighting (Go keywords)
	KeywordStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorInfo)).
			Bold(true)

	// StringStyle for syntax highlighting (strings)
	StringStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorSuccess))

	// CommentStyle for syntax highlighting (comments)
	CommentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorMuted)).
			Italic(true)

	// BreadcrumbStyle for navigation breadcrumbs
	BreadcrumbStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorMuted))

	// BreadcrumbActiveStyle for active breadcrumb item
	BreadcrumbActiveStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(ColorSuccess)).
				Bold(true)
)

// RenderHeader renders a styled header with the given text.
func RenderHeader(text string) string {
	return HeaderStyle.Render(text)
}

// RenderMuted renders text in muted style.
func RenderMuted(text string) string {
	return MutedStyle.Render(text)
}

// Simple wrapper functions that delegate to feedback.go functions
// These provide backwards compatibility for simple single-parameter calls

// RenderSuccess renders text in success style (simple wrapper).
func RenderSuccess(text string) string {
	return SuccessStyle.Render(text)
}

// RenderInfo renders text in info style (simple wrapper).
func RenderInfo(text string) string {
	return InfoStyle.Render(text)
}

// RenderWarning renders text in warning style (simple wrapper).
func RenderWarning(text string) string {
	return WarningStyle.Render(text)
}

// RenderError renders text in error style (simple wrapper).
func RenderError(text string) string {
	return ErrorStyle.Render(text)
}

// RenderSimpleInfo is an alias for RenderInfo (backward compatibility).
func RenderSimpleInfo(text string) string {
	return RenderInfo(text)
}

// RenderSimpleError is an alias for RenderError (backward compatibility).
func RenderSimpleError(text string) string {
	return RenderError(text)
}

// RenderBox renders text in a box.
func RenderBox(text string) string {
	return BoxStyle.Render(text)
}

// RenderTitleBox renders text in a prominent title box.
func RenderTitleBox(text string, width int) string {
	return TitleBoxStyle.Width(width).Render(text)
}

// RenderFooter renders footer help text.
func RenderFooter(text string) string {
	return FooterStyle.Render(text)
}

// RenderHighlight renders highlighted text.
func RenderHighlight(text string) string {
	return HighlightStyle.Render(text)
}

// Note: RenderSuccessBox, RenderErrorBox, RenderWarningBox, RenderInfoBox
// are defined in feedback.go with richer signatures.
// For simple box rendering, use the style directly: SuccessBox.Render(text)

// RenderGradientHeader renders a header with gradient background.
func RenderGradientHeader(text string, width int) string {
	return GradientHeaderStyle.Width(width).Align(lipgloss.Center).Render(text)
}

// RenderCode renders inline code snippet.
func RenderCode(text string) string {
	return CodeStyle.Render(text)
}

// RenderKey renders a keyboard shortcut.
func RenderKey(key string) string {
	return KeyStyle.Render(key)
}

// RenderBreadcrumb renders a breadcrumb navigation item.
func RenderBreadcrumb(text string, active bool) string {
	if active {
		return BreadcrumbActiveStyle.Render(text)
	}
	return BreadcrumbStyle.Render(text)
}

// RenderSeparator renders a horizontal separator line.
func RenderSeparator(width int) string {
	separator := ""
	for i := 0; i < width; i++ {
		separator += "─"
	}
	return SeparatorStyle.Render(separator)
}

// RenderProgressBar renders a simple progress bar.
// percent should be between 0.0 and 1.0
func RenderProgressBar(percent float64, width int) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 1 {
		percent = 1
	}

	filledWidth := int(float64(width) * percent)
	emptyWidth := width - filledWidth

	filled := ""
	for i := 0; i < filledWidth; i++ {
		filled += "█"
	}

	empty := ""
	for i := 0; i < emptyWidth; i++ {
		empty += "░"
	}

	return ProgressBarFilledStyle.Render(filled) + ProgressBarEmptyStyle.Render(empty)
}

// RenderListItem renders a list item (selected or not).
func RenderListItem(text string, selected bool) string {
	if selected {
		return SelectedListItemStyle.Render("▸ " + text)
	}
	return ListItemStyle.Render("  " + text)
}

// IsNoColorMode checks if NO_COLOR environment variable is set.
// When true, styles should fall back to plain text.
func IsNoColorMode() bool {
	return os.Getenv("NO_COLOR") != ""
}

// AdaptToTerminalWidth adjusts styles based on terminal width.
// For narrow terminals, reduce padding and margins.
func AdaptToTerminalWidth(width int) {
	if width < 60 {
		// Narrow terminal - reduce padding
		BoxStyle = BoxStyle.Padding(0, 1)
		TitleBoxStyle = TitleBoxStyle.Padding(0, 1)
		HeaderStyle = HeaderStyle.PaddingLeft(0)
	} else if width < 100 {
		// Medium terminal - moderate padding
		BoxStyle = BoxStyle.Padding(1, 1)
		TitleBoxStyle = TitleBoxStyle.Padding(0, 2)
	} else {
		// Wide terminal - comfortable padding
		BoxStyle = BoxStyle.Padding(1, 2)
		TitleBoxStyle = TitleBoxStyle.Padding(0, 3)
	}
}

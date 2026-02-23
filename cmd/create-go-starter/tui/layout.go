package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Terminal width breakpoints for responsive layouts
const (
	WidthNarrow  = 80
	WidthMedium  = 120
	WidthWide    = 160
	MinWidth     = 60
	DefaultWidth = 100
)

// AdaptiveLayout returns a style with width adapted to terminal width.
// Follows responsive breakpoints: narrow (< 80), medium (< 120), wide (>= 120).
func AdaptiveLayout(termWidth int) lipgloss.Style {
	if termWidth < MinWidth {
		// Very narrow terminal - minimal width
		return lipgloss.NewStyle().Width(termWidth - 4)
	} else if termWidth < WidthNarrow {
		// Narrow terminal (60-79 cols)
		return lipgloss.NewStyle().Width(termWidth - 4)
	} else if termWidth < WidthMedium {
		// Medium terminal (80-119 cols)
		return lipgloss.NewStyle().Width(WidthNarrow)
	} else if termWidth < WidthWide {
		// Wide terminal (120-159 cols)
		return lipgloss.NewStyle().Width(DefaultWidth)
	} else {
		// Extra wide terminal (>= 160 cols)
		return lipgloss.NewStyle().Width(WidthMedium)
	}
}

// CenterHorizontal centers text horizontally within the given width.
func CenterHorizontal(text string, width int) string {
	return lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Render(text)
}

// CenterVertical centers text vertically within the given height.
func CenterVertical(text string, height int) string {
	lines := strings.Split(text, "\n")
	linesCount := len(lines)
	if linesCount >= height {
		return text
	}

	paddingTop := (height - linesCount) / 2
	paddingBottom := height - linesCount - paddingTop

	var result strings.Builder
	for i := 0; i < paddingTop; i++ {
		result.WriteString("\n")
	}
	result.WriteString(text)
	for i := 0; i < paddingBottom; i++ {
		result.WriteString("\n")
	}

	return result.String()
}

// Center centers text both horizontally and vertically.
func Center(text string, width, height int) string {
	horizontal := CenterHorizontal(text, width)
	return CenterVertical(horizontal, height)
}

// AlignRight aligns text to the right within the given width.
func AlignRight(text string, width int) string {
	return lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Right).
		Render(text)
}

// AlignLeft aligns text to the left within the given width.
func AlignLeft(text string, width int) string {
	return lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Left).
		Render(text)
}

// TwoColumnLayout creates a two-column layout with the given content.
// leftRatio is the percentage of width for the left column (0.0 to 1.0).
func TwoColumnLayout(left, right string, totalWidth int, leftRatio float64) string {
	if leftRatio <= 0 || leftRatio >= 1 {
		leftRatio = 0.5 // Default to 50/50
	}

	leftWidth := int(float64(totalWidth) * leftRatio)
	rightWidth := totalWidth - leftWidth

	leftStyle := lipgloss.NewStyle().Width(leftWidth)
	rightStyle := lipgloss.NewStyle().Width(rightWidth)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftStyle.Render(left),
		rightStyle.Render(right),
	)
}

// ThreeColumnLayout creates a three-column layout with equal widths.
func ThreeColumnLayout(left, center, right string, totalWidth int) string {
	colWidth := totalWidth / 3

	leftStyle := lipgloss.NewStyle().Width(colWidth).Align(lipgloss.Left)
	centerStyle := lipgloss.NewStyle().Width(colWidth).Align(lipgloss.Center)
	rightStyle := lipgloss.NewStyle().Width(colWidth).Align(lipgloss.Right)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftStyle.Render(left),
		centerStyle.Render(center),
		rightStyle.Render(right),
	)
}

// PadContent adds padding around content based on terminal width.
func PadContent(content string, termWidth int) string {
	if termWidth < MinWidth {
		// No padding for very narrow terminals
		return content
	} else if termWidth < WidthNarrow {
		// Minimal padding for narrow terminals
		return lipgloss.NewStyle().Padding(0, 1).Render(content)
	} else if termWidth < WidthMedium {
		// Moderate padding for medium terminals
		return lipgloss.NewStyle().Padding(1, 2).Render(content)
	} else {
		// Comfortable padding for wide terminals
		return lipgloss.NewStyle().Padding(1, 3).Render(content)
	}
}

// WrapText wraps long text to fit within the given width.
func WrapText(text string, width int) string {
	if width <= 0 {
		return text
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	var lines []string
	var currentLine strings.Builder

	for _, word := range words {
		if currentLine.Len() == 0 {
			currentLine.WriteString(word)
		} else if currentLine.Len()+1+len(word) <= width {
			currentLine.WriteString(" ")
			currentLine.WriteString(word)
		} else {
			lines = append(lines, currentLine.String())
			currentLine.Reset()
			currentLine.WriteString(word)
		}
	}

	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return strings.Join(lines, "\n")
}

// BoxWithTitle creates a box with a title at the top.
func BoxWithTitle(title, content string, width int) string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(ColorSuccess)).
		Width(width - 4). // Account for border padding
		Align(lipgloss.Center)

	contentStyle := lipgloss.NewStyle().
		Width(width - 4).
		Align(lipgloss.Left)

	boxContent := lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render(title),
		RenderSeparator(width-4),
		contentStyle.Render(content),
	)

	return BoxStyle.Width(width).Render(boxContent)
}

// StackVertical stacks multiple strings vertically with spacing.
func StackVertical(spacing int, elements ...string) string {
	if len(elements) == 0 {
		return ""
	}

	var result strings.Builder
	for i, elem := range elements {
		result.WriteString(elem)
		if i < len(elements)-1 {
			for j := 0; j < spacing; j++ {
				result.WriteString("\n")
			}
		}
	}

	return result.String()
}

// GetContentWidth returns the recommended content width based on terminal width.
func GetContentWidth(termWidth int) int {
	if termWidth < MinWidth {
		return termWidth - 4
	} else if termWidth < WidthNarrow {
		return termWidth - 8
	} else if termWidth < WidthMedium {
		return WidthNarrow - 8
	} else if termWidth < WidthWide {
		return DefaultWidth
	} else {
		return WidthMedium
	}
}

// ResponsivePadding returns padding values based on terminal width.
// Returns (vertical, horizontal) padding.
func ResponsivePadding(termWidth int) (int, int) {
	if termWidth < MinWidth {
		return 0, 0
	} else if termWidth < WidthNarrow {
		return 0, 1
	} else if termWidth < WidthMedium {
		return 1, 2
	} else {
		return 1, 3
	}
}

// MaxWidth limits the content to a maximum width while keeping it centered.
func MaxWidth(content string, maxWidth, termWidth int) string {
	if termWidth <= maxWidth {
		return content
	}

	style := lipgloss.NewStyle().
		Width(maxWidth).
		Align(lipgloss.Center)

	centered := CenterHorizontal(style.Render(content), termWidth)
	return centered
}

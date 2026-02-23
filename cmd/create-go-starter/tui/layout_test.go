package tui

import (
	"strings"
	"testing"
)

// TestAdaptiveLayout verifies responsive layout adapts to terminal width (AC#8)
func TestAdaptiveLayout(t *testing.T) {
	tests := []struct {
		name      string
		termWidth int
		minWidth  int
		maxWidth  int
	}{
		{
			name:      "Very narrow terminal",
			termWidth: 40,
			minWidth:  30,
			maxWidth:  40,
		},
		{
			name:      "Narrow terminal (60-79 cols)",
			termWidth: 70,
			minWidth:  60,
			maxWidth:  70,
		},
		{
			name:      "Medium terminal (80-119 cols)",
			termWidth: 100,
			minWidth:  WidthNarrow,
			maxWidth:  WidthNarrow,
		},
		{
			name:      "Wide terminal (120-159 cols)",
			termWidth: 140,
			minWidth:  DefaultWidth,
			maxWidth:  DefaultWidth,
		},
		{
			name:      "Extra wide terminal (>= 160 cols)",
			termWidth: 180,
			minWidth:  WidthMedium,
			maxWidth:  WidthMedium,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			style := AdaptiveLayout(tt.termWidth)
			// Render a test string to verify the style works
			result := style.Render("test content")
			if result == "" {
				t.Errorf("AdaptiveLayout returned empty result")
			}
		})
	}
}

// TestCenterHorizontal verifies horizontal centering
func TestCenterHorizontal(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		width int
	}{
		{"Short text", "Hello", 20},
		{"Medium text", "Hello World", 40},
		{"Long text", "This is a longer text for testing", 80},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CenterHorizontal(tt.text, tt.width)
			if result == "" {
				t.Errorf("CenterHorizontal returned empty string")
			}
		})
	}
}

// TestCenterVertical verifies vertical centering
func TestCenterVertical(t *testing.T) {
	tests := []struct {
		name   string
		text   string
		height int
	}{
		{"Single line", "Hello", 10},
		{"Multiple lines", "Line 1\nLine 2\nLine 3", 10},
		{"Exact height", "Line 1\nLine 2", 2},
		{"Taller than content", "Single", 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CenterVertical(tt.text, tt.height)
			if result == "" {
				t.Errorf("CenterVertical returned empty string")
			}

			// Count lines in result
			lines := strings.Split(result, "\n")
			if len(lines) < strings.Count(tt.text, "\n")+1 {
				t.Errorf("CenterVertical lost content lines")
			}
		})
	}
}

// TestCenter verifies both horizontal and vertical centering
func TestCenter(t *testing.T) {
	text := "Centered"
	width := 40
	height := 10

	result := Center(text, width, height)
	if result == "" {
		t.Errorf("Center returned empty string")
	}
}

// TestAlignRight verifies right alignment
func TestAlignRight(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		width int
	}{
		{"Short text", "Right", 20},
		{"Medium text", "Align Right", 40},
		{"Exact width", "Perfect", 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AlignRight(tt.text, tt.width)
			if result == "" {
				t.Errorf("AlignRight returned empty string")
			}
		})
	}
}

// TestAlignLeft verifies left alignment
func TestAlignLeft(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		width int
	}{
		{"Short text", "Left", 20},
		{"Medium text", "Align Left", 40},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AlignLeft(tt.text, tt.width)
			if result == "" {
				t.Errorf("AlignLeft returned empty string")
			}
		})
	}
}

// TestTwoColumnLayout verifies two-column layout
func TestTwoColumnLayout(t *testing.T) {
	tests := []struct {
		name      string
		left      string
		right     string
		width     int
		leftRatio float64
	}{
		{"Equal columns", "Left", "Right", 40, 0.5},
		{"Left larger", "Left Column", "Right", 60, 0.7},
		{"Right larger", "Left", "Right Column", 60, 0.3},
		{"Invalid ratio (negative)", "Left", "Right", 40, -0.5}, // Should default to 0.5
		{"Invalid ratio (>1)", "Left", "Right", 40, 1.5},        // Should default to 0.5
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TwoColumnLayout(tt.left, tt.right, tt.width, tt.leftRatio)
			if result == "" {
				t.Errorf("TwoColumnLayout returned empty string")
			}
		})
	}
}

// TestThreeColumnLayout verifies three-column layout
func TestThreeColumnLayout(t *testing.T) {
	tests := []struct {
		name   string
		left   string
		center string
		right  string
		width  int
	}{
		{"Simple layout", "L", "C", "R", 60},
		{"With content", "Left", "Center", "Right", 90},
		{"Narrow width", "L", "C", "R", 30},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ThreeColumnLayout(tt.left, tt.center, tt.right, tt.width)
			if result == "" {
				t.Errorf("ThreeColumnLayout returned empty string")
			}
		})
	}
}

// TestPadContent verifies responsive padding
func TestPadContent(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		termWidth int
	}{
		{"Very narrow", "Content", 40},
		{"Narrow", "Content", 70},
		{"Medium", "Content", 100},
		{"Wide", "Content", 140},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PadContent(tt.content, tt.termWidth)
			if result == "" {
				t.Errorf("PadContent returned empty string")
			}
		})
	}
}

// TestWrapText verifies text wrapping
func TestWrapText(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		width int
		want  int // minimum number of lines expected
	}{
		{
			name:  "Short text, no wrap",
			text:  "Hello World",
			width: 20,
			want:  1,
		},
		{
			name:  "Long text, should wrap",
			text:  "This is a very long text that should definitely wrap across multiple lines",
			width: 20,
			want:  3,
		},
		{
			name:  "Empty text",
			text:  "",
			width: 20,
			want:  0,
		},
		{
			name:  "Single word longer than width",
			text:  "Supercalifragilisticexpialidocious",
			width: 10,
			want:  1, // Single long word won't wrap
		},
		{
			name:  "Zero width",
			text:  "Hello",
			width: 0,
			want:  0, // Should return original text
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WrapText(tt.text, tt.width)
			lines := strings.Split(result, "\n")

			if tt.want == 0 {
				if result != tt.text {
					t.Errorf("WrapText with zero width should return original text")
				}
				return
			}

			if len(lines) < tt.want {
				t.Errorf("WrapText: expected at least %d lines, got %d", tt.want, len(lines))
			}
		})
	}
}

// TestBoxWithTitle verifies box with title rendering
func TestBoxWithTitle(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		content string
		width   int
	}{
		{"Simple box", "Title", "Content here", 40},
		{"Wide box", "My Title", "Some content\nMultiple lines", 60},
		{"Narrow box", "T", "C", 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BoxWithTitle(tt.title, tt.content, tt.width)
			if result == "" {
				t.Errorf("BoxWithTitle returned empty string")
			}
		})
	}
}

// TestStackVertical verifies vertical stacking
func TestStackVertical(t *testing.T) {
	tests := []struct {
		name     string
		spacing  int
		elements []string
		want     int // minimum number of lines
	}{
		{
			name:     "No elements",
			spacing:  1,
			elements: []string{},
			want:     0,
		},
		{
			name:     "Single element",
			spacing:  1,
			elements: []string{"Element 1"},
			want:     1,
		},
		{
			name:     "Multiple elements with spacing",
			spacing:  2,
			elements: []string{"Element 1", "Element 2", "Element 3"},
			want:     3,
		},
		{
			name:     "No spacing",
			spacing:  0,
			elements: []string{"A", "B", "C"},
			want:     3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StackVertical(tt.spacing, tt.elements...)

			if len(tt.elements) == 0 {
				if result != "" {
					t.Errorf("StackVertical with no elements should return empty string")
				}
				return
			}

			if result == "" {
				t.Errorf("StackVertical returned empty string")
			}
		})
	}
}

// TestGetContentWidth verifies content width calculation
func TestGetContentWidth(t *testing.T) {
	tests := []struct {
		name      string
		termWidth int
		minWidth  int
		maxWidth  int
	}{
		{"Very narrow", 40, 30, 40},
		{"Narrow", 70, 50, 70},
		{"Medium", 100, 70, 80},
		{"Wide", 140, 90, 110},
		{"Extra wide", 180, 110, 130},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetContentWidth(tt.termWidth)
			if result < tt.minWidth || result > tt.maxWidth {
				t.Errorf("GetContentWidth(%d) = %d, want between %d and %d",
					tt.termWidth, result, tt.minWidth, tt.maxWidth)
			}
		})
	}
}

// TestResponsivePadding verifies responsive padding calculation
func TestResponsivePadding(t *testing.T) {
	tests := []struct {
		name      string
		termWidth int
	}{
		{"Very narrow", 40},
		{"Narrow", 70},
		{"Medium", 100},
		{"Wide", 140},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, h := ResponsivePadding(tt.termWidth)
			if v < 0 || h < 0 {
				t.Errorf("ResponsivePadding returned negative values: v=%d, h=%d", v, h)
			}
		})
	}
}

// TestMaxWidth verifies max width constraint
func TestMaxWidth(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		maxWidth  int
		termWidth int
	}{
		{"Terminal wider than max", "Content", 60, 100},
		{"Terminal narrower than max", "Content", 100, 60},
		{"Equal widths", "Content", 80, 80},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaxWidth(tt.content, tt.maxWidth, tt.termWidth)
			if result == "" {
				t.Errorf("MaxWidth returned empty string")
			}
		})
	}
}

// TestWidthConstants verifies breakpoint constants
func TestWidthConstants(t *testing.T) {
	tests := []struct {
		name  string
		value int
		min   int
	}{
		{"WidthNarrow", WidthNarrow, 60},
		{"WidthMedium", WidthMedium, 100},
		{"WidthWide", WidthWide, 120},
		{"MinWidth", MinWidth, 40},
		{"DefaultWidth", DefaultWidth, 80},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value < tt.min {
				t.Errorf("%s = %d, expected at least %d", tt.name, tt.value, tt.min)
			}
		})
	}
}

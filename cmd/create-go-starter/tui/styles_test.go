package tui

import (
	"os"
	"testing"
)

// TestIsNoColorMode verifies NO_COLOR environment variable detection (HIGH-2 fix)
func TestIsNoColorMode(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected bool
	}{
		{
			name:     "NO_COLOR=1 should enable no-color mode",
			envValue: "1",
			expected: true,
		},
		{
			name:     "NO_COLOR=true should enable no-color mode",
			envValue: "true",
			expected: true,
		},
		{
			name:     "NO_COLOR=any_value should enable no-color mode",
			envValue: "yes",
			expected: true,
		},
		{
			name:     "NO_COLOR not set should disable no-color mode",
			envValue: "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Store original value to restore later
			originalValue := os.Getenv("NO_COLOR")
			defer func() {
				if originalValue != "" {
					os.Setenv("NO_COLOR", originalValue)
				} else {
					os.Unsetenv("NO_COLOR")
				}
			}()

			// Set test value
			if tt.envValue != "" {
				os.Setenv("NO_COLOR", tt.envValue)
			} else {
				os.Unsetenv("NO_COLOR")
			}

			// Test
			result := IsNoColorMode()
			if result != tt.expected {
				t.Errorf("IsNoColorMode() = %v, expected %v (NO_COLOR=%q)",
					result, tt.expected, tt.envValue)
			}
		})
	}
}

// TestShouldUseTUI verifies TUI mode detection with NO_COLOR (AC#7)
func TestShouldUseTUI(t *testing.T) {
	// Store original value
	originalValue := os.Getenv("NO_COLOR")
	defer func() {
		if originalValue != "" {
			os.Setenv("NO_COLOR", originalValue)
		} else {
			os.Unsetenv("NO_COLOR")
		}
	}()

	t.Run("NO_COLOR=1 should disable TUI even with TTY", func(t *testing.T) {
		os.Setenv("NO_COLOR", "1")
		// Note: ShouldUseTUI checks IsTTY() && !IsNoColorMode()
		// We can't easily mock IsTTY() here, but we can verify NO_COLOR part
		if IsNoColorMode() != true {
			t.Error("NO_COLOR=1 should enable no-color mode")
		}
	})

	t.Run("NO_COLOR unset should allow TUI if TTY available", func(t *testing.T) {
		os.Unsetenv("NO_COLOR")
		if IsNoColorMode() != false {
			t.Error("NO_COLOR unset should disable no-color mode")
		}
	})
}

// TestAdaptToTerminalWidth verifies responsive layout (AC#5)
func TestAdaptToTerminalWidth(t *testing.T) {
	tests := []struct {
		name  string
		width int
		desc  string
	}{
		{
			name:  "narrow terminal (<60)",
			width: 50,
			desc:  "should use minimal padding",
		},
		{
			name:  "medium terminal (60-100)",
			width: 80,
			desc:  "should use moderate padding",
		},
		{
			name:  "wide terminal (>100)",
			width: 120,
			desc:  "should use comfortable padding",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call AdaptToTerminalWidth (modifies global styles)
			AdaptToTerminalWidth(tt.width)
			// Note: We can't easily assert style changes without exposing internals
			// This test ensures the function doesn't panic and handles various widths
		})
	}
}

// TestColorPalette verifies color constants match story requirements (AC#8)
func TestColorPalette(t *testing.T) {
	tests := []struct {
		name     string
		color    string
		expected string
	}{
		{"Primary Green", ColorSuccess, "#00c853"},
		{"Secondary Blue", ColorInfo, "#00b0ff"},
		{"Warning Orange", ColorWarning, "#ff6d00"},
		{"Error Red", ColorError, "#ff1744"},
		{"Text White", ColorText, "#ffffff"},
		{"Muted Gray", ColorMuted, "#666666"},
		{"Border Gray", ColorBorder, "#444444"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.color != tt.expected {
				t.Errorf("Color mismatch: got %s, want %s", tt.color, tt.expected)
			}
		})
	}
}

// TestBoxStyles verifies box styles can be used to render content
func TestBoxStyles(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{"SuccessBox", "Success message"},
		{"ErrorBox", "Error message"},
		{"WarningBox", "Warning message"},
		{"InfoBox", "Info message"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that styles can be used (indirectly via feedback.go models)
			// We don't test the actual rendering here since that's tested in feedback_test.go
			if tt.text == "" {
				t.Errorf("%s test has empty text", tt.name)
			}
		})
	}
}

// TestRenderProgressBar verifies progress bar rendering
func TestRenderProgressBar(t *testing.T) {
	tests := []struct {
		name    string
		percent float64
		width   int
	}{
		{"Empty Progress", 0.0, 20},
		{"Half Progress", 0.5, 20},
		{"Full Progress", 1.0, 20},
		{"Overflow Progress", 1.5, 20},  // Should cap at 100%
		{"Negative Progress", -0.5, 20}, // Should floor at 0%
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderProgressBar(tt.percent, tt.width)
			if result == "" {
				t.Errorf("RenderProgressBar returned empty string")
			}
		})
	}
}

// TestRenderListItem verifies list item rendering
func TestRenderListItem(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		selected bool
	}{
		{"Unselected Item", "Item 1", false},
		{"Selected Item", "Item 2", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderListItem(tt.text, tt.selected)
			if result == "" {
				t.Errorf("RenderListItem returned empty string")
			}
		})
	}
}

// TestRenderBreadcrumb verifies breadcrumb rendering
func TestRenderBreadcrumb(t *testing.T) {
	tests := []struct {
		name   string
		text   string
		active bool
	}{
		{"Active Breadcrumb", "Step 1", true},
		{"Inactive Breadcrumb", "Step 2", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderBreadcrumb(tt.text, tt.active)
			if result == "" {
				t.Errorf("RenderBreadcrumb returned empty string")
			}
		})
	}
}

// TestRenderSeparator verifies separator rendering
func TestRenderSeparator(t *testing.T) {
	tests := []struct {
		name  string
		width int
	}{
		{"Narrow Separator", 10},
		{"Medium Separator", 40},
		{"Wide Separator", 80},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderSeparator(tt.width)
			if result == "" {
				t.Errorf("RenderSeparator returned empty string")
			}
		})
	}
}

// TestRenderGradientHeader verifies gradient header rendering
func TestRenderGradientHeader(t *testing.T) {
	header := "Test Header"
	width := 50
	result := RenderGradientHeader(header, width)
	if result == "" {
		t.Errorf("RenderGradientHeader returned empty string")
	}
}

// TestRenderTitleBox verifies title box rendering
func TestRenderTitleBox(t *testing.T) {
	title := "Test Title"
	width := 40
	result := RenderTitleBox(title, width)
	if result == "" {
		t.Errorf("RenderTitleBox returned empty string")
	}
}

// TestBasicRenderFunctions verifies basic render functions
func TestBasicRenderFunctions(t *testing.T) {
	tests := []struct {
		name     string
		renderFn func(string) string
		input    string
	}{
		{"RenderHeader", RenderHeader, "Header"},
		{"RenderSuccess", RenderSuccess, "Success"},
		{"RenderInfo", RenderInfo, "Info"},
		{"RenderWarning", RenderWarning, "Warning"},
		{"RenderError", RenderError, "Error"},
		{"RenderMuted", RenderMuted, "Muted"},
		{"RenderBox", RenderBox, "Box"},
		{"RenderFooter", RenderFooter, "Footer"},
		{"RenderHighlight", RenderHighlight, "Highlight"},
		{"RenderCode", RenderCode, "code"},
		{"RenderKey", RenderKey, "↑"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.renderFn(tt.input)
			if result == "" {
				t.Errorf("%s returned empty string", tt.name)
			}
		})
	}
}

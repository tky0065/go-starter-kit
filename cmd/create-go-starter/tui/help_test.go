package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/key"
)

// TestDefaultKeyMap verifies that all keyboard shortcuts are properly defined.
func TestDefaultKeyMap(t *testing.T) {
	tests := []struct {
		name    string
		binding key.Binding
		wantKey string
	}{
		{"Up", DefaultKeyMap.Up, "↑/k"},
		{"Down", DefaultKeyMap.Down, "↓/j"},
		{"Left", DefaultKeyMap.Left, "←/h"},
		{"Right", DefaultKeyMap.Right, "→/l"},
		{"Back", DefaultKeyMap.Back, "b"},
		{"Next", DefaultKeyMap.Next, "n"},
		{"Help", DefaultKeyMap.Help, "?"},
		{"Quit", DefaultKeyMap.Quit, "Ctrl+C"},
		{"Enter", DefaultKeyMap.Enter, "Enter"},
		{"Space", DefaultKeyMap.Space, "Space"},
		{"PgUp", DefaultKeyMap.PgUp, "PgUp"},
		{"PgDown", DefaultKeyMap.PgDown, "PgDn"},
		{"Home", DefaultKeyMap.Home, "Home"},
		{"End", DefaultKeyMap.End, "End"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			help := tt.binding.Help()
			if !strings.Contains(help.Key, tt.wantKey) {
				t.Errorf("DefaultKeyMap.%s help key = %q, want to contain %q", tt.name, help.Key, tt.wantKey)
			}
		})
	}
}

// TestViewHelp_DisplaysCorrectHeader verifies that the help screen shows the correct header.
func TestViewHelp_DisplaysCorrectHeader(t *testing.T) {
	m := NewModel()
	m.state = StateHelp
	m.previousState = StateWelcome
	m.width = 80

	output := m.viewHelp()

	if !strings.Contains(output, "Keyboard Shortcuts") {
		t.Error("viewHelp() should contain 'Keyboard Shortcuts' header")
	}

	if !strings.Contains(output, "Global Shortcuts") {
		t.Error("viewHelp() should contain 'Global Shortcuts' section")
	}
}

// TestViewHelp_DisplaysGlobalShortcuts verifies that global shortcuts are always shown.
func TestViewHelp_DisplaysGlobalShortcuts(t *testing.T) {
	m := NewModel()
	m.state = StateHelp
	m.previousState = StateWelcome
	m.width = 80

	output := m.viewHelp()

	globalShortcuts := []string{
		"?",
		"Show this help screen",
		"Ctrl+C",
		"Quit application",
		"Esc",
		"Go back / Cancel",
	}

	for _, shortcut := range globalShortcuts {
		if !strings.Contains(output, shortcut) {
			t.Errorf("viewHelp() should contain global shortcut %q", shortcut)
		}
	}
}

// TestGetScreenTitle verifies that each screen state has a proper title.
func TestGetScreenTitle(t *testing.T) {
	tests := []struct {
		state AppState
		want  string
	}{
		{StateWelcome, "Welcome Screen"},
		{StateProjectName, "Project Name Entry"},
		{StateTemplateSelect, "Template Selection"},
		{StateDatabaseSelect, "Database Selection"},
		{StateObservabilitySelect, "Observability Selection"},
		{StateSummary, "Configuration Summary"},
		{StateGenerating, "Generation Progress"},
		{StateDone, "Completion Screen"},
		{StateHelp, "Current Screen"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := getScreenTitle(tt.state)
			if got != tt.want {
				t.Errorf("getScreenTitle(%v) = %q, want %q", tt.state, got, tt.want)
			}
		})
	}
}

// TestGetScreenShortcuts_WelcomeScreen verifies contextual shortcuts for welcome screen.
func TestGetScreenShortcuts_WelcomeScreen(t *testing.T) {
	shortcuts := getScreenShortcuts(StateWelcome)

	expectedShortcuts := []string{
		"↑ / ↓",
		"Navigate menu items",
		"Enter",
		"Select menu option",
		"Esc",
		"Quit application",
	}

	for _, expected := range expectedShortcuts {
		if !strings.Contains(shortcuts, expected) {
			t.Errorf("getScreenShortcuts(StateWelcome) should contain %q", expected)
		}
	}
}

// TestGetScreenShortcuts_ProjectNameScreen verifies contextual shortcuts for project name entry.
func TestGetScreenShortcuts_ProjectNameScreen(t *testing.T) {
	shortcuts := getScreenShortcuts(StateProjectName)

	expectedShortcuts := []string{
		"Enter",
		"Submit project name",
		"Esc",
		"Go back to welcome screen",
		"Type",
		"Enter your project name",
	}

	for _, expected := range expectedShortcuts {
		if !strings.Contains(shortcuts, expected) {
			t.Errorf("getScreenShortcuts(StateProjectName) should contain %q", expected)
		}
	}
}

// TestGetScreenShortcuts_ListScreens verifies contextual shortcuts for list-based screens.
func TestGetScreenShortcuts_ListScreens(t *testing.T) {
	listStates := []AppState{
		StateTemplateSelect,
		StateDatabaseSelect,
		StateObservabilitySelect,
	}

	expectedShortcuts := []string{
		"↑ / ↓",
		"Navigate list items",
		"Enter",
		"Select highlighted item",
		"Esc",
		"Go back to previous screen",
		"j / k",
		"Navigate (vim-style)",
	}

	for _, state := range listStates {
		t.Run(getScreenTitle(state), func(t *testing.T) {
			shortcuts := getScreenShortcuts(state)

			for _, expected := range expectedShortcuts {
				if !strings.Contains(shortcuts, expected) {
					t.Errorf("getScreenShortcuts(%v) should contain %q", state, expected)
				}
			}
		})
	}
}

// TestGetScreenShortcuts_SummaryScreen verifies contextual shortcuts for summary screen.
func TestGetScreenShortcuts_SummaryScreen(t *testing.T) {
	shortcuts := getScreenShortcuts(StateSummary)

	expectedShortcuts := []string{
		"Enter",
		"Confirm and start generation",
		"Esc",
		"Go back to edit configuration",
		"← / →",
		"Navigate back/forward",
		"b / n",
		"Back / Next (alternative)",
	}

	for _, expected := range expectedShortcuts {
		if !strings.Contains(shortcuts, expected) {
			t.Errorf("getScreenShortcuts(StateSummary) should contain %q", expected)
		}
	}
}

// TestGetScreenShortcuts_GeneratingScreen verifies contextual shortcuts for generation screen.
func TestGetScreenShortcuts_GeneratingScreen(t *testing.T) {
	shortcuts := getScreenShortcuts(StateGenerating)

	expectedMessages := []string{
		"No shortcuts available during generation",
		"Please wait for the process to complete",
	}

	for _, expected := range expectedMessages {
		if !strings.Contains(shortcuts, expected) {
			t.Errorf("getScreenShortcuts(StateGenerating) should contain %q", expected)
		}
	}
}

// TestGetScreenShortcuts_DoneScreen verifies contextual shortcuts for completion screen.
func TestGetScreenShortcuts_DoneScreen(t *testing.T) {
	shortcuts := getScreenShortcuts(StateDone)

	expectedShortcuts := []string{
		"Any key",
		"Exit application",
		"Enter",
	}

	for _, expected := range expectedShortcuts {
		if !strings.Contains(shortcuts, expected) {
			t.Errorf("getScreenShortcuts(StateDone) should contain %q", expected)
		}
	}
}

// TestViewHelp_DisplaysTips verifies that the help screen shows helpful tips.
func TestViewHelp_DisplaysTips(t *testing.T) {
	m := NewModel()
	m.state = StateHelp
	m.previousState = StateWelcome
	m.width = 80

	output := m.viewHelp()

	tips := []string{
		"Use arrow keys (↑↓) to navigate lists",
		"Press Esc at any time to go back",
		"Your configuration is saved as you progress",
		"Press ? at any time to see context-specific shortcuts",
	}

	for _, tip := range tips {
		if !strings.Contains(output, tip) {
			t.Errorf("viewHelp() should contain tip %q", tip)
		}
	}

	if !strings.Contains(output, "Tips") {
		t.Error("viewHelp() should contain 'Tips' header")
	}
}

// TestViewHelp_DisplaysFooter verifies that the help screen shows the footer.
func TestViewHelp_DisplaysFooter(t *testing.T) {
	m := NewModel()
	m.state = StateHelp
	m.previousState = StateWelcome
	m.width = 80

	output := m.viewHelp()

	if !strings.Contains(output, "Press any key to return") {
		t.Error("viewHelp() should contain footer instruction")
	}
}

// TestViewHelp_ContextualHelp verifies that help adapts to different screens.
func TestViewHelp_ContextualHelp(t *testing.T) {
	tests := []struct {
		name          string
		previousState AppState
		wantContain   []string
		wantNotContain []string
	}{
		{
			name:          "Welcome Screen Context",
			previousState: StateWelcome,
			wantContain:   []string{"Welcome Screen", "Navigate menu items"},
			wantNotContain: []string{"Submit project name", "Select highlighted item"},
		},
		{
			name:          "Project Name Context",
			previousState: StateProjectName,
			wantContain:   []string{"Project Name Entry", "Submit project name", "Enter your project name"},
			wantNotContain: []string{"Navigate list items", "Confirm and start generation"},
		},
		{
			name:          "Template Selection Context",
			previousState: StateTemplateSelect,
			wantContain:   []string{"Template Selection", "Navigate list items", "Select highlighted item"},
			wantNotContain: []string{"Submit project name", "Confirm and start generation"},
		},
		{
			name:          "Summary Context",
			previousState: StateSummary,
			wantContain:   []string{"Configuration Summary", "Confirm and start generation", "Go back to edit configuration"},
			wantNotContain: []string{"Submit project name", "Navigate list items"},
		},
		{
			name:          "Generating Context",
			previousState: StateGenerating,
			wantContain:   []string{"Generation Progress", "No shortcuts available during generation"},
			wantNotContain: []string{"Submit project name", "Navigate list items"},
		},
		{
			name:          "Done Context",
			previousState: StateDone,
			wantContain:   []string{"Completion Screen", "Exit application"},
			wantNotContain: []string{"Submit project name", "Navigate list items", "Confirm and start generation"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModel()
			m.state = StateHelp
			m.previousState = tt.previousState
			m.width = 80

			output := m.viewHelp()

			// Check that expected content is present
			for _, want := range tt.wantContain {
				if !strings.Contains(output, want) {
					t.Errorf("viewHelp() with previousState=%v should contain %q", tt.previousState, want)
				}
			}

			// Check that unwanted content is not present
			for _, notWant := range tt.wantNotContain {
				if strings.Contains(output, notWant) {
					t.Errorf("viewHelp() with previousState=%v should NOT contain %q", tt.previousState, notWant)
				}
			}
		})
	}
}

// TestRenderShortcut verifies the shortcut rendering function.
func TestRenderShortcut(t *testing.T) {
	output := renderShortcut("Enter", "Confirm selection")

	if !strings.Contains(output, "Enter") {
		t.Error("renderShortcut() should contain the key")
	}

	if !strings.Contains(output, "Confirm selection") {
		t.Error("renderShortcut() should contain the description")
	}
}

// TestRenderTip verifies the tip rendering function.
func TestRenderTip(t *testing.T) {
	output := renderTip("This is a helpful tip")

	if !strings.Contains(output, "This is a helpful tip") {
		t.Error("renderTip() should contain the tip text")
	}

	// Should have bullet point formatting
	if !strings.Contains(output, "•") {
		t.Error("renderTip() should contain bullet point")
	}
}

// TestViewHelp_WidthAdaptation verifies that help adapts to different terminal widths.
func TestViewHelp_WidthAdaptation(t *testing.T) {
	widths := []int{60, 80, 100, 120}

	for _, width := range widths {
		t.Run(string(rune(width)), func(t *testing.T) {
			m := NewModel()
			m.state = StateHelp
			m.previousState = StateWelcome
			m.width = width

			output := m.viewHelp()

			if output == "" {
				t.Errorf("viewHelp() should produce output for width %d", width)
			}

			// The output should contain the basic structure regardless of width
			if !strings.Contains(output, "Keyboard Shortcuts") {
				t.Errorf("viewHelp() should contain header for width %d", width)
			}
		})
	}
}

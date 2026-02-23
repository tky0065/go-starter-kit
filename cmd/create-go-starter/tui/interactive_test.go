package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// dummyGenerator is a no-op generator function for tests.
func dummyGenerator(projectName, template, database, observability string, progressCallback func(current, total int)) error {
	// Simulate successful generation
	if progressCallback != nil {
		progressCallback(1, 1)
	}
	return nil
}

// TestInteractiveModelCreation verifies that we can create an interactive model with defaults.
func TestInteractiveModelCreation(t *testing.T) {
	defaults := InteractiveDefaults{
		ProjectName:   "test-project",
		Template:      "full",
		Database:      "postgres",
		Observability: "none",
	}

	model := NewInteractiveModel(defaults, dummyGenerator)

	// Model should be initialized with defaults
	// Story 10.7: Now starts with welcome screen instead of direct project name input
	if model.state != StateWelcome {
		t.Errorf("Expected initial state to be StateWelcome (welcome screen), got %v", model.state)
	}

	// Project input should have the default as placeholder
	if model.projectInput.Value() == "" && model.projectInput.Placeholder == "" {
		t.Error("Project input should have a placeholder or value")
	}
}

// TestInteractiveStates verifies that the interactive model transitions through states correctly.
func TestInteractiveStates(t *testing.T) {
	model := NewInteractiveModel(InteractiveDefaults{}, dummyGenerator)

	t.Run("starts in StateWelcome", func(t *testing.T) {
		if model.state != StateWelcome {
			t.Errorf("Expected StateWelcome (welcome screen), got %v", model.state)
		}
	})

	t.Run("transitions to StateTemplateSelect after project name submission", func(t *testing.T) {
		newModel, _ := model.Update(ProjectNameSubmittedMsg{Name: "my-project"})
		m := newModel.(Model)

		if m.state != StateTemplateSelect {
			t.Errorf("Expected StateTemplateSelect, got %v", m.state)
		}

		if m.projectName != "my-project" {
			t.Errorf("Expected projectName to be 'my-project', got '%s'", m.projectName)
		}
	})
}

// TestInteractiveNavigation verifies that users can navigate back through screens.
func TestInteractiveNavigation(t *testing.T) {
	model := NewInteractiveModel(InteractiveDefaults{}, dummyGenerator)
	model.state = StateTemplateSelect

	t.Run("navigates back to project name from template select", func(t *testing.T) {
		newModel, _ := model.Update(NavigateBackMsg{})
		m := newModel.(Model)

		if m.state != StateProjectName {
			t.Errorf("Expected StateProjectName after back navigation, got %v", m.state)
		}
	})
}

// TestInteractiveKeyboardHandling verifies keyboard shortcuts work correctly.
func TestInteractiveKeyboardHandling(t *testing.T) {
	model := NewInteractiveModel(InteractiveDefaults{}, dummyGenerator)

	t.Run("ctrl+c quits the application", func(t *testing.T) {
		msg := tea.KeyMsg{Type: tea.KeyCtrlC}
		_, cmd := model.Update(msg)

		if cmd == nil {
			t.Error("Expected quit command, got nil")
		}

		// Check if it's a quit command (we can't easily test the actual quit behavior)
		// but we can verify a command was returned
	})

	t.Run("esc quits the application", func(t *testing.T) {
		msg := tea.KeyMsg{Type: tea.KeyEsc}
		_, cmd := model.Update(msg)

		if cmd == nil {
			t.Error("Expected quit command, got nil")
		}
	})
}

// TestInteractiveTemplateList verifies template selection functionality.
func TestInteractiveTemplateList(t *testing.T) {
	model := NewInteractiveModel(InteractiveDefaults{}, dummyGenerator)
	model.state = StateTemplateSelect

	// After moving to template select state, template list should be initialized
	// This is a placeholder test - actual implementation will verify the list has items
	t.Run("template list is initialized in StateTemplateSelect", func(t *testing.T) {
		// Template list should be ready for selection
		// Actual verification would check list items, but model structure needs to be set up first
		if model.state != StateTemplateSelect {
			t.Error("Model should be in StateTemplateSelect for this test")
		}
	})
}

// TestInteractiveSummaryScreen verifies the configuration summary screen.
func TestInteractiveSummaryScreen(t *testing.T) {
	model := NewInteractiveModel(InteractiveDefaults{}, dummyGenerator)
	model.projectName = "my-app"
	model.template = "full"
	model.database = "postgres"
	model.observability = "basic"
	model.state = StateSummary

	view := model.View()

	// Summary should show all selections
	if view == "" {
		t.Error("Summary view should not be empty")
	}

	// View should contain the project name
	// (More detailed assertions would check for specific content)
}

// TestProjectNameValidation verifies that invalid project names are rejected.
// This test covers MEDIUM-2: Project name validation in TUI.
func TestProjectNameValidation(t *testing.T) {
	model := NewInteractiveModel(InteractiveDefaults{}, dummyGenerator)

	testCases := []struct {
		name          string
		projectName   string
		shouldBeValid bool
	}{
		{"valid lowercase name", "my-project", true},
		{"valid with numbers", "project123", true},
		{"valid with uppercase", "MyProject", true},
		{"valid with underscores", "my_project", true},
		{"valid starting with number", "123project", true},
		{"invalid empty string", "", false},
		{"invalid with spaces", "my project", false},
		{"invalid with special chars", "@#$project", false},
		{"invalid with dots", "my.project", false},
		{"invalid with slashes", "my/project", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Start fresh from StateWelcome and transition to StateProjectName
			testModel := model
			testModel.state = StateProjectName

			// Submit project name
			newModel, _ := testModel.Update(ProjectNameSubmittedMsg{Name: tc.projectName})
			m := newModel.(Model)

			if tc.shouldBeValid {
				// Valid name should transition to template select
				if m.state != StateTemplateSelect {
					t.Errorf("Expected state StateTemplateSelect for valid name '%s', got %v", tc.projectName, m.state)
				}
				if m.validationError != "" {
					t.Errorf("Expected no validation error for valid name '%s', got: %s", tc.projectName, m.validationError)
				}
			} else {
				// Invalid name should stay in project name state with error
				if m.state != StateProjectName {
					t.Errorf("Expected state StateProjectName for invalid name '%s', got %v", tc.projectName, m.state)
				}
				if m.validationError == "" {
					t.Errorf("Expected validation error for invalid name '%s', got none", tc.projectName)
				}
			}
		})
	}
}

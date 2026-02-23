package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestModelImplementsTeaModel verifies that our main Model implements tea.Model interface.
func TestModelImplementsTeaModel(t *testing.T) {
	// This will fail to compile if Model doesn't implement tea.Model
	var _ tea.Model = Model{}
}

// TestModelInit verifies the Init() method returns the correct initial command.
func TestModelInit(t *testing.T) {
	m := Model{}
	cmd := m.Init()

	// Init should return nil or a valid command
	if cmd != nil {
		t.Logf("Init returned a command (expected for async initialization)")
	}
}

// TestModelUpdate verifies the Update() method handles messages correctly.
func TestModelUpdate(t *testing.T) {
	t.Run("handles quit message", func(t *testing.T) {
		m := Model{}

		// Send a KeyMsg with ctrl+c
		msg := tea.KeyMsg{Type: tea.KeyCtrlC}
		newModel, cmd := m.Update(msg)

		// Should return a quit command
		if cmd == nil {
			t.Error("Expected quit command, got nil")
		}

		// Model should be returned (possibly unchanged)
		if _, ok := newModel.(Model); !ok {
			t.Error("Update should return a Model")
		}
	})
}

// TestModelView verifies the View() method returns a string.
func TestModelView(t *testing.T) {
	m := Model{}
	view := m.View()

	// View should return a non-empty string
	if view == "" {
		t.Error("View() should return a non-empty string")
	}
}

// TestModelStates verifies that the Model supports different application states.
func TestModelStates(t *testing.T) {
	// Model should have a way to track current state/step
	// This is a structural test to ensure the Model has state management
	t.Run("model has state tracking", func(t *testing.T) {
		// The model should be instantiatable
		m := NewModel()
		if m.state != StateWelcome {
			t.Errorf("Expected initial state to be StateWelcome (welcome screen), got %v", m.state)
		}
	})
}

package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestWelcomeScreenCreation verifies the welcome screen can be created.
func TestWelcomeScreenCreation(t *testing.T) {
	model := NewInteractiveModel(InteractiveDefaults{}, dummyGenerator)

	// Initial state should be StateWelcome (not StateProjectName)
	if model.state != StateWelcome {
		t.Errorf("Expected initial state to be StateWelcome, got %v", model.state)
	}
}

// TestWelcomeScreenView verifies the welcome screen displays the logo and menu.
func TestWelcomeScreenView(t *testing.T) {
	model := NewInteractiveModel(InteractiveDefaults{}, dummyGenerator)
	model.state = StateWelcome

	view := model.View()

	// View should contain the logo/branding
	if !strings.Contains(view, "go-starter-kit") {
		t.Error("Welcome screen should display 'go-starter-kit' branding")
	}

	// View should contain menu options
	if !strings.Contains(view, "Create New Project") {
		t.Error("Welcome screen should have 'Create New Project' option")
	}
}

// TestWelcomeScreenNavigation verifies menu navigation works.
func TestWelcomeScreenNavigation(t *testing.T) {
	model := NewInteractiveModel(InteractiveDefaults{}, dummyGenerator)
	model.state = StateWelcome

	t.Run("Enter key starts project creation", func(t *testing.T) {
		msg := tea.KeyMsg{Type: tea.KeyEnter}
		newModel, _ := model.Update(msg)
		m := newModel.(Model)

		// Should transition to project name input
		if m.state != StateProjectName {
			t.Errorf("Expected state StateProjectName after Enter, got %v", m.state)
		}
	})
}

// TestWelcomeToProjectNameTransition verifies the flow from welcome to form.
func TestWelcomeToProjectNameTransition(t *testing.T) {
	model := NewInteractiveModel(InteractiveDefaults{}, dummyGenerator)
	model.state = StateWelcome

	// Simulate user selecting "Create New Project" and pressing Enter
	newModel, _ := model.Update(WelcomeMenuSelectedMsg{Action: "create"})
	m := newModel.(Model)

	if m.state != StateProjectName {
		t.Errorf("Expected transition to StateProjectName, got %v", m.state)
	}
}

// TestWelcomeMenuList verifies that the welcome screen has an interactive menu list.
// Story 10.7 Task 1.3: Menu interactif avec bubbles/list
func TestWelcomeMenuList(t *testing.T) {
	model := NewInteractiveModel(InteractiveDefaults{}, dummyGenerator)
	model.state = StateWelcome

	// Welcome menu list should be initialized
	if model.welcomeList.Items() == nil {
		t.Error("Welcome menu list should be initialized with items")
	}

	// Should have 3 menu options: Create, Help, Exit
	items := model.welcomeList.Items()
	if len(items) != 3 {
		t.Errorf("Expected 3 menu items, got %d", len(items))
	}

	// Verify menu items content
	itemTexts := []string{}
	for _, item := range items {
		if wItem, ok := item.(welcomeMenuItem); ok {
			itemTexts = append(itemTexts, wItem.title)
		}
	}

	expectedItems := []string{"Create New Project", "Help", "Exit"}
	for i, expected := range expectedItems {
		if i >= len(itemTexts) || !strings.Contains(itemTexts[i], expected) {
			t.Errorf("Expected menu item %d to contain '%s', got %v", i, expected, itemTexts)
		}
	}
}

// TestWelcomeMenuNavigation verifies arrow key navigation in welcome menu.
// Story 10.7 Task 1.3: Navigation avec ↑↓
func TestWelcomeMenuArrowNavigation(t *testing.T) {
	model := NewInteractiveModel(InteractiveDefaults{}, dummyGenerator)
	model.state = StateWelcome

	// Initially, first item should be selected
	initialIndex := model.welcomeList.Index()

	// Press Down arrow
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	newModel, _ := model.Update(downMsg)
	m := newModel.(Model)

	// Index should have incremented
	if m.welcomeList.Index() != initialIndex+1 {
		t.Errorf("Expected index to increment after Down key, got %d", m.welcomeList.Index())
	}

	// Press Up arrow
	upMsg := tea.KeyMsg{Type: tea.KeyUp}
	newModel2, _ := m.Update(upMsg)
	m2 := newModel2.(Model)

	// Index should have decremented back
	if m2.welcomeList.Index() != initialIndex {
		t.Errorf("Expected index to decrement after Up key, got %d", m2.welcomeList.Index())
	}
}

// TestWelcomeMenuSelection verifies Enter key selection works with the list.
// Story 10.7 Task 1.3: Sélection avec Enter
func TestWelcomeMenuSelectionWithList(t *testing.T) {
	model := NewInteractiveModel(InteractiveDefaults{}, dummyGenerator)
	model.state = StateWelcome

	// Navigate to "Create New Project" (should be first item)
	// Press Enter to select
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ := model.Update(enterMsg)
	m := newModel.(Model)

	// Should transition to StateProjectName
	if m.state != StateProjectName {
		t.Errorf("Expected StateProjectName after selecting 'Create New Project', got %v", m.state)
	}
}

// TestWelcomeMenuHelpSelection verifies Help option works.
// Story 10.7 Task 1.3 + Task 8: Help option
func TestWelcomeMenuHelpSelection(t *testing.T) {
	model := NewInteractiveModel(InteractiveDefaults{}, dummyGenerator)
	model.state = StateWelcome

	// Navigate to "Help" (second item)
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	newModel, _ := model.Update(downMsg)
	model = newModel.(Model)

	// Press Enter
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel2, _ := model.Update(enterMsg)
	m := newModel2.(Model)

	// Should transition to StateHelp
	if m.state != StateHelp {
		t.Errorf("Expected StateHelp after selecting 'Help', got %v", m.state)
	}
}

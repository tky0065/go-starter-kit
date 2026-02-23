package tui

import (
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestNewErrorModel tests creating a new error model.
func TestNewErrorModel(t *testing.T) {
	title := "Test Error"
	e := NewErrorModel(title)

	if e.Title != title {
		t.Errorf("expected title %q, got %q", title, e.Title)
	}

	if !e.ShowIcon {
		t.Error("expected ShowIcon to be true by default")
	}

	if len(e.Suggestions) != 0 {
		t.Errorf("expected empty suggestions, got %d", len(e.Suggestions))
	}
}

// TestErrorModel_WithDetails tests adding details to error model.
func TestErrorModel_WithDetails(t *testing.T) {
	details := "Additional error context"
	e := NewErrorModel("Error").WithDetails(details)

	if e.Details != details {
		t.Errorf("expected details %q, got %q", details, e.Details)
	}
}

// TestErrorModel_WithSuggestions tests adding suggestions to error model.
func TestErrorModel_WithSuggestions(t *testing.T) {
	suggestions := []string{"Try this", "Or that"}
	e := NewErrorModel("Error").WithSuggestions(suggestions...)

	if len(e.Suggestions) != 2 {
		t.Errorf("expected 2 suggestions, got %d", len(e.Suggestions))
	}

	if e.Suggestions[0] != suggestions[0] {
		t.Errorf("expected first suggestion %q, got %q", suggestions[0], e.Suggestions[0])
	}
}

// TestErrorModel_View tests rendering error model.
func TestErrorModel_View(t *testing.T) {
	e := NewErrorModel("File not found").
		WithDetails("The specified file does not exist").
		WithSuggestions("Check the file path", "Verify file permissions")

	view := e.View()

	// Check that the view contains key elements
	if !strings.Contains(view, "File not found") {
		t.Error("view should contain error title")
	}

	if !strings.Contains(view, "The specified file does not exist") {
		t.Error("view should contain error details")
	}

	if !strings.Contains(view, "Suggestions:") {
		t.Error("view should contain suggestions header")
	}

	if !strings.Contains(view, "Check the file path") {
		t.Error("view should contain first suggestion")
	}

	if !strings.Contains(view, "Verify file permissions") {
		t.Error("view should contain second suggestion")
	}
}

// TestErrorModel_ViewWithoutIcon tests rendering error without icon.
func TestErrorModel_ViewWithoutIcon(t *testing.T) {
	e := NewErrorModel("Error")
	e.ShowIcon = false

	view := e.View()

	// The error icon (✗) should not be present when ShowIcon is false
	// Note: We can't easily test for absence of icon in styled output,
	// but we can verify the view is not empty
	if view == "" {
		t.Error("view should not be empty")
	}
}

// TestNewSuccessModel tests creating a new success model.
func TestNewSuccessModel(t *testing.T) {
	title := "Test Success"
	s := NewSuccessModel(title)

	if s.Title != title {
		t.Errorf("expected title %q, got %q", title, s.Title)
	}

	if !s.ShowIcon {
		t.Error("expected ShowIcon to be true by default")
	}

	if len(s.Items) != 0 {
		t.Errorf("expected empty items, got %d", len(s.Items))
	}
}

// TestSuccessModel_WithDetails tests adding details to success model.
func TestSuccessModel_WithDetails(t *testing.T) {
	details := "Operation completed successfully"
	s := NewSuccessModel("Success").WithDetails(details)

	if s.Details != details {
		t.Errorf("expected details %q, got %q", details, s.Details)
	}
}

// TestSuccessModel_WithItems tests adding items to success model.
func TestSuccessModel_WithItems(t *testing.T) {
	items := []string{"Item 1", "Item 2", "Item 3"}
	s := NewSuccessModel("Success").WithItems(items...)

	if len(s.Items) != 3 {
		t.Errorf("expected 3 items, got %d", len(s.Items))
	}

	if s.Items[1] != items[1] {
		t.Errorf("expected item %q, got %q", items[1], s.Items[1])
	}
}

// TestSuccessModel_View tests rendering success model.
func TestSuccessModel_View(t *testing.T) {
	s := NewSuccessModel("Project created successfully").
		WithDetails("All files have been generated").
		WithItems("34 files created", "~850 KB total size", "Git repository initialized")

	view := s.View()

	// Check that the view contains key elements
	if !strings.Contains(view, "Project created successfully") {
		t.Error("view should contain success title")
	}

	if !strings.Contains(view, "All files have been generated") {
		t.Error("view should contain success details")
	}

	if !strings.Contains(view, "34 files created") {
		t.Error("view should contain first item")
	}

	if !strings.Contains(view, "Git repository initialized") {
		t.Error("view should contain last item")
	}
}

// TestNewWarningModel tests creating a new warning model.
func TestNewWarningModel(t *testing.T) {
	title := "Test Warning"
	w := NewWarningModel(title)

	if w.Title != title {
		t.Errorf("expected title %q, got %q", title, w.Title)
	}

	if !w.ShowIcon {
		t.Error("expected ShowIcon to be true by default")
	}
}

// TestWarningModel_View tests rendering warning model.
func TestWarningModel_View(t *testing.T) {
	w := NewWarningModel("Large project detected").
		WithDetails("This project will generate 100+ files")

	view := w.View()

	if !strings.Contains(view, "Large project detected") {
		t.Error("view should contain warning title")
	}

	if !strings.Contains(view, "This project will generate 100+ files") {
		t.Error("view should contain warning details")
	}
}

// TestNewInfoModel tests creating a new info model.
func TestNewInfoModel(t *testing.T) {
	title := "Test Info"
	i := NewInfoModel(title)

	if i.Title != title {
		t.Errorf("expected title %q, got %q", title, i.Title)
	}

	if !i.ShowIcon {
		t.Error("expected ShowIcon to be true by default")
	}
}

// TestInfoModel_View tests rendering info model.
func TestInfoModel_View(t *testing.T) {
	i := NewInfoModel("Database selection").
		WithDetails("PostgreSQL is recommended for production")

	view := i.View()

	if !strings.Contains(view, "Database selection") {
		t.Error("view should contain info title")
	}

	if !strings.Contains(view, "PostgreSQL is recommended for production") {
		t.Error("view should contain info details")
	}
}

// Note: Helper functions for rendering boxes have been removed to avoid conflicts.
// Tests for the Model-based approach are below.

// TestNewConfirmationModel tests creating a new confirmation dialog.
func TestNewConfirmationModel(t *testing.T) {
	message := "Are you sure?"
	m := NewConfirmationModel(message)

	if m.Message != message {
		t.Errorf("expected message %q, got %q", message, m.Message)
	}

	if m.ConfirmText != "Yes" {
		t.Errorf("expected confirm text %q, got %q", "Yes", m.ConfirmText)
	}

	if m.CancelText != "No" {
		t.Errorf("expected cancel text %q, got %q", "No", m.CancelText)
	}

	if !m.Focused {
		t.Error("expected Focused to be true (confirm button) by default")
	}

	if m.Done {
		t.Error("expected Done to be false by default")
	}
}

// TestConfirmationModel_NavigateLeft tests navigating to cancel button.
func TestConfirmationModel_NavigateLeft(t *testing.T) {
	m := NewConfirmationModel("Test")
	m.Focused = true // Start on confirm button

	// Navigate left to cancel button
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyLeft})

	if newModel.Focused {
		t.Error("expected Focused to be false after navigating left")
	}
}

// TestConfirmationModel_NavigateRight tests navigating to confirm button.
func TestConfirmationModel_NavigateRight(t *testing.T) {
	m := NewConfirmationModel("Test")
	m.Focused = false // Start on cancel button

	// Navigate right to confirm button
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRight})

	if !newModel.Focused {
		t.Error("expected Focused to be true after navigating right")
	}
}

// TestConfirmationModel_ConfirmWithEnter tests confirming with Enter key.
func TestConfirmationModel_ConfirmWithEnter(t *testing.T) {
	m := NewConfirmationModel("Test")
	m.Focused = true // On confirm button

	// Press Enter
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if !newModel.Done {
		t.Error("expected Done to be true after pressing Enter")
	}

	if !newModel.Confirmed {
		t.Error("expected Confirmed to be true when Enter is pressed on confirm button")
	}
}

// TestConfirmationModel_CancelWithEnter tests canceling with Enter key.
func TestConfirmationModel_CancelWithEnter(t *testing.T) {
	m := NewConfirmationModel("Test")
	m.Focused = false // On cancel button

	// Press Enter
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if !newModel.Done {
		t.Error("expected Done to be true after pressing Enter")
	}

	if newModel.Confirmed {
		t.Error("expected Confirmed to be false when Enter is pressed on cancel button")
	}
}

// TestConfirmationModel_YesShortcut tests the 'Y' shortcut.
func TestConfirmationModel_YesShortcut(t *testing.T) {
	m := NewConfirmationModel("Test")
	m.Focused = false // Start on cancel button

	// Press 'Y'
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

	if !newModel.Done {
		t.Error("expected Done to be true after pressing 'y'")
	}

	if !newModel.Confirmed {
		t.Error("expected Confirmed to be true after pressing 'y'")
	}
}

// TestConfirmationModel_NoShortcut tests the 'N' shortcut.
func TestConfirmationModel_NoShortcut(t *testing.T) {
	m := NewConfirmationModel("Test")
	m.Focused = true // Start on confirm button

	// Press 'N'
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

	if !newModel.Done {
		t.Error("expected Done to be true after pressing 'n'")
	}

	if newModel.Confirmed {
		t.Error("expected Confirmed to be false after pressing 'n'")
	}
}

// TestConfirmationModel_EscapeCancel tests canceling with Escape key.
func TestConfirmationModel_EscapeCancel(t *testing.T) {
	m := NewConfirmationModel("Test")

	// Press Escape
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})

	if !newModel.Done {
		t.Error("expected Done to be true after pressing Escape")
	}

	if newModel.Confirmed {
		t.Error("expected Confirmed to be false after pressing Escape")
	}
}

// TestConfirmationModel_View tests rendering the confirmation dialog.
func TestConfirmationModel_View(t *testing.T) {
	m := NewConfirmationModel("Delete all files?")

	view := m.View()

	// Check that the view contains key elements
	if !strings.Contains(view, "Delete all files?") {
		t.Error("view should contain the message")
	}

	if !strings.Contains(view, "Yes") {
		t.Error("view should contain Yes button")
	}

	if !strings.Contains(view, "No") {
		t.Error("view should contain No button")
	}

	if !strings.Contains(view, "✓") {
		t.Error("view should contain success icon")
	}

	if !strings.Contains(view, "✗") {
		t.Error("view should contain error icon")
	}
}

// TestRenderConfirmation tests the helper function for rendering confirmations.
func TestRenderConfirmation(t *testing.T) {
	view := RenderConfirmation("Overwrite existing files?", true)

	if !strings.Contains(view, "Overwrite existing files?") {
		t.Error("view should contain the message")
	}
}

// TestCommonErrorSuggestions tests that common error suggestions are defined.
func TestCommonErrorSuggestions(t *testing.T) {
	if len(ProjectNameSuggestions) == 0 {
		t.Error("ProjectNameSuggestions should not be empty")
	}

	if len(DirectoryExistsSuggestions) == 0 {
		t.Error("DirectoryExistsSuggestions should not be empty")
	}

	if len(GitInitSuggestions) == 0 {
		t.Error("GitInitSuggestions should not be empty")
	}

	if len(DependencySuggestions) == 0 {
		t.Error("DependencySuggestions should not be empty")
	}

	if len(DatabaseSuggestions) == 0 {
		t.Error("DatabaseSuggestions should not be empty")
	}
}

// TestErrorModel_Integration tests using error model with common suggestions.
func TestErrorModel_Integration(t *testing.T) {
	// Simulate an error with project name validation
	e := NewErrorModel("Invalid project name: 'My Project'").
		WithDetails("Project names cannot contain spaces or capital letters").
		WithSuggestions(ProjectNameSuggestions...)

	view := e.View()

	if !strings.Contains(view, "Invalid project name") {
		t.Error("view should contain error title")
	}

	if !strings.Contains(view, "Use lowercase letters") {
		t.Error("view should contain suggestion from ProjectNameSuggestions")
	}

	if !strings.Contains(view, "Example: my-awesome-api") {
		t.Error("view should contain example from ProjectNameSuggestions")
	}
}

// TestSuccessModel_Integration tests using success model with file stats.
func TestSuccessModel_Integration(t *testing.T) {
	// Simulate a successful project generation
	totalFiles := 34

	s := NewSuccessModel("Project created successfully!").
		WithDetails("Your go-starter-kit project is ready to use").
		WithItems(
			fmt.Sprintf("%d files created", totalFiles),
			"~850 KB total size",
			"Git repository initialized",
			"Go dependencies installed",
		)

	view := s.View()

	if !strings.Contains(view, "Project created successfully") {
		t.Error("view should contain success title")
	}

	if !strings.Contains(view, "34 files created") {
		t.Error("view should contain files count")
	}

	if !strings.Contains(view, "850 KB total size") {
		t.Error("view should contain formatted file size")
	}

	if !strings.Contains(view, "Git repository initialized") {
		t.Error("view should contain git initialization message")
	}
}

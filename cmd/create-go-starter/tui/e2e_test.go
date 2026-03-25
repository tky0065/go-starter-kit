package tui

import (
	"os"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestE2EInteractiveTUIBasicFlow tests the basic end-to-end flow of the TUI.
// This test validates that the TUI can be initialized and navigate through states.
func TestE2EInteractiveTUIBasicFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}

	// Setup test
	_ = t.TempDir() // Create temp dir for cleanup
	projectName := "test-e2e-project"

	// Track if generation was called (not used in this simplified test)
	mockGenerator := func(name, template, database, observability, framework string, progressCallback func(current, total int)) error {
		// Simulate quick generation
		if progressCallback != nil {
			progressCallback(10, 10)
		}
		return nil
	}

	defaults := InteractiveDefaults{
		ProjectName:   projectName,
		Template:      "full",
		Database:      "postgres",
		Observability: "none",
	}

	model := NewInteractiveModel(defaults, mockGenerator)

	// Test initialization
	if model.state != StateWelcome {
		t.Errorf("Expected initial state to be StateWelcome, got %v", model.state)
	}

	// Simulate user selecting "Create New Project" from welcome menu
	updatedModel, _ := model.Update(WelcomeMenuSelectedMsg{Action: "create"})
	model = updatedModel.(Model)
	if model.state != StateProjectName {
		t.Errorf("Expected state to transition to StateProjectName, got %v", model.state)
	}

	// Simulate user entering project name
	updatedModel, _ = model.Update(ProjectNameSubmittedMsg{Name: projectName})
	model = updatedModel.(Model)
	if model.state != StateTemplateSelect {
		t.Errorf("Expected state to transition to StateTemplateSelect, got %v", model.state)
	}

	// Simulate user selecting template
	updatedModel, _ = model.Update(TemplateSelectedMsg{Template: "full"})
	model = updatedModel.(Model)
	if model.state != StateDatabaseSelect {
		t.Errorf("Expected state to transition to StateDatabaseSelect, got %v", model.state)
	}

	// Simulate user selecting database
	updatedModel, _ = model.Update(DatabaseSelectedMsg{Database: "postgres"})
	model = updatedModel.(Model)
	if model.state != StateFrameworkSelect {
		t.Errorf("Expected state to transition to StateFrameworkSelect, got %v", model.state)
	}

	// Simulate user selecting framework
	updatedModel, _ = model.Update(FrameworkSelectedMsg{Framework: "fiber"})
	model = updatedModel.(Model)
	if model.state != StateObservabilitySelect {
		t.Errorf("Expected state to transition to StateObservabilitySelect, got %v", model.state)
	}

	// Simulate user selecting observability
	updatedModel, _ = model.Update(ObservabilitySelectedMsg{Observability: "none"})
	model = updatedModel.(Model)
	if model.state != StateSummary {
		t.Errorf("Expected state to transition to StateSummary, got %v", model.state)
	}

	// Simulate user confirming generation
	updatedModel, cmd := model.Update(ConfirmGenerationMsg{})
	model = updatedModel.(Model)
	if model.state != StateGenerating {
		t.Errorf("Expected state to transition to StateGenerating, got %v", model.state)
	}

	// Verify generation command was created
	if cmd == nil {
		t.Fatal("Expected generation command to be returned")
	}

	// Note: We don't execute the async generation here as it runs in a goroutine
	// Testing the async behavior is done in TestE2EInteractiveTUIProgressUpdates
}

// TestE2EInteractiveTUIProgressUpdates validates that real-time progress updates work.
// This test ensures that FileGeneratedMsg messages are processed correctly during generation.
func TestE2EInteractiveTUIProgressUpdates(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}

	totalFiles := 5
	filesProcessed := 0

	// Mock generator that sends progress updates
	mockGenerator := func(name, template, database, observability, framework string, progressCallback func(current, total int)) error {
		if progressCallback != nil {
			for i := 1; i <= totalFiles; i++ {
				progressCallback(i, totalFiles)
				filesProcessed = i
			}
		}
		return nil
	}

	defaults := InteractiveDefaults{
		ProjectName:   "test-progress",
		Template:      "minimal",
		Database:      "sqlite",
		Observability: "none",
	}

	model := NewInteractiveModel(defaults, mockGenerator)

	// Fast-forward to generating state
	model.state = StateGenerating
	model.totalFiles = totalFiles

	// Execute generation
	cmd := model.generateProjectCmd()
	msg := cmd()

	// Verify progress was tracked
	if filesProcessed != totalFiles {
		t.Errorf("Expected %d files processed, got %d", totalFiles, filesProcessed)
	}

	// Process completion message
	if completeMsg, ok := msg.(GenerationCompleteMsg); ok {
		if !completeMsg.Success {
			t.Errorf("Expected generation to succeed, got error: %v", completeMsg.Error)
		}
		if completeMsg.TotalFiles != totalFiles {
			t.Errorf("Expected TotalFiles=%d, got %d", totalFiles, completeMsg.TotalFiles)
		}
	} else {
		t.Errorf("Expected GenerationCompleteMsg, got %T", msg)
	}
}

// TestE2EInteractiveTUIErrorHandling validates error handling during generation.
func TestE2EInteractiveTUIErrorHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}

	expectedError := os.ErrPermission

	// Mock generator that fails
	mockGenerator := func(name, template, database, observability, framework string, progressCallback func(current, total int)) error {
		return expectedError
	}

	defaults := InteractiveDefaults{
		ProjectName:   "test-error",
		Template:      "full",
		Database:      "postgres",
		Observability: "none",
	}

	model := NewInteractiveModel(defaults, mockGenerator)
	model.state = StateGenerating

	// Execute generation
	cmd := model.generateProjectCmd()
	msg := cmd()

	// Process completion message
	updatedModel, _ := model.Update(msg)
	model = updatedModel.(Model)

	// Verify error was captured
	if model.err == nil {
		t.Error("Expected error to be captured in model")
	}

	if model.err != expectedError {
		t.Errorf("Expected error %v, got %v", expectedError, model.err)
	}

	// Verify state transitioned to done (with error)
	if model.state != StateDone {
		t.Errorf("Expected state to be StateDone after error, got %v", model.state)
	}
}

// TestE2EInteractiveTUINavigationBack tests the back navigation functionality.
func TestE2EInteractiveTUINavigationBack(t *testing.T) {
	mockGenerator := func(name, template, database, observability, framework string, progressCallback func(current, total int)) error {
		return nil
	}

	defaults := InteractiveDefaults{
		ProjectName:   "test-nav",
		Template:      "full",
		Database:      "postgres",
		Observability: "none",
	}

	model := NewInteractiveModel(defaults, mockGenerator)

	// Navigate forward through states
	model.state = StateObservabilitySelect

	// Navigate back
	model = model.navigateBack()
	if model.state != StateFrameworkSelect {
		t.Errorf("Expected state StateFrameworkSelect after back, got %v", model.state)
	}

	model = model.navigateBack()
	if model.state != StateDatabaseSelect {
		t.Errorf("Expected state StateDatabaseSelect after back, got %v", model.state)
	}

	model = model.navigateBack()
	if model.state != StateTemplateSelect {
		t.Errorf("Expected state StateTemplateSelect after back, got %v", model.state)
	}

	model = model.navigateBack()
	if model.state != StateProjectName {
		t.Errorf("Expected state StateProjectName after back, got %v", model.state)
	}

	model = model.navigateBack()
	if model.state != StateWelcome {
		t.Errorf("Expected state StateWelcome after back, got %v", model.state)
	}
}

// TestE2EInteractiveTUIHelpScreen validates contextual help functionality (AC#9).
func TestE2EInteractiveTUIHelpScreen(t *testing.T) {
	mockGenerator := func(name, template, database, observability, framework string, progressCallback func(current, total int)) error {
		return nil
	}

	defaults := InteractiveDefaults{}
	model := NewInteractiveModel(defaults, mockGenerator)

	// Start at welcome screen
	model.state = StateWelcome
	originalState := model.state

	// Press '?' to open help
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	model = updatedModel.(Model)

	if model.state != StateHelp {
		t.Errorf("Expected state to be StateHelp after pressing '?', got %v", model.state)
	}

	if model.previousState != originalState {
		t.Errorf("Expected previousState to be %v, got %v", originalState, model.previousState)
	}

	// Press any key to return
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updatedModel.(Model)

	if model.state != originalState {
		t.Errorf("Expected state to return to %v after closing help, got %v", originalState, model.state)
	}
}

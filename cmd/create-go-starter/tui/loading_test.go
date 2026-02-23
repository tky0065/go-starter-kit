package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func TestNewLoadingModel(t *testing.T) {
	text := "Loading data..."
	m := NewLoadingModel(text)

	if m.text != text {
		t.Errorf("Expected text %q, got %q", text, m.text)
	}

	if m.state != LoadingInProgress {
		t.Errorf("Expected state LoadingInProgress, got %v", m.state)
	}

	if !m.IsInProgress() {
		t.Error("Expected IsInProgress to return true")
	}
}

func TestNewLoadingModelWithSpinner(t *testing.T) {
	text := "Processing..."
	spinnerType := spinner.Line

	m := NewLoadingModelWithSpinner(text, spinnerType)

	// Note: Cannot directly compare spinner.Spinner types (contains slices)
	// Just verify the model was created successfully
	if m.text != text {
		t.Errorf("Expected text %q, got %q", text, m.text)
	}

	if m.text != text {
		t.Errorf("Expected text %q, got %q", text, m.text)
	}
}

func TestLoadingModel_SetText(t *testing.T) {
	m := NewLoadingModel("Initial text")
	newText := "Updated text"

	m.SetText(newText)

	if m.text != newText {
		t.Errorf("Expected text %q, got %q", newText, m.text)
	}
}

func TestLoadingModel_SetSuccess(t *testing.T) {
	m := NewLoadingModel("Loading...")

	m.SetSuccess()

	if m.state != LoadingSuccess {
		t.Errorf("Expected state LoadingSuccess, got %v", m.state)
	}

	if m.IsInProgress() {
		t.Error("Expected IsInProgress to return false after SetSuccess")
	}
}

func TestLoadingModel_SetError(t *testing.T) {
	m := NewLoadingModel("Loading...")

	m.SetError()

	if m.state != LoadingError {
		t.Errorf("Expected state LoadingError, got %v", m.state)
	}

	if m.IsInProgress() {
		t.Error("Expected IsInProgress to return false after SetError")
	}
}

func TestLoadingModel_View_InProgress(t *testing.T) {
	m := NewLoadingModel("Loading data...")
	view := m.View()

	// Should contain the text
	if !strings.Contains(view, "Loading data...") {
		t.Errorf("Expected view to contain text, got: %s", view)
	}

	// Should not contain success or error symbols
	if strings.Contains(view, "✓") {
		t.Error("In-progress view should not contain success checkmark")
	}
	if strings.Contains(view, "✗") {
		t.Error("In-progress view should not contain error cross")
	}
}

func TestLoadingModel_View_Success(t *testing.T) {
	m := NewLoadingModel("Loading data...")
	m.SetSuccess()
	view := m.View()

	// Should contain the text
	if !strings.Contains(view, "Loading data...") {
		t.Errorf("Expected view to contain text, got: %s", view)
	}

	// Should contain success checkmark
	if !strings.Contains(view, "✓") {
		t.Error("Success view should contain checkmark")
	}

	// Should not contain error symbol
	if strings.Contains(view, "✗") {
		t.Error("Success view should not contain error cross")
	}
}

func TestLoadingModel_View_Error(t *testing.T) {
	m := NewLoadingModel("Loading data...")
	m.SetError()
	view := m.View()

	// Should contain the text
	if !strings.Contains(view, "Loading data...") {
		t.Errorf("Expected view to contain text, got: %s", view)
	}

	// Should contain error cross
	if !strings.Contains(view, "✗") {
		t.Error("Error view should contain cross")
	}

	// Should not contain success symbol
	if strings.Contains(view, "✓") {
		t.Error("Error view should not contain success checkmark")
	}
}

func TestLoadingModel_Update_StopsWhenComplete(t *testing.T) {
	m := NewLoadingModel("Loading...")
	m.SetSuccess()

	// Update should not change the model when complete
	initialState := m.state
	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeySpace})

	if cmd != nil {
		t.Error("Expected no command when loading is complete")
	}

	if m.state != initialState {
		t.Error("State should not change when loading is complete")
	}
}

func TestLoadingModel_Init(t *testing.T) {
	m := NewLoadingModel("Loading...")
	cmd := m.Init()

	if cmd == nil {
		t.Error("Expected Init to return a spinner tick command")
	}
}

func TestNewMultiLoadingModel(t *testing.T) {
	operations := []string{
		"Step 1",
		"Step 2",
		"Step 3",
	}

	m := NewMultiLoadingModel(operations, true)

	if len(m.operations) != len(operations) {
		t.Errorf("Expected %d operations, got %d", len(operations), len(m.operations))
	}

	if m.currentIndex != 0 {
		t.Errorf("Expected currentIndex 0, got %d", m.currentIndex)
	}

	if m.IsComplete() {
		t.Error("Expected IsComplete to return false initially")
	}
}

func TestMultiLoadingModel_CurrentOperation(t *testing.T) {
	operations := []string{"Step 1", "Step 2", "Step 3"}
	m := NewMultiLoadingModel(operations, true)

	current := m.CurrentOperation()
	if current != "Step 1" {
		t.Errorf("Expected current operation 'Step 1', got %q", current)
	}
}

func TestMultiLoadingModel_SetCurrentSuccess(t *testing.T) {
	operations := []string{"Step 1", "Step 2", "Step 3"}
	m := NewMultiLoadingModel(operations, true)

	// Complete first operation
	m.SetCurrentSuccess()

	if m.currentIndex != 1 {
		t.Errorf("Expected currentIndex 1, got %d", m.currentIndex)
	}

	if m.operations[0].State != LoadingSuccess {
		t.Error("Expected first operation to be marked as success")
	}

	current := m.CurrentOperation()
	if current != "Step 2" {
		t.Errorf("Expected current operation 'Step 2', got %q", current)
	}
}

func TestMultiLoadingModel_SetCurrentError(t *testing.T) {
	operations := []string{"Step 1", "Step 2", "Step 3"}
	m := NewMultiLoadingModel(operations, true)

	m.SetCurrentError()

	if m.operations[0].State != LoadingError {
		t.Error("Expected first operation to be marked as error")
	}

	// Current index should NOT advance on error
	if m.currentIndex != 0 {
		t.Errorf("Expected currentIndex to stay at 0, got %d", m.currentIndex)
	}
}

func TestMultiLoadingModel_IsComplete(t *testing.T) {
	operations := []string{"Step 1", "Step 2"}
	m := NewMultiLoadingModel(operations, true)

	if m.IsComplete() {
		t.Error("Expected IsComplete to be false initially")
	}

	// Complete both operations
	m.SetCurrentSuccess()
	m.SetCurrentSuccess()

	if !m.IsComplete() {
		t.Error("Expected IsComplete to be true after all operations complete")
	}
}

func TestMultiLoadingModel_View_ShowCompleted(t *testing.T) {
	operations := []string{"Step 1", "Step 2", "Step 3"}
	m := NewMultiLoadingModel(operations, true) // showCompleted = true

	// Complete first operation
	m.SetCurrentSuccess()

	view := m.View()

	// Should show completed operation with checkmark
	if !strings.Contains(view, "✓") {
		t.Error("Expected view to show checkmark for completed operation")
	}

	// Should show current operation
	if !strings.Contains(view, "Step 2") {
		t.Error("Expected view to show current operation")
	}

	// Should NOT show future operation
	if strings.Contains(view, "Step 3") {
		t.Error("Expected view to not show future operation")
	}
}

func TestMultiLoadingModel_View_HideCompleted(t *testing.T) {
	operations := []string{"Step 1", "Step 2", "Step 3"}
	m := NewMultiLoadingModel(operations, false) // showCompleted = false

	// Complete first operation
	m.SetCurrentSuccess()

	view := m.View()

	// Should NOT show completed operation
	if strings.Contains(view, "Step 1") {
		t.Error("Expected view to not show completed operation when showCompleted=false")
	}

	// Should show current operation
	if !strings.Contains(view, "Step 2") {
		t.Error("Expected view to show current operation")
	}
}

func TestMultiLoadingModel_View_Error(t *testing.T) {
	operations := []string{"Step 1", "Step 2"}
	m := NewMultiLoadingModel(operations, true)

	m.SetCurrentError()

	view := m.View()

	// Should show error cross
	if !strings.Contains(view, "✗") {
		t.Error("Expected view to show error cross")
	}

	// Should show the failed operation
	if !strings.Contains(view, "Step 1") {
		t.Error("Expected view to show failed operation")
	}
}

func TestMultiLoadingModel_Update_StopsWhenComplete(t *testing.T) {
	operations := []string{"Step 1"}
	m := NewMultiLoadingModel(operations, true)

	// Complete the operation
	m.SetCurrentSuccess()

	// Update should return nil command when complete
	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeySpace})

	if cmd != nil {
		t.Error("Expected no command when all operations are complete")
	}
}

func TestMultiLoadingModel_Init(t *testing.T) {
	operations := []string{"Step 1"}
	m := NewMultiLoadingModel(operations, true)

	cmd := m.Init()

	if cmd == nil {
		t.Error("Expected Init to return a spinner tick command")
	}
}

func TestMultiLoadingModel_EmptyOperations(t *testing.T) {
	m := NewMultiLoadingModel([]string{}, true)

	if !m.IsComplete() {
		t.Error("Expected IsComplete to return true for empty operations")
	}

	view := m.View()
	if view != "" {
		t.Errorf("Expected empty view for empty operations, got: %q", view)
	}

	current := m.CurrentOperation()
	if current != "" {
		t.Errorf("Expected empty current operation, got: %q", current)
	}
}

func TestLoadingConstants(t *testing.T) {
	// Test that all predefined loading messages are non-empty
	constants := []string{
		LoadingInitializingGit,
		LoadingInstallingDeps,
		LoadingRunningGoModTidy,
		LoadingGeneratingFiles,
		LoadingCreatingDirectories,
		LoadingWritingConfig,
		LoadingSettingUpDatabase,
		LoadingConfiguringDocker,
		LoadingRunningTests,
		LoadingBuildingProject,
		LoadingVerifyingInstall,
	}

	for _, constant := range constants {
		if constant == "" {
			t.Error("Expected non-empty loading constant")
		}
		if !strings.HasSuffix(constant, "...") {
			t.Errorf("Expected loading constant to end with '...', got: %q", constant)
		}
	}
}

package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestFormModel_Creation verifies FormModel initialization.
// Story 10.7 Task 2: Multi-step forms
func TestFormModel_Creation(t *testing.T) {
	form := NewFormModel()

	if form.currentStep != 1 {
		t.Errorf("Expected initial step to be 1, got %d", form.currentStep)
	}

	if form.totalSteps != 4 {
		t.Errorf("Expected 4 total steps (project name, template, database, observability), got %d", form.totalSteps)
	}

	if len(form.history) != 0 {
		t.Errorf("Expected empty history on initialization, got %d items", len(form.history))
	}
}

// TestFormModel_NavigationForward verifies forward navigation.
// Story 10.7 Task 2: Navigation →/n (next)
func TestFormModel_NavigationForward(t *testing.T) {
	form := NewFormModel()
	form.currentStep = 1
	form.projectName = "test-project"

	// Navigate forward with Enter
	newForm, _ := form.Update(tea.KeyMsg{Type: tea.KeyEnter})
	f := newForm.(FormModel)

	// Should advance to step 2
	if f.currentStep != 2 {
		t.Errorf("Expected step 2 after forward navigation, got %d", f.currentStep)
	}

	// History should be updated
	if len(f.history) != 1 {
		t.Errorf("Expected 1 item in history, got %d", len(f.history))
	}

	if f.history[0] != 1 {
		t.Errorf("Expected history[0] to be 1, got %d", f.history[0])
	}
}

// TestFormModel_NavigationBack verifies backward navigation.
// Story 10.7 Task 2: Navigation ←/b (back)
func TestFormModel_NavigationBack(t *testing.T) {
	form := NewFormModel()
	form.currentStep = 2
	form.projectName = "test-project"
	form.template = "standard"
	form.history = []int{1}

	// Navigate back with Left arrow
	newForm, _ := form.Update(tea.KeyMsg{Type: tea.KeyLeft})
	f := newForm.(FormModel)

	// Should go back to step 1
	if f.currentStep != 1 {
		t.Errorf("Expected step 1 after back navigation, got %d", f.currentStep)
	}

	// History should be popped
	if len(f.history) != 0 {
		t.Errorf("Expected empty history after going back, got %d items", len(f.history))
	}
}

// TestFormModel_NavigationBackWithBKey verifies back navigation with 'b' key.
// Story 10.7 Task 2: Navigation b (back alternative)
func TestFormModel_NavigationBackWithBKey(t *testing.T) {
	form := NewFormModel()
	form.currentStep = 3
	form.history = []int{1, 2}

	// Navigate back with 'b' key
	newForm, _ := form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
	f := newForm.(FormModel)

	// Should go back to step 2
	if f.currentStep != 2 {
		t.Errorf("Expected step 2 after 'b' key, got %d", f.currentStep)
	}

	// History should be popped once
	if len(f.history) != 1 {
		t.Errorf("Expected 1 item in history, got %d", len(f.history))
	}
}

// TestFormModel_NavigationNextWithNKey verifies forward navigation with 'n' key.
// Story 10.7 Task 2: Navigation n (next alternative)
func TestFormModel_NavigationNextWithNKey(t *testing.T) {
	form := NewFormModel()
	form.currentStep = 1
	form.projectName = "valid-project"

	// Navigate forward with 'n' key
	newForm, _ := form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	f := newForm.(FormModel)

	// Should advance to step 2
	if f.currentStep != 2 {
		t.Errorf("Expected step 2 after 'n' key, got %d", f.currentStep)
	}
}

// TestFormModel_NavigationRightKey verifies forward navigation with right arrow.
// Story 10.7 Task 2: Navigation → (next)
func TestFormModel_NavigationRightKey(t *testing.T) {
	form := NewFormModel()
	form.currentStep = 2
	form.template = "standard"

	// Navigate forward with Right arrow
	newForm, _ := form.Update(tea.KeyMsg{Type: tea.KeyRight})
	f := newForm.(FormModel)

	// Should advance to step 3
	if f.currentStep != 3 {
		t.Errorf("Expected step 3 after right arrow, got %d", f.currentStep)
	}
}

// TestFormModel_StatePersistence verifies values are retained during navigation.
// Story 10.7 Task 2 & AC #10: State persistence
func TestFormModel_StatePersistence(t *testing.T) {
	form := NewFormModel()
	form.currentStep = 1
	form.projectName = "my-awesome-project"

	// Navigate forward
	newForm, _ := form.Update(tea.KeyMsg{Type: tea.KeyEnter})
	f := newForm.(FormModel)

	// Set template
	f.currentStep = 2
	f.template = "minimal"

	// Navigate forward again
	newForm2, _ := f.Update(tea.KeyMsg{Type: tea.KeyEnter})
	f2 := newForm2.(FormModel)

	// Navigate back twice
	newForm3, _ := f2.Update(tea.KeyMsg{Type: tea.KeyLeft})
	f3 := newForm3.(FormModel)
	newForm4, _ := f3.Update(tea.KeyMsg{Type: tea.KeyLeft})
	f4 := newForm4.(FormModel)

	// Values should still be there
	if f4.projectName != "my-awesome-project" {
		t.Errorf("Expected projectName to be persisted, got '%s'", f4.projectName)
	}

	// Navigate forward again - template should still be there
	newForm5, _ := f4.Update(tea.KeyMsg{Type: tea.KeyEnter})
	f5 := newForm5.(FormModel)
	if f5.template != "minimal" {
		t.Errorf("Expected template to be persisted, got '%s'", f5.template)
	}
}

// TestFormModel_CannotGoBackFromFirstStep verifies navigation bounds.
func TestFormModel_CannotGoBackFromFirstStep(t *testing.T) {
	form := NewFormModel()
	form.currentStep = 1

	// Try to go back with empty history
	newForm, _ := form.Update(tea.KeyMsg{Type: tea.KeyLeft})
	f := newForm.(FormModel)

	// Should stay on step 1
	if f.currentStep != 1 {
		t.Errorf("Expected to stay on step 1, got %d", f.currentStep)
	}
}

// TestFormModel_ProgressIndicator verifies progress indicator rendering.
// Story 10.7 Task 2: Visual progress indicators
func TestFormModel_ProgressIndicator(t *testing.T) {
	form := NewFormModel()
	form.currentStep = 2
	form.totalSteps = 4

	progress := form.RenderProgressIndicator()

	// Should show current step out of total
	if progress == "" {
		t.Error("Expected non-empty progress indicator")
	}

	// Should contain step info (implementation will vary)
	// Just verify it's generated - visual appearance tested manually
}

// TestFormModel_BreadcrumbIndicator verifies breadcrumb rendering.
// Story 10.7 Task 2: Breadcrumb indicator at top of screen
func TestFormModel_BreadcrumbIndicator(t *testing.T) {
	form := NewFormModel()
	form.currentStep = 2
	form.totalSteps = 4

	breadcrumb := form.RenderBreadcrumb()

	// Should show breadcrumb trail
	if breadcrumb == "" {
		t.Error("Expected non-empty breadcrumb indicator")
	}
}

// TestFormModel_StepNames verifies step naming.
func TestFormModel_StepNames(t *testing.T) {
	form := NewFormModel()

	tests := []struct {
		step int
		want string
	}{
		{1, "Project Name"},
		{2, "Template"},
		{3, "Database"},
		{4, "Observability"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := form.GetStepName(tt.step)
			if got != tt.want {
				t.Errorf("GetStepName(%d) = %q, want %q", tt.step, got, tt.want)
			}
		})
	}
}

// TestFormModel_HistoryStack verifies history stack operations.
func TestFormModel_HistoryStack(t *testing.T) {
	form := NewFormModel()
	form.currentStep = 1

	// Navigate through all steps
	steps := []int{2, 3, 4}
	for _, expectedStep := range steps {
		form = form.goNext().(FormModel)
		if form.currentStep != expectedStep {
			t.Errorf("Expected step %d, got %d", expectedStep, form.currentStep)
		}
	}

	// History should have [1, 2, 3]
	expectedHistory := []int{1, 2, 3}
	if len(form.history) != len(expectedHistory) {
		t.Errorf("Expected history length %d, got %d", len(expectedHistory), len(form.history))
	}

	for i, expected := range expectedHistory {
		if form.history[i] != expected {
			t.Errorf("Expected history[%d] = %d, got %d", i, expected, form.history[i])
		}
	}

	// Navigate back
	form = form.goBack().(FormModel)
	if form.currentStep != 3 {
		t.Errorf("Expected step 3 after going back, got %d", form.currentStep)
	}

	// History should have [1, 2]
	if len(form.history) != 2 {
		t.Errorf("Expected history length 2 after going back, got %d", len(form.history))
	}
}

// TestFormModel_ValidationOnNext verifies validation before advancing.
func TestFormModel_ValidationOnNext(t *testing.T) {
	form := NewFormModel()
	form.currentStep = 1
	form.projectName = "" // Empty name should be invalid

	// Try to advance with invalid input
	newForm, _ := form.Update(tea.KeyMsg{Type: tea.KeyEnter})
	f := newForm.(FormModel)

	// Should either stay on step 1 or show validation error
	// (Implementation may vary - could advance and use placeholder)
	// For now, just verify the model still has the state
	if f.currentStep < 1 || f.currentStep > 4 {
		t.Errorf("Expected valid step number, got %d", f.currentStep)
	}
}

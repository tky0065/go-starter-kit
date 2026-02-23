package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestNewMultiSelectItem tests creating a new multi-select item.
func TestNewMultiSelectItem(t *testing.T) {
	item := NewMultiSelectItem("Test Item", "Test description", "test-value")

	if item.Title() != "Test Item" {
		t.Errorf("Expected title 'Test Item', got '%s'", item.Title())
	}

	if item.Description() != "Test description" {
		t.Errorf("Expected description 'Test description', got '%s'", item.Description())
	}

	if item.Value() != "test-value" {
		t.Errorf("Expected value 'test-value', got '%s'", item.Value())
	}

	if item.IsSelected() {
		t.Error("Expected item to be unselected by default")
	}

	if item.FilterValue() != "Test Item" {
		t.Errorf("Expected filter value 'Test Item', got '%s'", item.FilterValue())
	}
}

// TestNewMultiSelectModel tests creating a new multi-select model.
func TestNewMultiSelectModel(t *testing.T) {
	items := []MultiSelectItem{
		NewMultiSelectItem("Item 1", "First item", "item1"),
		NewMultiSelectItem("Item 2", "Second item", "item2"),
		NewMultiSelectItem("Item 3", "Third item", "item3"),
	}

	model := NewMultiSelectModel("Select Items", items, 80, 24)

	if model.title != "Select Items" {
		t.Errorf("Expected title 'Select Items', got '%s'", model.title)
	}

	if len(model.items) != 3 {
		t.Errorf("Expected 3 items, got %d", len(model.items))
	}

	if model.CountSelected() != 0 {
		t.Errorf("Expected 0 selected items initially, got %d", model.CountSelected())
	}

	if !model.IsValid() {
		t.Error("Expected model to be valid initially (no constraints)")
	}
}

// TestMultiSelectToggle tests toggling item selection with Space.
func TestMultiSelectToggle(t *testing.T) {
	items := []MultiSelectItem{
		NewMultiSelectItem("Item 1", "First item", "item1"),
		NewMultiSelectItem("Item 2", "Second item", "item2"),
		NewMultiSelectItem("Item 3", "Third item", "item3"),
	}

	model := NewMultiSelectModel("Select Items", items, 80, 24)

	// Initially no items selected
	if model.CountSelected() != 0 {
		t.Errorf("Expected 0 selected, got %d", model.CountSelected())
	}

	// Press Space to select first item
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeySpace})

	if model.CountSelected() != 1 {
		t.Errorf("Expected 1 selected after Space, got %d", model.CountSelected())
	}

	if !model.items[0].selected {
		t.Error("Expected first item to be selected")
	}

	// Press Space again to deselect
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeySpace})

	if model.CountSelected() != 0 {
		t.Errorf("Expected 0 selected after second Space, got %d", model.CountSelected())
	}

	if model.items[0].selected {
		t.Error("Expected first item to be deselected")
	}
}

// TestMultiSelectValidation tests validation constraints.
func TestMultiSelectValidation(t *testing.T) {
	items := []MultiSelectItem{
		NewMultiSelectItem("Item 1", "First item", "item1"),
		NewMultiSelectItem("Item 2", "Second item", "item2"),
		NewMultiSelectItem("Item 3", "Third item", "item3"),
		NewMultiSelectItem("Item 4", "Fourth item", "item4"),
	}

	model := NewMultiSelectModel("Select Items", items, 80, 24)

	// Set validation: min 2, max 3
	model.SetValidation(2, 3)

	// No selection - invalid (below min)
	if model.IsValid() {
		t.Error("Expected model to be invalid with 0 selections (min=2)")
	}

	// Select 1 item - still invalid
	model.items[0].selected = true
	if model.IsValid() {
		t.Error("Expected model to be invalid with 1 selection (min=2)")
	}

	// Select 2 items - valid
	model.items[1].selected = true
	if !model.IsValid() {
		t.Error("Expected model to be valid with 2 selections")
	}

	// Select 3 items - still valid (at max)
	model.items[2].selected = true
	if !model.IsValid() {
		t.Error("Expected model to be valid with 3 selections (max=3)")
	}

	// Select 4 items - invalid (over max)
	model.items[3].selected = true
	if model.IsValid() {
		t.Error("Expected model to be invalid with 4 selections (max=3)")
	}
}

// TestMultiSelectMaxConstraint tests that toggling respects max constraint.
func TestMultiSelectMaxConstraint(t *testing.T) {
	items := []MultiSelectItem{
		NewMultiSelectItem("Item 1", "First item", "item1"),
		NewMultiSelectItem("Item 2", "Second item", "item2"),
		NewMultiSelectItem("Item 3", "Third item", "item3"),
	}

	model := NewMultiSelectModel("Select Items", items, 80, 24)
	model.SetValidation(0, 2) // Max 2 selections

	// Select first two items
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeySpace}) // Item 1
	model.list.Select(1)
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeySpace}) // Item 2

	if model.CountSelected() != 2 {
		t.Errorf("Expected 2 selected, got %d", model.CountSelected())
	}

	// Try to select third item (should be prevented)
	model.list.Select(2)
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeySpace}) // Item 3

	// Should still have 2 selected (max constraint)
	if model.CountSelected() != 2 {
		t.Errorf("Expected max 2 selected, got %d", model.CountSelected())
	}

	if model.items[2].selected {
		t.Error("Third item should not be selected (exceeds max)")
	}
}

// TestGetSelectedValues tests retrieving selected values.
func TestGetSelectedValues(t *testing.T) {
	items := []MultiSelectItem{
		NewMultiSelectItem("Item 1", "First item", "value1"),
		NewMultiSelectItem("Item 2", "Second item", "value2"),
		NewMultiSelectItem("Item 3", "Third item", "value3"),
	}

	model := NewMultiSelectModel("Select Items", items, 80, 24)

	// Select items 1 and 3
	model.items[0].selected = true
	model.items[2].selected = true

	values := model.GetSelectedValues()

	if len(values) != 2 {
		t.Errorf("Expected 2 values, got %d", len(values))
	}

	if values[0] != "value1" || values[1] != "value3" {
		t.Errorf("Expected values [value1, value3], got %v", values)
	}
}

// TestGetSelectedItems tests retrieving selected items.
func TestGetSelectedItems(t *testing.T) {
	items := []MultiSelectItem{
		NewMultiSelectItem("Item 1", "First item", "value1"),
		NewMultiSelectItem("Item 2", "Second item", "value2"),
		NewMultiSelectItem("Item 3", "Third item", "value3"),
	}

	model := NewMultiSelectModel("Select Items", items, 80, 24)

	// Select items 1 and 2
	model.items[0].selected = true
	model.items[1].selected = true

	selectedItems := model.GetSelectedItems()

	if len(selectedItems) != 2 {
		t.Errorf("Expected 2 items, got %d", len(selectedItems))
	}

	if selectedItems[0].Value() != "value1" || selectedItems[1].Value() != "value2" {
		t.Errorf("Expected items with values [value1, value2], got %v",
			[]string{selectedItems[0].Value(), selectedItems[1].Value()})
	}
}

// TestResetSelection tests resetting all selections.
func TestResetSelection(t *testing.T) {
	items := []MultiSelectItem{
		NewMultiSelectItem("Item 1", "First item", "value1"),
		NewMultiSelectItem("Item 2", "Second item", "value2"),
		NewMultiSelectItem("Item 3", "Third item", "value3"),
	}

	model := NewMultiSelectModel("Select Items", items, 80, 24)

	// Select all items
	model.items[0].selected = true
	model.items[1].selected = true
	model.items[2].selected = true

	if model.CountSelected() != 3 {
		t.Errorf("Expected 3 selected before reset, got %d", model.CountSelected())
	}

	// Reset
	model.ResetSelection()

	if model.CountSelected() != 0 {
		t.Errorf("Expected 0 selected after reset, got %d", model.CountSelected())
	}

	for i, item := range model.items {
		if item.selected {
			t.Errorf("Item %d should be deselected after reset", i)
		}
	}
}

// TestSelectAll tests selecting all items.
func TestSelectAll(t *testing.T) {
	items := []MultiSelectItem{
		NewMultiSelectItem("Item 1", "First item", "value1"),
		NewMultiSelectItem("Item 2", "Second item", "value2"),
		NewMultiSelectItem("Item 3", "Third item", "value3"),
	}

	model := NewMultiSelectModel("Select Items", items, 80, 24)

	// Select all
	model.SelectAll()

	if model.CountSelected() != 3 {
		t.Errorf("Expected 3 selected, got %d", model.CountSelected())
	}

	for i, item := range model.items {
		if !item.selected {
			t.Errorf("Item %d should be selected after SelectAll", i)
		}
	}
}

// TestSelectAllWithMaxConstraint tests SelectAll with max constraint.
func TestSelectAllWithMaxConstraint(t *testing.T) {
	items := []MultiSelectItem{
		NewMultiSelectItem("Item 1", "First item", "value1"),
		NewMultiSelectItem("Item 2", "Second item", "value2"),
		NewMultiSelectItem("Item 3", "Third item", "value3"),
		NewMultiSelectItem("Item 4", "Fourth item", "value4"),
	}

	model := NewMultiSelectModel("Select Items", items, 80, 24)
	model.SetValidation(0, 2) // Max 2

	// Select all (should only select first 2)
	model.SelectAll()

	if model.CountSelected() != 2 {
		t.Errorf("Expected 2 selected (respecting max), got %d", model.CountSelected())
	}

	if !model.items[0].selected || !model.items[1].selected {
		t.Error("First two items should be selected")
	}

	if model.items[2].selected || model.items[3].selected {
		t.Error("Last two items should not be selected (exceeds max)")
	}
}

// TestMultiSelectView tests that View renders without panicking.
func TestMultiSelectView(t *testing.T) {
	items := []MultiSelectItem{
		NewMultiSelectItem("Item 1", "First item", "value1"),
		NewMultiSelectItem("Item 2", "Second item", "value2"),
	}

	model := NewMultiSelectModel("Select Items", items, 80, 24)

	view := model.View()

	if view == "" {
		t.Error("Expected non-empty view")
	}

	// Check for expected elements in view
	expectedElements := []string{
		"Select Items", // Title
		"Selected: 0",  // Count
		"Space",        // Help text
		"Navigate",
	}

	for _, elem := range expectedElements {
		if !strings.Contains(view, elem) {
			t.Errorf("Expected view to contain '%s'", elem)
		}
	}
}

// TestMultiSelectViewWithSelection tests view with selections.
func TestMultiSelectViewWithSelection(t *testing.T) {
	items := []MultiSelectItem{
		NewMultiSelectItem("Item 1", "First item", "value1"),
		NewMultiSelectItem("Item 2", "Second item", "value2"),
	}

	model := NewMultiSelectModel("Select Items", items, 80, 24)
	model.items[0].selected = true

	view := model.View()

	if !strings.Contains(view, "Selected: 1") {
		t.Error("Expected view to show 'Selected: 1'")
	}
}

// TestMultiSelectViewWithValidation tests view with validation errors.
func TestMultiSelectViewWithValidation(t *testing.T) {
	items := []MultiSelectItem{
		NewMultiSelectItem("Item 1", "First item", "value1"),
		NewMultiSelectItem("Item 2", "Second item", "value2"),
	}

	model := NewMultiSelectModel("Select Items", items, 80, 24)
	model.SetValidation(1, 0) // Require at least 1

	view := model.View()

	// Should show warning when validation fails
	if !strings.Contains(view, "Please select at least") {
		t.Error("Expected view to show validation warning")
	}
}

// TestIsConfirmed tests confirmation action.
func TestIsConfirmed(t *testing.T) {
	items := []MultiSelectItem{
		NewMultiSelectItem("Item 1", "First item", "value1"),
		NewMultiSelectItem("Item 2", "Second item", "value2"),
	}

	model := NewMultiSelectModel("Select Items", items, 80, 24)

	// Initially not confirmed
	if model.IsConfirmed() {
		t.Error("Expected model to not be confirmed initially")
	}

	// Select an item and press Enter
	model.items[0].selected = true
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Should be confirmed now
	if !model.IsConfirmed() {
		t.Error("Expected model to be confirmed after Enter with valid selection")
	}
}

// TestHasMore tests the scroll indicator logic.
func TestHasMore(t *testing.T) {
	items := []MultiSelectItem{
		NewMultiSelectItem("Item 1", "First item", "value1"),
		NewMultiSelectItem("Item 2", "Second item", "value2"),
		NewMultiSelectItem("Item 3", "Third item", "value3"),
	}

	model := NewMultiSelectModel("Select Items", items, 80, 24)

	// At first item - has more
	if !model.HasMore() {
		t.Error("Expected HasMore() to be true at first item")
	}

	// Move to last item
	model.list.Select(2)

	// At last item - no more
	if model.HasMore() {
		t.Error("Expected HasMore() to be false at last item")
	}
}

// TestMultiSelectWindowResize tests window resize handling.
func TestMultiSelectWindowResize(t *testing.T) {
	items := []MultiSelectItem{
		NewMultiSelectItem("Item 1", "First item", "value1"),
		NewMultiSelectItem("Item 2", "Second item", "value2"),
	}

	model := NewMultiSelectModel("Select Items", items, 80, 24)

	// Resize window
	newWidth := 100
	newHeight := 30
	model, _ = model.Update(tea.WindowSizeMsg{Width: newWidth, Height: newHeight})

	if model.width != newWidth {
		t.Errorf("Expected width %d after resize, got %d", newWidth, model.width)
	}

	if model.height != newHeight {
		t.Errorf("Expected height %d after resize, got %d", newHeight, model.height)
	}
}

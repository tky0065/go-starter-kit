package tui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MultiSelectItem represents an item in a multi-select list.
// Story 10.7 Task 9: Multi-select for plugins, optional features
type MultiSelectItem struct {
	title       string
	description string
	value       string
	selected    bool
}

// NewMultiSelectItem creates a new MultiSelectItem.
func NewMultiSelectItem(title, description, value string) MultiSelectItem {
	return MultiSelectItem{
		title:       title,
		description: description,
		value:       value,
		selected:    false,
	}
}

// Title implements list.Item interface.
func (i MultiSelectItem) Title() string { return i.title }

// Description implements list.Item interface.
func (i MultiSelectItem) Description() string { return i.description }

// FilterValue implements list.Item interface.
func (i MultiSelectItem) FilterValue() string { return i.title }

// Value returns the item's value (used when retrieving selected items).
func (i MultiSelectItem) Value() string { return i.value }

// IsSelected returns whether the item is selected.
func (i MultiSelectItem) IsSelected() bool { return i.selected }

// MultiSelectModel manages a multi-select list with pagination.
// Story 10.7 Task 9: Advanced components (AC #2, #7)
type MultiSelectModel struct {
	list          list.Model
	items         []MultiSelectItem       // Track selected state separately
	delegate      multiSelectItemDelegate // Keep reference to delegate
	width         int
	height        int
	title         string
	minSelected   int // Minimum number of selections (0 = optional)
	maxSelected   int // Maximum number of selections (0 = unlimited)
	confirmAction bool
}

// multiSelectItemDelegate is a custom delegate for rendering multi-select items.
type multiSelectItemDelegate struct {
	items []MultiSelectItem // Reference to parent items for selection state
}

func (d multiSelectItemDelegate) Height() int                               { return 2 }
func (d multiSelectItemDelegate) Spacing() int                              { return 1 }
func (d multiSelectItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d multiSelectItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(MultiSelectItem)
	if !ok {
		return
	}

	// Check if this item is selected
	isSelected := false
	if index < len(d.items) {
		isSelected = d.items[index].selected
	}

	// Style based on focus and selection
	var titleStyle, descStyle lipgloss.Style
	var checkbox string

	isFocused := index == m.Index()

	// Checkbox indicator
	if isSelected {
		checkbox = "[✓] " // Checked
	} else {
		checkbox = "[ ] " // Unchecked
	}

	// Apply styles
	if isFocused {
		titleStyle = FocusedStyle.Bold(true)
		descStyle = InfoStyle
		if isSelected {
			checkbox = SuccessStyle.Render(checkbox)
		} else {
			checkbox = InfoStyle.Render(checkbox)
		}
	} else {
		if isSelected {
			titleStyle = SuccessStyle
			descStyle = SuccessStyle
			checkbox = SuccessStyle.Render(checkbox)
		} else {
			titleStyle = BlurredStyle
			descStyle = MutedStyle
			checkbox = MutedStyle.Render(checkbox)
		}
	}

	title := titleStyle.Render(checkbox + i.title)
	desc := descStyle.Render("    " + i.description)

	fmt.Fprintf(w, "%s\n%s", title, desc)
}

// NewMultiSelectModel creates a new MultiSelectModel.
// Story 10.7 Task 9: Multi-select with pagination
func NewMultiSelectModel(title string, items []MultiSelectItem, width, height int) MultiSelectModel {
	// Convert items to list.Item interface
	listItems := make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = item
	}

	// Create delegate with reference to items
	delegate := multiSelectItemDelegate{items: items}

	// Create list with pagination
	l := list.New(listItems, delegate, width-4, height-10)
	l.Title = title
	l.SetShowStatusBar(true) // Show pagination info
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false) // Custom help

	return MultiSelectModel{
		list:        l,
		items:       items,
		delegate:    delegate, // Store delegate reference
		width:       width,
		height:      height,
		title:       title,
		minSelected: 0,
		maxSelected: 0, // Unlimited
	}
}

// SetValidation sets minimum and maximum selection constraints.
func (m *MultiSelectModel) SetValidation(min, max int) {
	m.minSelected = min
	m.maxSelected = max
}

// Update handles messages for the multi-select list.
func (m MultiSelectModel) Update(msg tea.Msg) (MultiSelectModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetWidth(msg.Width - 4)
		m.list.SetHeight(msg.Height - 10)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case " ": // Space to toggle selection
			currentIndex := m.list.Index()
			if currentIndex >= 0 && currentIndex < len(m.items) {
				// Toggle selection
				m.items[currentIndex].selected = !m.items[currentIndex].selected

				// Update delegate's reference
				m.delegate.items = m.items
				m.list.SetDelegate(m.delegate)

				// Check if we've hit max selections
				if m.maxSelected > 0 && m.CountSelected() > m.maxSelected {
					// Deselect if over limit
					m.items[currentIndex].selected = false
					m.delegate.items = m.items
					m.list.SetDelegate(m.delegate)
				}

				return m, nil
			}

		case "enter":
			// Confirm selection (validate)
			if m.IsValid() {
				m.confirmAction = true
				return m, nil
			}
			// Invalid selection, don't confirm
			return m, nil
		}
	}

	// Update list for navigation
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View renders the multi-select list.
func (m MultiSelectModel) View() string {
	var b strings.Builder

	b.WriteString("\n")

	// Render list (includes title and pagination)
	b.WriteString(m.list.View())
	b.WriteString("\n\n")

	// Selection count
	selectedCount := m.CountSelected()
	countText := fmt.Sprintf("Selected: %d", selectedCount)
	if m.maxSelected > 0 {
		countText += fmt.Sprintf("/%d", m.maxSelected)
	}
	b.WriteString(InfoStyle.Render(countText))
	b.WriteString("\n\n")

	// Validation message
	if !m.IsValid() {
		if m.minSelected > 0 && selectedCount < m.minSelected {
			errMsg := fmt.Sprintf("⚠ Please select at least %d item(s)", m.minSelected)
			b.WriteString(WarningStyle.Render(errMsg))
			b.WriteString("\n\n")
		}
	}

	// Help footer
	footer := "Space: Toggle • ↑↓: Navigate • Enter: Confirm"
	if m.HasMore() {
		footer += " • ↓ More..."
	}
	b.WriteString(RenderFooter(footer))
	b.WriteString("\n")

	return b.String()
}

// CountSelected returns the number of selected items.
func (m MultiSelectModel) CountSelected() int {
	count := 0
	for _, item := range m.items {
		if item.selected {
			count++
		}
	}
	return count
}

// GetSelectedValues returns the values of all selected items.
func (m MultiSelectModel) GetSelectedValues() []string {
	var values []string
	for _, item := range m.items {
		if item.selected {
			values = append(values, item.value)
		}
	}
	return values
}

// GetSelectedItems returns all selected items.
func (m MultiSelectModel) GetSelectedItems() []MultiSelectItem {
	var selected []MultiSelectItem
	for _, item := range m.items {
		if item.selected {
			selected = append(selected, item)
		}
	}
	return selected
}

// IsValid checks if the current selection meets validation constraints.
func (m MultiSelectModel) IsValid() bool {
	count := m.CountSelected()

	// Check minimum
	if m.minSelected > 0 && count < m.minSelected {
		return false
	}

	// Check maximum
	if m.maxSelected > 0 && count > m.maxSelected {
		return false
	}

	return true
}

// HasMore returns true if there are more items below the current viewport.
// Story 10.7 Task 9: Scroll indicators
func (m MultiSelectModel) HasMore() bool {
	// Check if we're not at the bottom of the list
	totalItems := len(m.items)
	if totalItems == 0 {
		return false
	}

	// If we're at the last item, no more items
	return m.list.Index() < totalItems-1
}

// IsConfirmed returns true if the user confirmed the selection.
func (m MultiSelectModel) IsConfirmed() bool {
	return m.confirmAction
}

// ResetSelection deselects all items.
func (m *MultiSelectModel) ResetSelection() {
	for i := range m.items {
		m.items[i].selected = false
	}
	// Update delegate
	m.delegate.items = m.items
	m.list.SetDelegate(m.delegate)
}

// SelectAll selects all items (respecting maxSelected if set).
func (m *MultiSelectModel) SelectAll() {
	maxToSelect := len(m.items)
	if m.maxSelected > 0 && m.maxSelected < maxToSelect {
		maxToSelect = m.maxSelected
	}

	for i := 0; i < maxToSelect; i++ {
		m.items[i].selected = true
	}

	// Update delegate
	m.delegate.items = m.items
	m.list.SetDelegate(m.delegate)
}

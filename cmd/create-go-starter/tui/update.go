package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tky0065/go-starter-kit/pkg/utils"
)

// Update handles incoming messages and updates the model accordingly.
// This is the core of the Elm Architecture - a pure function that transforms state.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Handle keyboard input
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)

	// Handle window resize events
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Adapt styles to terminal width for responsiveness
		AdaptToTerminalWidth(m.width)

		// Update active list dimensions for proper responsiveness (MEDIUM-5 fix)
		switch m.state {
		case StateWelcome:
			m.welcomeList.SetSize(m.width-ListWidthMargin, 10) // Story 10.7: Welcome list sizing
		case StateTemplateSelect:
			m.templateList.SetSize(m.width, m.height-ListHeightOffset)
		case StateDatabaseSelect:
			m.databaseList.SetSize(m.width, m.height-ListHeightOffset)
		case StateFrameworkSelect:
			m.frameworkList.SetSize(m.width, m.height-ListHeightOffset)
		case StateObservabilitySelect:
			m.obsList.SetSize(m.width, m.height-ListHeightOffset)
		case StatePreview:
			// Update preview model viewport size (Story 10.7 Task 6)
			m.previewModel, _ = m.previewModel.Update(msg)
		}
		return m, nil

	// Handle custom application messages
	case WelcomeMenuSelectedMsg:
		// Handle welcome menu selection (Story 10.7 Task 1.3)
		switch msg.Action {
		case "create":
			m.state = StateProjectName
			return m, nil
		case "help":
			m.previousState = StateWelcome
			m.state = StateHelp
			return m, nil
		case "exit":
			return m, tea.Quit
		default:
			return m, nil
		}

	case ProjectNameSubmittedMsg:
		// Validate project name (MEDIUM-2 fix)
		if err := utils.ValidateGoModuleName(msg.Name); err != nil {
			// Invalid name - store error and stay in StateProjectName
			m.validationError = err.Error()
			// Keep the current state, don't transition
			return m, nil
		}

		// Valid name - clear any previous error and proceed
		m.validationError = ""
		m.projectName = msg.Name
		m.state = StateTemplateSelect
		m.templateList = initializeTemplateList(m.template)
		m.templateList.SetSize(m.width, m.height-ListHeightOffset)
		return m, nil

	case TemplateSelectedMsg:
		m.template = msg.Template
		m.state = StateDatabaseSelect
		m.databaseList = initializeDatabaseList(m.database)
		m.databaseList.SetSize(m.width, m.height-ListHeightOffset)
		return m, nil

	case DatabaseSelectedMsg:
		m.database = msg.Database
		m.state = StateFrameworkSelect
		m.frameworkList = initializeFrameworkList(m.framework)
		m.frameworkList.SetSize(m.width, m.height-ListHeightOffset)
		return m, nil

	case FrameworkSelectedMsg:
		m.framework = msg.Framework
		m.state = StateObservabilitySelect
		m.obsList = initializeObservabilityList(m.observability)
		m.obsList.SetSize(m.width, m.height-ListHeightOffset)
		return m, nil

	case ObservabilitySelectedMsg:
		m.observability = msg.Observability
		m.state = StateSummary
		return m, nil

	case ConfirmGenerationMsg:
		m.state = StateGenerating
		// Initialize generation stats (Story 10.7 Task 3)
		m.genStats = GenerationStats{
			FilesCreated: 0,
			TotalSize:    0,
			StartTime:    time.Now(),
			CurrentFile:  "",
			CurrentStep:  "Initializing...",
		}
		// Reset progress pulse animation (Story 10.7 Task 10)
		if m.progressPulse != nil {
			m.progressPulse.Reset()
		}
		// Start generation process and progress pulse animation
		return m, tea.Batch(
			m.generateProjectCmd(),
			ProgressPulseTick(50*time.Millisecond), // 50ms = 20 FPS
		)

	case FileGeneratedMsg:
		m.filesGenerated = msg.Index
		// Update generation stats (Story 10.7 Task 3)
		m.genStats.FilesCreated = msg.Index
		m.genStats.TotalSize += msg.Size
		m.genStats.CurrentFile = msg.FilePath
		if msg.Step != "" {
			m.genStats.CurrentStep = msg.Step
		}
		// Update progress bar percentage (MEDIUM-4 fix)
		if m.totalFiles > 0 {
			percent := float64(m.filesGenerated) / float64(m.totalFiles)
			m.progressBar.SetPercent(percent)
		}
		return m, nil

	case GenerationCompleteMsg:
		// Store results
		m.totalFiles = msg.TotalFiles
		m.filesGenerated = msg.TotalFiles // All files generated when complete
		m.err = msg.Error

		// Transition to done state - viewDone() handles both success and error display
		// No need for separate StateError - viewDone() checks m.err and renders accordingly
		m.state = StateDone
		return m, nil

	case NavigateBackMsg:
		return m.navigateBack(), nil

	case AnimationTickMsg:
		// Advance logo fade-in animation (Story 10.7 Task 1.5)
		// Only animate if we're on the welcome screen and animation is not done
		if m.state == StateWelcome && m.logoAnimation != nil && !m.logoAnimation.IsDone() {
			m.logoAnimation.Advance()
			// Continue ticking if animation is not done
			if !m.logoAnimation.IsDone() {
				return m, AnimationTick(msg.Elapsed)
			}
		}

		// Advance screen transition animation (Story 10.7 Task 10)
		if m.screenTransition != nil && !m.screenTransition.IsDone() {
			m.screenTransition.Advance()
			// Continue ticking if transition is not done
			if !m.screenTransition.IsDone() {
				return m, AnimationTick(msg.Elapsed)
			}
		}
		return m, nil

	case ProgressPulseTickMsg:
		// Advance progress bar pulse animation (Story 10.7 Task 10)
		// Only animate during generation
		if m.state == StateGenerating && m.progressPulse != nil {
			m.progressPulse.Advance()
			// Continue pulsing during generation
			return m, ProgressPulseTick(msg.Elapsed)
		}
		return m, nil

	default:
		// Update bubbles components based on current state
		return m.updateComponents(msg)
	}
}

// handleKeyMsg processes keyboard input based on the current state.
func (m Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle help key '?' globally (HIGH-1 fix - AC#4 compliance)
	if msg.String() == "?" && m.state != StateHelp && m.state != StateGenerating {
		m.previousState = m.state
		m.state = StateHelp
		return m, nil
	}

	// Handle return from help screen
	if m.state == StateHelp {
		// Any key returns from help
		m.state = m.previousState
		return m, nil
	}

	switch msg.Type {
	case tea.KeyCtrlC, tea.KeyEsc:
		// Always allow quitting with Ctrl+C or Esc
		return m, tea.Quit

	case tea.KeyEnter:
		// Handle Enter key based on current state
		switch m.state {
		case StateDone:
			// "Press any key to exit" — Enter also quits
			return m, tea.Quit

		case StateWelcome:
			// Get selected item from welcome menu list (Story 10.7 Task 1.3)
			selectedItem := m.welcomeList.SelectedItem()
			if item, ok := selectedItem.(welcomeMenuItem); ok {
				// Handle action based on selected menu item
				switch item.action {
				case "create":
					return m.Update(WelcomeMenuSelectedMsg{Action: "create"})
				case "help":
					return m.Update(WelcomeMenuSelectedMsg{Action: "help"})
				case "exit":
					return m, tea.Quit
				}
			}
			// Fallback to default behavior if item type assertion fails
			return m.Update(WelcomeMenuSelectedMsg{Action: "create"})

		case StateProjectName:
			// Submit project name
			name := m.projectInput.Value()
			if name == "" {
				name = m.projectInput.Placeholder
			}
			return m.Update(ProjectNameSubmittedMsg{Name: name})

		case StateTemplateSelect:
			// Submit template selection
			selectedItem := m.templateList.SelectedItem()
			if item, ok := selectedItem.(templateItem); ok {
				return m.Update(TemplateSelectedMsg{Template: item.name})
			}

		case StateDatabaseSelect:
			// Submit database selection
			selectedItem := m.databaseList.SelectedItem()
			if item, ok := selectedItem.(databaseItem); ok {
				return m.Update(DatabaseSelectedMsg{Database: item.name})
			}

		case StateFrameworkSelect:
			// Submit framework selection
			selectedItem := m.frameworkList.SelectedItem()
			if item, ok := selectedItem.(frameworkItem); ok {
				return m.Update(FrameworkSelectedMsg{Framework: item.name})
			}

		case StateObservabilitySelect:
			// Submit observability selection
			selectedItem := m.obsList.SelectedItem()
			if item, ok := selectedItem.(observabilityItem); ok {
				return m.Update(ObservabilitySelectedMsg{Observability: item.name})
			}

		case StateSummary:
			// Confirm generation
			return m.Update(ConfirmGenerationMsg{})
		}
	}

	// In StateDone, any key press should quit (matches "Press any key to exit" prompt)
	if m.state == StateDone {
		return m, tea.Quit
	}

	// Delegate to component-specific key handling
	return m.updateComponents(msg)
}

// updateComponents updates the active bubbles components based on the message.
func (m Model) updateComponents(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.state {
	case StateWelcome:
		// Update welcome menu list for arrow key navigation (Story 10.7 Task 1.3)
		m.welcomeList, cmd = m.welcomeList.Update(msg)

	case StateProjectName:
		m.projectInput, cmd = m.projectInput.Update(msg)

	case StateTemplateSelect:
		m.templateList, cmd = m.templateList.Update(msg)

	case StateDatabaseSelect:
		m.databaseList, cmd = m.databaseList.Update(msg)

	case StateFrameworkSelect:
		m.frameworkList, cmd = m.frameworkList.Update(msg)

	case StateObservabilitySelect:
		m.obsList, cmd = m.obsList.Update(msg)

	case StatePreview:
		// Update preview model viewport (Story 10.7 Task 6)
		m.previewModel, cmd = m.previewModel.Update(msg)

	case StateGenerating:
		// Update spinner animation
		m.loadSpinner, cmd = m.loadSpinner.Update(msg)
		// Progress bar returns tea.Model, so we need a type assertion
		progModel, progCmd := m.progressBar.Update(msg)
		m.progressBar = progModel.(progress.Model)
		// Batch both commands
		return m, tea.Batch(cmd, progCmd)
	}

	return m, cmd
}

// navigateBack moves to the previous state/screen.
func (m Model) navigateBack() Model {
	switch m.state {
	case StateProjectName:
		m.state = StateWelcome
	case StateTemplateSelect:
		m.state = StateProjectName
	case StateDatabaseSelect:
		m.state = StateTemplateSelect
	case StateFrameworkSelect:
		m.state = StateDatabaseSelect
	case StateObservabilitySelect:
		m.state = StateFrameworkSelect
	case StateSummary:
		m.state = StateObservabilitySelect
	case StatePreview:
		// Go back to summary from preview (Story 10.7 Task 6)
		m.state = StateSummary
	}
	return m
}

// generateProjectCmd creates a Cmd that runs the project generation.
// This runs in a goroutine and sends FileGeneratedMsg via p.Send() for real-time progress,
// then returns a GenerationCompleteMsg when done.
func (m *Model) generateProjectCmd() tea.Cmd {
	return func() tea.Msg {
		// Get the program reference for sending progress messages
		// The channel is buffered(1) and pre-filled by RunInteractiveTUI
		var p *tea.Program
		select {
		case p = <-m.programChan:
			// Got the program reference
		default:
			// No program available (e.g., in tests) — progress updates will be skipped
		}

		// Track total files for final report
		totalFiles := 0

		// Run generation synchronously within this Cmd goroutine
		err := m.generatorFunc(
			m.projectName,
			m.template,
			m.database,
			m.observability,
			m.framework,
			func(current, total int) {
				// Store total on first callback
				if totalFiles == 0 {
					totalFiles = total
				}
				// Send progress update to the Bubble Tea runtime via p.Send()
				// This is the correct pattern for goroutine-to-runtime communication
				if p != nil {
					p.Send(FileGeneratedMsg{
						FilePath: "",
						Index:    current,
						Size:     0,
						Step:     fmt.Sprintf("Generating file %d/%d...", current, total),
					})
				}
			},
		)

		// Return completion message (this is the final message from this Cmd)
		return GenerationCompleteMsg{
			TotalFiles: totalFiles,
			Success:    err == nil,
			Error:      err,
		}
	}
}

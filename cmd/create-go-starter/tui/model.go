package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// GeneratorFunc is a function that generates a project.
// It receives the project configuration and a progress callback.
// The progress callback is called with (current, total) for each file generated.
// The function should call progressCallback(current, total) for progress updates.
type GeneratorFunc func(projectName, template, database, observability string, progressCallback func(current, total int)) error

// UI dimension constants
const (
	ListHeightOffset = 4 // Space reserved for header/footer
	ListWidthMargin  = 4 // Side margins for content
)

// AppState represents the current state/screen of the application.
type AppState int

const (
	// StateWelcome is the initial welcome screen with logo and menu
	StateWelcome AppState = iota
	// StateProjectName is the state for collecting the project name
	StateProjectName
	// StateTemplateSelect is the state for selecting a template
	StateTemplateSelect
	// StateDatabaseSelect is the state for selecting a database
	StateDatabaseSelect
	// StateObservabilitySelect is the state for selecting observability level
	StateObservabilitySelect
	// StateSummary is the state for showing configuration summary
	StateSummary
	// StatePreview is the state for showing dry-run preview (Story 10.7 Task 6)
	StatePreview
	// StateGenerating is the state for showing generation progress
	StateGenerating
	// StateDone is the final state after generation completes
	StateDone
	// StateHelp is the state for showing the help screen (HIGH-1 fix)
	StateHelp
)

// Model is the main application model following Elm Architecture.
// It represents the complete state of the TUI application.
type Model struct {
	// Current state/screen
	state AppState

	// User selections
	projectName   string
	template      string
	database      string
	observability string

	// Bubbles components for different screens
	welcomeList  list.Model      // For welcome menu (Story 10.7 Task 1.3)
	projectInput textinput.Model // For project name input
	templateList list.Model      // For template selection
	databaseList list.Model      // For database selection
	obsList      list.Model      // For observability selection
	previewModel PreviewModel    // For dry-run preview (Story 10.7 Task 6)
	progressBar  progress.Model  // For generation progress
	loadSpinner  spinner.Model   // For loading indicator

	// UI state
	width           int    // Terminal width (from tea.WindowSizeMsg)
	height          int    // Terminal height (from tea.WindowSizeMsg)
	err             error  // Any error that occurred during generation
	validationError string // Validation error for current input (shown inline)

	// Animation state (Story 10.7 Task 1.5 & Task 10)
	logoAnimation    *AnimationState  // Logo fade-in animation
	progressPulse    *PulseState      // Progress bar gradient pulse animation (Task 10)
	screenTransition *TransitionState // Screen transition animation (Task 10)

	// Help screen state
	previousState AppState // State to return to after help (HIGH-1 fix)

	// Generation state
	filesGenerated int // Number of files generated
	totalFiles     int // Total files to generate

	// Generation statistics (Story 10.7 Task 3)
	genStats GenerationStats // Real-time statistics during generation

	// Generation function (dependency injection to avoid circular imports)
	generatorFunc GeneratorFunc // Function that performs actual project generation
}

// NewModel creates a new Model with initial state.
// This is the entry point for creating the TUI application.
func NewModel() Model {
	// Initialize text input for project name
	ti := textinput.New()
	ti.Placeholder = "my-project"
	ti.Focus()
	ti.CharLimit = 100

	// Initialize spinner for loading indicator
	sp := spinner.New()
	sp.Spinner = spinner.Dot

	// Initialize progress bar
	pb := progress.New(progress.WithDefaultGradient())

	// Initialize welcome menu list (Story 10.7 Task 1.3)
	wl := initializeWelcomeList()
	wl.SetSize(80-ListWidthMargin, 10) // Default size, will be adjusted on window resize

	// Initialize logo fade-in animation (Story 10.7 Task 1.5)
	// 20 frames at 50ms each = 1 second animation
	logoAnim := NewAnimationState(20)

	// Initialize progress bar pulse animation (Story 10.7 Task 10)
	// 30 frames at 50ms = 1.5 second pulse cycle, opacity range 0.7 to 1.0
	progressPulse := NewPulseState(30, 0.7, 1.0)

	return Model{
		state:         StateWelcome, // Start with welcome screen (Task 1 - Story 10.7)
		welcomeList:   wl,
		projectInput:  ti,
		loadSpinner:   sp,
		progressBar:   pb,
		logoAnimation: logoAnim,
		progressPulse: progressPulse,
		width:         80, // Default width
		height:        24, // Default height
	}
}

// Init initializes the model and returns the initial command.
// This is called once when the program starts.
func (m Model) Init() tea.Cmd {
	// Start text input cursor blinking, spinner animation, and logo fade-in animation
	// Story 10.7 Task 1.5: Logo animation
	return tea.Batch(
		textinput.Blink,
		m.loadSpinner.Tick,
		AnimationTick(50*time.Millisecond), // 50ms = 20 FPS
	)
}

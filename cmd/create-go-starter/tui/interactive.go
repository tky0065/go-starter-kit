package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// InteractiveDefaults holds default values for interactive mode prompts.
// This allows CLI flags to pre-populate the interactive TUI with defaults.
type InteractiveDefaults struct {
	ProjectName   string
	Template      string
	Database      string
	Observability string
	Framework     string
}

// NewInteractiveModel creates a new Model configured for interactive mode.
// It initializes all the necessary components and sets up default values.
// The generatorFunc is used to perform actual project generation (dependency injection).
func NewInteractiveModel(defaults InteractiveDefaults, generatorFunc GeneratorFunc) Model {
	// Initialize text input for project name
	ti := textinput.New()
	if defaults.ProjectName != "" {
		ti.SetValue(defaults.ProjectName)
		ti.Placeholder = defaults.ProjectName
	} else {
		ti.Placeholder = "my-project"
	}
	ti.Focus()
	ti.CharLimit = 100

	// Initialize spinner for loading indicator (MEDIUM-3 fix)
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

	// Initialize the model
	m := Model{
		state:         StateWelcome, // Start with welcome screen (Task 1 - Story 10.7)
		welcomeList:   wl,
		projectInput:  ti,
		loadSpinner:   sp,
		progressBar:   pb,
		logoAnimation: logoAnim,
		progressPulse: progressPulse,
		programChan:   make(chan *tea.Program, 1), // Buffered so Send() doesn't block
		width:         80,                         // Default width (will be updated by WindowSizeMsg)
		height:        24,                         // Default height (will be updated by WindowSizeMsg)
		generatorFunc: generatorFunc,
	}

	// Store defaults for later use when initializing lists
	m.projectName = defaults.ProjectName
	m.template = defaults.Template
	m.database = defaults.Database
	m.observability = defaults.Observability
	m.framework = defaults.Framework

	return m
}

// initializeTemplateList creates and initializes the template selection list.
func initializeTemplateList(defaultTemplate string) list.Model {
	items := []list.Item{
		templateItem{name: "full", desc: "Complete REST API with user auth, database, logging, Swagger"},
		templateItem{name: "minimal", desc: "Basic HTTP server with health check - lightweight start"},
		templateItem{name: "graphql", desc: "GraphQL API server with subscriptions and schema"},
	}

	delegate := list.NewDefaultDelegate()
	l := list.New(items, delegate, 0, 0)
	l.Title = "Select a Template"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	// Set default selection if provided
	if defaultTemplate != "" {
		for i, item := range items {
			if tpl, ok := item.(templateItem); ok && tpl.name == defaultTemplate {
				l.Select(i)
				break
			}
		}
	}

	return l
}

// initializeDatabaseList creates and initializes the database selection list.
func initializeDatabaseList(defaultDatabase string) list.Model {
	items := []list.Item{
		databaseItem{name: "postgres", desc: "PostgreSQL - Production-ready, advanced features"},
		databaseItem{name: "mysql", desc: "MySQL/MariaDB - Wide compatibility, shared hosting"},
		databaseItem{name: "sqlite", desc: "SQLite - Quick prototyping, embedded apps"},
	}

	delegate := list.NewDefaultDelegate()
	l := list.New(items, delegate, 0, 0)
	l.Title = "Select a Database"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	// Set default selection if provided
	if defaultDatabase != "" {
		for i, item := range items {
			if db, ok := item.(databaseItem); ok && db.name == defaultDatabase {
				l.Select(i)
				break
			}
		}
	}

	return l
}

// initializeObservabilityList creates and initializes the observability selection list.
func initializeObservabilityList(defaultObservability string) list.Model {
	items := []list.Item{
		observabilityItem{name: "none", desc: "No observability - minimal overhead"},
		observabilityItem{name: "basic", desc: "Basic - Enhanced /health endpoint only"},
		observabilityItem{name: "advanced", desc: "Advanced - Full Prometheus metrics + middleware"},
	}

	delegate := list.NewDefaultDelegate()
	l := list.New(items, delegate, 0, 0)
	l.Title = "Select Observability Level"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	// Set default selection if provided
	if defaultObservability != "" {
		for i, item := range items {
			if obs, ok := item.(observabilityItem); ok && obs.name == defaultObservability {
				l.Select(i)
				break
			}
		}
	}

	return l
}

// templateItem is a list item for template selection.
type templateItem struct {
	name string
	desc string
}

func (t templateItem) Title() string       { return t.name }
func (t templateItem) Description() string { return t.desc }
func (t templateItem) FilterValue() string { return t.name }

// databaseItem is a list item for database selection.
type databaseItem struct {
	name string
	desc string
}

func (d databaseItem) Title() string       { return d.name }
func (d databaseItem) Description() string { return d.desc }
func (d databaseItem) FilterValue() string { return d.name }

// observabilityItem is a list item for observability selection.
type observabilityItem struct {
	name string
	desc string
}

func (o observabilityItem) Title() string       { return o.name }
func (o observabilityItem) Description() string { return o.desc }
func (o observabilityItem) FilterValue() string { return o.name }

// initializeFrameworkList creates and initializes the framework selection list.
func initializeFrameworkList(defaultFramework string) list.Model {
	items := []list.Item{
		frameworkItem{name: "fiber", desc: "Fiber v2.52.10 - Fast HTTP framework inspired by Express (default)"},
		frameworkItem{name: "gin", desc: "Gin - High-performance HTTP web framework (planned for v2.0.0)"},
		frameworkItem{name: "echo", desc: "Echo - Minimalist high-performance HTTP framework (planned for v2.0.0)"},
	}

	delegate := list.NewDefaultDelegate()
	l := list.New(items, delegate, 0, 0)
	l.Title = "Select a Framework"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	// Set default selection if provided
	if defaultFramework != "" {
		for i, item := range items {
			if fw, ok := item.(frameworkItem); ok && fw.name == defaultFramework {
				l.Select(i)
				break
			}
		}
	}

	return l
}

// frameworkItem is a list item for framework selection.
type frameworkItem struct {
	name string
	desc string
}

func (f frameworkItem) Title() string       { return f.name }
func (f frameworkItem) Description() string { return f.desc }
func (f frameworkItem) FilterValue() string { return f.name }

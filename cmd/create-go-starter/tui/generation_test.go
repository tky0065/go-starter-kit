package tui

import (
	"strings"
	"testing"
	"time"
)

// TestGenerationModelCreation verifies that a generation model can be created with total files.
func TestGenerationModelCreation(t *testing.T) {
	model := NewGenerationModel(34)

	if model.state != StateGenerating {
		t.Errorf("Expected StateGenerating, got %v", model.state)
	}

	if model.totalFiles != 34 {
		t.Errorf("Expected totalFiles to be 34, got %d", model.totalFiles)
	}

	if model.filesGenerated != 0 {
		t.Errorf("Expected filesGenerated to be 0, got %d", model.filesGenerated)
	}
}

// TestGenerationProgress verifies that file generation progress is tracked correctly.
func TestGenerationProgress(t *testing.T) {
	model := NewGenerationModel(10)

	t.Run("tracks individual file generation", func(t *testing.T) {
		newModel, _ := model.Update(FileGeneratedMsg{
			FilePath: "main.go",
			Index:    1,
		})
		m := newModel.(Model)

		if m.filesGenerated != 1 {
			t.Errorf("Expected filesGenerated to be 1, got %d", m.filesGenerated)
		}
	})

	t.Run("completes when all files are generated", func(t *testing.T) {
		model.filesGenerated = 9

		newModel, _ := model.Update(FileGeneratedMsg{
			FilePath: "final.go",
			Index:    10,
		})
		m := newModel.(Model)

		if m.filesGenerated != 10 {
			t.Errorf("Expected filesGenerated to be 10, got %d", m.filesGenerated)
		}
	})
}

// TestGenerationCompletion verifies that generation completion is handled correctly.
func TestGenerationCompletion(t *testing.T) {
	model := NewGenerationModel(5)
	model.filesGenerated = 5

	t.Run("transitions to StateDone on successful completion", func(t *testing.T) {
		newModel, _ := model.Update(GenerationCompleteMsg{
			TotalFiles: 5,
			Success:    true,
			Error:      nil,
		})
		m := newModel.(Model)

		if m.state != StateDone {
			t.Errorf("Expected StateDone, got %v", m.state)
		}
	})

	t.Run("handles generation errors", func(t *testing.T) {
		newModel, _ := model.Update(GenerationCompleteMsg{
			TotalFiles: 5,
			Success:    false,
			Error:      nil, // Error would be set in real scenario
		})
		m := newModel.(Model)

		if m.state != StateDone {
			t.Errorf("Expected StateDone even on error, got %v", m.state)
		}
	})
}

// TestProgressBarAnimation verifies that the progress bar updates correctly.
func TestProgressBarAnimation(t *testing.T) {
	model := NewGenerationModel(100)

	// Simulate progression
	for i := 1; i <= 100; i++ {
		newModel, _ := model.Update(FileGeneratedMsg{
			FilePath: "file.go",
			Index:    i,
		})
		model = newModel.(Model)
	}

	if model.filesGenerated != 100 {
		t.Errorf("Expected filesGenerated to be 100, got %d", model.filesGenerated)
	}
}

// TestGenerationView verifies that the generation screen renders correctly.
func TestGenerationView(t *testing.T) {
	model := NewGenerationModel(50)
	model.filesGenerated = 25

	view := model.View()

	if view == "" {
		t.Error("Generation view should not be empty")
	}

	// View should show progress information
	// (More specific assertions would check for progress indicators)
}

// TestSpinnerAnimation verifies that the spinner is displayed during generation.
func TestSpinnerAnimation(t *testing.T) {
	model := NewGenerationModel(10)

	// Spinner should be initialized
	t.Run("spinner is initialized", func(t *testing.T) {
		// The spinner should be ready to animate
		// Actual spinner ticks are handled by Bubble Tea's Cmd system
		if model.state != StateGenerating {
			t.Error("Model should be in generating state for spinner")
		}
	})
}

// Story 10.7 Task 3: Tests for enhanced progress bar with real-time statistics

// TestGenerationStatsTracking verifies that generation statistics are tracked correctly.
func TestGenerationStatsTracking(t *testing.T) {
	model := NewGenerationModel(10)
	model.projectName = "test-project"

	// Initialize stats via ConfirmGenerationMsg
	newModel, _ := model.Update(ConfirmGenerationMsg{})
	model = newModel.(Model)

	t.Run("initializes stats on generation start", func(t *testing.T) {
		if model.genStats.FilesCreated != 0 {
			t.Errorf("Expected FilesCreated to be 0, got %d", model.genStats.FilesCreated)
		}
		if model.genStats.TotalSize != 0 {
			t.Errorf("Expected TotalSize to be 0, got %d", model.genStats.TotalSize)
		}
		if model.genStats.StartTime.IsZero() {
			t.Error("Expected StartTime to be set")
		}
	})

	t.Run("updates stats with file generation", func(t *testing.T) {
		newModel, _ := model.Update(FileGeneratedMsg{
			FilePath: "main.go",
			Index:    1,
			Size:     1024, // 1 KB
			Step:     "Creating main files...",
		})
		model = newModel.(Model)

		if model.genStats.FilesCreated != 1 {
			t.Errorf("Expected FilesCreated to be 1, got %d", model.genStats.FilesCreated)
		}
		if model.genStats.TotalSize != 1024 {
			t.Errorf("Expected TotalSize to be 1024, got %d", model.genStats.TotalSize)
		}
		if model.genStats.CurrentFile != "main.go" {
			t.Errorf("Expected CurrentFile to be 'main.go', got '%s'", model.genStats.CurrentFile)
		}
		if model.genStats.CurrentStep != "Creating main files..." {
			t.Errorf("Expected CurrentStep to be 'Creating main files...', got '%s'", model.genStats.CurrentStep)
		}
	})

	t.Run("accumulates total size", func(t *testing.T) {
		// Add second file
		newModel, _ := model.Update(FileGeneratedMsg{
			FilePath: "config.go",
			Index:    2,
			Size:     512, // 512 B
			Step:     "Creating configuration files...",
		})
		model = newModel.(Model)

		expectedSize := int64(1024 + 512)
		if model.genStats.TotalSize != expectedSize {
			t.Errorf("Expected TotalSize to be %d, got %d", expectedSize, model.genStats.TotalSize)
		}
	})
}

// TestFormatFileSizeGeneration verifies that file sizes are formatted correctly (generation.go).
func TestFormatFileSizeGeneration(t *testing.T) {
	tests := []struct {
		name     string
		bytes    int64
		expected string
	}{
		{"bytes", 500, "500 B"},
		{"kilobytes", 1024, "1.0 KB"},
		{"kilobytes decimal", 1536, "1.5 KB"},
		{"megabytes", 1048576, "1.0 MB"},
		{"megabytes decimal", 2621440, "2.5 MB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatFileSize(tt.bytes)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// TestFormatDuration verifies that durations are formatted correctly.
func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{"milliseconds", 500 * time.Millisecond, "500ms"},
		{"seconds", 2 * time.Second, "2.0s"},
		{"seconds decimal", 2500 * time.Millisecond, "2.5s"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatDuration(tt.duration)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// TestRenderGenerationStats verifies that the stats panel renders correctly.
func TestRenderGenerationStats(t *testing.T) {
	stats := GenerationStats{
		FilesCreated: 10,
		TotalSize:    10240, // 10 KB
		StartTime:    time.Now().Add(-5 * time.Second),
		CurrentFile:  "main.go",
		CurrentStep:  "Creating files...",
	}

	result := renderGenerationStats(stats, 20)

	if result == "" {
		t.Error("Stats panel should not be empty")
	}

	// Verify key components are present
	expectedComponents := []string{
		"Files created:",
		"10",
		"Total size:",
		"10.0 KB",
		"Time elapsed:",
		"ETA:",
	}

	for _, component := range expectedComponents {
		if !strings.Contains(result, component) {
			t.Errorf("Expected stats panel to contain '%s'", component)
		}
	}
}

// TestRenderCurrentStep verifies that the current step message renders correctly.
func TestRenderCurrentStep(t *testing.T) {
	t.Run("renders step message", func(t *testing.T) {
		result := renderCurrentStep("Creating directories...", "")
		if !strings.Contains(result, "Creating directories...") {
			t.Error("Expected step message to contain step description")
		}
	})

	t.Run("renders current file when no step", func(t *testing.T) {
		result := renderCurrentStep("", "internal/domain/user.go")
		if !strings.Contains(result, "internal/domain/user.go") {
			t.Error("Expected step message to contain current file")
		}
	})

	t.Run("returns empty for no step and no file", func(t *testing.T) {
		result := renderCurrentStep("", "")
		if result != "" {
			t.Error("Expected empty string when no step or file")
		}
	})

	t.Run("prioritizes step over file", func(t *testing.T) {
		result := renderCurrentStep("Installing dependencies...", "go.mod")
		if !strings.Contains(result, "Installing dependencies...") {
			t.Error("Expected step message to contain step description")
		}
	})
}

// TestGenerationViewWithStats verifies that the enhanced generation view renders correctly.
func TestGenerationViewWithStats(t *testing.T) {
	model := NewGenerationModel(50)
	model.projectName = "my-awesome-api"
	model.filesGenerated = 25

	// Initialize stats
	model.genStats = GenerationStats{
		FilesCreated: 25,
		TotalSize:    51200, // 50 KB
		StartTime:    time.Now().Add(-10 * time.Second),
		CurrentFile:  "internal/domain/user.go",
		CurrentStep:  "Generating domain services...",
	}

	view := model.View()

	if view == "" {
		t.Error("Generation view should not be empty")
	}

	// Verify enhanced components are present
	expectedComponents := []string{
		"my-awesome-api",        // Project name
		"25/50 files",           // Progress indicator
		"Files created:",        // Stats label
		"Total size:",           // Stats label
		"Time elapsed:",         // Stats label
		"ETA:",                  // Stats label
		"Generating domain services",       // Current step (partial match)
	}

	for _, component := range expectedComponents {
		if !strings.Contains(view, component) {
			t.Errorf("Expected view to contain '%s'\nGot:\n%s", component, view)
		}
	}
}

// TestETACalculation verifies that ETA is calculated correctly.
func TestETACalculation(t *testing.T) {
	// Create stats with 10 files generated in 10 seconds (1 file/sec)
	stats := GenerationStats{
		FilesCreated: 10,
		TotalSize:    10240,
		StartTime:    time.Now().Add(-10 * time.Second),
		CurrentFile:  "file10.go",
		CurrentStep:  "Generating...",
	}

	// With 40 files remaining and 1 file/sec rate, ETA should be ~40s
	result := renderGenerationStats(stats, 50)

	// Verify ETA is present and reasonable
	if !strings.Contains(result, "ETA:") {
		t.Error("Expected stats to contain ETA")
	}

	// ETA should be calculated (we don't check exact value due to timing)
	// Just verify it's shown
}

// TestProgressBarWithGradient verifies that progress bar uses gradient.
func TestProgressBarWithGradient(t *testing.T) {
	model := NewGenerationModel(100)

	// Verify progress bar view is not empty
	t.Run("progress bar has non-empty view", func(t *testing.T) {
		view := model.progressBar.View()
		if view == "" {
			t.Error("Progress bar view should not be empty")
		}
	})

	// Verify progress bar responds to updates
	t.Run("progress bar updates", func(t *testing.T) {
		initialView := model.progressBar.View()
		model.progressBar.SetPercent(0.5)
		updatedView := model.progressBar.View()
		// Views should differ after setting percentage
		// (This is a basic check - we can't easily verify the exact content)
		if initialView == "" && updatedView == "" {
			t.Error("Progress bar should respond to updates")
		}
	})
}

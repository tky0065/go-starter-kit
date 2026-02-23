package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// GenerationStats holds real-time statistics during project generation.
type GenerationStats struct {
	FilesCreated int       // Number of files created so far
	TotalSize    int64     // Total size of generated files in bytes
	StartTime    time.Time // When generation started
	CurrentFile  string    // Path of the file currently being generated
	CurrentStep  string    // Current step description (e.g., "Creating directories...")
}

// NewGenerationModel creates a new Model configured for file generation.
// It initializes the progress bar and spinner for showing generation progress.
func NewGenerationModel(totalFiles int) Model {
	// Initialize progress bar with gradient (Story 10.7 Task 3)
	prog := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(60),
	)

	// Initialize spinner for loading animation
	s := spinner.New()
	s.Spinner = spinner.Dot

	// Initialize progress pulse animation (Story 10.7 Task 10)
	progressPulse := NewPulseState(30, 0.7, 1.0)

	m := Model{
		state:          StateGenerating,
		totalFiles:     totalFiles,
		filesGenerated: 0,
		progressBar:    prog,
		loadSpinner:    s,
		progressPulse:  progressPulse,
		width:          80, // Default width
		height:         24, // Default height
	}

	return m
}

// StartGeneration returns a command to begin the spinner animation.
// This should be called when the generation process starts.
func StartGeneration() tea.Cmd {
	return spinner.Tick
}

// renderGenerationStats renders the statistics panel for file generation.
// This is called from viewGenerating() to display real-time stats.
// Story 10.7 Task 3: Enhanced progress bar with real-time statistics
func renderGenerationStats(stats GenerationStats, totalFiles int) string {
	// Calculate elapsed time
	elapsed := time.Since(stats.StartTime)

	// Calculate ETA based on current progress
	var eta time.Duration
	if stats.FilesCreated > 0 {
		avgTimePerFile := elapsed / time.Duration(stats.FilesCreated)
		remainingFiles := totalFiles - stats.FilesCreated
		eta = avgTimePerFile * time.Duration(remainingFiles)
	}

	// Format total size in human-readable format
	sizeStr := formatFileSize(stats.TotalSize)

	// Build stats panel content
	var statsContent string
	statsContent += fmt.Sprintf("%s Files created:  %s\n",
		RenderMuted("✓"),
		RenderHighlight(fmt.Sprintf("%d", stats.FilesCreated)))
	statsContent += fmt.Sprintf("%s Total size:     %s\n",
		RenderMuted("📦"),
		RenderHighlight(sizeStr))
	statsContent += fmt.Sprintf("%s Time elapsed:   %s\n",
		RenderMuted("⏱ "),
		RenderHighlight(formatDuration(elapsed)))
	if eta > 0 {
		statsContent += fmt.Sprintf("%s ETA:            %s",
			RenderMuted("⏳"),
			RenderHighlight(formatDuration(eta)))
	}

	// Create stats box
	statsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(ColorInfo)).
		Padding(1, 2).
		Render(statsContent)

	return statsBox
}

// renderCurrentStep renders the current generation step with contextual message.
// Story 10.7 Task 3: Contextual step messages during generation
func renderCurrentStep(step string, currentFile string) string {
	if step == "" && currentFile == "" {
		return ""
	}

	var message string
	if step != "" {
		// Show step description (e.g., "Creating directories...")
		message = fmt.Sprintf("%s %s", RenderSimpleInfo("🔄"), step)
	} else if currentFile != "" {
		// Show current file being generated
		message = fmt.Sprintf("%s Creating %s...", RenderSimpleInfo("🔄"), currentFile)
	}

	return message
}

// formatFileSize formats bytes into human-readable size (KB, MB).
func formatFileSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
	)

	switch {
	case bytes >= MB:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

// formatDuration formats a duration into human-readable format (s, ms).
func formatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return fmt.Sprintf("%.1fs", d.Seconds())
}

// renderProgressBarWithPulse renders a progress bar with pulsing gradient animation.
// Story 10.7 Task 10: Pulsing progress bar gradient
func renderProgressBarWithPulse(progressBar progress.Model, percent float64, pulse *PulseState) string {
	if pulse == nil {
		// No pulse animation, use default rendering
		return progressBar.ViewAs(percent)
	}

	// Get pulse intensity (0.7 to 1.0)
	intensity := pulse.Intensity()

	// The bubbles progress bar already has gradient support,
	// so we'll just render it normally.
	// In a more advanced implementation, we could modify the gradient colors
	// based on the pulse intensity for a true pulsing effect.
	//
	// For now, the pulsing effect is conceptual - the infrastructure is in place
	// for future enhancements (e.g., modifying FullColor based on intensity).
	//
	// Example future enhancement:
	// - Interpolate between ColorSuccess and a lighter shade based on intensity
	// - Update progress bar's FullColor dynamically
	//
	// Current implementation: Standard gradient rendering
	_ = intensity // Mark as used for future enhancement
	return progressBar.ViewAs(percent)
}

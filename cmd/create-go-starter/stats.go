package main

import (
	"fmt"
	"os"
	"time"
)

// GenerationStats tracks timing and file statistics during project generation.
type GenerationStats struct {
	StartTime     time.Time
	stepStarts    map[string]time.Time
	StepDurations map[string]time.Duration
	stepOrder     []string // insertion order for deterministic display
	FilesCount    int
	TotalBytes    int64
}

// NewGenerationStats creates a new stats tracker, starting the global timer now.
func NewGenerationStats() *GenerationStats {
	return &GenerationStats{
		StartTime:     time.Now(),
		stepStarts:    make(map[string]time.Time),
		StepDurations: make(map[string]time.Duration),
	}
}

// StartStep records the start time for a named step.
func (s *GenerationStats) StartStep(name string) {
	s.stepStarts[name] = time.Now()
	// Track insertion order for deterministic display
	for _, existing := range s.stepOrder {
		if existing == name {
			return
		}
	}
	s.stepOrder = append(s.stepOrder, name)
}

// EndStep records the elapsed time for a named step previously started with StartStep.
func (s *GenerationStats) EndStep(name string) {
	if start, ok := s.stepStarts[name]; ok {
		s.StepDurations[name] = time.Since(start)
	}
}

// AddFile increments the file counter and accumulates the file size via os.Stat.
func (s *GenerationStats) AddFile(path string) {
	s.FilesCount++
	if info, err := os.Stat(path); err == nil {
		s.TotalBytes += info.Size()
	}
}

// Display prints the generation statistics summary to stdout.
func (s *GenerationStats) Display() {
	total := time.Since(s.StartTime)
	sizeStr := formatBytes(s.TotalBytes)

	fmt.Printf("\n%s\n", Green("📊 Generation Statistics:"))
	fmt.Printf("   Files generated : %d\n", s.FilesCount)
	fmt.Printf("   Total size      : %s\n", sizeStr)
	fmt.Printf("   Total time      : %s\n", total.Round(time.Millisecond))

	if len(s.StepDurations) > 0 {
		fmt.Println("   Step breakdown  :")
		for _, step := range s.stepOrder {
			if dur, ok := s.StepDurations[step]; ok {
				fmt.Printf("     • %-25s %s\n", step+":", dur.Round(time.Millisecond))
			}
		}
	}
}

// formatBytes converts a byte count to a human-readable string (bytes, KB, or MB).
func formatBytes(bytes int64) string {
	switch {
	case bytes >= 1024*1024:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(1024*1024))
	case bytes >= 1024:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(1024))
	default:
		return fmt.Sprintf("%d bytes", bytes)
	}
}

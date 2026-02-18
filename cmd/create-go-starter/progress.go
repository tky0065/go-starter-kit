package main

import (
	"fmt"
	"os"
	"strings"
)

// ProgressBar is a simple terminal progress bar using only stdlib.
// It renders a visual bar like: [██████░░░░] 18/34 files
// When disabled (non-TTY or NO_COLOR), Update and Complete are no-ops.
type ProgressBar struct {
	total   int
	current int
	width   int
	enabled bool
}

// isTerminalFn is the terminal detection function. It is a package-level variable
// to allow overriding in tests for verifying TTY-dependent behavior.
var isTerminalFn = isTerminal

// NewProgressBar creates a new progress bar.
// The bar is automatically disabled when stdout is not a terminal or NO_COLOR is set.
func NewProgressBar(total, width int) *ProgressBar {
	return &ProgressBar{
		total:   total,
		width:   width,
		enabled: isTerminalFn() && os.Getenv("NO_COLOR") == "",
	}
}

// Update refreshes the progress bar display in-place using a carriage return.
func (pb *ProgressBar) Update(current int) {
	if !pb.enabled || pb.total <= 0 {
		return
	}
	pb.current = current
	percent := float64(current) / float64(pb.total)
	filled := int(percent * float64(pb.width))
	empty := pb.width - filled

	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
	fmt.Printf("\r  [%s] %d/%d files", bar, current, pb.total)
}

// Complete marks the progress bar as done (100%) and moves to the next line.
func (pb *ProgressBar) Complete() {
	if !pb.enabled {
		return
	}
	pb.Update(pb.total)
	fmt.Println()
}

// isTerminal returns true if stdout is connected to a terminal (not piped/redirected).
func isTerminal() bool {
	fileInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

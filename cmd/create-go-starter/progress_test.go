package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestProgressBarUpdate verifies that Update does not panic for various values.
func TestProgressBarUpdate(t *testing.T) {
	pb := &ProgressBar{total: 10, width: 10, enabled: false}
	pb.Update(0)
	pb.Update(5)
	pb.Update(10)
}

// TestProgressBarUpdateEnabled verifies Update works with the bar enabled.
func TestProgressBarUpdateEnabled(t *testing.T) {
	pb := &ProgressBar{total: 10, width: 10, enabled: true}
	// Should not panic; output goes to stdout
	pb.Update(3)
	pb.Update(10)
}

// TestProgressBarComplete verifies Complete does not panic when disabled.
func TestProgressBarComplete(t *testing.T) {
	pb := &ProgressBar{total: 10, width: 10, enabled: false}
	pb.Complete()
}

// TestProgressBarCompleteEnabled verifies Complete works with the bar enabled.
func TestProgressBarCompleteEnabled(t *testing.T) {
	pb := &ProgressBar{total: 5, width: 10, enabled: true}
	pb.Complete()
}

// TestNewProgressBarDisabledWhenNoColor verifies that NO_COLOR disables the bar.
func TestNewProgressBarDisabledWhenNoColor(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	pb := NewProgressBar(10, 30)
	if pb.enabled {
		t.Error("ProgressBar should be disabled when NO_COLOR=1")
	}
}

// TestNewProgressBarEnabledWhenTTYAndNoColorAbsent verifies the bar is enabled when
// stdout is a terminal and NO_COLOR is not set.
func TestNewProgressBarEnabledWhenTTYAndNoColorAbsent(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	// Override isTerminalFn to simulate TTY
	original := isTerminalFn
	isTerminalFn = func() bool { return true }
	defer func() { isTerminalFn = original }()

	pb := NewProgressBar(10, 30)
	if !pb.enabled {
		t.Error("ProgressBar should be enabled when isTerminal()=true and NO_COLOR is empty")
	}
}

// TestNewProgressBarDisabledWhenNotTTY verifies the bar is disabled when not a terminal.
func TestNewProgressBarDisabledWhenNotTTY(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	// Override isTerminalFn to simulate non-TTY
	original := isTerminalFn
	isTerminalFn = func() bool { return false }
	defer func() { isTerminalFn = original }()

	pb := NewProgressBar(10, 30)
	if pb.enabled {
		t.Error("ProgressBar should be disabled when isTerminal()=false")
	}
}

// TestIsTerminal verifies the function returns a bool without panicking.
func TestIsTerminal(t *testing.T) {
	// In test environments (CI, go test), stdout is typically not a TTY.
	// This is a smoke test to ensure no panic.
	result := isTerminal()
	// result can be true or false depending on the environment; just ensure no panic.
	_ = result
}

// TestFormatBytes verifies byte formatting at boundary values.
func TestFormatBytes(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{0, "0 bytes"},
		{500, "500 bytes"},
		{1023, "1023 bytes"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1024 * 1024, "1.0 MB"},
		{1536 * 1024, "1.5 MB"},
	}
	for _, tt := range tests {
		result := formatBytes(tt.bytes)
		if result != tt.expected {
			t.Errorf("formatBytes(%d) = %q, want %q", tt.bytes, result, tt.expected)
		}
	}
}

// TestGenerationStatsStartEndStep verifies that step timing is recorded correctly.
func TestGenerationStatsStartEndStep(t *testing.T) {
	stats := NewGenerationStats()
	stats.StartStep("test-step")
	time.Sleep(2 * time.Millisecond)
	stats.EndStep("test-step")

	dur, ok := stats.StepDurations["test-step"]
	if !ok {
		t.Fatal("expected 'test-step' duration to be recorded")
	}
	if dur < time.Millisecond {
		t.Errorf("expected duration >= 1ms, got %v", dur)
	}
}

// TestGenerationStatsEndStepWithoutStart verifies EndStep is safe without StartStep.
func TestGenerationStatsEndStepWithoutStart(t *testing.T) {
	stats := NewGenerationStats()
	stats.EndStep("nonexistent") // should not panic
	if _, ok := stats.StepDurations["nonexistent"]; ok {
		t.Error("should not record duration for step never started")
	}
}

// TestGenerationStatsAddFile verifies file counting and size accumulation.
func TestGenerationStatsAddFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	content := []byte("hello world")
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatal(err)
	}

	stats := NewGenerationStats()
	stats.AddFile(path)

	if stats.FilesCount != 1 {
		t.Errorf("expected FilesCount=1, got %d", stats.FilesCount)
	}
	if stats.TotalBytes == 0 {
		t.Error("expected TotalBytes > 0")
	}
	if stats.TotalBytes != int64(len(content)) {
		t.Errorf("expected TotalBytes=%d, got %d", len(content), stats.TotalBytes)
	}
}

// TestGenerationStatsAddFileMissing verifies AddFile is safe for non-existent files.
func TestGenerationStatsAddFileMissing(t *testing.T) {
	stats := NewGenerationStats()
	stats.AddFile("/nonexistent/path/file.txt") // should not panic
	if stats.FilesCount != 1 {
		t.Errorf("expected FilesCount=1 (counter incremented even if file missing), got %d", stats.FilesCount)
	}
	// TotalBytes should remain 0 since os.Stat failed
	if stats.TotalBytes != 0 {
		t.Errorf("expected TotalBytes=0 for missing file, got %d", stats.TotalBytes)
	}
}

// TestGenerationStatsDisplay verifies Display does not panic.
func TestGenerationStatsDisplay(t *testing.T) {
	stats := NewGenerationStats()
	stats.StartStep("Creating directories")
	time.Sleep(1 * time.Millisecond)
	stats.EndStep("Creating directories")
	stats.FilesCount = 5
	stats.TotalBytes = 1024
	// Display prints to stdout; just verify no panic
	stats.Display()
}

// TestWriteFilesWithProgress verifies that the onProgress callback is called correctly.
func TestWriteFilesWithProgress(t *testing.T) {
	dir := t.TempDir()
	files := []FileGenerator{
		{Path: filepath.Join(dir, "a.txt"), Content: "content a"},
		{Path: filepath.Join(dir, "b.txt"), Content: "content b"},
		{Path: filepath.Join(dir, "c.txt"), Content: "content c"},
	}

	var calls []int
	err := writeFiles(files, func(current, total int) {
		if total != len(files) {
			t.Errorf("expected total=%d, got %d", len(files), total)
		}
		calls = append(calls, current)
	})
	if err != nil {
		t.Fatalf("writeFiles failed: %v", err)
	}

	if len(calls) != len(files) {
		t.Errorf("expected %d progress calls, got %d", len(files), len(calls))
	}
	for i, v := range calls {
		if v != i+1 {
			t.Errorf("calls[%d] = %d, want %d", i, v, i+1)
		}
	}
}

// TestWriteFilesWithNilProgress verifies that nil callback does not panic.
func TestWriteFilesWithNilProgress(t *testing.T) {
	dir := t.TempDir()
	files := []FileGenerator{
		{Path: filepath.Join(dir, "file.txt"), Content: "content"},
	}
	if err := writeFiles(files, nil); err != nil {
		t.Fatalf("writeFiles with nil callback failed: %v", err)
	}
}

// TestProgressBarTotalZeroNoPanic verifies ProgressBar handles total=0 without dividing by zero.
func TestProgressBarTotalZeroNoPanic(t *testing.T) {
	// Test with disabled (no output)
	pb := &ProgressBar{total: 0, width: 10, enabled: false}
	pb.Update(0)
	pb.Complete()

	// Test with enabled — must not divide by zero or produce NaN
	pbEnabled := &ProgressBar{total: 0, width: 10, enabled: true}
	pbEnabled.Update(0)
	pbEnabled.Complete()
}

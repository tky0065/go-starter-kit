package main

import (
	"bytes"
	"strings"
	"testing"
)

// TestCheckGoVersionIntegration tests that Go is available in the test environment
func TestCheckGoVersionIntegration(t *testing.T) {
	result := checkGoVersion()
	if !result.Status {
		t.Errorf("Go should be available in test environment: %s", result.Message)
	}
	if result.Name != "Go" {
		t.Errorf("expected Name = 'Go', got %q", result.Name)
	}
	if result.Version == "" {
		t.Error("expected a non-empty Version string")
	}
}

// TestCheckGoVersionFields tests the CheckResult struct for Go
func TestCheckGoVersionFields(t *testing.T) {
	result := checkGoVersion()
	if result.Name != "Go" {
		t.Errorf("expected Name = 'Go', got %q", result.Name)
	}
	// When status is true, Fix should be empty
	if result.Status && result.Fix != "" {
		t.Errorf("expected Fix to be empty when Go check passes, got %q", result.Fix)
	}
}

// TestCheckGitIntegration tests that Git is available in the test environment
func TestCheckGitIntegration(t *testing.T) {
	result := checkGit()
	if result.Name != "Git" {
		t.Errorf("expected Name = 'Git', got %q", result.Name)
	}
	// Git should be available in CI/dev environment; if not, just verify no panic
	if !result.Status {
		// Git not available — verify Fix message is provided
		if result.Fix == "" {
			t.Error("expected Fix message when Git is not found")
		}
	} else {
		if result.Version == "" {
			t.Error("expected a non-empty Version when Git is found")
		}
	}
}

// TestCheckDockerIntegration verifies Docker check doesn't panic and returns valid result
func TestCheckDockerIntegration(t *testing.T) {
	result := checkDocker()
	if result.Name != "Docker" {
		t.Errorf("expected Name = 'Docker', got %q", result.Name)
	}
	// Docker may or may not be available — verify Fix is set when failed
	if !result.Status && result.Fix == "" {
		t.Error("expected Fix message when Docker check fails")
	}
}

// TestIsGoVersionSufficient tests version comparison logic with table-driven tests
func TestIsGoVersionSufficient(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		minMajor int
		minMinor int
		want     bool
	}{
		{"exact minimum", "1.21", 1, 21, true},
		{"above minimum", "1.25.5", 1, 21, true},
		{"below minimum minor", "1.20", 1, 21, false},
		{"below minimum major", "0.99", 1, 21, false},
		{"above minimum major", "2.0", 1, 21, true},
		{"pre-release rc", "1.22rc1", 1, 21, true},
		{"pre-release beta", "1.23beta2", 1, 21, true},
		{"pre-release at minimum", "1.21rc1", 1, 21, true},
		{"pre-release below minimum", "1.20rc1", 1, 21, false},
		{"single segment", "1", 1, 21, false},
		{"empty string", "", 1, 21, false},
		{"non-numeric", "abc.def", 1, 21, false},
		{"major only numeric", "1.abc", 1, 21, false},
		{"three segments", "1.25.5", 1, 21, true},
		{"minimum with patch", "1.21.0", 1, 21, true},
		{"just below with patch", "1.20.99", 1, 21, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isGoVersionSufficient(tt.version, tt.minMajor, tt.minMinor)
			if got != tt.want {
				t.Errorf("isGoVersionSufficient(%q, %d, %d) = %v, want %v",
					tt.version, tt.minMajor, tt.minMinor, got, tt.want)
			}
		})
	}
}

// TestStripNonNumericSuffix tests the helper function for cleaning version segments
func TestStripNonNumericSuffix(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"22rc1", "22"},
		{"5", "5"},
		{"23beta2", "23"},
		{"", ""},
		{"abc", ""},
		{"1alpha", "1"},
		{"100", "100"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := stripNonNumericSuffix(tt.input)
			if got != tt.want {
				t.Errorf("stripNonNumericSuffix(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// TestWriteDoctorReport tests the report display with known CheckResults using io.Writer
func TestWriteDoctorReport(t *testing.T) {
	results := []CheckResult{
		{Name: "Go", Status: true, Version: "go1.25.5 darwin/arm64"},
		{Name: "Git", Status: false, Message: "not found", Fix: "install git from https://git-scm.com"},
	}

	var buf bytes.Buffer
	writeDoctorReport(&buf, results)
	output := buf.String()

	// Verify Go check appears in output
	if !strings.Contains(output, "Go") {
		t.Error("expected 'Go' in doctor report output")
	}
	// Verify Git check appears in output
	if !strings.Contains(output, "Git") {
		t.Error("expected 'Git' in doctor report output")
	}
	// Verify Fix message appears
	if !strings.Contains(output, "install git") {
		t.Error("expected fix message in doctor report output")
	}
	// Verify version info for failed check with version
	if !strings.Contains(output, "not found") {
		t.Error("expected error message in doctor report output")
	}
}

// TestComputeExitCode tests that exit code is 0 when all pass, 1 when issues found
func TestComputeExitCode(t *testing.T) {
	// All passing
	allOK := []CheckResult{
		{Name: "Go", Status: true, Version: "go1.25.5"},
		{Name: "Git", Status: true, Version: "git version 2.43.0"},
	}
	code := computeExitCode(allOK)
	if code != 0 {
		t.Errorf("expected exit code 0 when all checks pass, got %d", code)
	}

	// One failing
	oneFail := []CheckResult{
		{Name: "Go", Status: true, Version: "go1.25.5"},
		{Name: "Git", Status: false, Message: "not found", Fix: "install git"},
	}
	code = computeExitCode(oneFail)
	if code != 1 {
		t.Errorf("expected exit code 1 when a check fails, got %d", code)
	}

	// All failing
	allFail := []CheckResult{
		{Name: "Go", Status: false, Message: "not found"},
		{Name: "Git", Status: false, Message: "not found"},
	}
	code = computeExitCode(allFail)
	if code != 1 {
		t.Errorf("expected exit code 1 when all checks fail, got %d", code)
	}

	// Empty results
	code = computeExitCode([]CheckResult{})
	if code != 0 {
		t.Errorf("expected exit code 0 when no checks, got %d", code)
	}
}

// TestVersionConstant tests that Version constant is defined and non-empty
func TestVersionConstant(t *testing.T) {
	if Version == "" {
		t.Error("Version constant should not be empty")
	}
	if !strings.HasPrefix(Version, "1.") {
		t.Errorf("Version should start with '1.', got %q", Version)
	}
}

// TestDoctorReportContainsVersion tests that version appears in doctor report
func TestDoctorReportContainsVersion(t *testing.T) {
	results := []CheckResult{
		{Name: "Go", Status: true, Version: "go1.25.5"},
	}

	var buf bytes.Buffer
	writeDoctorReport(&buf, results)
	output := buf.String()

	if !strings.Contains(output, Version) {
		t.Errorf("expected Version %q in doctor report output", Version)
	}
}

// TestDoctorAllPassedSummary tests the summary message when all checks pass
func TestDoctorAllPassedSummary(t *testing.T) {
	results := []CheckResult{
		{Name: "Go", Status: true, Version: "go1.25.5"},
		{Name: "Git", Status: true, Version: "git version 2.43.0"},
	}

	var buf bytes.Buffer
	writeDoctorReport(&buf, results)
	output := buf.String()

	if !strings.Contains(output, "All checks passed") {
		t.Errorf("expected 'All checks passed' in output when all checks pass, got:\n%s", output)
	}
}

// TestDoctorIssueFoundSummaryPlural tests the summary when multiple issues are found
func TestDoctorIssueFoundSummaryPlural(t *testing.T) {
	results := []CheckResult{
		{Name: "Go", Status: true, Version: "go1.25.5"},
		{Name: "Git", Status: false, Message: "not found", Fix: "install git"},
		{Name: "Docker", Status: false, Message: "not found", Fix: "install docker"},
	}

	var buf bytes.Buffer
	writeDoctorReport(&buf, results)
	output := buf.String()

	if !strings.Contains(output, "2 issues found") {
		t.Errorf("expected '2 issues found' in output when 2 checks fail, got:\n%s", output)
	}
}

// TestDoctorIssueFoundSummarySingular tests the summary when exactly 1 issue is found
func TestDoctorIssueFoundSummarySingular(t *testing.T) {
	results := []CheckResult{
		{Name: "Go", Status: true, Version: "go1.25.5"},
		{Name: "Git", Status: false, Message: "not found", Fix: "install git"},
	}

	var buf bytes.Buffer
	writeDoctorReport(&buf, results)
	output := buf.String()

	if !strings.Contains(output, "1 issue found") {
		t.Errorf("expected '1 issue found' (singular) in output when 1 check fails, got:\n%s", output)
	}
	// Must NOT contain "issues" (plural)
	if strings.Contains(output, "1 issues") {
		t.Errorf("expected singular 'issue' not 'issues' for 1 failure, got:\n%s", output)
	}
}

// TestDoctorReportFailedWithVersion tests that version + error message display correctly
func TestDoctorReportFailedWithVersion(t *testing.T) {
	results := []CheckResult{
		{Name: "Docker", Status: false, Version: "Docker version 24.0.0", Message: "Docker daemon is not running", Fix: "Start Docker Desktop"},
	}

	var buf bytes.Buffer
	writeDoctorReport(&buf, results)
	output := buf.String()

	// When Version is set on a failed check, both version and message should appear
	if !strings.Contains(output, "Docker version 24.0.0") {
		t.Error("expected version string in output for failed check with version")
	}
	if !strings.Contains(output, "Docker daemon is not running") {
		t.Error("expected error message in output for failed check")
	}
	if !strings.Contains(output, "Start Docker Desktop") {
		t.Error("expected fix message in output for failed check")
	}
}

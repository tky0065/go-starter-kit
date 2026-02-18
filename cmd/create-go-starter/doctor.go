package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// CheckResult holds the result of a single environment check.
type CheckResult struct {
	Name    string // "Go", "Git", "Docker"
	Status  bool   // true = OK, false = FAIL
	Version string // e.g., "go1.25.5 darwin/arm64"
	Message string // Error message when Status is false
	Fix     string // Suggested fix when Status is false
}

// runDoctor runs all environment checks and prints the report.
// It returns the exit code: 0 if all checks pass, 1 if any check fails.
func runDoctor() int {
	results := []CheckResult{
		checkGoVersion(),
		checkGit(),
		checkDocker(),
	}
	displayDoctorReport(results)
	return computeExitCode(results)
}

// checkGoVersion checks whether Go is installed and meets the minimum version (1.21).
func checkGoVersion() CheckResult {
	cmd := exec.Command("go", "version")
	output, err := cmd.Output()
	if err != nil {
		return CheckResult{
			Name:    "Go",
			Status:  false,
			Message: "Go not found in PATH",
			Fix:     "Download Go from https://go.dev/dl/ and follow installation instructions",
		}
	}

	versionStr := strings.TrimSpace(string(output))
	// Format: "go version go1.25.5 darwin/arm64"
	parts := strings.Fields(versionStr)
	if len(parts) < 3 {
		return CheckResult{
			Name:    "Go",
			Status:  false,
			Message: "Could not parse Go version output",
			Fix:     "Download Go from https://go.dev/dl/ and follow installation instructions",
		}
	}

	goVersion := strings.TrimPrefix(parts[2], "go") // "1.25.5"
	if !isGoVersionSufficient(goVersion, 1, 21) {
		return CheckResult{
			Name:    "Go",
			Status:  false,
			Version: versionStr,
			Message: fmt.Sprintf("Go version %s is below minimum required 1.21", goVersion),
			Fix:     "Update Go to 1.21+ from https://go.dev/dl/",
		}
	}

	return CheckResult{
		Name:    "Go",
		Status:  true,
		Version: versionStr,
	}
}

// isGoVersionSufficient returns true if goVersion (e.g. "1.25.5") is >= major.minor.
// It handles pre-release suffixes like "1.22rc1", "1.23beta2" by stripping non-numeric
// trailing characters from each segment before parsing.
func isGoVersionSufficient(goVersion string, minMajor, minMinor int) bool {
	segments := strings.Split(goVersion, ".")
	if len(segments) < 2 {
		return false
	}
	major, err := strconv.Atoi(stripNonNumericSuffix(segments[0]))
	if err != nil {
		return false
	}
	minor, err := strconv.Atoi(stripNonNumericSuffix(segments[1]))
	if err != nil {
		return false
	}
	if major > minMajor {
		return true
	}
	if major == minMajor && minor >= minMinor {
		return true
	}
	return false
}

// stripNonNumericSuffix returns the leading numeric portion of a string.
// For example: "22rc1" -> "22", "5" -> "5", "23beta2" -> "23", "" -> "".
func stripNonNumericSuffix(s string) string {
	for i, c := range s {
		if c < '0' || c > '9' {
			return s[:i]
		}
	}
	return s
}

// checkGit checks whether Git is installed and available in the PATH.
func checkGit() CheckResult {
	cmd := exec.Command("git", "--version")
	output, err := cmd.Output()
	if err != nil {
		return CheckResult{
			Name:    "Git",
			Status:  false,
			Message: "Git not found in PATH",
			Fix:     "Install Git: macOS: brew install git | Ubuntu: sudo apt install git | Windows: https://git-scm.com/download/win",
		}
	}
	return CheckResult{
		Name:    "Git",
		Status:  true,
		Version: strings.TrimSpace(string(output)),
	}
}

// checkDocker checks whether Docker binary is installed and the daemon is running.
func checkDocker() CheckResult {
	// Check if Docker binary exists
	binaryCmd := exec.Command("docker", "--version")
	binaryOutput, err := binaryCmd.Output()
	if err != nil {
		return CheckResult{
			Name:    "Docker",
			Status:  false,
			Message: "Docker binary not found in PATH",
			Fix:     "Install Docker Desktop: https://www.docker.com/products/docker-desktop/",
		}
	}

	// Check if Docker daemon is running
	daemonCmd := exec.Command("docker", "info")
	if err := daemonCmd.Run(); err != nil {
		return CheckResult{
			Name:    "Docker",
			Status:  false,
			Version: strings.TrimSpace(string(binaryOutput)),
			Message: "Docker daemon is not running",
			Fix:     "Start Docker: macOS/Windows: open Docker Desktop app | Linux: sudo systemctl start docker",
		}
	}

	return CheckResult{
		Name:    "Docker",
		Status:  true,
		Version: strings.TrimSpace(string(binaryOutput)),
	}
}

// computeExitCode returns 0 if all checks pass, 1 if any check fails.
func computeExitCode(results []CheckResult) int {
	for _, r := range results {
		if !r.Status {
			return 1
		}
	}
	return 0
}

// displayDoctorReport prints the formatted doctor report to stdout.
func displayDoctorReport(results []CheckResult) {
	writeDoctorReport(os.Stdout, results)
}

// writeDoctorReport writes the formatted doctor report to the given writer.
// This is separated from displayDoctorReport for testability.
func writeDoctorReport(w io.Writer, results []CheckResult) {
	border := "══════════════════════════════════════════════════════"

	fmt.Fprintln(w, border)
	fmt.Fprintf(w, "  create-go-starter doctor — Environment Check\n")
	fmt.Fprintln(w, border)
	fmt.Fprintf(w, "  Tool version: v%s\n", Version)
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Checking your development environment...")
	fmt.Fprintln(w)

	for _, r := range results {
		if r.Status {
			fmt.Fprintf(w, "  %s %-10s %s\n", Green("✅"), r.Name, r.Version)
		} else {
			versionInfo := r.Message
			if r.Version != "" {
				versionInfo = r.Version + " — " + r.Message
			}
			fmt.Fprintf(w, "  %s %-10s %s\n", Red("❌"), r.Name, versionInfo)
			if r.Fix != "" {
				fmt.Fprintf(w, "     %s\n", Yellow("Fix: "+r.Fix))
			}
		}
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, border)

	issueCount := 0
	for _, r := range results {
		if !r.Status {
			issueCount++
		}
	}

	if issueCount == 0 {
		fmt.Fprintf(w, "  %s\n", Green("✅ All checks passed! Your environment is ready."))
	} else {
		summary := fmt.Sprintf("⚠️  %d issue", issueCount)
		if issueCount > 1 {
			summary += "s"
		}
		summary += " found. Please resolve before using create-go-starter."
		fmt.Fprintf(w, "  %s\n", Yellow(summary))
	}

	fmt.Fprintln(w, border)
}

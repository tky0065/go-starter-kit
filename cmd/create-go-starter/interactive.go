package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/tky0065/go-starter-kit/pkg/utils"
)

// InteractiveDefaults holds default values for interactive mode prompts,
// allowing CLI flags to be pre-populated as defaults.
type InteractiveDefaults struct {
	ProjectName   string
	Template      string
	Database      string
	Observability string
}

// runInteractiveMode runs the interactive mode using os.Stdin with the given defaults.
func runInteractiveMode(defaults InteractiveDefaults) error {
	return runInteractiveModeWithReader(os.Stdin, defaults)
}

// runInteractiveModeWithReader runs the interactive mode reading from the given reader.
// This signature allows easy unit testing via strings.NewReader.
// The defaults parameter allows CLI flags to pre-populate prompt defaults.
func runInteractiveModeWithReader(r io.Reader, defaults InteractiveDefaults) error {
	reader := bufio.NewReader(r)

	// Print header
	fmt.Println()
	fmt.Println("══════════════════════════════════════════")
	fmt.Println("  create-go-starter — Interactive Mode")
	fmt.Println("══════════════════════════════════════════")
	fmt.Println()

	// 1. Collect project name
	defaultProjectName := defaults.ProjectName
	if defaultProjectName == "" {
		defaultProjectName = "my-project"
	}
	projectName, err := collectProjectName(reader, defaultProjectName)
	if err != nil {
		fmt.Println("Cancelled.")
		return err
	}

	// 2. Collect template
	templateOptionsOrdered := []string{TemplateFull, TemplateMinimal, TemplateGraphQL}
	templateDescs := []string{TemplateFullDesc, TemplateMinimalDesc, TemplateGraphQLDesc}
	// Strip "(default)" from descriptions since the menu adds its own marker
	for i, desc := range templateDescs {
		templateDescs[i] = strings.TrimSuffix(strings.TrimSuffix(desc, " (default)"), "(default)")
	}

	templateDefaultIdx := indexOfOption(templateOptionsOrdered, defaults.Template, 0)
	template, err := promptSelectFromReader(reader, "Template", templateOptionsOrdered, templateDescs, templateDefaultIdx)
	if err != nil {
		fmt.Println("Cancelled.")
		return err
	}

	// 3. Collect database
	databaseOptions := []string{DatabasePostgres, DatabaseMySQL, DatabaseSQLite}
	databaseDescs := []string{
		"PostgreSQL - Production-ready, advanced features",
		"MySQL/MariaDB - Wide compatibility, shared hosting",
		"SQLite - Quick prototyping, embedded apps",
	}

	databaseDefaultIdx := indexOfOption(databaseOptions, defaults.Database, 0)
	database, err := promptSelectFromReader(reader, "Database", databaseOptions, databaseDescs, databaseDefaultIdx)
	if err != nil {
		fmt.Println("Cancelled.")
		return err
	}

	// 4. Collect observability
	obsOptions := []string{ObservabilityNone, ObservabilityBasic, ObservabilityAdvanced}
	obsDescs := []string{
		"None - No observability",
		"Basic - Enhanced /health endpoint only",
		"Advanced - Full Prometheus metrics + middleware",
	}

	obsDefaultIdx := indexOfOption(obsOptions, defaults.Observability, 0)
	observability, err := promptSelectFromReader(reader, "Observability", obsOptions, obsDescs, obsDefaultIdx)
	if err != nil {
		fmt.Println("Cancelled.")
		return err
	}

	// 5. Validate observability + template combination
	if observability == ObservabilityAdvanced && template != TemplateFull {
		fmt.Println(Red(fmt.Sprintf("  --observability=advanced is only supported with template=full (got template=%s)", template)))
		fmt.Println("  Please re-run interactive mode with a compatible combination.")
		return fmt.Errorf("--observability=advanced is only supported with --template=full (got --template=%s)", template)
	}

	// 6. Display summary
	fmt.Println()
	fmt.Println("══════════════════════════════════════════")
	fmt.Println("Summary:")
	fmt.Printf("  Project:       %s\n", projectName)
	fmt.Printf("  Template:      %s\n", template)
	fmt.Printf("  Database:      %s\n", database)
	fmt.Printf("  Observability: %s\n", observability)
	fmt.Println("══════════════════════════════════════════")
	fmt.Println()

	// 7. Confirm
	confirmed, err := promptConfirmFromReader(reader, "Confirm?", true)
	if err != nil {
		fmt.Println("Cancelled.")
		return err
	}
	if !confirmed {
		fmt.Println("Cancelled.")
		return fmt.Errorf("cancelled by user")
	}

	// 8. Run project creation
	return run(projectName, template, database, observability)
}

// indexOfOption returns the index of value in the options slice, or defaultIdx if not found.
func indexOfOption(options []string, value string, defaultIdx int) int {
	for i, opt := range options {
		if opt == value {
			return i
		}
	}
	return defaultIdx
}

// collectProjectName prompts for a valid project name, re-prompting on invalid input.
func collectProjectName(reader *bufio.Reader, defaultName string) (string, error) {
	for {
		name, err := promptStringFromReader(reader, "Project name", defaultName)
		if err != nil {
			return "", err
		}
		if err := utils.ValidateGoModuleName(name); err != nil {
			fmt.Println(Red(fmt.Sprintf("  Invalid project name: %v", err)))
			fmt.Println("  Please try again.")
			continue
		}
		return name, nil
	}
}

// promptStringFromReader reads a string from the reader with a default value.
// Returns the default if the user just presses Enter. Returns an error on EOF.
func promptStringFromReader(r *bufio.Reader, prompt, defaultVal string) (string, error) {
	if defaultVal != "" {
		fmt.Printf("? %s [%s]: ", prompt, defaultVal)
	} else {
		fmt.Printf("? %s: ", prompt)
	}

	line, err := r.ReadString('\n')
	if err != nil {
		if err == io.EOF && line == "" {
			return "", fmt.Errorf("interrupted")
		}
		// If we got some input before EOF, use it
		if line != "" {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" {
				return defaultVal, nil
			}
			return trimmed, nil
		}
		return "", fmt.Errorf("interrupted")
	}

	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return defaultVal, nil
	}
	return trimmed, nil
}

// promptSelectFromReader displays a numbered selection menu and reads the user's choice.
// Returns the selected option string. Re-prompts on invalid input. Returns error on EOF.
func promptSelectFromReader(r *bufio.Reader, prompt string, options []string, descriptions []string, defaultIdx int) (string, error) {
	fmt.Printf("? %s: (use number to select)\n", prompt)
	for i, opt := range options {
		marker := ""
		if i == defaultIdx {
			marker = " (default)"
		}
		if i < len(descriptions) {
			fmt.Printf("  %d) %-9s - %s%s\n", i+1, opt, descriptions[i], marker)
		} else {
			fmt.Printf("  %d) %s%s\n", i+1, opt, marker)
		}
	}

	for {
		fmt.Printf("  Select [%d]: ", defaultIdx+1)

		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF && line == "" {
				return "", fmt.Errorf("interrupted")
			}
			if line != "" {
				// Process partial input before EOF
				trimmed := strings.TrimSpace(line)
				if trimmed == "" {
					return options[defaultIdx], nil
				}
				idx, parseErr := strconv.Atoi(trimmed)
				if parseErr == nil && idx >= 1 && idx <= len(options) {
					return options[idx-1], nil
				}
			}
			return "", fmt.Errorf("interrupted")
		}

		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			return options[defaultIdx], nil
		}

		idx, parseErr := strconv.Atoi(trimmed)
		if parseErr != nil || idx < 1 || idx > len(options) {
			fmt.Printf("  Invalid selection. Please enter a number between 1 and %d.\n", len(options))
			continue
		}

		return options[idx-1], nil
	}
}

// promptConfirmFromReader reads a yes/no confirmation from the reader.
// If defaultYes is true, empty input returns true (displayed as [Y/n]).
// If defaultYes is false, empty input returns false (displayed as [y/N]).
// Returns true for "y", "Y", "yes", "YES"; false for "n", "N", "no", "NO".
func promptConfirmFromReader(r *bufio.Reader, prompt string, defaultYes bool) (bool, error) {
	hint := "[y/N]"
	if defaultYes {
		hint = "[Y/n]"
	}
	fmt.Printf("? %s %s: ", prompt, hint)

	line, err := r.ReadString('\n')
	if err != nil {
		if err == io.EOF && line == "" {
			return false, fmt.Errorf("interrupted")
		}
		if line != "" {
			trimmed := strings.TrimSpace(strings.ToLower(line))
			if trimmed == "" {
				return defaultYes, nil
			}
			return trimmed == "y" || trimmed == "yes", nil
		}
		return false, fmt.Errorf("interrupted")
	}

	trimmed := strings.TrimSpace(strings.ToLower(line))
	if trimmed == "" {
		return defaultYes, nil
	}
	return trimmed == "y" || trimmed == "yes", nil
}

package utils

import (
	"fmt"
	"regexp"
)

// Valid Go module name pattern
var ValidGoModuleNamePattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]*$`)

// ValidateGoModuleName validates that a module name is valid for Go modules.
// Valid names must:
// - Start with a letter or number
// - Contain only letters, numbers, hyphens, or underscores
// - Not be empty
func ValidateGoModuleName(name string) error {
	if name == "" {
		return fmt.Errorf("module name cannot be empty")
	}

	if !ValidGoModuleNamePattern.MatchString(name) {
		return fmt.Errorf("invalid module name '%s': must start with a letter or number and contain only letters, numbers, hyphens, or underscores", name)
	}

	return nil
}

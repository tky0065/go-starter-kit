package main

import (
	"fmt"
	"strings"
	"unicode"
)

// FieldDefinition represents a parsed model field with its type and GORM modifiers.
type FieldDefinition struct {
	Name     string   // PascalCase (e.g., "Title")
	Type     string   // Go type (e.g., "string")
	GORMTags []string // Modifiers (e.g., ["unique", "not_null"])
	JSONName string   // snake_case (e.g., "title")
}

// supportedTypes maps user-friendly type names to Go types.
var supportedTypes = map[string]string{
	"string":  "string",
	"int":     "int",
	"uint":    "uint",
	"float64": "float64",
	"bool":    "bool",
	"time":    "time.Time",
}

// validGORMModifiers contains the allowed GORM tag modifiers.
var validGORMModifiers = map[string]bool{
	"unique":   true,
	"not_null": true,
	"index":    true,
}

// toPascalCase converts a string to PascalCase.
func toPascalCase(s string) string {
	if s == "" {
		return ""
	}
	// If already PascalCase, return as-is
	if unicode.IsUpper(rune(s[0])) {
		return s
	}
	// Capitalize first letter
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// toSnakeCase converts a PascalCase or camelCase string to snake_case.
// Handles acronyms properly: "ID" → "id", "UserID" → "user_id", "APIKey" → "api_key".
func toSnakeCase(s string) string {
	if s == "" {
		return ""
	}

	var result []rune
	runes := []rune(s)

	for i := 0; i < len(runes); i++ {
		r := runes[i]

		if unicode.IsUpper(r) {
			// Add underscore if:
			// 1. Not the first character, AND
			// 2. (Previous char is lowercase OR previous char is digit) OR next char is lowercase (end of acronym)
			if i > 0 {
				prevIsLower := unicode.IsLower(runes[i-1])
				prevIsDigit := unicode.IsDigit(runes[i-1])
				nextIsLower := i+1 < len(runes) && unicode.IsLower(runes[i+1])

				if prevIsLower || prevIsDigit || nextIsLower {
					result = append(result, '_')
				}
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

// parseFields parses a fields string like "title:string,completed:bool:not_null"
// into a slice of FieldDefinition.
func parseFields(fieldsStr string) ([]FieldDefinition, error) {
	if fieldsStr == "" {
		return nil, fmt.Errorf("fields cannot be empty")
	}

	parts := strings.Split(fieldsStr, ",")
	var fields []FieldDefinition

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		segments := strings.Split(part, ":")
		if len(segments) < 2 {
			return nil, fmt.Errorf("invalid field format %q: expected name:type[:modifier...]", part)
		}

		name := strings.TrimSpace(segments[0])
		typeName := strings.TrimSpace(segments[1])

		if name == "" {
			return nil, fmt.Errorf("field name cannot be empty in %q", part)
		}

		// Validate type
		goType, ok := supportedTypes[typeName]
		if !ok {
			validTypes := make([]string, 0, len(supportedTypes))
			for k := range supportedTypes {
				validTypes = append(validTypes, k)
			}
			return nil, fmt.Errorf("unsupported type %q for field %q: valid types are: %s",
				typeName, name, strings.Join(validTypes, ", "))
		}

		// Parse GORM modifiers
		var gormTags []string
		for _, seg := range segments[2:] {
			mod := strings.TrimSpace(seg)
			if mod == "" {
				continue
			}
			if !validGORMModifiers[mod] {
				validMods := make([]string, 0, len(validGORMModifiers))
				for k := range validGORMModifiers {
					validMods = append(validMods, k)
				}
				return nil, fmt.Errorf("unsupported GORM modifier %q for field %q: valid modifiers are: %s",
					mod, name, strings.Join(validMods, ", "))
			}
			gormTags = append(gormTags, mod)
		}

		pascalName := toPascalCase(name)
		fields = append(fields, FieldDefinition{
			Name:     pascalName,
			Type:     goType,
			GORMTags: gormTags,
			JSONName: toSnakeCase(pascalName),
		})
	}

	if len(fields) == 0 {
		return nil, fmt.Errorf("no valid fields found")
	}

	return fields, nil
}

// validateModelName checks that the model name is a valid Go identifier in PascalCase.
func validateModelName(name string) error {
	if name == "" {
		return fmt.Errorf("model name is required")
	}
	if !unicode.IsUpper(rune(name[0])) {
		return fmt.Errorf("model name %q must start with an uppercase letter (PascalCase)", name)
	}
	for _, r := range name {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return fmt.Errorf("model name %q contains invalid character %q: only letters and digits are allowed", name, string(r))
		}
	}
	return nil
}

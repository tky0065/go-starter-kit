package main

import (
	"testing"
)

// TestParseFieldsValid tests parsing valid field definitions (AC: 1)
func TestParseFieldsValid(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []FieldDefinition
	}{
		{
			name:  "single string field",
			input: "title:string",
			expected: []FieldDefinition{
				{Name: "Title", Type: "string", GORMTags: nil, JSONName: "title"},
			},
		},
		{
			name:  "multiple fields",
			input: "title:string,completed:bool",
			expected: []FieldDefinition{
				{Name: "Title", Type: "string", GORMTags: nil, JSONName: "title"},
				{Name: "Completed", Type: "bool", GORMTags: nil, JSONName: "completed"},
			},
		},
		{
			name:  "field with modifiers",
			input: "email:string:unique:not_null",
			expected: []FieldDefinition{
				{Name: "Email", Type: "string", GORMTags: []string{"unique", "not_null"}, JSONName: "email"},
			},
		},
		{
			name:  "time type",
			input: "createdAt:time",
			expected: []FieldDefinition{
				{Name: "CreatedAt", Type: "time.Time", GORMTags: nil, JSONName: "created_at"},
			},
		},
		{
			name:  "all supported types",
			input: "name:string,count:int,id:uint,price:float64,active:bool,date:time",
			expected: []FieldDefinition{
				{Name: "Name", Type: "string", GORMTags: nil, JSONName: "name"},
				{Name: "Count", Type: "int", GORMTags: nil, JSONName: "count"},
				{Name: "Id", Type: "uint", GORMTags: nil, JSONName: "id"},
				{Name: "Price", Type: "float64", GORMTags: nil, JSONName: "price"},
				{Name: "Active", Type: "bool", GORMTags: nil, JSONName: "active"},
				{Name: "Date", Type: "time.Time", GORMTags: nil, JSONName: "date"},
			},
		},
		{
			name:  "field with index modifier",
			input: "email:string:index",
			expected: []FieldDefinition{
				{Name: "Email", Type: "string", GORMTags: []string{"index"}, JSONName: "email"},
			},
		},
		{
			name:  "PascalCase field name preserved",
			input: "BlogPost:string",
			expected: []FieldDefinition{
				{Name: "BlogPost", Type: "string", GORMTags: nil, JSONName: "blog_post"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields, err := parseFields(tt.input)
			if err != nil {
				t.Fatalf("parseFields(%q) returned error: %v", tt.input, err)
			}

			if len(fields) != len(tt.expected) {
				t.Fatalf("parseFields(%q) returned %d fields, want %d", tt.input, len(fields), len(tt.expected))
			}

			for i, f := range fields {
				exp := tt.expected[i]
				if f.Name != exp.Name {
					t.Errorf("field[%d].Name = %q, want %q", i, f.Name, exp.Name)
				}
				if f.Type != exp.Type {
					t.Errorf("field[%d].Type = %q, want %q", i, f.Type, exp.Type)
				}
				if f.JSONName != exp.JSONName {
					t.Errorf("field[%d].JSONName = %q, want %q", i, f.JSONName, exp.JSONName)
				}
				if len(f.GORMTags) != len(exp.GORMTags) {
					t.Errorf("field[%d].GORMTags = %v, want %v", i, f.GORMTags, exp.GORMTags)
				}
			}
		})
	}
}

// TestParseFieldsInvalid tests parsing invalid field definitions (AC: 1, 3)
func TestParseFieldsInvalid(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		errContains string
	}{
		{"empty string", "", "fields cannot be empty"},
		{"no type", "title", "invalid field format"},
		{"unsupported type", "title:date", "unsupported type"},
		{"invalid modifier", "title:string:unknown", "unsupported GORM modifier"},
		{"empty field name", ":string", "field name cannot be empty"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseFields(tt.input)
			if err == nil {
				t.Fatalf("parseFields(%q) returned nil error, want error containing %q", tt.input, tt.errContains)
			}
			if !containsString(err.Error(), tt.errContains) {
				t.Errorf("parseFields(%q) error = %q, want error containing %q", tt.input, err.Error(), tt.errContains)
			}
		})
	}
}

// TestValidateModelName tests model name validation (AC: 3)
func TestValidateModelName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid PascalCase", "Todo", false},
		{"valid multi-word", "BlogPost", false},
		{"valid single char", "T", false},
		{"empty name", "", true},
		{"lowercase start", "todo", true},
		{"contains space", "Blog Post", true},
		{"contains hyphen", "Blog-Post", true},
		{"contains underscore", "Blog_Post", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateModelName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateModelName(%q) error = %v, wantErr = %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

// TestToPascalCase tests PascalCase conversion
func TestToPascalCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"title", "Title"},
		{"Title", "Title"},
		{"name", "Name"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := toPascalCase(tt.input)
			if got != tt.want {
				t.Errorf("toPascalCase(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// TestToSnakeCase tests snake_case conversion
func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Title", "title"},
		{"BlogPost", "blog_post"},
		{"CreatedAt", "created_at"},
		{"ID", "id"},
		{"UserID", "user_id"},
		{"APIKey", "api_key"},
		{"HTTPSConnection", "https_connection"},
		{"OAuth2Token", "o_auth2_token"},
		{"V2API", "v2_api"},
		{"Product2Name", "product2_name"},
		{"name", "name"},
		{"TODO", "todo"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := toSnakeCase(tt.input)
			if got != tt.want {
				t.Errorf("toSnakeCase(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstr(s, substr))
}

func containsSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

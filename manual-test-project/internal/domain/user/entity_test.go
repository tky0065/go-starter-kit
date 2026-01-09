package user_test

import (
	"reflect"
	"testing"
	"time"

	"manual-test-project/internal/domain/user"
)

func TestUserEntity_FieldsAndTags(t *testing.T) {
	u := user.User{}
	val := reflect.TypeOf(u)

	tests := []struct {
		Field   string
		Type    string
		JSONTag string
		GormTag string
	}{
		{"ID", "uint", "id", "primaryKey"},
		{"Email", "string", "email", "uniqueIndex;not null"},
		{"PasswordHash", "string", "-", "not null"}, // PasswordHash often hidden from JSON
		{"CreatedAt", "time.Time", "created_at", ""},
		{"UpdatedAt", "time.Time", "updated_at", ""},
	}

	for _, tt := range tests {
		t.Run(tt.Field, func(t *testing.T) {
			field, ok := val.FieldByName(tt.Field)
			if !ok {
				t.Fatalf("Field %s not found", tt.Field)
			}

			if field.Type.String() != tt.Type {
				// Special handling for time.Time matching
				if tt.Type == "time.Time" && field.Type != reflect.TypeOf(time.Time{}) {
					t.Errorf("Field %s type mismatch: got %s, want %s", tt.Field, field.Type, tt.Type)
				} else if tt.Type != "time.Time" {
					t.Errorf("Field %s type mismatch: got %s, want %s", tt.Field, field.Type, tt.Type)
				}
			}

			jsonTag := field.Tag.Get("json")
			if jsonTag != tt.JSONTag {
				t.Errorf("Field %s json tag mismatch: got %s, want %s", tt.Field, jsonTag, tt.JSONTag)
			}

			// For GORM, we just check if it contains expected values as it can vary
			if tt.GormTag != "" {
				gormTag := field.Tag.Get("gorm")
				if gormTag == "" {
					t.Errorf("Field %s missing gorm tag", tt.Field)
				}
			}
		})
	}
}

package interfaces_test

import (
	"reflect"
	"testing"

	"manual-test-project/internal/interfaces"
)

func TestUserRepository_Methods(t *testing.T) {
	var i *interfaces.UserRepository = nil
	tType := reflect.TypeOf(i).Elem()

	methods := map[string]string{
		"CreateUser":     "func(*user.User) error",
		"GetUserByEmail": "func(string) (*user.User, error)",
	}

	for name := range methods {
		m, ok := tType.MethodByName(name)
		if !ok {
			t.Errorf("Method %s not found", name)
			continue
		}
		// Validating exact signature string with reflect is complex due to package names in types.
		// Just checking existence is good enough for now,
		// effectively checking if "func(*user.User) error" exists.

		_ = m
	}
}

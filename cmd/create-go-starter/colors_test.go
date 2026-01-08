package main

import (
	"testing"
)

func TestColors(t *testing.T) {
	msg := "test"

	// Test Green
	expectedGreen := "\033[32mtest\033[0m"
	if got := Green(msg); got != expectedGreen {
		t.Errorf("Green() = %q, want %q", got, expectedGreen)
	}

	// Test Red
	expectedRed := "\033[31mtest\033[0m"
	if got := Red(msg); got != expectedRed {
		t.Errorf("Red() = %q, want %q", got, expectedRed)
	}
}

package tui

import (
	"testing"
	"time"
)

// TestAnimationStateCreation verifies AnimationState can be created.
// Story 10.7 Task 1.5: Animations fade-in du logo
func TestAnimationStateCreation(t *testing.T) {
	anim := NewAnimationState(10) // 10 frames

	if anim.frame != 0 {
		t.Errorf("Expected initial frame to be 0, got %d", anim.frame)
	}

	if anim.maxFrames != 10 {
		t.Errorf("Expected maxFrames to be 10, got %d", anim.maxFrames)
	}

	if anim.done {
		t.Error("Expected animation to not be done initially")
	}
}

// TestAnimationStateProgress verifies animation progresses through frames.
func TestAnimationStateProgress(t *testing.T) {
	anim := NewAnimationState(3) // 3 frames

	// Initially at frame 0
	if anim.Opacity() != 0.0 {
		t.Errorf("Expected opacity 0.0 at start, got %f", anim.Opacity())
	}

	// Advance to frame 1
	anim.Advance()
	if anim.frame != 1 {
		t.Errorf("Expected frame 1 after advance, got %d", anim.frame)
	}
	expectedOpacity := 1.0 / 3.0
	if anim.Opacity() < expectedOpacity-0.01 || anim.Opacity() > expectedOpacity+0.01 {
		t.Errorf("Expected opacity ~%f, got %f", expectedOpacity, anim.Opacity())
	}

	// Advance to frame 2
	anim.Advance()
	if anim.frame != 2 {
		t.Errorf("Expected frame 2, got %d", anim.frame)
	}

	// Advance to frame 3 (max)
	anim.Advance()
	if anim.frame != 3 {
		t.Errorf("Expected frame 3, got %d", anim.frame)
	}
	if !anim.IsDone() {
		t.Error("Expected animation to be done after reaching maxFrames")
	}
	if anim.Opacity() != 1.0 {
		t.Errorf("Expected opacity 1.0 at end, got %f", anim.Opacity())
	}
}

// TestAnimationStateDoesNotExceedMax verifies animation stops at maxFrames.
func TestAnimationStateDoesNotExceedMax(t *testing.T) {
	anim := NewAnimationState(2)

	// Advance beyond max
	anim.Advance() // frame 1
	anim.Advance() // frame 2 (max)
	anim.Advance() // should stay at frame 2

	if anim.frame != 2 {
		t.Errorf("Expected frame to stay at 2, got %d", anim.frame)
	}

	if !anim.IsDone() {
		t.Error("Expected animation to be done")
	}
}

// TestAnimationTickMessage verifies tick message can be created.
func TestAnimationTickMessage(t *testing.T) {
	msg := AnimationTickMsg{Elapsed: 50 * time.Millisecond}

	if msg.Elapsed != 50*time.Millisecond {
		t.Errorf("Expected elapsed 50ms, got %v", msg.Elapsed)
	}
}

// TestModelWithAnimation verifies Model can track animation state.
func TestModelWithAnimation(t *testing.T) {
	model := NewInteractiveModel(InteractiveDefaults{}, dummyGenerator)

	// Model should have animation state initialized
	if model.logoAnimation == nil {
		t.Error("Expected logoAnimation to be initialized in Model")
	}

	// Animation should not be done initially
	if model.logoAnimation.IsDone() {
		t.Error("Expected logo animation to not be done initially")
	}
}

// TestTransitionStateCreation verifies TransitionState can be created.
// Story 10.7 Task 10: Screen transitions
func TestTransitionStateCreation(t *testing.T) {
	// Test fade-in transition
	fadeIn := NewTransition(TransitionFadeIn, 10)

	if fadeIn == nil {
		t.Fatal("Expected transition to be created")
	}

	if fadeIn.Type != TransitionFadeIn {
		t.Errorf("Expected TransitionFadeIn, got %d", fadeIn.Type)
	}

	if fadeIn.IsDone() {
		t.Error("Expected transition to not be done initially")
	}

	// Test fade-out transition
	fadeOut := NewTransition(TransitionFadeOut, 10)

	if fadeOut.Type != TransitionFadeOut {
		t.Errorf("Expected TransitionFadeOut, got %d", fadeOut.Type)
	}
}

// TestTransitionFadeIn verifies fade-in opacity progression.
func TestTransitionFadeIn(t *testing.T) {
	transition := NewTransition(TransitionFadeIn, 4)

	// Initially at opacity 0.0 (transparent)
	if transition.Opacity() != 0.0 {
		t.Errorf("Expected opacity 0.0 at start of fade-in, got %f", transition.Opacity())
	}

	// Advance to 25%
	transition.Advance()
	expected := 0.25
	if transition.Opacity() < expected-0.01 || transition.Opacity() > expected+0.01 {
		t.Errorf("Expected opacity ~%f, got %f", expected, transition.Opacity())
	}

	// Advance to 50%
	transition.Advance()
	expected = 0.5
	if transition.Opacity() < expected-0.01 || transition.Opacity() > expected+0.01 {
		t.Errorf("Expected opacity ~%f, got %f", expected, transition.Opacity())
	}

	// Complete the transition
	transition.Advance()
	transition.Advance()

	// Should reach opacity 1.0 (fully visible)
	if transition.Opacity() != 1.0 {
		t.Errorf("Expected opacity 1.0 at end of fade-in, got %f", transition.Opacity())
	}

	if !transition.IsDone() {
		t.Error("Expected transition to be done")
	}
}

// TestTransitionFadeOut verifies fade-out opacity progression.
func TestTransitionFadeOut(t *testing.T) {
	transition := NewTransition(TransitionFadeOut, 4)

	// Initially at opacity 1.0 (fully visible)
	if transition.Opacity() != 1.0 {
		t.Errorf("Expected opacity 1.0 at start of fade-out, got %f", transition.Opacity())
	}

	// Advance to 75%
	transition.Advance()
	expected := 0.75
	if transition.Opacity() < expected-0.01 || transition.Opacity() > expected+0.01 {
		t.Errorf("Expected opacity ~%f, got %f", expected, transition.Opacity())
	}

	// Advance to 50%
	transition.Advance()
	expected = 0.5
	if transition.Opacity() < expected-0.01 || transition.Opacity() > expected+0.01 {
		t.Errorf("Expected opacity ~%f, got %f", expected, transition.Opacity())
	}

	// Complete the transition
	transition.Advance()
	transition.Advance()

	// Should reach opacity 0.0 (fully transparent)
	if transition.Opacity() != 0.0 {
		t.Errorf("Expected opacity 0.0 at end of fade-out, got %f", transition.Opacity())
	}

	if !transition.IsDone() {
		t.Error("Expected transition to be done")
	}
}

// TestPulseStateCreation verifies PulseState can be created.
// Story 10.7 Task 10: Pulsing progress bar gradient
func TestPulseStateCreation(t *testing.T) {
	pulse := NewPulseState(20, 0.5, 1.0)

	if pulse == nil {
		t.Fatal("Expected pulse to be created")
	}

	if pulse.maxFrames != 20 {
		t.Errorf("Expected maxFrames 20, got %d", pulse.maxFrames)
	}

	if pulse.minOpacity != 0.5 {
		t.Errorf("Expected minOpacity 0.5, got %f", pulse.minOpacity)
	}

	if pulse.maxOpacity != 1.0 {
		t.Errorf("Expected maxOpacity 1.0, got %f", pulse.maxOpacity)
	}

	// Should start at minimum intensity
	if pulse.Intensity() != 0.5 {
		t.Errorf("Expected initial intensity 0.5, got %f", pulse.Intensity())
	}
}

// TestPulseProgression verifies pulse cycles through intensity values.
func TestPulseProgression(t *testing.T) {
	pulse := NewPulseState(10, 0.5, 1.0)

	// Initially at minimum (0.5)
	if pulse.Intensity() != 0.5 {
		t.Errorf("Expected intensity 0.5 at start, got %f", pulse.Intensity())
	}

	// Advance 5 frames (halfway = max intensity)
	for i := 0; i < 5; i++ {
		pulse.Advance()
	}

	// Should be at midpoint (0.75)
	expected := 0.75
	if pulse.Intensity() < expected-0.01 || pulse.Intensity() > expected+0.01 {
		t.Errorf("Expected intensity ~%f at midpoint, got %f", expected, pulse.Intensity())
	}

	// Advance to max (frame 10)
	for i := 0; i < 5; i++ {
		pulse.Advance()
	}

	// Should be at maximum (1.0)
	if pulse.Intensity() != 1.0 {
		t.Errorf("Expected intensity 1.0 at peak, got %f", pulse.Intensity())
	}

	// Continue pulsing - should reverse direction
	pulse.Advance()

	// Should start decreasing
	if pulse.Intensity() >= 1.0 {
		t.Errorf("Expected intensity to decrease after peak, got %f", pulse.Intensity())
	}
}

// TestPulseCycles verifies pulse continues cycling indefinitely.
func TestPulseCycles(t *testing.T) {
	pulse := NewPulseState(4, 0.0, 1.0)

	// First cycle: 0.0 -> 0.25 -> 0.5 -> 0.75 -> 1.0
	for i := 0; i < 4; i++ {
		pulse.Advance()
	}

	if pulse.Intensity() != 1.0 {
		t.Errorf("Expected intensity 1.0 at first peak, got %f", pulse.Intensity())
	}

	// Second cycle: 1.0 -> 0.75 -> 0.5 -> 0.25 -> 0.0
	for i := 0; i < 4; i++ {
		pulse.Advance()
	}

	if pulse.Intensity() != 0.0 {
		t.Errorf("Expected intensity 0.0 at trough, got %f", pulse.Intensity())
	}

	// Third cycle: 0.0 -> 0.25 -> 0.5 -> 0.75 -> 1.0
	for i := 0; i < 4; i++ {
		pulse.Advance()
	}

	if pulse.Intensity() != 1.0 {
		t.Errorf("Expected intensity 1.0 at second peak, got %f", pulse.Intensity())
	}
}

// TestPulseReset verifies pulse can be reset to initial state.
func TestPulseReset(t *testing.T) {
	pulse := NewPulseState(10, 0.3, 0.9)

	// Advance the pulse
	for i := 0; i < 7; i++ {
		pulse.Advance()
	}

	// Intensity should have changed
	if pulse.Intensity() == 0.3 {
		t.Error("Expected intensity to have changed after advancing")
	}

	// Reset the pulse
	pulse.Reset()

	// Should be back at minimum
	if pulse.Intensity() != 0.3 {
		t.Errorf("Expected intensity 0.3 after reset, got %f", pulse.Intensity())
	}

	if pulse.frame != 0 {
		t.Errorf("Expected frame 0 after reset, got %d", pulse.frame)
	}

	if pulse.direction != 1 {
		t.Errorf("Expected direction 1 after reset, got %d", pulse.direction)
	}
}

// TestProgressPulseTickMessage verifies pulse tick message can be created.
func TestProgressPulseTickMessage(t *testing.T) {
	msg := ProgressPulseTickMsg{Elapsed: 50 * time.Millisecond}

	if msg.Elapsed != 50*time.Millisecond {
		t.Errorf("Expected elapsed 50ms, got %v", msg.Elapsed)
	}
}

// TestModelWithPulse verifies Model can track pulse animation state.
func TestModelWithPulse(t *testing.T) {
	model := NewInteractiveModel(InteractiveDefaults{}, dummyGenerator)

	// Model should have pulse state initialized
	if model.progressPulse == nil {
		t.Error("Expected progressPulse to be initialized in Model")
	}

	// Pulse should start at minimum intensity
	if model.progressPulse.Intensity() != 0.7 {
		t.Errorf("Expected pulse intensity 0.7, got %f", model.progressPulse.Intensity())
	}
}

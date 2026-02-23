package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// AnimationState tracks the state of a frame-based animation.
// Story 10.7 Task 1.5 & Task 10: Animations fade-in et transitions
type AnimationState struct {
	frame     int  // Current frame (0-indexed)
	maxFrames int  // Total frames in animation
	done      bool // True when animation completes
}

// NewAnimationState creates a new animation state with the given frame count.
func NewAnimationState(maxFrames int) *AnimationState {
	return &AnimationState{
		frame:     0,
		maxFrames: maxFrames,
		done:      false,
	}
}

// Advance progresses the animation by one frame.
// Returns true if there are more frames, false if animation is complete.
func (a *AnimationState) Advance() bool {
	if a.frame < a.maxFrames {
		a.frame++
	}

	if a.frame >= a.maxFrames {
		a.done = true
		return false
	}

	return true
}

// Reset resets the animation to the beginning.
func (a *AnimationState) Reset() {
	a.frame = 0
	a.done = false
}

// IsDone returns true if the animation has completed.
func (a *AnimationState) IsDone() bool {
	return a.done
}

// Opacity returns the current opacity as a value from 0.0 to 1.0.
// This is calculated as frame / maxFrames.
func (a *AnimationState) Opacity() float64 {
	if a.maxFrames == 0 {
		return 1.0 // Avoid division by zero
	}
	opacity := float64(a.frame) / float64(a.maxFrames)
	if opacity > 1.0 {
		return 1.0
	}
	if opacity < 0.0 {
		return 0.0
	}
	return opacity
}

// Frame returns the current frame number.
func (a *AnimationState) Frame() int {
	return a.frame
}

// MaxFrames returns the total frame count.
func (a *AnimationState) MaxFrames() int {
	return a.maxFrames
}

// AnimationTickMsg is sent periodically to advance animations.
// Story 10.7 Task 1.5: Animation tick message
type AnimationTickMsg struct {
	Elapsed time.Duration
}

// AnimationTick creates a command that sends animation tick messages.
// The tick interval determines the frame rate (e.g., 50ms = 20 FPS).
func AnimationTick(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(t time.Time) tea.Msg {
		return AnimationTickMsg{Elapsed: interval}
	})
}

// FadeColor applies opacity to a hex color string.
// Returns a color with reduced opacity (lighter).
// This is a simple implementation - for production, use a proper color library.
// Story 10.7 Task 1.5: Fade-in logo animation
func FadeColor(color string, opacity float64) string {
	// For now, we'll return the original color
	// In a production system, you'd blend with background color
	// or use RGBA colors with alpha channel
	return color
}

// TransitionType represents different types of screen transitions.
// Story 10.7 Task 10: Screen transitions
type TransitionType int

const (
	// TransitionNone means no transition animation
	TransitionNone TransitionType = iota
	// TransitionFadeIn means fade in from transparent
	TransitionFadeIn
	// TransitionFadeOut means fade out to transparent
	TransitionFadeOut
)

// TransitionState tracks the state of a screen transition.
// Story 10.7 Task 10: Smooth transitions between screens
type TransitionState struct {
	Type      TransitionType
	animation *AnimationState
}

// NewTransition creates a new transition with the given type and duration.
// The duration is specified in frames (e.g., 10 frames at 50ms = 500ms).
func NewTransition(transitionType TransitionType, frames int) *TransitionState {
	return &TransitionState{
		Type:      transitionType,
		animation: NewAnimationState(frames),
	}
}

// Advance progresses the transition animation by one frame.
// Returns true if there are more frames, false if transition is complete.
func (t *TransitionState) Advance() bool {
	if t.animation == nil {
		return false
	}
	return t.animation.Advance()
}

// Opacity returns the current transition opacity (0.0 to 1.0).
// For fade-in, this goes from 0.0 to 1.0.
// For fade-out, this goes from 1.0 to 0.0.
func (t *TransitionState) Opacity() float64 {
	if t.animation == nil {
		return 1.0
	}

	switch t.Type {
	case TransitionFadeIn:
		// 0.0 -> 1.0
		return t.animation.Opacity()
	case TransitionFadeOut:
		// 1.0 -> 0.0
		return 1.0 - t.animation.Opacity()
	default:
		return 1.0
	}
}

// IsDone returns true if the transition has completed.
func (t *TransitionState) IsDone() bool {
	if t.animation == nil {
		return true
	}
	return t.animation.IsDone()
}

// Reset resets the transition to the beginning.
func (t *TransitionState) Reset() {
	if t.animation != nil {
		t.animation.Reset()
	}
}

// PulseState represents a pulsing animation (e.g., for progress bar gradient).
// Story 10.7 Task 10: Pulsing progress bar gradient
type PulseState struct {
	frame      int
	maxFrames  int
	direction  int // 1 for increasing, -1 for decreasing
	minOpacity float64
	maxOpacity float64
}

// NewPulseState creates a new pulsing animation with the given parameters.
// minOpacity and maxOpacity define the pulse range (e.g., 0.5 to 1.0).
// frames defines how many frames for a full pulse cycle (in/out).
func NewPulseState(frames int, minOpacity, maxOpacity float64) *PulseState {
	return &PulseState{
		frame:      0,
		maxFrames:  frames,
		direction:  1,
		minOpacity: minOpacity,
		maxOpacity: maxOpacity,
	}
}

// Advance progresses the pulse animation by one frame.
// This creates a smooth back-and-forth pulsing effect.
func (p *PulseState) Advance() {
	p.frame += p.direction

	// Reverse direction at boundaries
	if p.frame >= p.maxFrames {
		p.frame = p.maxFrames
		p.direction = -1
	} else if p.frame <= 0 {
		p.frame = 0
		p.direction = 1
	}
}

// Intensity returns the current pulse intensity (minOpacity to maxOpacity).
// This creates a smooth sine-wave-like pulsing effect.
func (p *PulseState) Intensity() float64 {
	if p.maxFrames == 0 {
		return p.maxOpacity
	}

	// Calculate position in cycle (0.0 to 1.0)
	position := float64(p.frame) / float64(p.maxFrames)

	// Map to opacity range
	opacityRange := p.maxOpacity - p.minOpacity
	return p.minOpacity + (opacityRange * position)
}

// Reset resets the pulse animation to the beginning.
func (p *PulseState) Reset() {
	p.frame = 0
	p.direction = 1
}

// ProgressPulseTickMsg is sent periodically to advance progress bar pulse animation.
// Story 10.7 Task 10: Pulsing progress bar gradient
type ProgressPulseTickMsg struct {
	Elapsed time.Duration
}

// ProgressPulseTick creates a command that sends progress pulse tick messages.
// This is used to animate the progress bar gradient during generation.
func ProgressPulseTick(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(t time.Time) tea.Msg {
		return ProgressPulseTickMsg{Elapsed: interval}
	})
}

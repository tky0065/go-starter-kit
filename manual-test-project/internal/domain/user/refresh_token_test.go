package user

import (
	"testing"
	"time"
)

func TestRefreshToken_IsExpired(t *testing.T) {
	tests := []struct {
		name      string
		expiresAt time.Time
		want      bool
	}{
		{
			name:      "expired token",
			expiresAt: time.Now().Add(-1 * time.Hour),
			want:      true,
		},
		{
			name:      "valid token",
			expiresAt: time.Now().Add(1 * time.Hour),
			want:      false,
		},
		{
			name:      "just expired",
			expiresAt: time.Now().Add(-1 * time.Second),
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := &RefreshToken{
				ExpiresAt: tt.expiresAt,
			}
			if got := rt.IsExpired(); got != tt.want {
				t.Errorf("RefreshToken.IsExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRefreshToken_IsRevoked(t *testing.T) {
	tests := []struct {
		name    string
		revoked bool
		want    bool
	}{
		{
			name:    "revoked token",
			revoked: true,
			want:    true,
		},
		{
			name:    "active token",
			revoked: false,
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := &RefreshToken{
				Revoked: tt.revoked,
			}
			if got := rt.IsRevoked(); got != tt.want {
				t.Errorf("RefreshToken.IsRevoked() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRefreshToken_IsValid(t *testing.T) {
	tests := []struct {
		name      string
		revoked   bool
		expiresAt time.Time
		want      bool
	}{
		{
			name:      "valid token",
			revoked:   false,
			expiresAt: time.Now().Add(1 * time.Hour),
			want:      true,
		},
		{
			name:      "revoked token",
			revoked:   true,
			expiresAt: time.Now().Add(1 * time.Hour),
			want:      false,
		},
		{
			name:      "expired token",
			revoked:   false,
			expiresAt: time.Now().Add(-1 * time.Hour),
			want:      false,
		},
		{
			name:      "expired and revoked",
			revoked:   true,
			expiresAt: time.Now().Add(-1 * time.Hour),
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := &RefreshToken{
				Revoked:   tt.revoked,
				ExpiresAt: tt.expiresAt,
			}
			if got := rt.IsValid(); got != tt.want {
				t.Errorf("RefreshToken.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

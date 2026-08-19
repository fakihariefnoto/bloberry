package api

import "testing"

func TestIsInlineSafeType(t *testing.T) {
	tests := []struct {
		ct     string
		inline bool
	}{
		{"image/png", true},
		{"image/jpeg", true},
		{"image/webp", true},
		{"image/svg+xml", false}, // svg can embed scripts
		{"image/svg", false},
		{"audio/mpeg", true},
		{"video/mp4", true},
		{"text/html", false},
		{"text/javascript", false},
		{"application/pdf", false},
		{"application/octet-stream", false},
		{"", false},
	}
	for _, tt := range tests {
		if got := isInlineSafeType(tt.ct); got != tt.inline {
			t.Errorf("isInlineSafeType(%q) = %v, want %v", tt.ct, got, tt.inline)
		}
	}
}

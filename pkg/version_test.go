package release

import (
	"testing"
)

func TestUpdateVersion(t *testing.T) {
	tests := []struct {
		name     string
		input    Version
		changeType semverChange
		expected Version
	}{
		{"Major", Version{1, 1, 1}, major, Version{2, 0, 0}},
		{"Minor", Version{1, 1, 1}, minor, Version{1, 2, 0}},
		{"Patch", Version{1, 1, 1}, patch, Version{1, 1, 2}},
		{"Noop",  Version{1, 1, 1}, noop,  Version{1, 1, 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := updateVersion(tt.input, tt.changeType)
			if result != tt.expected {
				t.Errorf("Expected: %+v, Got: %+v", tt.expected, result)
			}
		})
	}
}

package utils

import (
	"testing"

	"github.com/google/uuid"
)

func TestIsValidUUID(t *testing.T) {
	tests := []struct {
		name     string
		input    uuid.UUID
		expected bool
	}{
		{
			name:     "Valid UUID",
			input:    uuid.New(),
			expected: true,
		},
		{
			name:     "Nil UUID",
			input:    uuid.Nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidUUID(tt.input)
			if result != tt.expected {
				t.Errorf("IsValidUUID(%v) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

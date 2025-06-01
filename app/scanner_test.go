package main

import (
	"testing"
)

func TestMatching(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		pattern  string
		expected bool
	}{
		{
			"parentheses",
			"(()",
			"",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Scan(tt.input)
			// if result != tt.expected {
			// 	t.Errorf("incorrect result for %v, %v", string(tt.input), tt.pattern)
			// }
		})
	}
}

package scan

import (
	"testing"
)

func TestSingleCharacterTokens(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		pattern string
	}{
		{
			"parentheses",
			"(()",
			"",
		},
		{
			"brackets",
			"{{}}",
			"",
		},
		{
			"other-tokens",
			"({*.,+*})",
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := Scan(tt.input)
			if len(scanner.tokens) != len(tt.input)+1 {
				t.Errorf("incorrect result for %v, %v", string(tt.input), tt.pattern)
			}
		})
	}
}

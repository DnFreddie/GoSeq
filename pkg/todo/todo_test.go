package todo

import (
	"testing"
)

func TestContainsPattern(t *testing.T) {
	testCases := []struct {
		name      string
		line      string
		pattern   Pattern
		shouldErr bool
		message   string
	}{
		{"Test TODO", "TODO!! This is an urgent task", TODO, false, "! This is an urgent task"},
		{"Test IDEA", "IDEA: This is an Idea", IDEA, false, ": This is an Idea"},
		{"Test NOTE", "NOTE! This is the Note!", NOTE, false, "This is the Note!"},
		{"Test no pattern", "Random text without patterns", TODO, true, ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ContainsPattern(tc.line, 1, tc.pattern)

			if tc.shouldErr {
				if result != nil {
					t.Errorf("Expected result to be nil, but got: %v", result)
				}
			} else {
				if result == nil {
					t.Errorf("Expected a valid result, but got nil")
				} else if result.Title != tc.message {
					t.Errorf("Expected result title to be '%s', but got '%s'", tc.message, result.Title)
				}
			}
		})
	}
}

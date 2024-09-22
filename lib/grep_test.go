package lib

import (
	"regexp"
	"testing"
)

func TestSearchFunctions(t *testing.T) {
	testCases := []struct {
		name              string
		line              string
		pattern           string
		flag              GrepFlag
		expectMatch       bool
		expectedHighlight string
	}{
		{"ToLower: Matches lowercase", "TODO!! This is an urgent task", "todo", ToLower, true, "\033[31mTODO\033[0m!! This is an urgent task"},
		{"ToLower: No match uppercase", "TODO!! This is an urgent task", "URGENT", ToLower, true, "TODO!! This is an \033[31murgent\033[0m task"},
		{"Normal: Exact match", "IDEA: This is an Idea", "IDEA", 0, true, "\033[31mIDEA\033[0m: This is an Idea"},
		{"Normal: No match", "IDEA: This is an Idea", "NOTE", 0, false, ""},
		{"Regex: Matches pattern", "NOTE! This is the Note!", "Note!$", Regex, true, "NOTE! This is the \033[31mNote!\033[0m"},
		{"Regex: No match", "Random text without patterns", "pattern$", Regex, false, ""},
		{"Regex and ToLower: Case-insensitive match", "HELLO World", "hello.*LD", Regex | ToLower, true, "\033[31mHELLO World\033[0m"},
		{"Regex and ToLower: No match", "HELLO World", "hello.*XYZ", Regex | ToLower, false, ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var matched bool
			var highlightedMatch string

			switch tc.flag {
			case ToLower:
				matched, highlightedMatch = searchToLower(tc.line, tc.pattern)
			case Regex:
				re := regexp.MustCompile(tc.pattern)
				matched, highlightedMatch = searchRegex(tc.line, re)
			case Regex | ToLower:
				re := regexp.MustCompile("(?i)" + tc.pattern)
				matched, highlightedMatch = searchRegex(tc.line, re)
			default:
				matched, highlightedMatch = searchNormal(tc.line, tc.pattern)
			}

			if matched != tc.expectMatch {
				t.Errorf("Expected match: %v, got: %v", tc.expectMatch, matched)
			}

			if tc.expectMatch && highlightedMatch != tc.expectedHighlight {
				t.Errorf("Expected highlight: %q, got: %q", tc.expectedHighlight, highlightedMatch)
			}
		})
	}
}

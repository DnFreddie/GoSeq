package lib

import "strings"

type Pattern uint

const (
	NONE Pattern = 0
	TODO Pattern = 1 << iota
	IDEA
	NOTE
	ALL Pattern = TODO | IDEA | NOTE
)

type Todo struct {
	Keyword       string
	Urgency       int
	ID            *string
	Filename      string
	Line          int
	BodySeparator string
	Title         string `json:"title"`
	Body          string `json:"body"`
	Pattern       Pattern
}

func ContainsPattern(line string, lineIndex int, patterns Pattern) *Todo {

	if patterns&TODO != 0 {
		index := strings.Index(line, "TODO")
		if index != -1 {
			return processMatch(line, lineIndex, index, "TODO", TODO)
		}
	}

	if patterns&IDEA != 0 {
		index := strings.Index(line, "IDEA")
		if index != -1 {
			return processMatch(line, lineIndex, index, "IDEA", IDEA)
		}
	}

	if patterns&NOTE != 0 {
		index := strings.Index(line, "NOTE")
		if index != -1 {
			return processMatch(line, lineIndex, index, "NOTE", NOTE)
		}
	}

	return nil
}

func processMatch(line string, lineIndex int, index int, keyword string, pattern Pattern) *Todo {
	var titleIndex int
	var urgency int

	lastChar := keyword[len(keyword)-1]

	startIndex := index + len(keyword)

	for i := startIndex; i < len(line); i++ {
		if line[i] == lastChar {
			urgency += 1
		} else if line[i] == '!' {
			titleIndex = i + 1
			break
		} else if line[i] != ' ' {
			titleIndex = i
			break
		}
	}

	if titleIndex == 0 {
		titleIndex = startIndex
	}

	title := strings.TrimSpace(line[titleIndex:])
	if title == "" {
		return nil
	}

	return &Todo{
		Keyword: keyword,
		Urgency: urgency,
		Title:   title,
		Line:    lineIndex,
		Pattern: pattern,
	}
}

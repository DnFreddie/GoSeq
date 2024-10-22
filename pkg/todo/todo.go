package todo

import (
	"fmt"
	"strings"
)

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

var patternKeywords = map[Pattern]string{
	TODO: "TODO",
	IDEA: "IDEA",
	NOTE: "NOTE",
}

func ContainsPattern(line string, lineIndex int, patterns Pattern) *Todo {
	for pattern, keyword := range patternKeywords {
		if patterns&pattern != 0 {
			if index := strings.Index(line, keyword); index != -1 {
				return processMatch(line, lineIndex, index, keyword, pattern)
			}
		}
	}
	return nil
}

func processMatch(line string, lineIndex int, index int, keyword string, pattern Pattern) *Todo {
	var titleIndex, urgency int
	lastChar := keyword[len(keyword)-1]
	startIndex := index + len(keyword)

	for i := startIndex; i < len(line); i++ {
		if line[i] == lastChar {
			urgency++
		} else if line[i] == '!' || line[i] != ' ' {
			titleIndex = i
			if line[i] == '!' {
				titleIndex++
			}
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

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
)

func colorize(color, text string) string {
	return color + text + colorReset
}

func (todo *Todo) PrettyPrintTodo() {
	fmt.Printf("%s: %s ", colorize(colorCyan, "Keyword"), colorize(colorYellow, todo.Keyword))
	fmt.Printf("%s: %s\n", colorize(colorCyan, "Urgency"), colorize(colorRed, fmt.Sprintf("%d", todo.Urgency)))

	if todo.ID != nil {
		fmt.Printf("%s: %s\n", colorize(colorCyan, "ID"), *todo.ID)
	}

	fmt.Printf("%s: %s ", colorize(colorCyan, "Filename"), colorize(colorGreen, todo.Filename))
	fmt.Printf("%s: %s\n", colorize(colorCyan, "Line"), colorize(colorGreen, fmt.Sprintf("%d", todo.Line)))

	fmt.Printf("%s: %s\n", colorize(colorCyan, "BodySeparator"), todo.BodySeparator)

	fmt.Printf("%s: %s\n", colorize(colorCyan, "Title"), colorize(colorBlue, todo.Title))
	fmt.Printf("%s:\n%s\n", colorize(colorCyan, "Body"), colorize(colorBlue, strings.TrimSpace(todo.Body)))

	patternName, exists := patternKeywords[todo.Pattern]
	if !exists {
		patternName = "UNKNOWN"
	}
	fmt.Printf("%s: %s\n", colorize(colorCyan, "Pattern"), colorize(colorYellow, patternName))

	fmt.Println(strings.Repeat("-", 50))
}

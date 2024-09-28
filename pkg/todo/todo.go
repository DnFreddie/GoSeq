package todo

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

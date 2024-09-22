package lib

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type GrepFlag uint

const (
	Regex GrepFlag = 1 << iota
	ToLower
)

type GrepMatch struct {
	Line     int64
	Match    string
	Location string
}

func GrepFile(file string, pat []byte, flag GrepFlag) ([]GrepMatch, error) {
	var matches []GrepMatch
	var index int64

	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	pattern := string(pat)

	for scanner.Scan() {
		index++
		line := scanner.Text()

		var matched bool
		var highlightedMatch string
		switch flag {

		case Regex:
			matched, highlightedMatch = searchRegex(line, pattern, flag)
		case ToLower:
			matched, highlightedMatch = searchToLower(line, pattern)
		default:
			matched, highlightedMatch = searchNormal(line, pattern)
		}

		if matched {
			match := GrepMatch{
				Line:     index,
				Match:    highlightedMatch,
				Location: f.Name(),
			}
			matches = append(matches, match)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("no matches found")
	}

	return matches, nil
}

func highlightMatch(text, match string) string {
	redColor := string(Red)
	resetColor := string(Reset)
	return strings.Replace(text, match, redColor+match+resetColor, 1)
}

func searchRegex(line, pattern string, flag GrepFlag) (bool, string) {
	re, err := regexp.CompilePOSIX(pattern)
	if err != nil {
		return false, ""
	}
	if flag&ToLower != 0 {
		re = regexp.MustCompile("(?i)" + re.String())
	}
	match := re.FindString(line)
	if match != "" {
		return true, highlightMatch(line, match)
	}
	return false, ""
}

func searchNormal(line, pattern string) (bool, string) {
	if strings.Contains(line, pattern) {
		return true, highlightMatch(line, pattern)
	}
	return false, ""
}

func searchToLower(line, pattern string) (bool, string) {
	lowerLine := strings.ToLower(line)
	lowerPattern := strings.ToLower(pattern)
	index := strings.Index(lowerLine, lowerPattern)
	if index != -1 {
		match := line[index : index+len(pattern)]
		return true, highlightMatch(line, match)
	}
	return false, ""
}

func ProcessUserInput(matchArray []map[string][]GrepMatch) error {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Choose the note to open:")
	fmt.Print("#? ")

	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			break
		}

		i, err := strconv.Atoi(text)
		if err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			fmt.Print("#? ")
			continue
		}

		if i < 1 || i > len(matchArray) {
			fmt.Println("Unable to choose  a note")
			fmt.Print("#? ")
			continue
		}

		for k := range matchArray[i-1] {
			if err := Edit(k); err != nil {
				return fmt.Errorf("error editing file %s: %w", k, err)
			}
		}
		break
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}

	return nil
}

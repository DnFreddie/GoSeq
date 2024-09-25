package lib

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/sync/semaphore"
)

type GrepFlag uint

const (
	Regex GrepFlag = 1 << iota
	ToLower
)

type GrepMatch struct {
	Line     int64
	Match    string
	Formated string
}

type Searchable interface {
	GetPath() string
	Format() (string, error)
}




func GrepFile[T Note](searchF T, pat []byte, flag GrepFlag) ([]GrepMatch, error) {
	var matches []GrepMatch
	var index int64
	f, err := os.Open(searchF.GetPath()) 
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	pattern := string(pat)

	var re *regexp.Regexp
	if flag&Regex != 0 {
		var err error
		if flag&ToLower != 0 {
			re, err = regexp.Compile("(?i)" + pattern)
		} else {
			re, err = regexp.Compile(pattern)
		}
		if err != nil {
			return nil, fmt.Errorf("invalid regex pattern: %w", err)
		}
	}

	for scanner.Scan() {
		index++
		line := scanner.Text()
		var matched bool
		var highlightedMatch string

		if flag&Regex != 0 {
			matched, highlightedMatch = searchRegex(line, re)
		} else if flag&ToLower != 0 {
			matched, highlightedMatch = searchToLower(line, pattern)
		} else {
			matched, highlightedMatch = searchNormal(line, pattern)
		}

		if matched {
			formated, err := searchF.Format() 
			if err != nil {
				formated = searchF.GetPath() 
			}
			match := GrepMatch{
				Line:     index,
				Match:    highlightedMatch,
				Formated: formated,
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

func searchRegex(line string, re *regexp.Regexp) (bool, string) {
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

func GrepMulti [T Note](searches []T, toParse string, flag GrepFlag) ([]map[string][]GrepMatch, error) {
	var wg sync.WaitGroup
	sem := semaphore.NewWeighted(10)
	results := make([]map[string][]GrepMatch, 0)
	resultsMutex := &sync.Mutex{}

	for _, s := range searches {
		pattern := []byte(toParse)
		wg.Add(1)
		go func(search Searchable) {
			defer wg.Done()
			ctx := context.Background()
			if err := sem.Acquire(ctx, 1); err != nil {
				return
			}
			defer sem.Release(1)
			matches, err := GrepFile(s, pattern, flag)
			if err != nil {
				return
			}
			if len(matches) > 0 {
				resultsMutex.Lock()
				results = append(results, map[string][]GrepMatch{(s).GetPath(): matches})
				resultsMutex.Unlock()
			}
		}(s)
	}
	wg.Wait()
	if len(results) == 0 {
		return results, fmt.Errorf("No results found")
	}
	return results, nil
}

func OpenMatched(matchArray *[]map[string][]GrepMatch) error {
	formatMatches(matchArray)
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
		if i < 1 || i > len(*matchArray) {
			fmt.Println("Unable to choose a note")
			fmt.Print("#? ")
			continue
		}
		for k := range (*matchArray)[i-1] {
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

func formatMatches(notes *[]map[string][]GrepMatch) {
	for i, note := range *notes {
		for _, matches := range note {
			for _, match := range matches {
				InColors(Green, fmt.Sprintf("%d. ", i+1))
				InColors(Blue, match.Formated+"\n")
				fmt.Printf("Line:%d %s\n", match.Line, match.Match)
			}
		}
	}
}

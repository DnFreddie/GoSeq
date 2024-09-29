package terminal

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"golang.org/x/term"
)

type Key int

const (
	Unknown = iota
	CtrlC
	Backspace
	Enter
	Escape
	UpArrow
	DownArrow
	Other
)

type EscapeCode string

const (
	// Escape codes
	Clear       EscapeCode = "\033[H\033[2J\033[H" // Clear screen and reset cursor
	ResetCursor EscapeCode = "\033[0G"             // Move cursor to the beginning of the line
	HideCursor  EscapeCode = "\033[?25l"
	ShowCursor  EscapeCode = "\033[?25h"
)

func clearTerminal() {
	print(Clear)
}

type Color string

const ( // Colors
	Red    Color = "\033[31m"
	Reset  Color = "\033[0m"
	Green  Color = "\033[32m"
	Blue   Color = "\033[34m"
	Cyan   Color = "\033[36m"
	Yellow Color = "\033[33m"
)

func InColors(c Color, s string) {
	fmt.Print(c, s, Reset)
}

type Term interface {
	Start()
	Close()
	Clear()
}

func Quit(t Term) {
	t.Close()
	t.Clear()
	defer os.Exit(0)
}

func NewTerm() Term {
	newTerm := &Terminal{}
	newTerm.Start()
	return newTerm
}

type Terminal struct {
	oldState *term.State
}

func (t *Terminal) Start() {
	t.startRawMode()
}

func (t *Terminal) startRawMode() {
	var err error
	t.oldState, err = term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(fmt.Sprintf("Failed to set raw mode: %v", err))
	}
	fmt.Print(HideCursor)
}

func (t *Terminal) Close() {
	fmt.Print(ShowCursor)
	t.stopRawMode()
}

func (t *Terminal) stopRawMode() {
	if t.oldState != nil {
		if err := term.Restore(int(os.Stdin.Fd()), t.oldState); err != nil {
			panic(fmt.Sprintf("Failed to restore terminal: %v", err))
		}
		t.oldState = nil
	}
}

func (t *Terminal) Clear() {
	clearTerminal()
}

func read() (Key, rune) {
	buf := make([]byte, 3)
	n, err := os.Stdin.Read(buf)
	if err != nil {
		panic(fmt.Sprintf("Failed to read input: %v", err))
	}
	if n == 0 {
		return Unknown, 0
	}

	switch buf[0] {
	case 3:
		return CtrlC, 0
	case 127:
		return Backspace, 0
	case 13:
		return Enter, 0
	case 27:
		if n > 1 && buf[1] == '[' {
			switch buf[2] {
			case 'A':
				return UpArrow, 0
			case 'B':
				return DownArrow, 0
			}
		}
		return Escape, 0
	default:
		return Other, rune(buf[0])
	}
}

func RunTerm[T any](maps []map[string]T) (map[string]T, error) {
	var input string
	var selectionIndex int

	combinedItems := make(map[string]T)
	for _, m := range maps {
		for k, v := range m {
			combinedItems[k] = v
		}
	}

	if len(combinedItems) == 0 {
		return nil, fmt.Errorf("No items available.")
	}

	term := NewTerm()
	defer term.Close()

	for {
		term.Clear()
		fmt.Printf("> %s\n\n", input)

		filteredItems := filterItems(combinedItems, input)

		//check for index out of range panic
		if len(filteredItems) == 0 {
			selectionIndex = 0
		} else if selectionIndex >= len(filteredItems) {
			selectionIndex = len(filteredItems) - 1
		}
		displayResults(filteredItems, selectionIndex)

		key, r := read()

		switch key {
		case CtrlC:
			Quit(term)

		case Backspace:
			if len(input) > 0 {
				input = input[:len(input)-1]
			}
		case Enter:
			if len(filteredItems) > 0 {
				selected := filteredItems[selectionIndex]
				return map[string]T{selected: combinedItems[selected]}, nil
			}
		case Escape:
			Quit(term)

		case UpArrow:
			if selectionIndex > 0 {
				selectionIndex--
			}
		case DownArrow:
			if selectionIndex < len(filteredItems)-1 {
				selectionIndex++
			}
		case Other:
			input += string(r)
		}
	}
}

func displayResults(filteredItems []string, selectionIndex int) {
	if len(filteredItems) == 0 {
		fmt.Println("No results found.")
	} else {
		fmt.Print(ResetCursor)
		for i, item := range filteredItems {
			if i == selectionIndex {
				InColors(Blue, fmt.Sprintf("> %v\n", item))
				fmt.Print(ResetCursor)
			} else {
				fmt.Printf("  %v\n", item)
				fmt.Print(ResetCursor)
			}
		}
	}
}

func filterItems[T any](items map[string]T, input string) []string {
	filtered := make(map[string]T)
	inputLower := strings.ToLower(input)

	for key := range items {
		if strings.Contains(strings.ToLower(key), inputLower) {
			filtered[key] = items[key]
		}
	}
	keys := make([]string, 0, len(filtered))
	for key := range filtered {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	return keys
}

package lib

import (
	"fmt"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"golang.org/x/term"
)

type EscapeCode string
type Signal int
type Color string

const (
	// Escape codes
	Clear        EscapeCode = "\033[H\033[2J\033[H" // Clear screen and reset cursor
	ResetCursor  EscapeCode = "\033[0G"             // Move cursor to the beginning of the line
	HideCursor   EscapeCode = "\033[?25l"
	ShowCursor   EscapeCode = "\033[?25h"

	// Signals
	CtrlC     Signal = 3
	Backspace Signal = 127
	Enter     Signal = 13
	Escape    Signal = 27

	// Colors
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

func clearTerminal() {
	fmt.Print(Clear)
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

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	// Save terminal state
	oldState, err := term.GetState(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Error getting terminal state:", err)
		return nil, err
	}

	_, err = term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Error setting raw mode:", err)
		return nil, err
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	go func() {
		<-sigs
		term.Restore(int(os.Stdin.Fd()), oldState)
		os.Exit(1)
	}()

	buf := make([]byte, 3) 

	for {
		clearTerminal()
		fmt.Print("> ", input, "\n\n")
		fmt.Print(HideCursor)
		fmt.Print(ResetCursor)

		filteredItems := filterItems(combinedItems, input)

		displayResults(filteredItems, selectionIndex)

		fmt.Print(ResetCursor)

		n, err := os.Stdin.Read(buf)
		if err != nil {
			fmt.Println("Error reading input:", err)
			return nil, err
		}

		if n > 0 {
			switch buf[0] {
			case byte(CtrlC):
				fmt.Print(Clear)
				term.Restore(int(os.Stdin.Fd()), oldState)
				fmt.Print(ShowCursor)
				os.Exit(0)
			case byte(Backspace):
				if len(input) > 0 {
					input = input[:len(input)-1]
				}
			case byte(Enter):
				term.Restore(int(os.Stdin.Fd()), oldState)
				clearTerminal()
				if len(filteredItems) > 0 {
					selectedKey := filteredItems[selectionIndex]
					return map[string]T{selectedKey: combinedItems[selectedKey]}, nil
				}
			case byte(Escape):
				if n > 1 && buf[1] == '[' {
					switch buf[2] {
					case 'A': // Up arrow
						if selectionIndex > 0 {
							selectionIndex--
						}
					case 'B': // Down arrow
						if selectionIndex < len(filteredItems)-1 {
							selectionIndex++
						}
					}
				}
			default:
				input += string(buf[0])
			}
		}
	}
}

func displayResults(filteredItems []string, selectionIndex int) {
	if len(filteredItems) == 0 {
		fmt.Println("No results found.")
	} else {
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


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
	//Escape codes
	Clear        EscapeCode = "\033[H\033[2J\033[H"
	ResetCursor EscapeCode = "\033[0G" // Move cursor to the beginning of the line
	HideCursor   EscapeCode = "\033[?25l"
	ShowCursor   EscapeCode = "\033[?25h"
	//Signals
	CtrlC     Signal = 3 //In ASCI
	Backspace Signal = 127
	Enter     Signal = 13
)

const ( //Corlors
	Red    Color = "\033[31m"
	Reset  Color = "\033[0m"
	Green  Color = "\033[32m"
	Blue   Color = "\033[34m"
	Cyan   Color = "\033[36m"
	Yellow Color = "\033[33m"
)

func inColors(c Color, s string) {
	fmt.Print(c, s, Reset)
}

func clearTerminal() {
	fmt.Print(Clear)
}

func RunTerm[T any](items []map[string]T) (map[string]T, error) {
	var input string

	if len(items) == 0 {
		return nil, fmt.Errorf("No notes, create one")
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

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

	buf := make([]byte, 1)

	for {
		clearTerminal()
		fmt.Print("> ", input, "\n\n")
		fmt.Print(HideCursor)

		filteredItems := filterItems(items, input)
		displayResults(filteredItems, input)

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
				fmt.Print(ResetCursor)
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
					return filteredItems, nil
				}

			default:
				input += string(buf[0])
			}
		}
	}
}

func displayResults[T any](filteredItems map[string]T, input string) {
	if len(filteredItems) == 0 {
		fmt.Println("No results found.")
	} else {
		index := 0
		for key , _ := range filteredItems {
			fmt.Print(ResetCursor)
			if index == 0 {
				inColors(Blue, fmt.Sprintf(">%v \n", key))
			} else {
				fmt.Printf("%v\n", key)
			}
			index++
		}
	}
}

// Filter items based on the input adn sort them 
func filterItems[T any](items []map[string]T, input string) map[string]T {
	filtered := make(map[string]T)
	inputLower := strings.ToLower(input)

	for _, itemMap := range items {
		for key, item := range itemMap {
			if strings.Contains(strings.ToLower(key), inputLower) {
				filtered[key] = item
			}
		}
	}

	keys := make([]string, 0, len(filtered))
	for key := range filtered {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	sortedFiltered := make(map[string]T)
	for _, key := range keys {
		sortedFiltered[key] = filtered[key]
	}

	return sortedFiltered
}


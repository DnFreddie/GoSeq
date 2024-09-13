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
	ResetCoursor EscapeCode = "\033[0G" // Move cursor to the beginning of the line
	HideCursor   EscapeCode = "\033[?25l"
	ShowCurosr   EscapeCode = "\033[?25h"
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

func RunTerm(items []string) (string, error) {
	var input string

	if 1 >len(items)  {
		return "", fmt.Errorf("No notes, create one")
	}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	oldState, err := term.GetState(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Error getting terminal state:", err)
		return "", err
	}

	_, err = term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Error setting raw mode:", err)
		return "", err
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

		fmt.Print(ResetCoursor)

		n, err := os.Stdin.Read(buf)
		if err != nil {
			fmt.Println("Error reading input:", err)
			return "", err
		}

		if n > 0 {
			switch buf[0] {
			case byte(CtrlC):
				fmt.Print(Clear)
				term.Restore(int(os.Stdin.Fd()), oldState)
				fmt.Print(ResetCoursor)
				fmt.Print(ShowCurosr)
				os.Exit(0)
			case byte(Backspace):
				if len(input) > 0 {
					input = input[:len(input)-1]
				}
			case byte(Enter):
				term.Restore(int(os.Stdin.Fd()), oldState)
				clearTerminal()
				if len(filteredItems) > 0 {
					choice := filteredItems[0]
					fmt.Print(ShowCurosr)
					return choice, nil
				}

			default:
				input += string(buf[0])
			}
		}
	}
}

func displayResults(filteredItems []string, input string) {
	if len(filteredItems) == 0 {
		fmt.Println("No results found.")
	} else {
		for index, item := range filteredItems {
			fmt.Print(ResetCoursor)
			if index == 0 {
				inColors(Blue, fmt.Sprintf(">%v \n", item))
				continue
			}
			fmt.Printf("%v\n", item)
		}
	}
}

func filterItems(items []string, input string) []string {
	var filtered []string
	inputLower := strings.ToLower(input)
	for _, item := range items {
		if strings.Contains(strings.ToLower(item), inputLower) {
			filtered = append(filtered, item)
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		return strings.HasPrefix(strings.ToLower(filtered[i]), inputLower) && !strings.HasPrefix(strings.ToLower(filtered[j]), inputLower) ||
			(strings.HasPrefix(strings.ToLower(filtered[i]), inputLower) == strings.HasPrefix(strings.ToLower(filtered[j]), inputLower) && len(filtered[i]) < len(filtered[j]))
	})

	return filtered
}

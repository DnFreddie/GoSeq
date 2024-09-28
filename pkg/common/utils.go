package common

import (
	"github.com/DnFreddie/goseq/pkg/grep"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func Edit(fPath string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return fmt.Errorf("Failed to foudn $EDITOR")
	}
	cmd := exec.Command(editor, fPath)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		return err
	}

	return nil
}

func OpenMatched(matchArray *[]map[string][]grep.GrepMatch) error {

	grep.FormatMatches(matchArray)
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

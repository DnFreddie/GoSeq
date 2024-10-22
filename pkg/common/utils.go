package common

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/DnFreddie/goseq/pkg/grep"
	"github.com/rogpeppe/go-internal/lockedfile"
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

// Tries Creating  a Locked File with timeout 3s
func CreteFLocked(path string) (*lockedfile.File, error) {
	done := make(chan struct{})
	var f *lockedfile.File
	var openErr error

	go func() {
		f, openErr = lockedfile.Create(path)
		close(done)
	}()

	select {
	case <-time.After(3 * time.Second):
		return nil, FLockedErr{}

	case <-done:
		if openErr != nil {
			return nil, fmt.Errorf("error opening joined file: %w", openErr)
		}
	}
	return f, nil
}

// Find traverses a directory and applies a condition function to each DirEntry.
// It returns a slice of items of type T that match the condition and a combined error if any errors occurred.
func Find[T any](dir string, condition func(fs.DirEntry) (T, bool)) ([]T, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("error getting absolute path: %w", err)
	}

	var results []T
	var errorsArr []error

	walk := func(path string, d os.DirEntry, err error) error {
		if err != nil {
			errorsArr = append(errorsArr, fmt.Errorf("error accessing path %s: %w", path, err))
			return nil
		}

		if item, ok := condition(d); ok {
			results = append(results, item)
		}
		return nil
	}

	err = filepath.WalkDir(absDir, walk)
	if err != nil {
		errorsArr = append(errorsArr, err)
	}
	return results, errors.Join(errorsArr...)
}

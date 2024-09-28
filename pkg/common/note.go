package common

import (
	"github.com/DnFreddie/goseq/pkg/grep"
	"github.com/DnFreddie/goseq/pkg/interfaces"
	"errors"
	"fmt"
	"sort"
	"sync"
)

type SearchResult struct {
	Note      interfaces.Note
	Matches   []grep.GrepMatch
	Formatted string
}

func Search[T interfaces.Note](notes []T, toParse string, flag grep.GrepFlag) error {
	if len(notes) == 0 {
		return fmt.Errorf("No items found")
	}
	if toParse == "" {
		return fmt.Errorf("No pattern to look for ")
	}
	matches, err := grep.GrepMulti(notes, toParse, flag)
	if err != nil {
		return err
	}

	OpenMatched(&matches)

	if err != nil {
		return err
	}

	return nil
}
func SortNotes[T interfaces.Note](notes []T) {
	sort.Slice(notes, func(i, j int) bool {
		return notes[i].GetDate().Before(notes[j].GetDate())
	})
}

func ScanJoined[T interfaces.Note](scanner interfaces.NoteScanner[T]) error {
	var wg sync.WaitGroup
	nCh := make(chan T)
	errCh := make(chan error, 1)
	var errs []error

	wg.Add(1)
	go func() {
		defer wg.Done()
		for scanner.Scan() {
			note := scanner.Note()
			nCh <- note
		}
		if err := scanner.Err(); err != nil {
			errCh <- err
		}
	}()

	go func() {
		wg.Wait()
		close(nCh)
		close(errCh)
	}()

	for note := range nCh {
		if err := note.Write(); err != nil {
			errs = append(errs, fmt.Errorf("error writing note: %w", err))
		}
	}

	if err := <-errCh; err != nil {
		errs = append(errs, fmt.Errorf("error during scanning: %w", err))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

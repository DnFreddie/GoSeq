package lib

import (
	"errors"
	"fmt"
	"io"
	"sync"
)

type NoteManager[T Note] interface {
	GetNotes(period Period) ([]T, error)
	JoinNotesByTitle(notes *[]T) (io.Reader,error)
	JoinNotesWithContents(notes *[]T) (io.Reader, error)
	Scan(r io.Reader, scanner NoteScanner[T]) ([]T, error)
	DelteByTitle(r io.Reader, n *[]Note)error
}

type NoteScanner[T Note] interface {
	Scan() bool
	Note() T
	Err() error
}

func ScanJoined[T Note](scanner NoteScanner[T]) error {
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

package project

import (
	"bufio"
	"io"
)

type ProjectScanner struct {
	scanner     *bufio.Scanner
	currentNote Project
	err         error
}

func NewDNoteScanner(r io.Reader) *ProjectScanner {
	return &ProjectScanner{
		scanner: bufio.NewScanner(r),
	}
}

func (s *ProjectScanner) Note() Project {
	return s.currentNote
}

func (s *ProjectScanner) Err() error {
	if s.err != nil {
		return s.err
	}
	return s.scanner.Err()
}

func (s *ProjectScanner) Scan() bool {
	return false
}

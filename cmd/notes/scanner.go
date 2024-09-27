package notes

import (
	"bufio"
	"bytes"
	"io"
	"time"

)

type DNoteScanner struct {
	scanner     *bufio.Scanner
	currentNote DNote
	err         error
}

func NewDNoteScanner(r io.Reader) *DNoteScanner {
	return &DNoteScanner{
		scanner: bufio.NewScanner(r),
	}
}

func (s *DNoteScanner) Note() DNote {
	return s.currentNote
}

func (s *DNoteScanner) Err() error {
	if s.err != nil {
		return s.err
	}
	return s.scanner.Err()
}

func (s *DNoteScanner) Scan() bool {
	isCollecting := false
	layout := string(FullDate)

	for s.scanner.Scan() {
		line := s.scanner.Text()
		trimmedLine := bytes.TrimSpace([]byte(line))

		if parsedTime, err := time.Parse(layout, string(trimmedLine)); err == nil {
			if isCollecting {
				return true
			}
			s.currentNote = DNote{Date: parsedTime }

			isCollecting = true
			continue
		}

		if checkSeparator(string(trimmedLine)) {
			if isCollecting {
				s.currentNote.Contents = bytes.TrimRight(s.currentNote.Contents, "\n")
				return true
			}
			continue
		}

		if isCollecting {
			s.currentNote.Contents = append(s.currentNote.Contents, line...)
			s.currentNote.Contents = append(s.currentNote.Contents, '\n')
		}
	}

	if isCollecting {
		s.currentNote.Contents = bytes.TrimRight(s.currentNote.Contents, "\n")
		return true
	}

	return false
}

func checkSeparator(line string) bool {
	if len(line) < 4 || line[0] != '#' {
		return false
	}
	var hyphenCount int
	for i := 1; i < len(line); i++ {
		if line[i] == '-' {
			hyphenCount++
		} else {
			break
		}
	}
	return hyphenCount >= 3
}

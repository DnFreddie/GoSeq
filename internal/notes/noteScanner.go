package notes

import (
	"bufio"
	"bytes"
	"io"
	"strings"
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
	for s.scanner.Scan() {
		line := s.scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		if date, ok := parseDateFromSeparator(trimmedLine); ok {
			if isCollecting {
				return true
			}
			s.currentNote = DNote{Date: date}
			isCollecting = true
			continue
		}

		if checkSeparator(trimmedLine) {
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
	return strings.HasPrefix(line, "#-") && strings.HasSuffix(line, "-")
}

func parseDateFromSeparator(line string) (time.Time, bool) {
	if !strings.HasPrefix(line, "#-") || !strings.HasSuffix(line, "-") {
		return time.Time{}, false
	}

	dateStr := strings.Trim(line, "#-")
	dateStr = strings.TrimSpace(dateStr)

	layouts := []string{
		"January 2 2006",
		"January 03 2006",
		"Jan 2 2006",
		"Jan 02 2006",
	}

	for _, layout := range layouts {
		if date, err := time.Parse(layout, dateStr); err == nil {
			return date, true
		}
	}

	return time.Time{}, false
}

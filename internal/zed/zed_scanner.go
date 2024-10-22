package zed

import "bufio"

type ZedNoteScanner struct {
	scanner     *bufio.Scanner
	currentNote ZedNote
	err         error
}

func (s *ZedNoteScanner) Scan() bool {
	panic("not implemented") // TODO: Implement
}

func (s *ZedNoteScanner) Note() ZedNote {
	panic("not implemented") // TODO: Implement
}

func (s *ZedNoteScanner) Err() error {
	panic("not implemented") // TODO: Implement
}

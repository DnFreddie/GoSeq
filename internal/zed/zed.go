package zed

import (
	"time"
)

// ZedNote implements the Note interface
type ZedNote struct {
	Title   string
	Content string
	Path    string
	Date    time.Time
}

func (zednote *ZedNote) Format() (string, error) {
	panic("not implemented") // TODO: Implement
}

func (zednote *ZedNote) GetPath() string {
	panic("not implemented") // TODO: Implement
}

func (zednote *ZedNote) GetDate() time.Time {
	panic("not implemented") // TODO: Implement
}

func (zednote *ZedNote) Delete() error {
	panic("not implemented") // TODO: Implement
}

func (zednote *ZedNote) Write() error {
	panic("not implemented") // TODO: Implement
}

// ZedNoteManager implements the NoteManager interface

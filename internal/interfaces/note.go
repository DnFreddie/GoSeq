package interfaces

import (
	"io"
	"time"
)

type Note interface {
	Format() (string, error)
	GetPath() string
	GetDate() time.Time
	Delete() error
	Write() error
}

type NoteManager[T Note, E any] interface {
	GetNotes(conditon E) ([]T, error)
	JoinNotesByTitle(notes *[]T) (io.Reader, error)
	JoinNotesWithContents(notes *[]T) (io.Reader, error)
	Scan(r io.Reader, scanner NoteScanner[T]) ([]T, error)
	DelteByTitle(r io.Reader, n *[]Note) error
}

type NoteScanner[T Note] interface {
	Scan() bool
	Note() T
	Err() error
}


package dnotes

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

type Nnote interface {
	GetTitle() string
	GetContent() []byte
	GetPath() string
	GetDate() time.Time
	Format() (string, error)
	Write() error
}
type NoteFormatter interface {
	FormatNote(note *BasicNote)
	ParseNote(content string) (BasicNote, error)
}

// NoteManager handles operations on collections of notes
type NNoteManager interface {
	GetNotes(filter any) ([]BasicNote, error)
	SaveNote(note BasicNote) error
	DeleteNote(note BasicNote) error
	JoinNotes(notes []BasicNote) (io.Reader, error)
	Sort() error
}

// NoteScanner reads notes from a source
type NNoteScanner interface {
	Scan() bool
	NNote() BasicNote
	Err() error
}

// BasicNote provides a standard implementation

// implements [Nnote] interface
type BasicNote struct {
	Title   string
	Content []byte
	Path    string
	Date    time.Time
	Manager *DnoteManager
}

func NewBasicNote(notePath string) (*BasicNote, error) {

	noteName := path.Base(notePath)
	rawDate := strings.Replace(noteName, ".md", "", -1)
	date, err := time.Parse(string(FileDate), rawDate)
	if err != nil {
		return nil, err
	}
	note := BasicNote{
		Path:  notePath,
		Date:  date,
		Title: noteName,
	}

	return &note, nil
}

func (n *BasicNote) GetTitle() string   { return n.Title }
func (n *BasicNote) GetContent() []byte { return n.Content }
func (n *BasicNote) GetPath() string    { return n.Path }
func (n *BasicNote) Delete() error      { return nil }
func (n *BasicNote) GetDate() time.Time { return n.Date }
func (n *BasicNote) Format() (string, error) {
	if n.Manager == nil {
		return "", fmt.Errorf("no manager set for note")
	}
	return "", nil
}
func (n *BasicNote) Write() error {
	if n.Manager == nil {
		return fmt.Errorf("no manager set for note")
	}
	return n.Manager.SaveNote(n)

}

func (n *BasicNote) Read() error {
	f, err := os.Open(n.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	n.Content, err = io.ReadAll(f)
	if err != nil {
		return err
	}
	return nil
}

type StandardNoteScanner struct {
	scanner     *bufio.Scanner
	formatter   NoteFormmatter
	currentNote BasicNote
	err         error
}

func NewNoteScanner(r io.Reader, formatter NoteFormmatter) *StandardNoteScanner {
	return &StandardNoteScanner{
		scanner:   bufio.NewScanner(r),
		formatter: formatter,
	}
}
func sanitizeFileName(name string) string {
	// Replace invalid characters with underscores
	re := regexp.MustCompile(`[^a-zA-Z0-9-_]`)
	return re.ReplaceAllString(name, "_")
}

type NoteFormmatter struct {
	separatorChar  string
	separatorWidth int
}

func NewMarkdownFormatter() *NoteFormmatter {
	return &NoteFormmatter{
		separatorChar:  "-",
		separatorWidth: 30,
	}
}

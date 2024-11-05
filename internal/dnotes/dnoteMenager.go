package dnotes

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/DnFreddie/goseq/internal/common"
)

type DnoteManager struct {
	formatter NoteFormatter
	basePath  string
	Notes     []BasicNote
}

// NewDNoteManager initializes a new DnoteManager with a formatter and base path.
func NewDNoteManager(formatter NoteFormatter) *DnoteManager {
	return &DnoteManager{
		formatter: formatter,
	}
}

// NewBNoteFormatter initializes a new BNoteFormatter with default settings.
func NewBNoteFormatter() *BNoteFormatter {
	return &BNoteFormatter{
		separator: "-",
		sWidth:    30,
		buff:      bytes.Buffer{},
	}
}

// BNoteFormatter is used to format notes with a specified separator and width.
type BNoteFormatter struct {
	separator string
	sWidth    int
	buff      bytes.Buffer
}

func (fm *BNoteFormatter) Clear() {
	fm.buff = bytes.Buffer{}

}

func (fm *BNoteFormatter) ParseNote(content string) (BasicNote, error) {
	return BasicNote{}, nil

}

func (fm *BNoteFormatter) FormatNote(note *BasicNote) {
	content := note.GetContent()
	sep := strings.Repeat(fm.separator, fm.sWidth)
	header := fmt.Sprintf("%s%s%s\n", sep, note.Title, sep)
	footer := fmt.Sprintf("%s\n\n", strings.Repeat("-", len(header)-1))

	fmt.Fprintf(&fm.buff, "%s%s\n\n%s", header, content, footer)
}

// SaveNote saves a note (implementation pending).
func (m *DnoteManager) SaveNote(note *BasicNote) error {
	return nil
}

// GetNotes retrieves notes based on a filter (implementation pending).
func (c *DnoteManager) GetNotes(period common.Period) error {
	notes, err := XGetNotes(period)
	c.Notes = notes
	if err != nil {
		return err
	}
	return nil
}

// DeleteNote deletes a specified note (implementation pending).
func (c *DnoteManager) DeleteNote(note BasicNote) error {
	panic("not implemented") // TODO: Implement
}
func (c *DnoteManager) JoinNotes() (io.Reader, error) {
	return XjoinNotes(&c.Notes)
}

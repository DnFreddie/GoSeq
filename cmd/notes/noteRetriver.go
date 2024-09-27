package notes

import (
	"DnFreddie/goseq/lib"
	"io"
)

type DailyNoteManager struct{}

func NewDailyNoteManager() *DailyNoteManager {
	return &DailyNoteManager{}
}

// Returns the Error of proccesing notes or No notes error if agenda not founded
func (d *DailyNoteManager) GetNotes(p lib.Period) ([]DNote, error) {
	return getNotes(p)
}

func (d *DailyNoteManager) DeleteNotes(notes []DNote) error {
	return nil
}

func (d *DailyNoteManager) JoinNotesWithContents(notes *[]DNote) (io.Reader,error) {
	return joinNotes(notes)
}

func (d *DailyNoteManager) JoinNotesByTitle(notes *[]DNote) ([]DNote, error) {
	return []DNote{}, nil
}

func (d *DailyNoteManager) Scan(r io.Reader, scanner DNoteScanner) ([]DNote, error) {
	var notes []DNote
	for scanner.Scan() {
		notes = append(notes, scanner.Note())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return notes, nil
}

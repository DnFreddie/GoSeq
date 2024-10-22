package zed

import (
	"fmt"
	"io"
	"io/fs"

	"github.com/DnFreddie/goseq/pkg/common"
	"github.com/spf13/viper"
)

type Filter interface {
	common.Period | string
}

type ZedNoteManager[T Filter] struct{}

func NewProjectManager[T Filter]() *ZedNoteManager[T] {
	retriever := ZedNoteManager[T]{}
	return &retriever
}

// GetNotes retrieves notes based on a provided filter
//
// Filter is either common.Period or string
func (m *ZedNoteManager[T]) GetNotes(filter T) ([]ZedNote, error) {
	notesDir := viper.GetString("NOTES")
	if notesDir == "" {
		return nil, fmt.Errorf("notes directory is not configured")
	}

	condition := func(d fs.DirEntry) (ZedNote, bool) {
		switch v := any(filter).(type) {
		case string:
			if d.Name() == v {
				return ZedNote{}, true
			}
		case common.Period:
		}
		return ZedNote{}, false
	}

	zedNotes, err := common.Find(notesDir, condition)
	if err != nil {
		fmt.Printf("error indexing...\n%v", err)
		return nil, err
	}

	if len(zedNotes) == 0 {
		return zedNotes, common.NoNotesFoundErr{}
	}

	return zedNotes, nil
}

func (m *ZedNoteManager[T]) JoinNotesByTitle(notes *[]ZedNote) (io.Reader, error) {
	panic("not implemented") // TODO: Implement
}

func (m *ZedNoteManager[T]) JoinNotesWithContents(notes *[]ZedNote) (io.Reader, error) {
	panic("not implemented") // TODO: Implement
}

func (m *ZedNoteManager[T]) DeleteByTitle(r io.Reader, notes *[]ZedNote) error {
	panic("not implemented") // TODO: Implement
}

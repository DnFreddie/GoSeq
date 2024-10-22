package notes

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unicode"

	"github.com/DnFreddie/goseq/pkg/common"
	"github.com/DnFreddie/goseq/pkg/terminal"
	"github.com/spf13/viper"
)

const (
	JOINED        = "/tmp/.go_seq_notes_joined.md"
	JOINED_DELETE = "/tmp/.go_seq_notes_delete_joined.md"
)

type DailyNoteManager struct{}
type separator string

func NewDailyNoteManager() *DailyNoteManager {
	return &DailyNoteManager{}
}

// Returns the Error of proccesing notes or No notes error if agenda not founded
func (d *DailyNoteManager) GetNotes(p common.Period) ([]DNote, error) {
	return getNotes(p)
}

func (d *DailyNoteManager) DeleteByTitle(r io.Reader, n *[]DNote) error {

	return deleteByTitle(r, n)
}

func (d *DailyNoteManager) JoinNotesWithContents(notes *[]DNote) (io.Reader, error) {
	return joinNotes(notes)
}
func (d *DailyNoteManager) JoinNotesByTitle(notes *[]DNote) (io.Reader, error) {

	return joinByTitle(notes)
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

func getNotes(pr common.Period) ([]DNote, error) {
	var errMessages []string
	noteArray := []DNote{}

	AGENDA := viper.GetString("AGENDA")

	entries, err := os.ReadDir(AGENDA)
	if err != nil {

		return noteArray, &common.NoNotesFoundErr{}
	}

	now := time.Now()

	for _, v := range entries {
		if !v.IsDir() {
			rawDate := strings.Replace(v.Name(), ".md", "", -1)
			date, err := time.Parse(string(FileDate), rawDate)
			if err != nil {
				errMessages = append(errMessages, fmt.Sprintf("failed to parse file %s: %v", v.Name(), err))
				continue
			}

			if !common.DateInRange(now, pr, date) {
				continue
			}

			note := DNote{
				Path: path.Join(AGENDA, v.Name()),
				Date: date,
			}
			noteArray = append(noteArray, note)
		}
	}

	common.SortNotes(noteArray)

	if len(errMessages) > 0 {
		return noteArray, fmt.Errorf(strings.Join(errMessages, "; "))
	}

	return noteArray, nil
}

func joinNotes(notes *[]DNote) (io.Reader, error) {
	if len(*notes) == 0 {
		return nil, common.NoNotesFoundErr{}
	}

	f, err := common.CreteFLocked(JOINED)
	if err != nil {
		return nil, err
	}

	defer common.CleanupFileHandler(f, JOINED)

	var fullBuffer bytes.Buffer
	for _, v := range *notes {
		if err := v.Read(); err != nil {
			log.Printf("Error reading note: %v", err)
			continue
		}
		formated, err := v.Format()
		if err != nil {
			log.Printf("Error formatting note: %v", err)
			continue
		}
		sep := fmt.Sprintf("%v", strings.Repeat("-", 30))
		header := fmt.Sprintf("%v%v%v\n", sep, formated, sep)
		footer := fmt.Sprintf("%s\n\n", strings.Repeat("-", len(header)-1))
		fmt.Fprintf(&fullBuffer, "%s%s\n\n%s",
			header,
			v.Contents,
			footer)
		v.Contents = nil
	}

	if _, err := f.Write(fullBuffer.Bytes()); err != nil {
		return nil, fmt.Errorf("error writing to joined file: %w", err)
	}

	if err := common.Edit(JOINED); err != nil {
		return nil, fmt.Errorf("error editing file: %w", err)
	}

	readFile, err := os.Open(JOINED)
	if err != nil {
		return nil, fmt.Errorf("error opening edited file: %w", err)
	}

	reader := &trimReader{
		reader: bufio.NewReader(readFile),
		file:   readFile,
	}
	return reader, nil
}

type trimReader struct {
	reader *bufio.Reader
	file   *os.File
}

func (tr *trimReader) Close() error {
	return tr.file.Close()
}

func (tr *trimReader) Read(p []byte) (n int, err error) {
	line, err := tr.reader.ReadBytes('\n')
	if err != nil && err != io.EOF {
		return 0, err
	}

	// Trim trailing whitespace
	line = bytes.TrimRightFunc(line, unicode.IsSpace)
	if len(line) > 0 || err == nil {
		line = append(line, '\n')
	}

	n = copy(p, line)
	if n < len(line) {
		err = nil
	}
	return n, err
}

func deleteByTitle(r io.Reader, notes *[]DNote) error {
	var titles []string
	var wasDeleted atomic.Bool
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		titles = append(titles, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	errChan := make(chan error, len(*notes))
	var wg sync.WaitGroup
	for _, note := range *notes {
		formatted, err := note.Format()
		if err != nil {
			formatted = note.GetPath()
		}
		if !slices.Contains(titles, formatted) {
			wg.Add(1)
			go func(n DNote) {
				defer wg.Done()
				if err := n.Delete(); err != nil {
					errChan <- err
				} else {
					wasDeleted.Store(true)
				}
			}(note)
		}
	}
	go func() {
		wg.Wait()
		close(errChan)
	}()
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}
	if len(errors) > 0 {
		return fmt.Errorf("multiple errors occurred: %v", errors)
	}
	if !wasDeleted.Load() {
		terminal.InColors(terminal.Red, "Nothing to delete ...\n")
	}
	return nil
}

func joinByTitle(notes *[]DNote) (io.Reader, error) {
	if len(*notes) == 0 {
		return nil, common.NoNotesFoundErr{}
	}

	f, err := common.CreteFLocked(JOINED_DELETE)
	if err != nil {
		return nil, err
	}

	defer common.CleanupFileHandler(f, JOINED_DELETE)

	var titles []string
	for _, note := range *notes {
		formattedName, err := note.Format()
		if err != nil {
			formattedName = note.GetPath()
		}
		titles = append(titles, formattedName)
	}
	joinedTitles := strings.Join(titles, "\n")

	if _, err := f.Write([]byte(joinedTitles)); err != nil {
		return nil, err
	}

	if err := common.Edit(JOINED_DELETE); err != nil {
		return nil, err
	}

	updatedContent, err := os.ReadFile(JOINED_DELETE)
	if err != nil {
		return nil, fmt.Errorf("failed to read updated file: %w", err)
	}

	return bytes.NewReader(updatedContent), nil
}

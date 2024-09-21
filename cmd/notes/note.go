package notes

import (
	"DnFreddie/goseq/lib"
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Note struct {
	Date     time.Time
	Contents []byte
	Path     string
}

type DateLayout string

const (
	FileDate DateLayout = "2006-01-02"
	FullDate DateLayout = "January 2 2006"
)

func (n *Note) writeNote() error {
	AGENDA := viper.GetString("AGENDA")
	if n.Path == "" && !n.Date.IsZero() {

		n.Path = n.Date.Format("2006-01-02.md")

	}

	absPath := path.Join(AGENDA, n.Path)
	if _, err := os.Stat(absPath); errors.Is(err, os.ErrNotExist) {

		return err

	}
	f, err := os.Create(absPath)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := bufio.NewWriter(f)

	_, err = writer.WriteString(n.parseDate() + "\n")
	if err != nil {
		return err
	}

	_, err = writer.Write(n.Contents)
	if err != nil {
		return err
	}

	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}

func (n *Note) read() error {
	f, err := os.Open(n.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	n.Contents, err = io.ReadAll(f)
	if err != nil {
		return err
	}
	return nil
}

func (n *Note) parseDate() string {
	//layout := "Mon Jan 2 15:04:05 PM MST 2006"
	formated_Date := n.Date.Format(string(FullDate))

	return formated_Date
}

func sortNotes(notes []Note) {

	sort.Slice(notes, func(i, j int) bool {
		return notes[i].Date.Before(notes[j].Date)
	})
}

func dailyNote() error {
	agenda := checkAgenda()
	now := time.Now()
	date := now.Format(string(FileDate))
	formattedTime := now.Format(string(FullDate))

	dailyNote := path.Join(agenda, date+".md")
	if _, err := os.Stat(dailyNote); errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(dailyNote)
		defer f.Close()
		if err != nil {
			return err
		}
		f.Write([]byte(formattedTime))
		f.Write([]byte("\n" + strings.Repeat("-", len(formattedTime))))

	}
	err := lib.Edit(dailyNote)
	if err != nil {
		return err
	}

	return nil
}

func checkAgenda() string {

	HOME := viper.GetString("HOME")
	agenda := path.Join(HOME, "Documents/Agenda")

	if err := os.MkdirAll(agenda, os.FileMode(0755)); err != nil {
		log.Fatal("Failed to create an Agenda Directory", err)
	}

	return agenda
}
func ChoseNote() error {

	AGENDA := viper.GetString("AGENDA")
	entries, err := os.ReadDir(AGENDA)

	if err != nil {
		return err
	}

	var names []map[string]time.Time
	for _, entry := range entries {
		rawDate, err := isNote(entry)

		if err != nil {
			continue
		}

		dateMap := make(map[string]time.Time)
		fmtTime := rawDate.Format(string(FullDate))
		dateMap[fmtTime] = rawDate
		names = append(names, dateMap)

	}
	if len(names) == 0 {
		fmt.Errorf("No DailyNotes found try to create one with goseq new")
		return err
	}

	choice, err := lib.RunTerm(names)
	if err != nil {
		return err

	}
	chosenNote := &Note{}
	for _, v := range choice {
		*chosenNote = Note{
			Path: path.Join(AGENDA, v.Format(string(FileDate)+".md")),
			Date: v,
		}
	}
	err = lib.Edit(chosenNote.Path)
	if err != nil {
		return nil
	}
	return nil
}
func isNote(entry os.DirEntry) (time.Time, error) {
	var rawDate time.Time
	if entry.IsDir() {
		return rawDate, fmt.Errorf("Not a file so not a note")
	}
	dateStirng := strings.Replace(entry.Name(), ".md", "", -1)
	rawDate, err := time.Parse(string(FileDate), dateStirng)

	if err != nil {
		return rawDate, err
	}
	return rawDate, nil
}

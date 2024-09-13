package lib

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"sort"
	"strings"
)

type Token string

const (
	START Token = "START\n"
	END   Token = "END\n\n"
	LINE  Token = "\n-----------------------------------\n"
)

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

func JoinNotes(entries *[]fs.DirEntry) error {
	agenda := path.Join(AGENDA, ".joined_test.md")
	notes := GetNotes(entries)
	sortNotes(notes)

	f, err := os.OpenFile(agenda, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, v := range notes {

		err := v.read()
		if err != nil {
			fmt.Println(err)
			continue
		}
		var buffer bytes.Buffer
		buffer.Write(v.Contents)
		buffer.Write([]byte("END"))
		buffer.Write([]byte("\n\n"))

		_, err = f.Write(buffer.Bytes())

		if err != nil {
			log.Println(err)
		}
		v.Contents = nil
	}

	err = edit(agenda)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
func sortNotes(notes []Note) {

	sort.Slice(notes, func(i, j int) bool {
		return notes[i].Date.Before(notes[j].Date)
	})
}
func GetNotes(e *[]os.DirEntry) []Note {
	var noteArray []Note
	for _, v := range *e {
		if !v.IsDir() {
			raw_date := strings.Replace(v.Name(), ".md", "", -1)
			date, _, err := parseTimeNote(raw_date)
			if err != nil {
				continue
			}
			note := Note{
				Path: path.Join(AGENDA, v.Name()),
				Date: date,
			}
			noteArray = append(noteArray, note)

		}

	}

	return noteArray

}

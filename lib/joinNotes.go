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
	"time"

	"github.com/spf13/viper"
)

type DateRange int

const (
	Day   DateRange = 1
	Week  DateRange = 7
	Month DateRange = 30
	Year  DateRange = 365
	All   DateRange = 0
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

func JoinNotes(entries *[]fs.DirEntry, period DateRange) error {
	join := path.Join(JOINED)
	notes := GetNotes(entries, period)

	f, err := os.OpenFile(join, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
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

	err = Edit(join)
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

func GetNotes(e *[]os.DirEntry, dr DateRange) []Note {
	var noteArray []Note
	for _, v := range *e {
		if !v.IsDir() {
			raw_date := strings.Replace(v.Name(), ".md", "", -1)
			date, err := time.Parse(string(FileDate), raw_date)
			if err != nil {
				continue
			}
			home := viper.GetString("HOME")
			note := Note{

				Path: path.Join(home, AGENDA, v.Name()),
				Date: date,
			}
			noteArray = append(noteArray, note)

		}

	}
	sortNotes(noteArray)

	if dr == All || int(dr) > len(noteArray) {
		return noteArray
	}

	startIndex := len(noteArray) - int(dr)
	return noteArray[startIndex:]
}

package notes

import (
	"DnFreddie/goseq/lib"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"strings"
	"sync"
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

	JOINED = "/tmp/.go_seq_joined.md"
)

func parseDateRange(input string) DateRange {
	switch strings.ToLower(input) {
	case "day":
		return Day
	case "week":
		return Week
	case "month":
		return Month
	case "year":
		return Year
	case "all":
		return All
	default:
		fmt.Printf("Invalid date range: %s. Defaulting to 'all'.\n", input)
		return Week
	}
}
func JoinNotes(entries *[]fs.DirEntry, period DateRange) error {
	join := path.Join(JOINED)
	notes := getNotes(entries, period)

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

	err = lib.Edit(join)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func ScanEverything() {
	var wg sync.WaitGroup
	ch := make(chan Note)
	f, err := os.Open(JOINED)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	wg.Add(1)

	go func() {
		defer wg.Done()
		err := ScanAgenda(f, ch)
		if err != nil {
			log.Fatal("Error scanning agenda:", err)
		}
	}()

	go func() {
		wg.Wait()
		close(ch)
	}()

	for note := range ch {
		wg.Add(1)
		note := note

		go func(n Note) {
			defer wg.Done()

			err := n.writeNote()
			if err != nil {
				fmt.Printf("The %v note has errored: %v\n", n.Path, err)
			}
		}(note)
	}

	wg.Wait()

}

func ScanAgenda(contents io.Reader, ch chan<- Note) error {
	var currentNote Note
	isCollecting := false

	layout := string(FullDate)
	s := bufio.NewScanner(contents)
	for s.Scan() {
		line := s.Text()
		trimmedLine := strings.TrimSpace(line)

		parsedTime, err := time.Parse(layout, trimmedLine)
		if err == nil {
			if !isCollecting {
				currentNote = Note{Date: parsedTime}
				isCollecting = true
			}
			continue
		}

		if strings.EqualFold(trimmedLine, "END") {
			if isCollecting {
				ch <- currentNote
				isCollecting = false
			}
			continue
		}

		if isCollecting {
			currentNote.Contents = append(currentNote.Contents, []byte(line+"\n")...)
		}
	}

	if err := s.Err(); err != nil {
		return err
	}

	return nil
}

func getNotes(e *[]os.DirEntry, dr DateRange) []Note {
	var noteArray []Note
	AGENDA := viper.GetString("AGENDA")
	for _, v := range *e {
		if !v.IsDir() {
			raw_date := strings.Replace(v.Name(), ".md", "", -1)
			date, err := time.Parse(string(FileDate), raw_date)
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
	sortNotes(noteArray)

	if dr == All || int(dr) > len(noteArray) {
		return noteArray
	}

	startIndex := len(noteArray) - int(dr)
	return noteArray[startIndex:]
}

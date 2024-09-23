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

type separator string

const (
	EOF separator = "#------------------------------"
)

type DateRange int

const (
	Day       DateRange = 1
	Week      DateRange = 7
	Month     DateRange = 30
	Year      DateRange = 365
	All       DateRange = 0
	Yesterday DateRange = 2
	JOINED              = "/tmp/.go_seq_joined.md"
)

type Period struct {
	Range  DateRange
	Amount int
}

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

func JoinNotes(entries *[]fs.DirEntry, period Period) error {
	join := path.Join(JOINED)
	notes := getNotes(entries, period)
	if len(notes) == 0 {
		return fmt.Errorf("No DailyNotes found try to create one with goseq new")
	}
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
		buffer.WriteString("\n\n")
		buffer.WriteString(string(EOF))
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

	content, err := os.ReadFile(join)
	if err != nil {
		return err
	}
	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t")
	}
	strippedContent := strings.Join(lines, "\n")
	return os.WriteFile(join, []byte(strippedContent), 0644)
}

func checkSeparator(line string) bool {
	if len(line) < 4 || line[0] != '#' {
		return false
	}
	var hyphenCount int
	for i := 1; i < len(line); i++ {
		if line[i] == '-' {
			hyphenCount++
		} else {
			break
		}
	}
	return hyphenCount >= 3
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
			if isCollecting {
				ch <- currentNote
			}
			currentNote = Note{Date: parsedTime}
			isCollecting = true
			continue
		}
		if checkSeparator(trimmedLine) {
			if isCollecting {
				currentNote.Contents = bytes.TrimRight(currentNote.Contents, "\n")
				ch <- currentNote
				isCollecting = false
			}
			continue
		}
		if isCollecting {
			currentNote.Contents = append(currentNote.Contents, line...)
			currentNote.Contents = append(currentNote.Contents, '\n')
		}
	}
	if isCollecting {
		// Handle the last if there's no separator
		currentNote.Contents = bytes.TrimRight(currentNote.Contents, "\n")
		ch <- currentNote
	}
	return s.Err()
}

func getNotes(e *[]os.DirEntry, pr Period) []Note {
	var noteArray []Note
	AGENDA := viper.GetString("AGENDA")
	now := time.Now()
	for _, v := range *e {
		if !v.IsDir() {
			raw_date := strings.Replace(v.Name(), ".md", "", -1)
			date, err := time.Parse(string(FileDate), raw_date)

			if err != nil {
				continue
			}
			if !dateInRange(now, pr, date) {
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

	return noteArray
}

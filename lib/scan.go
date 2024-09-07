package lib

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

const AGENDA = "/home/rocky/Documents/Agenda/"
const JOINED = "/home/rocky/Documents/Agenda/.joined.md"

type Note struct {
	Date     string
	Contents []string
	Path     string
}

func (n *Note) writeNote() error {
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

	_, err = writer.WriteString(n.Date + "\n")
	if err != nil {
		return err
	}

	for _, content := range n.Contents {
		_, err = writer.WriteString(content + "\n")
		if err != nil {
			return err
		}
	}

	err = writer.Flush()
	if err != nil {
		return err
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
				fmt.Printf("The %v note has errored: %v\n",n.Path,err)
			}
		}(note)
	}

	wg.Wait()

}

func ScanAgenda(contents io.Reader, ch chan<- Note) error {
	s := bufio.NewScanner(contents)
	var currentNote Note
	isCollecting := false

	for s.Scan() {
		line := s.Text()

		if strings.Contains(line, "START") {
			isCollecting = true
			continue
		}

		if isCollecting {
			if currentNote.Date == "" {
				trimmedLine := strings.TrimSpace(line)
				currentNote.Date = trimmedLine
				err := currentNote.parseDate()
				if err != nil {
					fmt.Printf("Parising failed fo %v due to %v\n", currentNote.Path, err)
					continue
				}
			} else if line != "END" {
				currentNote.Contents = append(currentNote.Contents, line)
			}
		}

		if line == "END" && isCollecting {
			ch <- currentNote
			currentNote = Note{}
			isCollecting = false
		}
	}

	if err := s.Err(); err != nil {
		return err
	}

	return nil
}

func (n *Note) parseDate() error {
	layout := "Mon Jan 2 15:04:05 PM MST 2006"
	parsedTime, err := time.Parse(layout, n.Date)
	if err != nil {
		return err
	}
	n.Path = parsedTime.Format("2006-01-02.md")
	return nil
}


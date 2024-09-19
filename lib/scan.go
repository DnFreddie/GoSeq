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

	"github.com/spf13/viper"
)

const (
AGENDA = "/Documents/Agenda/"
JOINED = "/tmp/.go_seq_joined.md"
PROJECTS = "/Documents/Agenda/projects"
)


type Note struct {
	Date     time.Time
	Contents []byte
	Path     string
}

func (n *Note) writeNote() error {
	home := viper.GetString("HOME")
	if n.Path == "" && !n.Date.IsZero() {

		n.Path = n.Date.Format("2006-01-02.md")

	}

	absPath := path.Join(home,AGENDA, n.Path)
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

func (n *Note) parseDate() string {
	//layout := "Mon Jan 2 15:04:05 PM MST 2006"
	formated_Date := n.Date.Format(string(FullDate))


	return formated_Date
}

func ScanAgenda(contents io.Reader, ch chan<- Note) error {
	var currentNote Note
	isCollecting := false

	layout := string(FullDate)
	s := bufio.NewScanner(contents)
for  s.Scan() {
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


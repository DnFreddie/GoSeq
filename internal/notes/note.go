package notes

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/DnFreddie/goseq/pkg/common"
	"github.com/DnFreddie/goseq/pkg/terminal"

	"github.com/spf13/viper"
)

type DateLayout string

const (
	FileDate DateLayout = "2006-01-02"
	FullDate DateLayout = "January 2 2006"
)

type DNote struct {
	Date     time.Time
	Contents []byte
	Path     string
}

func (d DNote) Format() (string, error) {

	rawDate := path.Base(strings.TrimSuffix(d.Path, ".md"))

	date, err := time.Parse("2006-01-02", rawDate)
	if err != nil {
		return "", err
	}
	return date.Format(string(FullDate)), nil
}

func (d DNote) GetPath() string {
	return d.Path
}
func (d DNote) GetDate() time.Time {
	return time.Time{}
}

func (d DNote) Delete() error {
	formated, err := d.Format()
	if err != nil {
		formated = d.GetPath()
	}
	if err := os.Remove(d.GetPath()); err != nil {
		return err
	}
	fmt.Println("Successfully deleted the Note", formated)
	return nil

}

func (n *DNote) read() error {
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

func (n *DNote) parseDate() string {
	//layout := "Mon Jan 2 15:04:05 PM MST 2006"
	formated_Date := n.Date.Format(string(FullDate))

	return formated_Date
}

func (n DNote) Write() error {
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

func DailyNote() error {
	agenda := checkAgenda()
	now := time.Now()
	date := now.Format(string(FileDate))
	formattedTime := now.Format(string(FullDate))
	var buffer bytes.Buffer
	buffer.Write([]byte(fmt.Sprintf("# %v", formattedTime)))

	dailyNote := path.Join(agenda, date+".md")
	if _, err := os.Stat(dailyNote); errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(dailyNote)
		defer f.Close()
		if err != nil {
			return err
		}

		f.Write(buffer.Bytes())
	}
	err := common.Edit(dailyNote)

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
func ChoseNote(notesArray *[]DNote) error {

	var names []map[string]*DNote
	if len(*notesArray) == 0 {
		return fmt.Errorf("No DailyNotes found try to create one with goseq new")
	}
	for _, entry := range *notesArray {

		dateMap := make(map[string]*DNote)
		fmtTime, err := entry.Format()
		if err != nil {
			fmtTime = entry.GetPath()
		}
		dateMap[fmtTime] = &entry
		names = append(names, dateMap)

	}

	choice, err := terminal.RunTerm(names)
	if err != nil {

		return err
	}

	chosenNote := &DNote{}
	for _, v := range choice {
		*chosenNote = *v
	}

	err = common.Edit(chosenNote.Path)

	if err != nil {

		return err
	}
	return nil
}

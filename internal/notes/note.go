package notes

import (
	"DnFreddie/goseq/pkg/common"
	"DnFreddie/goseq/pkg/terminal"
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

func DailyNote() error {
	agenda := checkAgenda()
	now := time.Now()
	date := now.Format(string(FileDate))
	formattedTime := now.Format(string(FullDate))
	var buffer bytes.Buffer
	buffer.Write([]byte(formattedTime))
	buffer.Write([]byte("\n" + strings.Repeat("-", len(formattedTime))))

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
func ChoseNote() error {

	AGENDA := viper.GetString("AGENDA")
	entries, err := os.ReadDir(AGENDA)

	if err != nil {
		return fmt.Errorf("No DailyNotes found try to create one with goseq new")
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

	choice, err := terminal.RunTerm(names)
	if err != nil {
		return err

	}
	chosenNote := &DNote{}
	for _, v := range choice {
		*chosenNote = DNote{
			Path: path.Join(AGENDA, v.Format(string(FileDate)+".md")),
			Date: v,
		}
	}
	err = common.Edit(chosenNote.Path)
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

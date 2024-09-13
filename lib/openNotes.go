package lib

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

func checkAgenda() string {

	HOME, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	agenda := path.Join(HOME, "Documents/Agenda")

	err = os.MkdirAll(agenda, os.FileMode(0755))

	if err != nil {
		log.Fatal(err)
	}

	return agenda
}

func DailyNote() error {
	agenda := checkAgenda()
	layout := "Mon Jan  2 15:04:05 PM MST 2006"
	now := time.Now()
	date := now.Format("2006-01-02")
	formattedTime := now.Format(layout)

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
	err := edit(dailyNote)
	if err != nil {
		return err
	}

	return nil
}


func parseTimeNote(d string) (time.Time, string, error) {
    layoutISO := "2006-01-02"
    layoutRequested := "January 2 2006"

    var date time.Time
    var err error
    var formattedDate string

    if strings.Contains(d, "-") {
        date, err = time.Parse(layoutISO, d) 
        if err != nil {
            return time.Time{}, "", err
        }
        formattedDate = date.Format(layoutRequested) 
    } else {
        date, err = time.Parse(layoutRequested, d) 
        if err != nil {
            return time.Time{}, "", err
        }
        formattedDate = date.Format(layoutISO) 
    }

    return date, formattedDate, nil
}


func ChoseNote() error {

	entries, err := os.ReadDir(AGENDA)

	if err != nil {
		return err
	}

	var names []string
	for _, entry := range entries {
		if !entry.IsDir() {

			raw_date := strings.Replace(entry.Name(), ".md", "", -1)
			_,fmtDate, err := parseTimeNote(raw_date)
			if err != nil {
				continue
			}

			names = append(names, fmtDate)
		}
	}

	choice, err := RunTerm(names)
	if err != nil {
		return err
	}

	date,fmtDate, err := parseTimeNote(choice)
	if err != nil {
		return err
	}

	chosenNote := Note{
		Path: path.Join(AGENDA, fmtDate+".md"),
		Date: date,
	}
	err = edit(chosenNote.Path)
	if err != nil {
		return nil
	}
	return nil
}

func edit(fPath string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return fmt.Errorf("Failed to foudn $EDITOR")
	}
	cmd := exec.Command(editor, fPath)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		return err
	}

	return nil
}

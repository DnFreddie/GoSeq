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
	err := edit(dailyNote)
	if err != nil {
		return err
	}

	return nil
}




func ChoseNote() error {

	entries, err := os.ReadDir(AGENDA)

	if err != nil {
		return err
	}

	var names []string
	for _, entry := range entries {
		if !entry.IsDir() {
			dateStirng := strings.Replace(entry.Name(), ".md", "", -1)
			fmtDate,err := time.Parse(string(FileDate),dateStirng)

			if err != nil {
				continue
			}

			names = append(names, fmtDate.Format(string(FullDate)))
		}
	}

	choice, err := RunTerm(names)
	if err != nil {
		return err
	}

	date,err := time.Parse(string(FullDate),choice) 
	if err != nil {
		return err
	}

	chosenNote := Note{
		Path: path.Join(AGENDA, date.Format(string(FileDate)+".md")),
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

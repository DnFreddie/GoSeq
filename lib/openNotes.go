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

	"github.com/spf13/viper"
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
	err := Edit(dailyNote)
	if err != nil {
		return err
	}

	return nil
}

func ChoseNote() error {

	home := viper.GetString("HOME")
	entries, err := os.ReadDir(path.Join(home,AGENDA))

	if err != nil {
		return err
	}

	var names []map[string]time.Time
	for _, entry := range entries {
		if !entry.IsDir() {
			dateStirng := strings.Replace(entry.Name(), ".md", "", -1)
			rawDate, err := time.Parse(string(FileDate), dateStirng)

			if err != nil {
				continue
			}

			dateMap := make(map[string]time.Time)
			fmtTime := rawDate.Format(string(FullDate))
			dateMap[fmtTime] = rawDate
			names = append(names, dateMap)

		}
	}

	choice, err := RunTerm(names)
	if err != nil {
		return err

	}
	if err != nil {
		return err
	}
	chosenNote := &Note{}
	for _, v := range choice {
		*chosenNote = Note{
			Path: path.Join(home,AGENDA, v.Format(string(FileDate)+".md")),
			Date: v,
		}
	}
	err = Edit(chosenNote.Path)
	if err != nil {
		return nil
	}
	return nil
}

func Edit(fPath string) error {
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

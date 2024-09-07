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

func CheckAgenda() string {

	HOME, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	agenda := path.Join(HOME, "TEST_AGENDA")

	err = os.MkdirAll(agenda, os.FileMode(0755))

	if err != nil {
		log.Fatal(err)
	}

	return agenda
}

func DailyNote() error {
	agenda := CheckAgenda()
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
		f.Write([]byte( "\n" + strings.Repeat("-", len(formattedTime))))

	}
	err := edit(dailyNote)
	if err != nil {
		return err
	}

	return nil
}

func edit(fPath string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return fmt.Errorf("Failed to foudn $EDITOR")
	}
	cmd := exec.Command("nvim", fPath)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		return err
	}

	fmt.Println("Neovim has exited.")
	return nil
}

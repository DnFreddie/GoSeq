package notes

import (
	"DnFreddie/goseq/lib"
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/sync/semaphore"
)

func SearchNotes(phrase string, flag lib.GrepFlag) error {
	agendaDir := viper.GetString("AGENDA")
	entries, err := os.ReadDir(agendaDir)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	sem := semaphore.NewWeighted(10)

	var matchArray []map[string][]lib.GrepMatch
	for _, entry := range entries {
		if _, err := isNote(entry); err != nil {
			continue
		}

		pattern := []byte(phrase)
		filePath := path.Join(agendaDir, entry.Name())

		wg.Add(1)
		go func(fPath string) {
			defer wg.Done()
			ctx := context.Background()
			if err := sem.Acquire(ctx, 1); err != nil {
				return
			}
			defer sem.Release(1)

			matches, err := lib.GrepFile(fPath, pattern, flag)
			if err != nil {
				return
			}

			m := make(map[string][]lib.GrepMatch)

			m[fPath] = matches

			matchArray = append(matchArray, m)

		}(filePath)
	}

	wg.Wait()

	if len(matchArray) != 0 {
		printNotes(matchArray)
		lib.ProcessUserInput(matchArray)

		return nil
	}

	lib.InColors(lib.Red, "No results found\n")
	return nil
}

func printNotes(notes []map[string][]lib.GrepMatch) {
	for i, note := range notes {
		for key, matches := range note {
			fileName := path.Base(key)
			rawDate := strings.TrimSuffix(fileName, ".md")
			date, err := time.Parse(string(FileDate), rawDate)
			if err != nil {
				lib.InColors(lib.Blue, fileName+"\n")
			} else {
				lib.InColors(lib.Green, fmt.Sprintf("%d. ", i+1))
				lib.InColors(lib.Blue, date.Format(string(FullDate))+"\n")
			}

			for _, match := range matches {
				fmt.Printf("Line:%d %s\n", match.Line, match.Match)
			}
		}
	}
}

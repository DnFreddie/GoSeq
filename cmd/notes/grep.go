package notes

import (
	"DnFreddie/goseq/lib"
	"context"
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/spf13/viper"
	"golang.org/x/sync/semaphore"
)

func SearchNotes(phrase string) error {
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

			matches, err := lib.GrepFile(fPath, pattern, lib.Regex)
			if err != nil {
				return
			}

			m := make(map[string][]lib.GrepMatch)

			m[fPath] = matches

			matchArray = append(matchArray, m)

		}(filePath)
	}

	wg.Wait()
	fmt.Println(matchArray)
	return nil
}

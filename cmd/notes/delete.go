/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package notes

import (
	"DnFreddie/goseq/lib"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deltes daily note from the file",
	Long: `Join notes and deltes the ones that are beeing 
deleted by the user
`,
	Run: func(cmd *cobra.Command, args []string) {
		period := lib.Period{
			Range:  lib.All,
			Amount: 0,
		}
		noteManager := NewDailyNoteManager()
		notes, err := noteManager.GetNotes(period)

		locker :=  lib.NewFileLocker(lib.DeleteLock,"Delete Notes")
		if err != nil {
			if errors.Is(err, lib.NoNotesError{}) {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Println(err)

		}
		if err := locker.Lock(); err!= nil{
			fmt.Println(err)
			os.Exit(1)
		}
		defer locker.Unlock()

		reader, err := noteManager.JoinNotesByTitle(&notes)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := noteManager.DeleteByTitle(reader, &notes); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	},
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
func deleteByTitle(r io.Reader, notes *[]DNote) error {
	var titles []string
	var wasDeleted atomic.Bool
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		titles = append(titles, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	errChan := make(chan error, len(*notes))
	var wg sync.WaitGroup
	for _, note := range *notes {
		formatted, err := note.Format()
		if err != nil {
			formatted = note.GetPath()
		}
		if !slices.Contains(titles, formatted) {
			wg.Add(1)
			go func(n DNote) {
				defer wg.Done()
				if err := n.Delete(); err != nil {
					errChan <- err
				} else {
					wasDeleted.Store(true)
				}
			}(note)
		}
	}
	go func() {
		wg.Wait()
		close(errChan)
	}()
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}
	if len(errors) > 0 {
		return fmt.Errorf("multiple errors occurred: %v", errors)
	}
	if !wasDeleted.Load() {
		lib.InColors(lib.Red, "Nothing to delete ...\n")
	}
	return nil
}

func joinByTitle(notes *[]DNote) (io.Reader, error) {
	f, err := os.OpenFile(JOINED_DELETE, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	var titles []string
	for _, note := range *notes {
		formattedName, err := note.Format()
		if err != nil {
			formattedName = note.GetPath()
		}
		titles = append(titles, formattedName)
	}
	joinedTitles := strings.Join(titles, "\n")

	if _, err := f.Write([]byte(joinedTitles)); err != nil {
		return nil, err
	}

	if err := lib.Edit(JOINED_DELETE); err != nil {
		return nil, err
	}

	updatedContent, err := os.ReadFile(JOINED_DELETE)
	if err != nil {
		return nil, fmt.Errorf("failed to read updated file: %w", err)
	}

	return bytes.NewReader(updatedContent), nil
}

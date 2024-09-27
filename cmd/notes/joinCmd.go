/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package notes

import (
	"DnFreddie/goseq/lib"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"time"
	"unicode"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// joinCmd represents the join command

const (
	JOINED        = "/tmp/.go_seq_notes_joined.md"
	JOINED_DELETE = "/tmp/.go_seq_notes_delete_joined.md"
)

var periodVarCmd string
var dateRangeVar int
var JoinCmd = &cobra.Command{
	Use:   "join",
	Short: "Joins notes in one ",
	Long:  `Join notes any changes to the notes will be applaied to the notes (by defult from one week last 7 notes) `,
	Run: func(cmd *cobra.Command, args []string) {

		period := lib.Period{
			Range:  lib.ParseDateRange(periodVarCmd),
			Amount: dateRangeVar,
		}

		noteManager := NewDailyNoteManager()
		notes, err := noteManager.GetNotes(period)
		reader, err := noteManager.JoinNotesWithContents(&notes)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		scanner := NewDNoteScanner(reader)
		lib.ScanJoined(scanner)
	},
}

func init() {
	// Here you will define your flags and configuration settings.
	//JoinCmd.Flags().StringVar(&periodVarCmd, "range", "r", "week", "Date range (day, week, month, year, all)")
	JoinCmd.Flags().StringVarP(&periodVarCmd, "range", "r", "week", "Specify a time unit (week, year, day)Default 1 week ")
	JoinCmd.Flags().IntVarP(&dateRangeVar, "times", "t", 1, "Specify home many times ago(3 weeeks ago)")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// joinCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// joinCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type separator string

const (
	EOF separator = "#------------------------------"
)

type DateRange int

func joinNotes(notes *[]DNote) (io.Reader, error) {
	if len(*notes) == 0 {
		return nil, fmt.Errorf("no DailyNotes found; try to create one with goseq new")
	}

	f, err := os.OpenFile(JOINED, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	for _, v := range *notes {
		if err := v.read(); err != nil {
			log.Printf("Error reading note: %v", err)
			continue
		}

		var buffer bytes.Buffer
		buffer.Write(v.Contents)
		buffer.WriteString("\n\n")
		buffer.WriteString(string(EOF))
		buffer.WriteString("\n\n")

		if _, err := f.Write(buffer.Bytes()); err != nil {
			log.Printf("Error writing to file: %v", err)
		}

		v.Contents = nil
	}

	if err := lib.Edit(JOINED); err != nil {
		return nil, fmt.Errorf("error editing file: %w", err)
	}

	readFile, err := os.Open(JOINED)
	if err != nil {
		return nil, fmt.Errorf("error opening edited file: %w", err)
	}

	reader := bufio.NewReader(readFile)

	return &trimReader{reader: reader, file: readFile}, nil
}

type trimReader struct {
	reader *bufio.Reader
	file   *os.File
}

func (tr *trimReader) Close() error {
	return tr.file.Close()
}

func (tr *trimReader) Read(p []byte) (n int, err error) {
	line, err := tr.reader.ReadBytes('\n')
	if err != nil && err != io.EOF {
		return 0, err
	}

	// Trim trailing whitespace
	line = bytes.TrimRightFunc(line, unicode.IsSpace)
	if len(line) > 0 || err == nil {
		line = append(line, '\n')
	}

	n = copy(p, line)
	if n < len(line) {
		err = nil
	}
	return n, err
}

func getNotes(pr lib.Period) ([]DNote, error) {
	var errMessages []string
	noteArray := []DNote{}

	AGENDA := viper.GetString("AGENDA")

	entries, err := os.ReadDir(AGENDA)
	if err != nil {

		return noteArray, &lib.NoNotesError{}
	}

	now := time.Now()

	for _, v := range entries {
		if !v.IsDir() {
			rawDate := strings.Replace(v.Name(), ".md", "", -1)
			date, err := time.Parse(string(FileDate), rawDate)
			if err != nil {
				errMessages = append(errMessages, fmt.Sprintf("failed to parse file %s: %v", v.Name(), err))
				continue
			}

			if !lib.DateInRange(now, pr, date) {
				continue
			}

			note := DNote{
				Path: path.Join(AGENDA, v.Name()),
				Date: date,
			}
			noteArray = append(noteArray, note)
		}
	}

	lib.SortNotes(noteArray)

	if len(errMessages) > 0 {
		return noteArray, fmt.Errorf(strings.Join(errMessages, "; "))
	}

	return noteArray, nil
}


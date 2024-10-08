/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package notes

import (
	"fmt"

	"github.com/DnFreddie/goseq/internal/notes"
	"github.com/DnFreddie/goseq/pkg/common"
	"github.com/DnFreddie/goseq/pkg/locker"

	"github.com/spf13/cobra"
)

// joinCmd represents the join command

const (
	JoinLock locker.LockFile = "/tmp/.goseq_join.lock"
)

var periodVarCmd string
var dateRangeVar int
var JoinCmd = &cobra.Command{
	Use:   "join",
	Short: "Joins notes in one ",
	Long:  `Join notes any changes to the notes will be applaied to the notes (by defult from one week last 7 notes) `,
	Run: func(cmd *cobra.Command, args []string) {

		period := common.Period{
			Range:  common.ParseDateRange(periodVarCmd),
			Amount: dateRangeVar,
		}

		locker := locker.NewFileLocker(JoinLock, "Join Notes")
		noteManager := notes.NewDailyNoteManager()
		notesArray, err := noteManager.GetNotes(period)
		if err := locker.Lock(); err != nil {
			fmt.Println(err)

			return
		}
		defer locker.Unlock()

		reader, err := noteManager.JoinNotesWithContents(&notesArray)
		if err != nil {
			fmt.Println(err)
			return
		}

		scanner := notes.NewDNoteScanner(reader)
		common.ScanJoined(scanner)
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

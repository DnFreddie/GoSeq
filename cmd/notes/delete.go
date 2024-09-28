/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package notes

import (
	"DnFreddie/goseq/internal/notes"
	"DnFreddie/goseq/pkg/common"
	"DnFreddie/goseq/pkg/locker"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const DeleteLock locker.LockFile = "/tmp/.goseq_delete.lock"

// deleteCmd represents the delete command
var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deltes daily note from the file",
	Long: `Join notes and deltes the ones that are beeing 
deleted by the user
`,
	Run: func(cmd *cobra.Command, args []string) {
		period := common.Period{
			Range:  common.All,
			Amount: 0,
		}
		noteManager := notes.NewDailyNoteManager()
		notes, err := noteManager.GetNotes(period)

		locker := locker.NewFileLocker(DeleteLock, "Delete Notes")
		if err != nil {
			if errors.Is(err, common.NoNotesError{}) {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Println(err)

		}
		if err := locker.Lock(); err != nil {
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

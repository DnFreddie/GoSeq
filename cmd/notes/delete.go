/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package notes

import (
	"errors"
	"fmt"
	"time"

	"github.com/DnFreddie/goseq/internal/common"
	"github.com/DnFreddie/goseq/internal/dnotes"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deltes daily note from the file",
	Long: `Join dnotes and deltes the ones that are beeing 
deleted by the user
`,
	Run: func(cmd *cobra.Command, args []string) {
		period := common.Period{
			Range:  common.All,
			Amount: 0,
			Today:  time.Now(),
		}
		noteManager := dnotes.NewDailyNoteManager()
		dnotes, err := noteManager.GetNotes(period)

		if err != nil {

			if errors.Is(err, common.NoNotesFoundErr{}) {
				fmt.Println(err)
				return
			}

			fmt.Println(err)

		}

		reader, err := noteManager.JoinNotesByTitle(&dnotes)
		if err != nil {
			fmt.Println(err)
			return
		}

		if err := noteManager.DeleteByTitle(reader, &dnotes); err != nil {
			fmt.Println(err)
			return
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

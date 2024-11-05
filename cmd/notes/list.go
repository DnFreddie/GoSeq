/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package notes

import (
	"fmt"
	"os"
	"time"

	"github.com/DnFreddie/goseq/internal/common"
	"github.com/DnFreddie/goseq/internal/dnotes"

	"github.com/spf13/cobra"
)

// ListCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List notes and pick the one you want",
	Long: `
	List all notes in agenda and lets you chose the one u want then it opens it and applay the changes 
`,
	Run: func(cmd *cobra.Command, args []string) {
		now := time.Now()
		notesArray, err := dnotes.NewDailyNoteManager().GetNotes(common.Period{Range: common.All, Amount: 0, Today: now})

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if err = dnotes.ChoseNote(&notesArray); err != nil {

			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

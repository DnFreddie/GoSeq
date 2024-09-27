/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package git

import (
	"DnFreddie/goseq/lib"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a projects from the list",
	Long: `Get a list of the project in a file 
Chagnes to the file will delete the associated projects.
 `,
	Run: func(cmd *cobra.Command, args []string) {

		period := lib.Period{
			Range:  lib.All,
			Amount: 0,
		}
		projectManager := NewProjectManager()
		projects, err := projectManager.GetNotes(period)
		if err != nil {
			if errors.Is(err, lib.NoNotesError{}) {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Println(err)

		}
		reader, err := projectManager.JoinNotesByTitle(&projects)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := projectManager.DeleteByTitle(reader, &projects); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	},
}

func init() {
	GitCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

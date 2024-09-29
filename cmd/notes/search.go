/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package notes

import (
	"errors"
	"fmt"
	"github.com/DnFreddie/goseq/internal/notes"
	"github.com/DnFreddie/goseq/pkg/common"
	"github.com/DnFreddie/goseq/pkg/grep"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var iname bool
var regex bool
var SearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for a pattern in the notes and open the ones that are relevant.",
	Long: `Accept the pattern along with the grep flags -i or -E to search for matches.
It will then display the matches and allow you to open the desired note.`,
	Run: func(cmd *cobra.Command, args []string) {

		var re grep.GrepFlag
		var insencitive grep.GrepFlag

		userInput := strings.Join(args, " ")

		if userInput == "" {
			fmt.Println("No pattern to look for ")
			os.Exit(1)
		}

		period := common.Period{
			Range:  common.All,
			Amount: 0,
		}
		if iname {
			insencitive = grep.ToLower

		}
		if regex {
			re = grep.Regex
		}

		noteManager := notes.NewDailyNoteManager()
		notes, err := noteManager.GetNotes(period)
		if err != nil {
			if !errors.Is(err, common.NoNotesError{}) {
				fmt.Println(err)
				os.Exit(1)
			} else {
				fmt.Println(err)
			}
		}
		if err := common.Search(notes, userInput, re|insencitive); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {

	// Here you will define your flags and configuration settings.

	SearchCmd.Flags().BoolVarP(&iname, "iname", "i", false, "Case Insensitive Search")
	SearchCmd.Flags().BoolVarP(&regex, "regex", "E", false, "Accepts Posix Regex Search")
	//(&iname, "iname", "i", , "Description of the iname flag")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

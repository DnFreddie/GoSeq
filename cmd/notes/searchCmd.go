/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package notes

import (
	"DnFreddie/goseq/lib"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

		var re lib.GrepFlag
		var insencitive lib.GrepFlag
		if iname {
			insencitive = lib.ToLower

		}
		if regex {
			re = lib.Regex
		}

		if err := SearchT(strings.Join(args, " "), re|insencitive); err != nil {
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

func SearchT(toParse string, flag lib.GrepFlag) error {
	agendaDir := viper.GetString("AGENDA")
	entries, err := os.ReadDir(agendaDir)
	if err != nil {
		return err
	}

	period := Period{
		Range:  All,
		Amount: 0,
	}
	notes := getNotes(&entries, period)

	if len(notes) == 0 {
		return fmt.Errorf("No notes found\n")
	}

	locations := make([]string, len(notes))

	for i, n := range notes {
		locations[i] = n.Path
	}


	matches, err := lib.GrepMulti(locations, toParse, flag)
	if err != nil {
		return err
	}
	lib.OpenNotes(&matches, searchNoteForrmater)


	return nil
}
func searchNoteForrmater(s string) (string, error) {
	rawDate := strings.TrimSuffix(s, ".md")
	date, err := time.Parse(string(FileDate), rawDate)
	if err != nil {
		return "", err
	}
	return date.Format(string(FullDate)), nil
}

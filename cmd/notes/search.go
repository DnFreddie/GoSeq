/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package notes

import (
	"DnFreddie/goseq/lib"
	"log"
	"strings"

	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var iname bool
var regex bool
var SearchCmd = &cobra.Command{
	Use:   "search",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		var re lib.GrepFlag
		var insencitive lib.GrepFlag
		if iname {
			insencitive = lib.ToLower

		}
		if regex {
			re = lib.Regex
		}

		if err := SearchNotes(strings.Join(args, " "), re|insencitive); err != nil {
			log.Fatal(err)
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

/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package git

import (
	"DnFreddie/GoSeq/lib/github"
	"log"

	"github.com/spf13/cobra"
)

// testCmd represents the test command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List added Projects and chose one to edit",
	Long: `It lists the projects that was previously added and lets you chose one.
The paths are located in $HOME/Documents/Agenda/projects/.PROJECTS_META.json `,
	Run: func(cmd *cobra.Command, args []string) {
		if err := github.ReadRecent(true);err!=nil{
			log.Fatal(err)
		}
	},
}

func init() {
	GitCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

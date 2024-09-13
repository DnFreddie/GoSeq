/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"DnFreddie/GoSeq/lib"
	"log"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List notes and pick the one you want",
	Long: `
	List all notes in agenda and lets you chose the one u want then it opens it and applay the changes 
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := lib.ChoseNote()
		if err != nil{

		log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

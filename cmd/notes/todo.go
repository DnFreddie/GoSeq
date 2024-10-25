/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package notes

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/DnFreddie/goseq/internal/dnotes"
	"github.com/DnFreddie/goseq/internal/common"
	"github.com/DnFreddie/goseq/pkg/todo"
	"github.com/spf13/cobra"
)

// todoCmd represents the todo command
var TodoCmd = &cobra.Command{
	Use:   "todo",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		n, err := notes.NewDailyNoteManager().GetNotes(common.Period{
			Range:  common.All,
			Amount: 0,
		})
		if err != nil {
			fmt.Println("Error retrieving notes:", err)
			return
		}

		for _, v := range n {
			err := v.Read()
			if err != nil {
				fmt.Println("Error reading note:", err)
				continue
			}

			scanner := bufio.NewScanner(strings.NewReader(string(v.Contents)))
			lineNumber := 0
			for scanner.Scan() {
				line := scanner.Text()
				lineNumber++
				if t := todo.ContainsPattern(line, lineNumber, todo.ALL); t != nil {
					t.PrettyPrintTodo()
				}

			}

			if err := scanner.Err(); err != nil {
				fmt.Println("Error scanning note contents:", err)
			}
		}
	},
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// todoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// todoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

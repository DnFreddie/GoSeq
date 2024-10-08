/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package notes

import (
	"github.com/DnFreddie/goseq/internal/notes"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// newCmd represents the new command
var NewCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new daily note ",
	Long: `Create a new daily note or open an exsisitng one for today
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := notes.DailyNote()
		if err != nil{

			fmt.Println(err)
			os.Exit(1)
		}

	},
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

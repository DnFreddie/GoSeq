/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"DnFreddie/goseq/cmd/git"
	"DnFreddie/goseq/cmd/notes"
	"DnFreddie/goseq/config"
	"log"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "goseq",
	Short: "Goseq diary for the devleoper and their projectes all in one binary",
	Long: `Goseq provides a way to connect your daily notes with your project notes.
It also allows you to seamlessly upload GitHub issues written in the code.
All in one binary.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// err := lib.DailyNote()

	// if err != nil {
	// fmt.Println(err)
	// 	os.Exit(1)
	// }
	//lib.RunTerm()

	///	errror  := lib.ChoseNote()
	err := RootCmd.Execute()

	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.goseq.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	cobra.OnInitialize(config.LoadConfig)
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}
	viper.Set("HOME", homeDir)
	agenda := path.Join(homeDir,"/Documents/Agenda/")
	viper.Set("AGENDA",agenda)
	viper.Set("PROJECTS",agenda+"/projects")
	addSubcommandsPallet()
}

func addSubcommandsPallet() {
	RootCmd.AddCommand(git.GitCmd)
	RootCmd.AddCommand(notes.JoinCmd)
	RootCmd.AddCommand(notes.NewCmd)
	RootCmd.AddCommand(notes.SearchCmd)
	RootCmd.AddCommand(notes.ListCmd)
}

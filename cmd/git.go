/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"DnFreddie/GoSeq/lib/github"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// gitCmd represents the git command
var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		//github.WalkProject("/home/rocky/github.com/DnFreddie/rlbl")
		
		err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

		//val := os.Getenv("GITHUB_TOKEN")
		p,err  := github.ProjectInit("/home/rocky/github.com/DnFreddie/rlbl/")


		if err != nil{

			log.Fatal(err)
		}
		err = p.WalkProject()
		if err != nil{
			log.Fatal(err)
		}
		fmt.Println(p.Name)
		p.Read()
		//fmt.Println(p.FetchGitHubIssues(val))
	},
}

func init() {
	rootCmd.AddCommand(gitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

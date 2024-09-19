/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"DnFreddie/GoSeq/lib"
	"DnFreddie/GoSeq/lib/github"
	"io"
	"log"
	"log/slog"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var projectPath string

// gitCmd represents the git command
var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "Open a note for a specyfied repo",
	Long:  `Opens a note for the project if paht not specyfied it finds the recent one`,
	Run: func(cmd *cobra.Command, args []string) {

		//github.WalkProject("/home/rocky/github.com/DnFreddie/rlbl")

		home := viper.GetString("HOME")
		if projectPath != "" {
			pAth := github.PickProject(projectPath)

			location := path.Join(home, lib.PROJECTS, pAth)
			err := os.WriteFile(lib.ENV_VAR, []byte(location), 0644)
			if err != nil {
				slog.Warn("Failed to read to the TMP file", "err", err)
			}

		} else {
			log.Println("No project path provided.")
			f, err := os.Open(lib.ENV_VAR)
			defer f.Close()
			if err != nil {
				log.Fatalf("No recent Projects found: %v", err)
			}
			p, err := io.ReadAll(f)
			if err != nil {
				log.Fatalf("No recent Projects found: %v", err)
			}
			lib.Edit(string(p) + ".md")

		}

		//fmt.Println(p.FetchGitHubIssues(val))
	},
}

func init() {
	rootCmd.AddCommand(gitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:

	gitCmd.PersistentFlags().StringVar(&projectPath, "path", "", "A path to your project/dir where you store them")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

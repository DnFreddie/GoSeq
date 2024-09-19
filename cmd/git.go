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
)

var projectPath string

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

		if projectPath != "" {
			pAth := github.PickProject(projectPath)

			location := path.Join(lib.AGENDA, "projects", pAth)
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
			lib.Edit(string(p))

		}

		//fmt.Println(p.FetchGitHubIssues(val))
	},
}

func init() {
	rootCmd.AddCommand(gitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	gitCmd.Flags().StringVarP(&projectPath, "path", "p","", "A path to your project/dir where u store them")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

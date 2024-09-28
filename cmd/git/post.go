/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package git

import (
	"github.com/DnFreddie/goseq/internal/project"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// postCmd represents the post command
var PostCmd = &cobra.Command{
	Use:   "post",
	Short: "Post new issues to the github",
	Long:  `It scans the given direcories for githubissues  and post the ones that don't exist to the github`,
	Run: func(cmd *cobra.Command, args []string) {

		tokenValue := viper.Get("token")
		token, ok := tokenValue.(string)

		if !ok || token == "" {
			log.Fatal("Failed to found the Github Api token")
		}
		if ProjectPath != "" {
			p, err := project.NewProject(ProjectPath)
			if err != nil {
				log.Fatal(err)
			}
			if err := p.WalkProject(); err != nil {
				log.Fatal(err)
			}

			gitIssues, err := p.FetchGitHubIssues(token)
			if err != nil {
				log.Fatal(err)

			}
			p.ApplayIssues(token, &gitIssues)

		}

	},
}

func init() {
	GitCmd.AddCommand(PostCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// postCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// postCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

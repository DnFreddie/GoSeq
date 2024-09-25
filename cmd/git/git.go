/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package git

import (
	"fmt"
	"log/slog"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ProjectPath string
var recent bool

// gitCmd represents the git command
var GitCmd = &cobra.Command{
	Use:   "git",
	Short: "Open a note for a specyfied repo",
	Long: `Opens a note for the project if paths  not specyfied it finds the recent one
And the for the saved.You can see what Projects did u saved via the git list command`,
	Run: func(cmd *cobra.Command, args []string) {
		//github.WalkProject("/home/rocky/github.com/DnFreddie/rlbl")
		PROJECTS := viper.GetString("PROJECTS")
		if ProjectPath != "" {
			p := PickProject(ProjectPath)
			pAth := path.Join(p.Owner, p.Name)

			location := path.Join(PROJECTS, pAth)
			err := os.WriteFile(ENV_VAR, []byte(location), 0644)
			if err != nil {
				slog.Warn("Failed to read to the TMP file", "err", err)
			}
			return
		}
		if recent  {
			if err := ReadRecent(false);err !=nil{
				fmt.Println(err)
				os.Exit(1)
			}

			
		}
		if recent == false&& ProjectPath == "" {
			cmd.Help()
		}

	},
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:

	GitCmd.PersistentFlags().StringVarP(&ProjectPath, "path", "p", "", "A path to your project/dir where you store them")
	GitCmd.Flags().BoolVarP(&recent, "recent", "r", false, "Open Recent Project")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

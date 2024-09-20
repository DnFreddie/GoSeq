/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package git

import (
	"DnFreddie/goseq/lib/github"
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/spf13/cobra"
	"golang.org/x/sync/semaphore"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scans todos from the given path",
	Long: `Findes the project and prints all the todos founded inside the project/s `,
	Run: func(cmd *cobra.Command, args []string) {
		if projectPath != "" {
			prArray, err := github.ListRepos(projectPath)
			if err != nil {
				log.Fatal(err)
			}

			var wg sync.WaitGroup
			ctx := context.Background()
			sem := semaphore.NewWeighted(int64(10))

			for _, v := range prArray {
				wg.Add(1)
				go func(repo github.Project) {
					sem.Acquire(ctx, 1)
					defer wg.Done();sem.Release(1)
					if err := repo.WalkProject(); err != nil {
						fmt.Printf("Failed to scan %v, %v\n", repo.Name, err)
					}
					repo.PrintTodos()
				}(v)
			}

			wg.Wait()
		} else {
			cmd.Help()
		}
	},
}

func init() {
	GitCmd.AddCommand(scanCmd)

	// Here you will define your flags and configuration settings.
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}


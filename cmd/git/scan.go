/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package git

import (
	"github.com/DnFreddie/goseq/internal/project"
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/spf13/cobra"
	"golang.org/x/sync/semaphore"
)
var add bool
var ScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scans todos from the given path",
	Long:  `Findes the project and prints all the todos founded inside the project/s `,
	Run: func(cmd *cobra.Command, args []string) {
		if ProjectPath != "" {
			prArray, err := project.ListProjects(ProjectPath)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			var wg sync.WaitGroup
			ctx := context.Background()
			sem := semaphore.NewWeighted(int64(10))

			for _, v := range prArray {
				wg.Add(1)
				go func(repo project.Project) {
					sem.Acquire(ctx, 1)
					defer wg.Done()
					sem.Release(1)
					if err := repo.WalkProject(); err != nil {
						fmt.Printf("Failed to scan %v, %v\n", repo.Name, err)
					}
					repo.PrintTodos()
				}(v)
			}

			wg.Wait()
			if len(prArray)!= 0 && add {
				for _,p:= range prArray{
					p.SaveProject()
				}
			}

			
		} else {
			cmd.Help()
		}
	},
}

func init() {
	GitCmd.AddCommand(ScanCmd)

	ScanCmd.Flags().BoolVarP(&add, "add", "a", false, "Add project to known projects")
	// Here you will define your flags and configuration settings.
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}

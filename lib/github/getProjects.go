package github

import (
	"DnFreddie/GoSeq/lib"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/sync/semaphore"
)

func FindRepos(pt string) {
	var wg sync.WaitGroup
	ctx := context.Background()

	var projects []Project
	repoChan := make(chan *Project, 100)

	var sem = semaphore.NewWeighted(int64(20))
	go func() {
		for repo := range repoChan {
			projects = append(projects,  *repo)
			//fmt.Println("Absolute path of Git repository:", repo)
			
		}
	}()

	err := filepath.WalkDir(pt, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() && d.Name() == ".git" {
			repoPath := filepath.Dir(path)
			absPath, err := filepath.Abs(repoPath)
			if err != nil {
				return err
			}

			sem.Acquire(ctx, 1)
			wg.Add(1)
			go func(repoPath string) {
				defer wg.Done()

				defer sem.Release(1)
				p, err := ProjectInit(repoPath)
				if err != nil {
					return
				}
				repoChan <- p
			}(absPath)
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error walking the directory:", err)
	}

	wg.Wait()

	var options []map[string]*Project
	
for _, v := range projects {
    newOption := make(map[string]*Project)
    newOption[v.Name] = &v
    options = append(options, newOption)
}
	choice,err  := lib.RunTerm(options)
	if err != nil{
		log.Fatal(err)
	}
	fmt.Println(choice)
	

}

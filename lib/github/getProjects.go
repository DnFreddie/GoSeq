package github

import (
	"DnFreddie/GoSeq/lib"
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"

	"golang.org/x/sync/semaphore"
)

func ListRepos(pt string) ([]Project, error) {
	var wg sync.WaitGroup
	ctx := context.Background()

	var projects []Project
	repoChan := make(chan *Project, 100)

	var sem = semaphore.NewWeighted(int64(20))
	go func() {
		for repo := range repoChan {
			projects = append(projects, *repo)

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

	return projects, nil

}

func ChoseProject(pr *[]Project) *Project {
	var options []map[string]*Project
	var p  *Project
	for _, v := range *pr {
		newOption := make(map[string]*Project)
		newOption[fmt.Sprintf("%v/%v",v.Owner,v.Name)] = &v
		options = append(options, newOption)
	}
	choice, err := lib.RunTerm(options)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range choice{
		p = v


	}
	return p
}

func (p *Project) EditProject() {
	pDir := path.Join(lib.AGENDA, "projects",p.Owner)
	err := os.MkdirAll(pDir, 0755)

	if err != nil {
		log.Fatal(err)
	}
	project := path.Join(pDir, p.Name)
	lib.Edit(project)

}

func PickProject(pPath string) string {

	pr, err := ListRepos(pPath)
	for _, v := range pr {
		fmt.Println("------------------------------")
		fmt.Println(v.Location)
		fmt.Println(v.DefaultBranch)
		fmt.Println(v.Name)
		fmt.Println(v.Owner)
		fmt.Println(v.Url)
		fmt.Println("------------------------------")
		
	}
	if err != nil {

		log.Fatal(err)
	}

	gp := ChoseProject(&pr)
	gp.EditProject()
	joinded := path.Join(gp.Owner,gp.Name)
	return  joinded
}

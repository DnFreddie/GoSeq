package github

import (
	"DnFreddie/GoSeq/lib"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/spf13/viper"
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
	var p *Project
	for _, v := range *pr {
		newOption := make(map[string]*Project)
		newOption[fmt.Sprintf("%v/%v", v.Owner, v.Name)] = &v
		options = append(options, newOption)
	}
	choice, err := lib.RunTerm(options)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range choice {
		p = v

	}
	return p
}
func (p *Project) printProperites() string {

	properties := fmt.Sprintf(`------------------------------
Repo: %v/%v
Branch: %v
Url: %v
------------------------------`, p.Owner, p.Name, p.DefaultBranch, p.Url)
	return properties
}

func (p *Project) EditProject() {
	home :=viper.GetString("HOME")
	pDir := path.Join(home, lib.PROJECTS, p.Owner)
	err := os.MkdirAll(pDir, 0755)

	if err != nil {
		log.Fatal(err)
	}
	project := path.Join(pDir, p.Name+".md")

	if _, err := os.Stat(project); errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(project)
		if err != nil {
			log.Fatal(err)
		}

		props := p.printProperites()
		if _, err = f.Write([]byte(props)); err != nil {
			log.Fatal(err)
		}

	}

	lib.Edit(project)

}

func PickProject(pPath string) string {

	prArray, err := ListRepos(pPath)
	var p *Project
	if err != nil {

		log.Fatal(err)
	}

	if len(prArray) == 1 {
		p = &prArray[0]

	} else {
		p = ChoseProject(&prArray)
	}
	p.EditProject()
	joinded := path.Join(p.Owner, p.Name)
	return joinded
}

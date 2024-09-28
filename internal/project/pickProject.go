package project

import (
	"DnFreddie/goseq/pkg/common"
	"DnFreddie/goseq/pkg/terminal"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/sync/semaphore"
)

const (
	ENV_VAR       = "/tmp/GO_SEQ_PROJECT.txt"
	PROJECTS_META = ".PROJECTS_META.json"
)

func PickProject(pPath string) *Project {

	prArray, err := ListProjects(pPath)
	var p *Project
	if err != nil {

		log.Fatal(err)
	}

	if len(prArray) == 1 {
		p = &prArray[0]

	} else {
		p = choseProject(&prArray)
	}
	p.EditProject()
	return p
}

func choseProject(pr *[]Project) *Project {
	var options []map[string]*Project
	var p *Project
	for _, v := range *pr {
		newOption := make(map[string]*Project)
		newOption[fmt.Sprintf("%v/%v", v.Owner, v.Name)] = &v
		options = append(options, newOption)
	}
	choice, err := terminal.RunTerm(options)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range choice {
		p = v

	}
	return p
}

func ListProjects(pt string) ([]Project, error) {
	var wg sync.WaitGroup
	ctx := context.Background()

	var projects []Project
	repoChan := make(chan *Project, 100)

	var sem = semaphore.NewWeighted(int64(30))
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
				p, err := NewProject(repoPath)
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
	if len(projects) == 0 {
		return projects, fmt.Errorf("No Projects Found")
	}

	return projects, nil

}
func ReadRecent(list bool) error {
    if !list {
        p, err := readRecentProject()
        if err != nil {
            fmt.Println("No recent projects found, listing added projects instead")
            return ReadRecent(true)
        }
        return common.Edit(string(p) + ".md")
    }
    
    projects, err := getSavedProjects()
    if err != nil {
        return err
    }
    return choseProject(&projects).EditProject()
}

func readRecentProject() ([]byte, error) {
    if _, err := os.Stat(ENV_VAR); os.IsNotExist(err) {
        return nil, err
    }
    
    f, err := os.Open(ENV_VAR)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    
    return io.ReadAll(f)
}


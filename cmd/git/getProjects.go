package git

import (
	"DnFreddie/goseq/lib"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/spf13/viper"
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
	choice, err := lib.RunTerm(options)
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
	if len(projects)==0{
		return projects,fmt.Errorf("No Projects Found")
	}

	return projects, nil

}


func getSavedProjects()([]Project,error){
	PROJECTS := viper.GetString("PROJECTS")
	var projecArray []Project
	f, err := os.Open(path.Join(PROJECTS, PROJECTS_META))

	if err != nil {
		return projecArray, fmt.Errorf("The meta file is empty add the project to fix this\n\ngit -p <path/to/project/\n")
	}

	contents, err := io.ReadAll(f)
	if err != nil {
		return projecArray,err
	}
	err = json.Unmarshal(contents, &projecArray)
	if err != nil {
		return projecArray,err
	}
	if len(projecArray) == 0 {
		return projecArray, lib.NoNotesError{}
	}
	return projecArray,nil
}


func ReadRecent(list bool) error {

	if !list {
		f, err := os.Open(ENV_VAR)
		if err != nil {
			fmt.Println("No recent Projects found, listing added projects instead")
			return ReadRecent(true)
		}
		defer f.Close()

		p, err := io.ReadAll(f)
		if err != nil {
			log.Println("Failed to read recent project, listing recent projects instead")
			return ReadRecent(true)
		}

		lib.Edit(string(p) + ".md")
		return nil
	}

	projecArray,err:= getSavedProjects()
	if err !=  nil{
		return err
	}

	pr := choseProject(&projecArray)
	pr.EditProject()
	return nil

}

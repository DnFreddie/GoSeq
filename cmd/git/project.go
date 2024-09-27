package git

import (
	"DnFreddie/goseq/lib"
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/sync/semaphore"
)


type Project struct {
	Name          string `json:"name"`
	Owner         string `json:"owner"`
	DefaultBranch string `json:"default_branch"`
	Url           string `json:"repo_url"`
	Issues        []map[string][]lib.Todo
	Location      string `json:"location"`
	NotePath      string `json:"note_path"`
}

func (p Project) GetPath() string {
	return p.NotePath
}

func (p Project) Write() error {

	return nil
}

func (p Project) Format() (string, error) {
	return path.Join(p.Owner, p.Name), nil
}
func (p Project) GetDate() time.Time {
	time := time.Now()
	return time
}

func (p Project)Delete(){

}
type ProjectRetriver struct{}

func NewDRetriver() *ProjectRetriver {
	retriver := ProjectRetriver{}
	return &retriver
}

func (d *ProjectRetriver) GetNotes(p lib.Period) ([]Project, error) {
	notes, err := getSavedProjects()
	fmt.Println(notes)
	if err != nil {
		return notes, err
	}
	return notes, nil
}

func (d *ProjectRetriver)JoinNotes(p lib.Period)error{

	return nil

}

func (d *ProjectRetriver)Delete(p lib.Period)error{

	return nil

}


func ProjectInit(localPath string) (*Project, error) {
	absoluteP, err := makeAbsolute(localPath)
	if err != nil {
		slog.Error("Doesn't exist:", "path", path.Base(localPath))
		return nil, err
	}

	HEAD := filepath.Join(absoluteP, ".git/HEAD")
	CONFIG := filepath.Join(absoluteP, ".git/config")
	reBranch := regexp.MustCompile(`refs/heads/(\w+)`)
	reUrl := regexp.MustCompile(`url = (.+\.git)$`)

	var defaultBranch string
	var gitURL string
	var repoName string
	var rOwner string

	headFile, err := os.Open(HEAD)
	if err != nil {
		slog.Warn("Failed to open HEAD file:", "error", err)
	} else {
		defer headFile.Close()
		reader := bufio.NewReader(headFile)
		defaultBranch, err = extractMatch(reader, reBranch)
		if err != nil {
			slog.Warn("Failed to extract default branch:", "error", err)
		}
	}

	configFile, err := os.Open(CONFIG)
	if err != nil {
		slog.Warn("Failed to open config file:", "error", err)
	} else {
		defer configFile.Close()
		configReader := bufio.NewReader(configFile)
		gitURL, err = extractMatch(configReader, reUrl)
		if err != nil {
			slog.Warn("Failed to extract URL:", "error", err)
		}
	}

	if gitURL != "" {
		parts := strings.Split(gitURL, "/")
		if len(parts) > 0 {
			rOwner = parts[len(parts)-2]
			repoName = strings.TrimSuffix(parts[len(parts)-1], ".git")
		}
	}

	if repoName == "" {
		return nil, fmt.Errorf("failed to determine repository name from URL: %s", gitURL)
	}

	return &Project{
		Name:          repoName,
		Location:      absoluteP,
		DefaultBranch: defaultBranch,
		Url:           gitURL,
		Owner:         rOwner,
	}, nil
}

func (p *Project) PrintTodos() {
	message := fmt.Sprintf("%v/%v\n", p.Owner, p.Name)
	lib.InColors(lib.Cyan, message)
	if len(p.Issues) == 0 {
		lib.InColors(lib.Red, "No TODOS found\n")

	} else {

		fmt.Println("------------------------------")
		for _, issueMap := range p.Issues {
			for issueKey, todos := range issueMap {

				printSortedTodos(path.Base(issueKey), todos)
			}
		}
	}

	fmt.Println("------------------------------")
}

func printSortedTodos(issueKey string, todos []lib.Todo) {

	sort.Slice(todos, func(i, j int) bool {
		return todos[i].Urgency > todos[j].Urgency

	})

	lib.InColors(lib.Blue, fmt.Sprintf("Location: %s\n", issueKey))
	for _, todo := range todos {
		title := fmt.Sprintf("TODO: %v\n", todo.Title)
		lib.InColors(lib.Green, title)
		fmt.Printf("Line: %d\nUrgency: %d\n\n", todo.Line, todo.Urgency)
	}

}
func extractMatch(reader io.Reader, re *regexp.Regexp) (string, error) {
	s := bufio.NewScanner(reader)
	for s.Scan() {
		line := s.Text()
		matches := re.FindStringSubmatch(line)
		if matches != nil {
			captureGroup := matches[1]
			return captureGroup, nil
		}
	}

	if err := s.Err(); err != nil {
		return "", err
	}

	return "", nil
}
func makeAbsolute(fPath string) (string, error) {
	var dest string
	if !filepath.IsAbs(fPath) {
		pwd, err := os.Getwd()
		if err != nil {
			log.Fatal("Can't get the current working directory", err)
		}
		dest = filepath.Join(pwd, fPath)

	} else {
		dest = fPath
	}
	_, err := os.Stat(dest)
	if os.IsNotExist(err) {
		return "", err
	}

	return dest, nil
}

func (p *Project) saveProject() error {
	PROJECTS := viper.GetString("PROJECTS")
	var projects []Project

	metaPath := path.Join(PROJECTS, PROJECTS_META)

	if _, err := os.Stat(metaPath); err == nil {
		f, err := os.Open(metaPath)
		if err != nil {
			return err
		}
		defer f.Close()

		contents, err := io.ReadAll(f)
		if err != nil {
			return err
		}

		if len(contents) > 0 {
			err = json.Unmarshal(contents, &projects)
			if err != nil {
				return err
			}
		}
	} else if os.IsNotExist(err) {
		projects = make([]Project, 0)
	} else {
		return err
	}

	if projectExists(projects, p) {
		return nil
	}

	projects = append(projects, *p)

	jsonProjects, err := json.Marshal(projects)
	if err != nil {
		return err
	}

	tempFilePath := metaPath + ".tmp"
	tempFile, err := os.OpenFile(tempFilePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0664)
	if err != nil {
		return err
	}
	defer tempFile.Close()

	if _, err := tempFile.Write(jsonProjects); err != nil {
		return err
	}

	if err := os.Rename(tempFilePath, metaPath); err != nil {
		return err
	}

	return nil
}

func projectExists(projects []Project, newProject *Project) bool {
	for _, existingProject := range projects {
		if existingProject.Owner == newProject.Owner && existingProject.Name == newProject.Name {
			return true
		}
	}
	return false
}

func (p *Project) EditProject() {
	PROJECTS := viper.GetString("PROJECTS")
	pDir := path.Join(PROJECTS, p.Owner)
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
	p.NotePath = project
	if err := p.saveProject(); err != nil {
		fmt.Errorf("Failed to save the project u have to rerwrite to be able to open the noptes err:%v\n", err)
		time.Sleep(3 * time.Second)
	}

	lib.Edit(project)

}

func (p *Project) printProperites() string {

	properties := fmt.Sprintf(`------------------------------
Repo: %v/%v
Branch: %v
Url: %v
------------------------------`, p.Owner, p.Name, p.DefaultBranch, p.Url)
	return properties
}

func (pr *Project) WalkProject() error {
	if pr.Location == "" {
		log.Fatal("Failed to find the path to the Project")
	}
	// Change the directory because else git-ls will fail
	err := os.Chdir(pr.Location)
	if err != nil {
		fmt.Printf("Error changing directory: %v\n", err)
		return nil
	}

	cmd := exec.Command("git", "ls-files")
	var outb bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &outErr

	err = cmd.Run()
	if err != nil {
		if outErr.Len() > 0 {
			fmt.Println(outErr.String())
		} else {
			fmt.Printf("Error running command: %v\n", err)
		}
		return err
	}

	var wg sync.WaitGroup
	ctx := context.Background()
	var sem = semaphore.NewWeighted(int64(20))
	ch := make(chan map[string][]lib.Todo, 10)

	go func() {
		wg.Wait()
		close(ch)
	}()

	for scanner := bufio.NewScanner(&outb); scanner.Scan(); {
		filepath := scanner.Text()
		abFilepath := path.Join(pr.Location, filepath)
		wg.Add(1)
		go func(ab string) {
			sem.Acquire(ctx, 1)
			defer wg.Done()
			defer sem.Release(1)
			todoArray := walkFile(ab)

			if todoArray != nil && len(todoArray) > 0 {
				todosMap := make(map[string][]lib.Todo)
				todosMap[ab] = todoArray
				ch <- todosMap
			}
		}(abFilepath)
	}

	for v := range ch {
		pr.Issues = append(pr.Issues, v)
	}

	return nil
}

func walkFile(p string) []lib.Todo {
	info, err := os.Stat(p)
	if err != nil {
		//fmt.Println(err)
		return nil
	}
	if info.IsDir() {
		fmt.Println(info.Name(), "Probably a submodule")
		return nil
	}

	f, err := os.Open(p)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer f.Close()

	ch := make(chan lib.Todo)
	var wg sync.WaitGroup
	var TODOS []lib.Todo

	go func() {
		for todo := range ch {

			todo.Filename = path.Base(p)
			TODOS = append(TODOS, todo)

		}

	}()

	lineIndex := 0
	for s := bufio.NewScanner(f); s.Scan(); {
		line := s.Text()
		lineIndex++

		wg.Add(1)
		go func(line string, index int) {
			defer wg.Done()
			todo := lib.ContainsPattern(line, index, lib.TODO)
			if todo != nil {
				ch <- *todo
			}
		}(line, lineIndex)
	}

	wg.Wait()
	close(ch)

	if len(TODOS) != 0 {
		return TODOS
	}
	return nil
}

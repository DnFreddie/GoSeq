package github

import (
	"DnFreddie/goseq/lib"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/spf13/viper"
)

type Project struct {
	Name          string `json:"name"`
	Owner         string `json:"owner"`
	DefaultBranch string `json:"default_branch"`
	Url           string `json:"repo_url"`
	Issues        []map[string][]Todo
	Location      string `json:"location"`
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

	// Attempt to read the HEAD file
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

	// Attempt to read the config file
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

				printTodos(path.Base(issueKey), todos)
			}
		}
	}

	fmt.Println("------------------------------")
}

func printTodos(issueKey string, todos []Todo) {

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
	home := viper.GetString("HOME")
	var projects []Project

	metaPath := path.Join(home, lib.PROJECTS_META)

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

	// Create a temp file for to not lose data
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

// Helper function to check if the project already exists
func projectExists(projects []Project, newProject *Project) bool {
	for _, existingProject := range projects {
		if existingProject.Owner == newProject.Owner && existingProject.Name == newProject.Name {
			return true
		}
	}
	return false
}

func ReadRecent(list bool) error {
	home := viper.GetString("HOME")

	if !list {
		f, err := os.Open(lib.ENV_VAR)
		if err != nil {
			log.Println("No recent Projects found, listing recent projects instead")
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

	f, err := os.Open(path.Join(home, lib.PROJECTS_META))

	if err != nil {
		return fmt.Errorf("The meta file is empty add the project to fix this\n")
	}
	var projecArray []Project

	contents, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	err = json.Unmarshal(contents, &projecArray)
	if err != nil {
		return err
	}
	if len(projecArray) == 0 {
		return fmt.Errorf("The meta file is empty add the project to fix this\n")
	}
	pr := ChoseProject(&projecArray)
	pr.EditProject()
	return nil

}

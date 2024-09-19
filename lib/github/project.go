package github

import (
	"DnFreddie/GoSeq/lib"
	"bufio"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

type Project struct {
	Name          string
	Owner         string
	DefaultBranch string `json:"default_branch"`
	Url           string `json:"repo_url"`
	Issues        []map[string][]Todo
	Location      string
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

func (p *Project) Read() {

	if len(p.Issues) == 0 {
		lib.InColors(lib.Blue, "No TODOS found")
	}
	for _, issueMap := range p.Issues {
		for issueKey, todos := range issueMap {

			printTodos(path.Base(issueKey), todos)
		}
	}
}

func printTodos(issueKey string, todos []Todo) {
	lib.InColors(lib.Blue, fmt.Sprintf("Issue: %s\n", issueKey))
	for _, todo := range todos {
		fmt.Printf("Line:%d TODO:%s\n", todo.Line, todo.Title)
	}
}

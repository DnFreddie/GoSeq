package project

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/DnFreddie/goseq/internal/common"
)

type GitFileReader interface {
	ReadHead() (io.ReadCloser, error)
	ReadConfig() (io.ReadCloser, error)
}

// FSGitReader implements GitFileReader using the actual filesystem
type FSGitReader struct {
	projectPath string
}

func NewFSGitReader(path string) *FSGitReader {
	return &FSGitReader{projectPath: path}
}

func (fr *FSGitReader) ReadHead() (io.ReadCloser, error) {
	headPath := filepath.Join(fr.projectPath, ".git/HEAD")
	return os.Open(headPath)
}

func (fr *FSGitReader) ReadConfig() (io.ReadCloser, error) {
	configPath := filepath.Join(fr.projectPath, ".git/config")
	return os.Open(configPath)
}

type gitConfig struct {
	defaultBranch string
	url           string
	owner         string
	repoName      string
}

func extractGitConfig(reader GitFileReader) (gitConfig, error) {
	var config gitConfig

	// Extract branch from HEAD file
	headReader, err := reader.ReadHead()
	if err != nil {
		slog.Warn("Failed to read HEAD file", "error", err)
	} else {
		defer headReader.Close()
		config.defaultBranch, err = extractMatch(
			bufio.NewReader(headReader),
			regexp.MustCompile(`refs/heads/(\w+)`),
		)
		if err != nil {
			slog.Warn("Failed to find defutl branch", "err", err)
		}
	}

	// Extract URL from config file
	configReader, err := reader.ReadConfig()
	if err != nil {
		slog.Warn("Failed to read config file", "error", err)
	} else {
		defer configReader.Close()
		config.url, err = extractMatch(
			bufio.NewReader(configReader),
			regexp.MustCompile(`^\s*url\s*=\s*(https?://.+\.git)\s*$`),
		)
		if err != nil {
			slog.Warn("Failed to fund Url ", "err", err)
		}

	}

	// Get the owner and repo name from URL
	if config.url != "" {
		config.owner, config.repoName = parseGitURL(config.url)
	}

	if config.repoName == "" {
		slog.Error("Failed to determin URL", "err",
			common.NoMatchErr{
				Matched:  config.url,
				Expected: "repo name",
				Location: "Gitconfig lookup",
				Message:  "No url provided",
			})
		return config, fmt.Errorf("failed to determine repository name from URL: %s", config.url)
	}

	return config, nil
}

func extractMatch(reader *bufio.Reader, regex *regexp.Regexp) (string, error) {
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return "", err
		}

		if match := regex.FindStringSubmatch(line); len(match) > 1 {
			return match[1], nil
		}

		if err == io.EOF {
			break
		}
	}
	return "", common.NoMatchErr{
		Condition: regex.String(),
		Matched:   "nothing",
		Expected:  "something",
		Location:  "Gitconfig lookup",
		Message:   "Exact match expected",
	}
}

func parseGitURL(url string) (owner, repoName string) {
	parts := strings.Split(url, "/")
	if len(parts) < 2 {
		return "", ""
	}
	owner = parts[len(parts)-2]
	repoName = strings.TrimSuffix(parts[len(parts)-1], ".git")
	return owner, repoName

}

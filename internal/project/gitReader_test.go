package project

import (
	"io"
	"strings"
	"testing"
)

type MockGitReader struct {
	headContents   string
	configContents string
}

func (m *MockGitReader) ReadHead() (io.ReadCloser, error) {
	return io.NopCloser(strings.NewReader(m.headContents)), nil
}

func (m *MockGitReader) ReadConfig() (io.ReadCloser, error) {
	return io.NopCloser(strings.NewReader(m.configContents)), nil
}
func TestExtractGitConfig(t *testing.T) {
	mockReader := &MockGitReader{
		headContents: "ref: refs/heads/main\n",
		configContents: `[remote "origin"]
    url = https://github.com/owner/repo.git
    fetch = +refs/heads/*:refs/remotes/origin/*`,
	}

	config, err := extractGitConfig(mockReader)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if config.defaultBranch != "main" {
		t.Errorf("Expected branch 'main', got %s", config.defaultBranch)
	}
	if config.owner != "owner" {
		t.Errorf("Expected owner 'owner', got %s", config.owner)
	}
	if config.repoName != "repo" {
		t.Errorf("Expected repo 'repo', got %s", config.repoName)
	}
}

package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type gitIssue struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func sendRequest(method, url, token string, body io.Reader) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if method == "POST" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected response status: %s, body: %s", resp.Status, string(body))
	}

	return io.ReadAll(resp.Body)
}

func (p *Project) FetchGitHubIssues(token string) ([]gitIssue, error) {
	var issues []gitIssue
	url := fmt.Sprintf("https://api.github.com/repos/%v/%v/issues", p.Owner, p.Name)
	body, err := sendRequest("GET", url, token, nil)
	if err != nil {
		return issues, err
	}
	if err := json.Unmarshal(body, &issues); err != nil {
		return issues, fmt.Errorf("error unmarshaling issues: %v", err)
	}

	return issues, nil
}

func (t *Todo) createGitHubIssue(token string, owner string, repo string) error {
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	url := fmt.Sprintf("https://api.github.com/repos/%v/%v/issues", owner, repo)
	data, err := json.Marshal(t)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	if _, err = sendRequest("POST", url, token, bytes.NewBuffer(data)); err != nil {
		return err
	}

	return nil
}

func (p *Project) SearchIssueByTitle(token string, existIssue *[]gitIssue) {

	var wg sync.WaitGroup
	var mu sync.Mutex

	successCount := 0
	failedIssues := make(map[string]error)

	for _, issueMap := range p.Issues {
		for _, todos := range issueMap {
			for _, todo := range todos {
				isNewIssue := true
				for _, existingIssue := range *existIssue {
					// If the title matches, it's not a new issue
					if strings.EqualFold(todo.Title, existingIssue.Title) {
						isNewIssue = false
						fmt.Println("This is an old todo:", todo.Title)
						break
					}
				}

				if isNewIssue {
					wg.Add(1)
					go func(todo Todo) {
						defer wg.Done()
						err := todo.createGitHubIssue(token, p.Owner, p.Name)
						mu.Lock()
						defer mu.Unlock()
						if err != nil {
							failedIssues[todo.Title] = err
						} else {
							successCount++
						}
					}(todo)
				}
			}
		}
	}

	wg.Wait()

	fmt.Printf("\nSummary:\n")
	fmt.Printf("Successfully posted issues: %d\n", successCount)
	if len(failedIssues) > 0 {
		fmt.Printf("Failed to post %d issues:\n", len(failedIssues))
		for title, err := range failedIssues {
			fmt.Printf(" - %s: %v\n", title, err)
		}
	} else {
		fmt.Println("No issues failed to post.")
	}

}

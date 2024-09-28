package project

import (
	"github.com/DnFreddie/goseq/pkg/terminal"
	"github.com/DnFreddie/goseq/pkg/todo"
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

func (p *Project) ApplayIssues(token string, existIssue *[]gitIssue) {
	var newIssues []todo.Todo
	var oldIssues []todo.Todo

	for _, issueMap := range p.Issues {
		for _, todos := range issueMap {
			for _, todo := range todos {
				isNewIssue := true
				for _, existingIssue := range *existIssue {
					if strings.EqualFold(todo.Title, existingIssue.Title) {
						isNewIssue = false
						oldIssues = append(oldIssues, todo)
						break
					}
				}
				if isNewIssue {
					newIssues = append(newIssues, todo)
				}
			}
		}
	}

	if err := goseqPlan(newIssues, oldIssues); err != nil {

		return

	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	successCount := 0
	failedIssues := make(map[string]error)

	for _, v := range newIssues {
		wg.Add(1)
		go func( t todo.Todo) {
			defer wg.Done()
			err := createGitHubIssue[todo.Todo](t, token, p.Owner, p.Name)

			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				failedIssues[t.Title] = err
			} else {
				successCount++
			}
		}(v)
	}
	wg.Wait()

	// Summary
	fmt.Printf("\nExecution Summary:\n")
	terminal.InColors(terminal.Green, fmt.Sprintf("Successfully posted issues: %d\n", successCount))
	if len(failedIssues) > 0 {
		fmt.Printf("Failed to post %d issues:\n", len(failedIssues))
		for title, err := range failedIssues {
			terminal.InColors(terminal.Red, fmt.Sprintf(" - %s: %v\n", title, err))
		}
	} else {
		fmt.Println("No issues failed to post.")
	}
}

func goseqPlan(newIssues []todo.Todo, oldIssues []todo.Todo) error {
	fmt.Println("\nPlan:")
	fmt.Printf("New issues to be created: %d\n", len(newIssues))
	for _, issue := range newIssues {
		terminal.InColors(terminal.Green, fmt.Sprintf(" + %s\n", issue.Title))
	}
	fmt.Printf("\nExisting issues (no changes): %d\n", len(oldIssues))
	for _, issue := range oldIssues {
		terminal.InColors(terminal.Red, fmt.Sprintf(" = %s\n", issue.Title))
	}
	if len(newIssues) == 0 {
		fmt.Println("No new issues exiting...")
		return fmt.Errorf("No new issues")
	}
	fmt.Print("\nDo you want to proceed with this plan? (yes/no): ")
	var response string
	fmt.Scanln(&response)

	if strings.ToLower(response) != "yes" {
		fmt.Println("Operation cancelled.")

		return fmt.Errorf("GoSeq Plan cancelled")
	}
	return nil
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

func createGitHubIssue[T any](todo T, token string, owner string, repo string) error {
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	url := fmt.Sprintf("https://api.github.com/repos/%v/%v/issues", owner, repo)
	data, err := json.Marshal(todo)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	if _, err = sendRequest("POST", url, token, bytes.NewBuffer(data)); err != nil {
		return err
	}

	return nil
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

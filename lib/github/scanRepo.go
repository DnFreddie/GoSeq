package github

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
)

func WalkFile(p string) []Todo {
	info, err := os.Stat(p)
	if err != nil {
		fmt.Println(err)
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

	ch := make(chan Todo)
	var wg sync.WaitGroup
	var TODOS []Todo

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
		go func(s string, index int) {
			defer wg.Done()
			todo := containsTODO(s, index)
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

func containsTODO(line string, lineIndex int) *Todo {
	// Find the index of "TODO"
	index := strings.Index(line, "TODO")
	var titleIndex int

	var pririoryty int
	if index == -1 {
		return nil
	}

	for i := index + 4; i < len(line); i++ {
		if line[i] == '!' {
			titleIndex = i + 1
			break

		} else if line[i] != 'O' {
			titleIndex = i
			break

		}
		pririoryty += 1
	}
	title := line[titleIndex:]

	if title == "" {
		return nil
	}
	return &Todo{
		Urgency: pririoryty,
		Title:   strings.TrimSpace(title),
		Line:    lineIndex,
	}
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
	ch := make(chan map[string][]Todo)

	go func() {
		wg.Wait()
		close(ch)
	}()

	for scanner := bufio.NewScanner(&outb); scanner.Scan(); {
		filepath := scanner.Text()
		abFilepath := path.Join(pr.Location, filepath)
		wg.Add(1)
		go func(ab string) {
			defer wg.Done()
			todoArray := WalkFile(ab)

			if todoArray != nil && len(todoArray) > 0 {
				todosMap := make(map[string][]Todo)
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

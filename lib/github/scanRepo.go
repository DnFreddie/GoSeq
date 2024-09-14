package github

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
)


func WalkFile(p string) {
	info, err := os.Stat(p)
	if err != nil {
		fmt.Println(err)
		return
	}
	if info.IsDir() {
		fmt.Println(info.Name(), "Probably a submodule")
		return
	}

	f, err := os.Open(p)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close() 

	ch := make(chan Todo)
	var wg sync.WaitGroup

	go func() {
		for todo := range ch {
			fmt.Println(todo)
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

func WalkProject(p string) {

	//U have to change the dir becouse else it git-ls  will fail
	err := os.Chdir(p)
	if err != nil {
		fmt.Printf("Error changing directory: %v\n", err)
		return
	}
	cmd := exec.Command("git", "ls-files")
	var outb bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &outErr

	err = cmd.Run()
	if err != nil {
		fmt.Println(outErr)
	}

	var wg sync.WaitGroup
	for scanner := bufio.NewScanner(&outb); scanner.Scan(); {
		filepath := scanner.Text()
		abFilepaht := path.Join(p, filepath)
		wg.Add(1)
		go func(ab string) {
			defer wg.Done()
			WalkFile(abFilepaht)

		}(abFilepaht)

	}
	wg.Wait()

}

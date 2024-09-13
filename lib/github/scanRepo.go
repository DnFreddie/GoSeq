package github

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"sync"
)

func WalkFile(p string) {

	info, err := os.Stat(p)
	if err != nil {
		fmt.Println(err)
		return
	}
	if info.IsDir() {
		fmt.Println(info.Name(),"Propably a submodule ")
		return
	}
	fmt.Println(info.Name())

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
		go func (ab string){
		defer wg.Done()
        WalkFile(abFilepaht)


		}(abFilepaht)

	}
	wg.Wait()

}

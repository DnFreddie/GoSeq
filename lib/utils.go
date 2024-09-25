package lib

import (
	"fmt"
	"os"
	"os/exec"
)


const (

)

func Map[T any, U any](input []T, fn func(T) U) []U {
	output := make([]U, len(input))
	for i, v := range input {
		output[i] = fn(v)
	}
	return output
}
func Edit(fPath string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return fmt.Errorf("Failed to foudn $EDITOR")
	}
	cmd := exec.Command(editor, fPath)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		return err
	}

	return nil
}





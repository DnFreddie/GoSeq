package lib

import (
	"fmt"

)

func TesContainsPattern() {
	line := "TODO!! This is an urgent task"
	line2 := "IDEA: This is an Idea"
	line3 := "NOTE! This is the Note!"
//	line4 := "Random text without patterns"

	todoItem := ContainsPattern(line, 1, TODO)
	if todoItem != nil {
		fmt.Printf("Found: %+v\n", *todoItem)
	} else {
		fmt.Println("No pattern found.")
	}

	todoItem2 := ContainsPattern(line2, 1, TODO|IDEA)
	if todoItem2 != nil {
		fmt.Printf("Found: %+v\n", *todoItem2)
	} else {
		fmt.Println("No pattern found.")
	}

	todoItem3 := ContainsPattern(line3, 1, ALL)
	if todoItem3 != nil {
		fmt.Printf("Found (ALL): %+v\n", *todoItem3)
	} else {
		fmt.Println("No pattern found with ALL.")
	}
}

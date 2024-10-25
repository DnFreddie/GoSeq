package quizz

import (
	"bufio"
	"fmt"
	"strings"
)

// Parse Header and Get the Name and Level

// Get the tree of headers
func printHeaders(headers []Header, indent int) {
	for _, h := range headers {
		prefix := "Header: "
		if indent > 0 {
			prefix = "Child: "
		}
		fmt.Printf("%s%sLevel=%d, Name=%q, Range=(%d, %d)\n",
			strings.Repeat("  ", indent), prefix, h.Level, h.Name, h.Range.Start, h.Range.End)
		printHeaders(h.Children, indent+1)
	}
}

func ExampleHeaderScanner_Scan_basic() {
	content := `# Main Topic
## Subtopic A
### Deep Topic
## Subtopic B
# Another Main Topic`

	reader := strings.NewReader(content)
	headerScanner := &HeaderScanner{
		Scanner: bufio.NewScanner(reader),
	}

	headers, err := headerScanner.Scan()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print headers with proper formatting
	printHeaders(headers, 0)

	// Output:
	// Header: Level=1, Name="Main Topic", Range=(0, 3)
	//   Child: Level=2, Name="Subtopic A", Range=(1, 2)
	//     Child: Level=3, Name="Deep Topic", Range=(2, 2)
	//   Child: Level=2, Name="Subtopic B", Range=(3, 3)
	// Header: Level=1, Name="Another Main Topic", Range=(4, 4)
}

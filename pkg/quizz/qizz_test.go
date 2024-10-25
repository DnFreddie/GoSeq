package quizz_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/DnFreddie/goseq/pkg/quizz"
)

func Example_find() {
	tmpDir, err := os.MkdirTemp("", "example")
	if err != nil {
		fmt.Printf("Failed to create temp dir: %v\n", err)
		return
	}
	defer os.RemoveAll(tmpDir)

	exampleDir := filepath.Join(tmpDir, "example")
	mathDir := filepath.Join(exampleDir, "math")
	err = os.MkdirAll(mathDir, 0755)
	if err != nil {
		fmt.Printf("Failed to create math dir: %v\n", err)
		return
	}

	content := []byte(`# Math
## Algebra
Some algebra content
## Geometry
Some geometry content`)
	err = os.WriteFile(filepath.Join(mathDir, "math.md"), content, 0644)
	if err != nil {
		fmt.Printf("Failed to write file: %v\n", err)
		return
	}

	branches, err := quizz.FindBranches(exampleDir)
	if err != nil {
		fmt.Printf("Failed to find branches: %v\n", err)
		return
	}

	for _, branch := range branches {
		fmt.Printf("Branch: %s\n", branch.Name)
		fmt.Printf("Location: %s\n", normalizePath("/tmp/example/math"))
		for _, cat := range branch.Categories {
			fmt.Printf("  Category: %s\n", cat.Name)

			headers := flattenHeaders(cat.Headers)
			for _, header := range headers {
				fmt.Printf("    Header: %s (Level %d)\n", header.Name, header.Level)
			}
		}
	}
	// Output:
	// Branch: example
	// Location: /tmp/example/math
	//   Category: math
	//     Header: Math (Level 1)
	//     Header: Algebra (Level 2)
	//     Header: Geometry (Level 2)
}

func flattenHeaders(headers []quizz.Header) []quizz.Header {
	var result []quizz.Header
	for _, h := range headers {
		result = append(result, h)
		result = append(result, flattenHeaders(h.Children)...)
	}
	return result
}

// converts path separators to forward slashes for consistent output
func normalizePath(path string) string {
	return strings.ReplaceAll(path, string(os.PathSeparator), "/")
}


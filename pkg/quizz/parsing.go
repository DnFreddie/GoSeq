package quizz

import (
	"bufio"
	"io"
	"strings"
)

func ParseHeader(s string) *Header {
	level := 0

	for i, char := range s {
		if char != '#' {
			break
		}
		level = i + 1
	}
	if level == 0 || len(s) <= level || s[level] != ' ' {
		return nil
	}

	name := strings.TrimSpace(s[level+1:])
	//  Check empty header with all spaces
	if name == "" {
		return nil
	}

	return &Header{
		Level: level,
		Name:  name,
	}
}

type HeaderScanner struct {
	*bufio.Scanner
}

func NewHeaderScanner(r io.Reader) *HeaderScanner {
	return &HeaderScanner{
		Scanner: bufio.NewScanner(r),
	}
}

func (hs *HeaderScanner) Scan() ([]Header, error) {
	var headers []Header
	var currentPath []int // Tracks indices of current hierarchy
	lineNumber := 0

	for hs.Scanner.Scan() {
		line := strings.TrimSpace(hs.Scanner.Text())
		header := ParseHeader(line)

		if header != nil {
			header.Range = Range{Start: lineNumber, End: lineNumber}

			for len(currentPath) > 0 {
				parentIdx := currentPath[len(currentPath)-1]
				var parentHeader *Header
				if len(currentPath) == 1 {
					parentHeader = &headers[parentIdx]
				} else {
					temp := &headers[currentPath[0]]
					for i := 1; i < len(currentPath)-1; i++ {
						temp = &temp.Children[currentPath[i]]
					}
					parentHeader = &temp.Children[parentIdx]
				}

				if parentHeader.Level < header.Level {
					parentHeader.Children = append(parentHeader.Children, *header)
					currentPath = append(currentPath, len(parentHeader.Children)-1)
					break
				}
				// Remove levels >= current header level
				currentPath = currentPath[:len(currentPath)-1]
			}

			if len(currentPath) == 0 {
				// Top-level header
				headers = append(headers, *header)
				currentPath = []int{len(headers) - 1}
			}
		}
		lineNumber++
	}

	// Update end ranges
	updateEndRanges(headers, lineNumber-1)

	if err := hs.Scanner.Err(); err != nil {
		return nil, err
	}
	return headers, nil
}

func updateEndRanges(headers []Header, lastLine int) {
	for i := range headers {
		if len(headers[i].Children) > 0 {
			// If header has children, end range is just before next sibling
			if i < len(headers)-1 {
				headers[i].Range.End = headers[i+1].Range.Start - 1
			} else {
				headers[i].Range.End = lastLine
			}
			updateEndRanges(headers[i].Children, headers[i].Range.End)
		} else {
			// If no children, end range extends to next header or end of file
			if i < len(headers)-1 {
				headers[i].Range.End = headers[i+1].Range.Start - 1
			} else {
				headers[i].Range.End = lastLine
			}
		}
	}
}

package common

import (
	"fmt"
)

type NoNotesFoundErr struct{}

func (e NoNotesFoundErr) Error() string {
	return "No notes available ..."
}

type FLockedErr struct{}

func (e FLockedErr) Error() string {
	return "Timeout waiting for file lock - another instance may be running"

}

type NoMatchErr struct {
	Condition string
	Matched   string // (empty if no match found)
	Expected  string
	Message   string // An optional custom message to provide context
	Location  string // A location or function name where the error occurred
}

func (e NoMatchErr) Error() string {
	return fmt.Sprintf("No match found at %s: Condition: %s, Expected: %s, Message: %s, Matched: %s",
		e.Location, e.Condition, e.Expected, e.Message, e.Matched)
}

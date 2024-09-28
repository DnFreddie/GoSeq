package common
type NoNotesError struct{}

func (e NoNotesError) Error() string {
	return "No notes available ..."
}


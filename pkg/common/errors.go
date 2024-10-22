package common

type NoNotesFoundErr struct{}

func (e NoNotesFoundErr) Error() string {
	return "No notes available ..."
}

type FLockedErr struct{}

func (e FLockedErr) Error() string {
	return "Timeout waiting for file lock - another instance may be running"
}

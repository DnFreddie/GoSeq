package lib

type DateLayout string

const (
	ENV_VAR = "/tmp/GO_SEQ_PROJECT.txt"
	PROJECTS_META = PROJECTS + "/.PROJECTS_META.json"

	FileDate DateLayout = "2006-01-02"
	FullDate DateLayout = "January 2 2006"
)


func Map[T any, U any](input []T, fn func(T) U) []U {
	output := make([]U, len(input))
	for i, v := range input {
		output[i] = fn(v)
	}
	return output
}

package github

import "fmt"

type Todo struct {
	Keyword       string
	Urgency       int
	ID            *string
	Filename      string
	Line          int
	Title         string
	Body          []string
	BodySeparator string
}
type GitCreds struct {
	Token string
}

func (gc GitCreds) query() error {

	return nil
}

func (gc GitCreds) postIssue(repo string, todo Todo) (Todo, error) {
	fmt.Println(repo)
	err := gc.query()
	
	if err != nil {
		return todo, err
	}

	//id := "#" + strconv.Itoa(int(json["number"].(float64)))
	//todo.ID = &id

	return todo, err
}




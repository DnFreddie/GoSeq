package project

import (
	"github.com/DnFreddie/goseq/pkg/common"
	"github.com/DnFreddie/goseq/pkg/terminal"
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"slices"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

const (
	JOINED_DELETE = "/tmp/.go_seq_projects_joined.md"
	JOINED        = "/tmp/.go_seq_projects_joined.md"
)

type ProjectManager struct{}

func NewProjectManager() *ProjectManager {
	retriver := ProjectManager{}
	return &retriver
}

func (pm *ProjectManager) GetNotes(p common.Period) ([]Project, error) {

	return getSavedProjects()
}

func (p *ProjectManager) JoinNotesByTitle(notes *[]Project) (io.Reader, error) {

	return joinByTitle(notes)
}
func (p *ProjectManager) JoinNotesWithContents(notes *[]Project) (io.Reader, error) {
	return nil, nil
}

func (pm *ProjectManager) Scan(r io.Reader, scanner ProjectScanner) ([]Project, error) {

	return nil, nil
}

func (pm *ProjectManager) DeleteByTitle(r io.Reader, n *[]Project) error {

	return deleteByTitle(r, n)

}

func getSavedProjects()([]Project,error){
	PROJECTS := viper.GetString("PROJECTS")
	var projecArray []Project
	f, err := os.Open(path.Join(PROJECTS, PROJECTS_META))

	if err != nil {
		return projecArray, fmt.Errorf("The meta file is empty add the project to fix this\n\ngit -p <path/to/project/\n")
	}

	contents, err := io.ReadAll(f)
	if err != nil {
		return projecArray,err
	}
	err = json.Unmarshal(contents, &projecArray)
	if err != nil {
		return projecArray,err
	}
	if len(projecArray) == 0 {
		return projecArray, common.NoNotesError{}
	}
	return projecArray,nil
}

func joinByTitle(notes *[]Project) (io.Reader, error) {
	f, err := os.OpenFile(JOINED_DELETE, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	var titles []string
	for _, note := range *notes {
		formattedName, err := note.Format()
		if err != nil {
			formattedName = note.GetPath()
		}
		titles = append(titles, formattedName)
	}
	joinedTitles := strings.Join(titles, "\n")

	if _, err := f.Write([]byte(joinedTitles)); err != nil {
		return nil, err
	}

	if err := common.Edit(JOINED_DELETE); err != nil {
		return nil, err
	}

	updatedContent, err := os.ReadFile(JOINED_DELETE)
	if err != nil {
		return nil, fmt.Errorf("failed to read updated file: %w", err)
	}

	return bytes.NewReader(updatedContent), nil
}

func deleteByTitle(r io.Reader, notes *[]Project) error {
	titles, err := readTitles(r)
	if err != nil {
		return err
	}

	updatedProjectsData, errChan := processNotes(*notes, titles)

	if err := collectErrors(errChan); err != nil {
		return err
	}

	if len(updatedProjectsData) == len(*notes) {
		terminal.InColors(terminal.Red, "Nothing to delete ...\n")
		return nil
	}

	return updatedSavedJson(updatedProjectsData)
}

func readTitles(r io.Reader) ([]string, error) {
	var titles []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		titles = append(titles, scanner.Text())
	}
	return titles, scanner.Err()
}

func processNotes(notes []Project, titles []string) ([]*Project, chan error) {
	updatedProjectsData := make([]*Project, 0, len(notes))
	errChan := make(chan error, len(notes))
	var wg sync.WaitGroup

	for _, note := range notes {
		formatted, err := note.Format()
		if err != nil {
			note.SaveProject()
			formatted = note.GetPath()
		}
		if !slices.Contains(titles, formatted) {
			wg.Add(1)
			go func(n Project) {
				defer wg.Done()
				if err := note.Delete(); err != nil {
					errChan <- err
				}
			}(note)
		} else {
			updatedProjectsData = append(updatedProjectsData, &note)
		}
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	return updatedProjectsData, errChan
}

func collectErrors(errChan <-chan error) error {
	var errs []error
	for err := range errChan {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func updatedSavedJson(updatedProjectsData []*Project) error {
	PROJECTS := viper.GetString("PROJECTS")
	metaPath := path.Join(PROJECTS, PROJECTS_META)
	updatedContent, err := json.Marshal(updatedProjectsData)
	if err != nil {
		return fmt.Errorf("error marshaling updated JSON: %v", err)
	}
	err = os.WriteFile(metaPath, updatedContent, 0644)
	if err != nil {
		return fmt.Errorf("error writing updated ENV_VAR file: %v", err)
	}
	return nil
}

package quizz

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

const (
	DEFAULT_INDEX_DIR  = ".qizz"
	DEFAULT_INDEX_FILE = ".index"
)

type Range struct {
	Start int
	End   int
}

type Header struct {
	Name       string
	Level      int
	Range      Range
	Flashcards []int
	Children   []Header
}

type Category struct {
	Name     string
	Location string
	Headers  []Header
	ModTime  time.Time
}

type Branch struct {
	Name       string
	Location   string
	Categories []Category
}

type FlashManager struct {
	Branches []Branch
}

func NewFlashManager() *FlashManager {
	return &FlashManager{
		Branches: make([]Branch, 0),
	}
}

func (fm *FlashManager) Index(fPath string) error {
	branches, err := FindBranches(fPath)
	if err != nil {
		return fmt.Errorf("error finding categories and creating branches: %w", err)
	}

	if len(branches) == 0 {
		return fmt.Errorf("no branches found")
	}

	fm.Branches = branches
	return nil
}

func FindBranches(dir string) ([]Branch, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("error getting absolute path: %w", err)
	}

	w := &walker{
		branchMap: make(map[string]*Branch),
		errorsArr: make([]error, 0),
	}

	err = filepath.WalkDir(absDir, w.filter)
	if err != nil {
		w.errorsArr = append(w.errorsArr, err)
	}

	if len(w.errorsArr) > 0 {
		return nil, errors.Join(w.errorsArr...)
	}

	// Convert map to slice
	branches := make([]Branch, 0, len(w.branchMap))
	for _, branch := range w.branchMap {
		branches = append(branches, *branch)
	}

	logBranches(branches)
	return branches, nil
}

type walker struct {
	branchMap map[string]*Branch
	errorsArr []error
}

func (w *walker) filter(path string, d os.DirEntry, err error) error {
	if err != nil {
		w.errorsArr = append(w.errorsArr, fmt.Errorf("error accessing path %s: %w", path, err))
		return nil
	}

	if category, ok := ParseToCategory(d, path); ok {
		branchName := filepath.Base(filepath.Dir(category.Location))
		if branchName == "." || branchName == "/" {
			branchName = filepath.Base(filepath.Dir(filepath.Dir(category.Location)))
		}

		if branch, exists := w.branchMap[branchName]; exists {
			branch.Categories = append(branch.Categories, category)
		} else {
			w.branchMap[branchName] = &Branch{
				Name:       branchName,
				Location:   filepath.Dir(category.Location),
				Categories: []Category{category},
			}
		}
	}
	return nil
}

func ParseToCategory(dirE fs.DirEntry, currentPath string) (Category, bool) {
	if dirE.IsDir() {
		return Category{}, false
	}

	fp := filepath.Join(filepath.Dir(currentPath), dirE.Name())
	f, err := os.Open(fp)
	if err != nil {
		return Category{}, false
	}
	defer f.Close()

	hs := NewHeaderScanner(f)
	headers, err := hs.Scan()
	if err != nil {
		slog.Warn("Scan had failures", "failures", err)
	}

	location := filepath.Dir(fp)
	name := filepath.Base(location)

	fstats, err := os.Stat(fp)
	if err != nil {
		slog.Error("Failed to get file stats", "error", err)
		return Category{}, false
	}

	return Category{
		Name:     name,
		Location: location,
		Headers:  headers,
		ModTime:  fstats.ModTime(),
	}, true
}

func logBranches(branches []Branch) {
	branchNames := make([]string, len(branches))
	for i, branch := range branches {
		branchNames[i] = branch.Name
	}
	slog.Info("Created branches",
		"count", len(branches),
		"names", branchNames,
	)
}


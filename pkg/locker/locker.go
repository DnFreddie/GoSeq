package locker

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Locker interface {
	Lock() error
	Unlock() error
	IsLocked() bool
}

type LockFile string


type FileLocker struct {
	LockFile string
	Service  string
}

func NewFileLocker(lockFile LockFile, service string) *FileLocker {
	return &FileLocker{
		LockFile: string(lockFile),
		Service:  service,
	}
}
func (f *FileLocker) Lock() error {
	err := createLockFile(f.LockFile)
	if err != nil {
		return fmt.Errorf("%s failed to acquire lock: %w", f.Service, err)
	}
	return nil
}

func (f *FileLocker) Unlock() error {
	err := os.Remove(f.LockFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("%s  lock file does not exist: %w", f.Service, err)
		}
		return fmt.Errorf("%s  failed to remove lock file: %w", f.Service, err)
	}
	return nil
}

func (f *FileLocker) IsLocked() bool {
	_, err := os.Stat(f.LockFile)
	return err == nil
}

func createLockFile(fpath string) error {
	dir := filepath.Dir(fpath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory for lock file: %w", err)
	}
	file, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return fmt.Errorf("lock file already exists")
		}
		return fmt.Errorf("failed to create lock file: %w", err)
	}
	defer file.Close()

	return nil
}


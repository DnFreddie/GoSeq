package common

import (
	"log"
	"os"

	"github.com/rogpeppe/go-internal/lockedfile"
)

// Closes and Cleanups the temporary  file
func CleanupFileHandler(f *lockedfile.File, path string) {
	if closeErr := f.Close(); closeErr != nil {
		log.Printf("Error closing temporary file: %v", closeErr)
	}

	if removeErr := os.Remove(path); removeErr != nil {
		log.Printf("Error removing temporary file: %v", removeErr)

	}
}

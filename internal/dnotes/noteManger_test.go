package dnotes

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/DnFreddie/goseq/internal/common"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestEnvironment(t *testing.T) string {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "notes-test-*")
	require.NoError(t, err)

	// Set up viper config for test
	viper.Set("DAILIES", tmpDir)

	// Create test files
	testFiles := map[string]string{
		"2024-01-01.md": "Test content 1",
		"2024-01-15.md": "Test content 2",
		"2024-02-01.md": "Test content 3",
		"2024-03-01.md": "Test content 4",
		"invalid.txt":   "Invalid file",
		"2024-13-01.md": "Invalid date",
	}

	for filename, content := range testFiles {
		err := os.WriteFile(filepath.Join(tmpDir, filename), []byte(content), 0644)
		require.NoError(t, err)
	}

	return tmpDir
}

func TestXGetNotes(t *testing.T) {
	tmpDir := setupTestEnvironment(t)
	defer os.RemoveAll(tmpDir)

	baseTime := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)

	t.Run("XGetNotesLastMonth", func(t *testing.T) {
		period := common.Period{
			Range:  common.Month,
			Amount: 1,
			Today:  baseTime,
		}

		notes, err := XGetNotes(period)
		assert.NoError(t, err)
		assert.Len(t, notes, 2, "should only find February notes")
		verifyNotes(t, notes, baseTime, period)
	})

	t.Run("XGetNotesLastThreeMonths", func(t *testing.T) {
		period := common.Period{
			Range:  common.Month,
			Amount: 3,
			Today:  baseTime,
		}

		notes, err := XGetNotes(period)
		assert.NoError(t, err)
		assert.Len(t, notes, 4, "should find all valid notes from Jan to March")
		verifyNotes(t, notes, baseTime, period)
	})

	t.Run("XGetNotesNoNotesInRange", func(t *testing.T) {
		oldTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
		period := common.Period{
			Range:  common.Month,
			Amount: 1,
			Today:  oldTime,
		}

		notes, err := XGetNotes(period)
		assert.Error(t, err)
		assert.ErrorIs(t, err, &common.NoNotesFoundErr{})
		assert.Nil(t, notes)
	})

	t.Run("XGetNotesSingleDay", func(t *testing.T) {
		specificDay := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		period := common.Period{
			Range:  common.Day,
			Amount: 1,
			Today:  specificDay,
		}

		notes, err := XGetNotes(period)
		assert.NoError(t, err)
		assert.Len(t, notes, 1, "should only find 2024-01-15.md")
		verifyNotes(t, notes, specificDay, period)
	})

	t.Run("XGetNotesAll", func(t *testing.T) {
		period := common.Period{
			Range: common.All,
			Today: baseTime,
		}

		notes, err := XGetNotes(period)
		assert.NoError(t, err)
		assert.Len(t, notes, 4, "should find all valid .md files")
		verifyNotes(t, notes, baseTime, period)
	})
}

func verifyNotes(t *testing.T, notes []BasicNote, today time.Time, period common.Period) {
	for _, note := range notes {
		assert.NotEmpty(t, note.Path)
		assert.NotEmpty(t, note.Title)
		assert.True(t, strings.HasSuffix(note.Title, ".md"))
		assert.False(t, note.Date.IsZero())

		// Verify the date is within the expected range
		if period.Range != common.All {
			assert.True(t, common.DateInRange(period, note.Date),
				"note date %v should be within range for period %v from %v",
				note.Date, period, today)
		}
	}
}

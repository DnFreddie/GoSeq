package common

import (
	"testing"
	"time"
)

func TestDateInRange(t *testing.T) {
	// This is our reference "today" date for all tests
	fixedToday := time.Date(2024, 9, 22, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		r        DateRange
		amount   int
		testDate time.Time // The date we're testing against
		expect   bool
	}{
		{
			name:     "1 day ago",
			r:        Day,
			amount:   1,
			testDate: time.Date(2024, 9, 21, 0, 0, 0, 0, time.UTC),
			expect:   true,
		},
		{
			name:     "3 day ago",
			r:        Day,
			amount:   3,
			testDate: time.Date(2024, 9, 19, 0, 0, 0, 0, time.UTC),
			expect:   true,
		},
		{
			name:     "1 week ago",
			r:        Week,
			amount:   1,
			testDate: time.Date(2024, 9, 15, 0, 0, 0, 0, time.UTC),
			expect:   true,
		},
		{
			name:     "1 month ago",
			r:        Month,
			amount:   1,
			testDate: time.Date(2024, 8, 22, 0, 0, 0, 0, time.UTC),
			expect:   true,
		},
		{
			name:     "1 year ago",
			r:        Year,
			amount:   1,
			testDate: time.Date(2023, 9, 22, 0, 0, 0, 0, time.UTC),
			expect:   true,
		},
		{
			name:     "today",
			r:        Day,
			amount:   0, // Changed to 0 for today
			testDate: time.Date(2024, 9, 22, 0, 0, 0, 0, time.UTC),
			expect:   true,
		},
		{
			name:     "in the future",
			r:        Week,
			amount:   1,
			testDate: time.Date(2024, 9, 23, 0, 0, 0, 0, time.UTC),
			expect:   false,
		},
		{
			name:     "yesterday",
			r:        Yesterday,
			amount:   0,
			testDate: time.Date(2024, 9, 21, 0, 0, 0, 0, time.UTC), // Fixed date for yesterday
			expect:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := Period{
				Range:  test.r,
				Amount: test.amount,
				Today:  fixedToday, // Using fixedToday as our reference point
			}

			result := DateInRange(p, test.testDate)
			if result != test.expect {
				t.Errorf("DateInRange(Period{Range: %v, Amount: %d, Today: %v}, %v) = %v, want %v",
					test.r, test.amount, fixedToday, test.testDate, result, test.expect)
			}
		})
	}
}

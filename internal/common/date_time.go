package common

import (
	"fmt"
	"strings"
	"time"
)

type DateRange int

const (
	Day       DateRange = 1
	Week      DateRange = 7
	Month     DateRange = 30
	Year      DateRange = 365
	All       DateRange = 0
	Yesterday DateRange = 2
)

type Period struct {
	Range  DateRange
	Amount int
	Today  time.Time
}

func ParseDateRange(input string) DateRange {

	switch strings.ToLower(input) {
	case "day":
		return Day
	case "week":
		return Week
	case "month":
		return Month
	case "year":
		return Year
	case "all":
		return All
	default:
		fmt.Printf("Invalid date range: %s. Defaulting to 'all'.\n", input)
		return Week
	}
}

func DateInRange(period Period, date time.Time) bool {
	if period.Range == All {
		return true
	}

	// Don't allow future dates
	if date.After(period.Today) {
		return false
	}

	var startDate time.Time

	switch period.Range {
	case Day:
		if period.Amount == 0 {
			// Special case for today
			startDate = period.Today.Truncate(24 * time.Hour)
		} else {
			startDate = period.Today.AddDate(0, 0, -period.Amount)
		}
	case Yesterday:
		startDate = period.Today.AddDate(0, 0, -1)
	case Week:
		startDate = period.Today.AddDate(0, 0, -period.Amount*7)
	case Month:
		startDate = period.Today.AddDate(0, -period.Amount, 0)
	case Year:
		startDate = period.Today.AddDate(-period.Amount, 0, 0)
	default:
		return false
	}

	// Normalize times to start of day for comparison
	startDate = startDate.Truncate(24 * time.Hour)
	compareDate := date.Truncate(24 * time.Hour)
	endDate := period.Today.Truncate(24 * time.Hour)

	return (compareDate.After(startDate) || compareDate.Equal(startDate)) &&
		(compareDate.Before(endDate) || compareDate.Equal(endDate))
}

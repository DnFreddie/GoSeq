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
func DateInRange(today time.Time, r Period, date time.Time) bool {
	var searchPattern time.Time

	switch r.Range {
	case Day:
		searchPattern = today.AddDate(0, 0, -r.Amount)
	case Yesterday:
		searchPattern = today.AddDate(0, 0, -2)
	case Week:
		searchPattern = today.AddDate(0, 0, -r.Amount*7)
	case Month:
		searchPattern = today.AddDate(0, -r.Amount, 0)
	case Year:
		searchPattern = today.AddDate(-r.Amount, 0, 0)
	case All:
		return true

	default:
		return false
	}

	return !date.Before(searchPattern) && !date.After(today)

}

package notes

import (
	"time"
)


func dateInRange(today time.Time, r Period , date time.Time) bool {
    var searchPattern time.Time

    switch r.Range {
    case Day:
        searchPattern = today.AddDate(0, 0, -r.Amount)
	case Yesterday:
	searchPattern = today.AddDate(0,0,-2)
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

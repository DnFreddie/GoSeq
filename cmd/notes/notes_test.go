package notes

import (
	"testing"
	"time"
)


func TestFindDate(t *testing.T) {
    fixedDate := time.Date(2024, 9, 22, 0, 0, 0, 0, time.UTC)

    tests := []struct {
        name   string
        r      DateRange
        amount int
        date   time.Time
        expect bool
    }{
        {"1 day ago", Day, 1, time.Date(2024, 9, 21, 0, 0, 0, 0, time.UTC), true},
        {"3 day ago", Day, 3, time.Date(2024, 9, 19, 0, 0, 0, 0, time.UTC), true},
        {"1 week ago", Week, 1, time.Date(2024, 9, 15, 0, 0, 0, 0, time.UTC), true},
        {"1 month ago", Month, 1, time.Date(2024, 8, 22, 0, 0, 0, 0, time.UTC), true},
        {"1 year ago", Year, 1, time.Date(2023, 9, 22, 0, 0, 0, 0, time.UTC), true},
        {"today", Day, 1, time.Date(2024, 9, 22, 0, 0, 0, 0, time.UTC), true},
        {"in the future", Week, 1, time.Date(2024, 9, 23, 0, 0, 0, 0, time.UTC), false},
        {"yesterday", Yesterday, 0, time.Date(2024, 9, 20, 0, 0, 0, 0, time.UTC), true},


    }

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
			p:=  Period{
				Range: test.r,
				Amount: test.amount,
			}
            result := dateInRange(fixedDate,p, test.date)
            if result != test.expect {
                t.Errorf("findDate(%v, %d, %v) = %v, want %v", test.r, test.amount, test.date, result, test.expect)
            }
        })
    }
}


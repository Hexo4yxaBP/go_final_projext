package datetime

import (
	"errors"
	"time"
)

// DateDiff returns the difference between a and b expressed in the requested unit.
// unit may be one of: "days", "weeks", "months", "years".
// The result is signed: positive if b is after a, negative if b is before a.
//
// Implementation notes:
// - days/weeks are calculated using whole calendar days (UTC midnight) to avoid DST issues.
// - months/years are computed by calendar differences: full months between dates,
//   decrementing when the day-of-month of the later date is smaller than the earlier.

func DateDiff(a, b time.Time, unit string) (int, error) {
	switch unit {
	case "days":
		return int(b.Sub(a).Hours() / 24), nil
	case "weeks":
		return int(b.Sub(a).Hours()/24) / 7, nil
	case "months", "years":
		// calendar-based months/years difference
		sign := 1
		left, right := a, b
		if left.After(right) {
			sign = -1
			left, right = right, left
		}
		y1, m1, d1 := left.Date()
		y2, m2, d2 := right.Date()
		months := int((y2-y1))*12 + int(m2-m1)
		if d2 < d1 {
			months--
		}
		if unit == "months" {
			return sign * months, nil
		}
		return sign * (months / 12), nil
	default:
		return 0, errors.New("unsupported unit; use days|weeks|months|years")
	}
}

// DaysInMonth returns the number of days for the specified month in the given year.
// The parameter i represents the year (e.g., 2024) and nextMonth is a time.Month.
// It accounts for months with 31 and 30 days and uses IsLeapYear(i) to determine
// whether February has 28 or 29 days. If nextMonth is not a valid month (outside 1–12),
// the function returns 0.
func DaysInMonth(i int, nextMonth time.Month) int {
	switch nextMonth {
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	case 4, 6, 9, 11:
		return 30
	case 2:
		if IsLeapYear(i) {
			return 29
		}
		return 28
	default:
		return 0
	}
}

// IsLeapYear reports whether the given year i is a leap year according to the Gregorian calendar.
// It returns true when i is divisible by 4, except when divisible by 100 (unless also divisible by 400).
func IsLeapYear(i int) bool {
	if i%4 != 0 {
		return false
	}
	if i%100 != 0 {
		return true
	}
	if i%400 != 0 {
		return false
	}
	return true
}

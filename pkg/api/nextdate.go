package api

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Hexo4yxaBP/go_final_projext/pkg/datetime"
)

// GetNextDateHandler handles HTTP requests to calculate the next occurrence of a date
// according to a start date and a repetition rule.
//
// Expected input (via form values or URL query parameters):
//   - "now"   (optional): reference date in "YYYYMMDD" format (Go layout "20060102").
//     If omitted, the current server time (time.Now()) is used.
//   - "date"  (required for nextDate): start date string passed to nextDate.
//   - "repeat" (required for nextDate): repetition rule string passed to nextDate.
//
// Behavior:
//   - Parses the "now" parameter if provided; returns HTTP 400 Bad Request with the
//     message "Invalid 'now' date format" if parsing fails.
//   - Forwards "date" and "repeat" along with the resolved reference time to nextDate.
//   - If nextDate returns an error, responds with HTTP 500 Internal Server Error and
//     the message "Error calculating next date: <err>".
//   - On success, writes the returned result string to the response body (status 200).
//
// Notes:
//   - The handler does not set a specific Content-Type header; it writes the raw result
//     bytes to the response.
//   - r.FormValue is used to read parameters, so values can come from URL query parameters
//     or form-encoded request bodies.
func GetNextDateHandler(w http.ResponseWriter, r *http.Request) {

	now := r.FormValue("now")

	var nowT time.Time
	if now == "" {
		nowT = time.Now()
	} else {
		var err error
		nowT, err = time.Parse("20060102", now)
		if err != nil {
			http.Error(w, "Invalid 'now' date format", http.StatusBadRequest)
			return
		}
	}

	dstart := r.FormValue("date")
	repeat := r.FormValue("repeat")
	res, err := nextDate(nowT, dstart, repeat)
	if err != nil {
		http.Error(w, "Error calculating next date: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(res))
}

// nextDate returns the next occurrence date (formatted as "20060102") for an event
// that began on dstart and repeats according to the repeat specification. The
// computation is done relative to the provided now time.
//
// Parameters:
//   - now: reference time used to compute the next occurrence (time zone of the
//     resulting date is taken from now).
//   - dstart: start date of the event, expected in "YYYYMMDD" format (layout
//     "20060102"). If this cannot be parsed an error is returned.
//   - repeat: repetition specification parsed by parseRepeat; the first character
//     (date prefix) must be one of:
//   - "d" — daily repetition
//   - "w" — weekly repetition
//   - "m" — monthly repetition
//   - "y" — yearly repetition
//     parseRepeat is expected to return (datePrefix, countPeriod, optMonths, err)
//     where countPeriod and optMonths drive behavior described below.
//
// Behavior summary:
//   - If now is strictly before the parsed start date, the function uses the
//     start date as the reference moment; otherwise it uses now.
//   - All returned dates are constructed at 00:00:00 in the reference time's
//     location and formatted as "YYYYMMDD".
//
// Per-prefix semantics:
//
//	d (daily):
//	  - countPeriod[0] is treated as the interval in days.
//	  - The function computes the number of days elapsed since start and
//	    advances to the next occurrence by adding the minimal positive number of
//	    days required to reach the next scheduled date. Note: in the current
//	    implementation, when the reference time already falls exactly on a
//	    scheduled day the result advances by a full period (i.e., returns the
//	    next occurrence after the reference).
//
//	w (weekly):
//	  - countPeriod contains sorted weekday codes (expected as integers 1..7).
//	  - The function selects the next weekday value from countPeriod that is
//	    strictly greater than the reference weekday; if none is greater it wraps
//	    to the first element in countPeriod (advancing into the next week).
//	  - Values in countPeriod outside 1..7 cause an error.
//
//	m (monthly):
//	  - countPeriod contains day-of-month values. Positive values are absolute
//	    days (1..31). Negative values are offsets from the end of the month
//	    (e.g. -1 denotes the last day).
//	  - The function picks the first day in countPeriod that comes after the
//	    reference day in the current month; if none is found it advances to the
//	    next month and uses the first value from countPeriod.
//	  - If the chosen day does not exist in the candidate month (e.g. 30 in
//	    February), the algorithm advances to the next month.
//	  - optMonths, if non-empty, restricts allowed months: the algorithm will
//	    choose the next month from optMonths that is >= the candidate month; if
//	    none is >= candidate it wraps to the first optMonths entry and advances
//	    the year.
//
//	y (yearly):
//	  - The function schedules occurrences on the same month and day as dstart.
//	  - It returns the occurrence in the current year if it is strictly after
//	    the reference; otherwise it returns the occurrence in the next year.
//
// Errors:
//   - Invalid dstart format (parsing failure).
//   - Any error returned by parseRepeat.
//   - For weekly repetition, weekday values outside the allowed range return an
//     error.
//   - An unsupported repeat prefix returns an error indicating supported values.
//
// Return values:
//   - (string, error): formatted next date "YYYYMMDD" and nil error on success;
//     otherwise an empty string and a descriptive error.
func nextDate(now time.Time, dstart string, repeat string) (string, error) {
	startT, err := time.Parse("20060102", dstart)

	if err != nil {
		return "", err
	}

	var resTime time.Time
	if now.Before(startT) {
		resTime = startT
	} else {
		resTime = now
	}

	var datePrefix string
	var countPeriod []int
	var optMonths []int

	datePrefix, countPeriod, optMonths, err = parseRepeat(repeat)
	if err != nil {
		return "", err
	}

	switch datePrefix {
	case "d":
		daysDiff, err := datetime.DateDiff(startT, resTime, "days")

		if err != nil {
			return "", err
		}

		//----- test logic
		/*if daysDiff == 0 {
			return startT.Format("20060102"), nil
		}*/
		//-----
		daysToAdd := countPeriod[0] - (daysDiff % countPeriod[0])
		nextDate := resTime.AddDate(0, 0, daysToAdd)

		return nextDate.Format("20060102"), nil
	case "w":

		dOw := resTime.Weekday()
		nextDoW := 0

		for _, v := range countPeriod {

			if v < 1 || v > 7 {
				return "", fmt.Errorf("day of week out of range: %d", v)
			}

			// reminder: countPeriod is sorted, so we can use the construction below
			if int(dOw) < v {
				nextDoW = v
				break
			}

		}

		// reminder: countPeriod is sorted, so we can use the construction below
		if nextDoW == 0 {
			nextDoW = countPeriod[0]
		}

		daysToAdd := nextDoW - int(dOw)

		if daysToAdd < 0 {
			daysToAdd += 7
		}

		nextDate := resTime.AddDate(0, 0, daysToAdd)

		return nextDate.Format("20060102"), nil
	case "m":

		nowDoM := resTime.Day()
		nowMonth := resTime.Month()
		nextMonth := nowMonth

		// find next day in month to set
		foundDay := false
		var dayToSet int

		for _, d := range countPeriod {
			if nowDoM < d || (d < 0 && nowDoM < datetime.DaysInMonth(resTime.Year(), nextMonth)+d+1) {
				if d > 0 {
					dayToSet = d
				} else {
					dayToSet = datetime.DaysInMonth(resTime.Year(), nextMonth) + d + 1
				}
				foundDay = true
				break
			}
		}

		if !foundDay {
			dayToSet = countPeriod[0]
			nextMonth += 1
			if nextMonth > 12 {
				nextMonth = 1
				resTime = resTime.AddDate(1, 0, 0)
			}
		} else if dayToSet > datetime.DaysInMonth(resTime.Year(), nextMonth) {
			nextMonth += 1
		}

		// handle optional months list
		if len(optMonths) > 0 {
			foundMonth := false

			for _, m := range optMonths {
				if int(nextMonth) <= m {
					nextMonth = time.Month(m)
					foundMonth = true
					break
				}
			}

			if !foundMonth {
				nextMonth = time.Month(optMonths[0])
				resTime = resTime.AddDate(1, 0, 0)
			}
		}

		nextDate := time.Date(resTime.Year(), nextMonth, dayToSet, 0, 0, 0, 0, resTime.Location())

		return nextDate.Format("20060102"), nil
	case "y":

		nextDate := time.Date(resTime.Year(), startT.Month(), startT.Day(), 0, 0, 0, 0, resTime.Location())
		if !nextDate.After(resTime) {
			nextDate = nextDate.AddDate(1, 0, 0)
		}
		return nextDate.Format("20060102"), nil
	default:
		return "", errors.New("unsupported repeat prefix; use d|w|m|y")
	}

}

func parseRepeat(repeat string) (string, []int, []int, error) {
	repeat = strings.TrimSpace(repeat)
	if repeat == "" {
		return "", nil, nil, errors.New("empty repeat")
	}

	parts := strings.Fields(repeat)
	if len(parts) < 2 && parts[0] != "y" {
		return "", nil, nil, errors.New("invalid repeat format; expected: <prefix> <days[,days...]> [months[,months...]]")
	}

	datePrefix := parts[0]
	var days []int
	var months []int
	if len(parts) > 1 {

		// days list is mandatory
		daysStr := parts[1]
		dayTokens := strings.Split(daysStr, ",")

		for _, tok := range dayTokens {
			tok = strings.TrimSpace(tok)
			if tok == "" {
				continue
			}
			v, err := strconv.Atoi(tok)
			if err != nil {
				return "", nil, nil, fmt.Errorf("invalid day %q: %w", tok, err)
			}
			if !((v >= 1 && v <= 31) || v == -1 || v == -2) {
				return "", nil, nil, fmt.Errorf("day out of range: %d", v)
			}
			days = append(days, v)
		}
		if len(days) == 0 {
			return "", nil, nil, errors.New("no valid days provided")
		}

		// optional months list
		if len(parts) >= 3 {
			monthsStr := parts[2]
			monthTokens := strings.Split(monthsStr, ",")
			for _, tok := range monthTokens {
				tok = strings.TrimSpace(tok)
				if tok == "" {
					continue
				}
				v, err := strconv.Atoi(tok)
				if err != nil {
					return "", nil, nil, fmt.Errorf("invalid month %q: %w", tok, err)
				}
				if v < 1 || v > 12 {
					return "", nil, nil, fmt.Errorf("month out of range: %d", v)
				}
				months = append(months, v)
			}
		}
	}

	//sort slices
	days = sortPN(days)
	sort.Ints(months)

	return datePrefix, days, months, nil
}

// slice special sorting function
// sorts the slice in ascending order, positive elements first, then negative elements
func sortPN(src []int) []int {

	// create two slices
	var nonNegative []int
	var negative []int

	for _, num := range src {
		if num >= 0 {
			nonNegative = append(nonNegative, num)
		} else {
			negative = append(negative, num)
		}
	}

	// sort non-negative numbers in ascending order
	sort.Ints(nonNegative)
	// sort negatives in ascending order (from smallest to largest)
	sort.Ints(negative)

	// concatenate: first nonNegative, then negative
	return append(nonNegative, negative...)
}

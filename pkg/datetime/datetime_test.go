package datetime

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func parseDate(s string) time.Time {
	t, _ := time.Parse("2006-01-02", s)
	return t
}

func TestDiffBasic(t *testing.T) {
	tests := []struct {
		a, b string
		unit string
		want int
	}{
		{"2025-01-01", "2025-01-10", "days", 9},
		{"2025-01-01", "2025-01-15", "weeks", 2},
		{"2025-01-15", "2025-04-14", "months", 2},
		{"2019-03-01", "2022-03-01", "years", 3},
		{"2025-01-10", "2025-01-01", "days", -9},
	}

	for _, tt := range tests {
		a := parseDate(tt.a)
		b := parseDate(tt.b)
		got, err := DateDiff(a, b, tt.unit)
		if err != nil {
			t.Fatalf("unexpected error for %v: %v", tt, err)
		}
		if got != tt.want {
			t.Fatalf("Diff(%s,%s,%s) = %d; want %d", tt.a, tt.b, tt.unit, got, tt.want)
		}
	}
}

func TestDiffInvalidUnit(t *testing.T) {
	_, err := DateDiff(parseDate("2025-01-01"), parseDate("2025-01-02"), "hours")
	if err == nil {
		t.Fatalf("expected error for unsupported unit")
	}
}

func TestIsLeapYear(t *testing.T) {
	tests := []struct {
		year int
		want bool
	}{
		{2000, true},  // divisible by 400
		{1900, false}, // divisible by 100 but not 400
		{1996, true},  // divisible by 4
		{1999, false}, // common year
		{2400, true},  // divisible by 400
		{2100, false}, // divisible by 100 but not 400
		{0, true},     // year 0 is divisible by 400
		{-4, true},    // negative year divisible by 4
		{-100, false}, // negative year divisible by 100 but not 400
	}

	for _, tc := range tests {
		got := IsLeapYear(tc.year)
		if got != tc.want {
			t.Fatalf("isLeapYear(%d) = %v; want %v", tc.year, got, tc.want)
		}
	}
}
func TestDaysInMonth(t *testing.T) {
	tests := []struct {
		year int
		mon  time.Month
		want int
	}{
		// 31-day months
		{2025, 1, 31},
		{2025, 3, 31},
		{2025, 5, 31},
		{2025, 7, 31},
		{2025, 8, 31},
		{2025, 10, 31},
		{2025, 12, 31},
		// 30-day months
		{2025, 4, 30},
		{2025, 6, 30},
		{2025, 9, 30},
		{2025, 11, 30},
		// February non-leap
		{2023, 2, 28},
		// February leap
		{2024, 2, 29},
		// year edge cases (negative / zero year are supported by isLeapYear)
		{0, 2, 29},    // year 0 divisible by 400 -> leap
		{-4, 2, 29},   // negative year divisible by 4 -> leap
		{-100, 2, 28}, // negative year divisible by 100 but not 400 -> not leap
	}

	for _, tc := range tests {
		name := tc.mon.String() + "_" + strconv.Itoa(tc.year)
		t.Run(name, func(t *testing.T) {
			got := DaysInMonth(tc.year, tc.mon)
			if got != tc.want {
				t.Fatalf("daysInMonth(%d, %d) = %d; want %d", tc.year, tc.mon, got, tc.want)
			}
		})
	}
}

func TestDaysInMonthInvalidMonth(t *testing.T) {
	tests := []struct {
		year int
		mon  time.Month
	}{
		{2025, 0},
		{2025, 13},
		{2025, time.Month(255)},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("m=%d", tc.mon), func(t *testing.T) {
			got := DaysInMonth(tc.year, tc.mon)
			if got != 0 {
				t.Fatalf("daysInMonth(%d, %d) = %d; want 0 for invalid month", tc.year, tc.mon, got)
			}
		})
	}
}

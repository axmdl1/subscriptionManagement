package utils

import (
	"errors"
	"fmt"
	"time"
)

const monthYearLayout = "01-2006"

func ParseMonthYear(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, errors.New("empty date")
	}

	t, err := time.Parse(monthYearLayout, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format, expected MM-YYYY: %w", err)
	}

	return t, nil
}

func FormatMonthYear(t time.Time) string {
	return t.Format(monthYearLayout)
}

func MonthsBetween(start, end time.Time) int {
	if end.Before(start) {
		return 0
	}

	yearDiff := end.Year() - start.Year()
	monthDiff := int(end.Month()) - int(start.Month())

	return yearDiff*12 + monthDiff + 1
}

func OverlapPeriod(aStart, aEnd, bStart, bEnd time.Time) (time.Time, time.Time, bool) {
	if aEnd.Before(bStart) || bEnd.Before(aStart) {
		return time.Time{}, time.Time{}, false
	}

	start := maxTime(aStart, bStart)
	end := minTime(aEnd, bEnd)

	return start, end, true
}

func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func minTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

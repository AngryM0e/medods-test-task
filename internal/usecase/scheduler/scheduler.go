package scheduler

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

func NextDate(now time.Time, dateStr string, repeatType string) (string, error) {
	if repeatType == "" {
		return "", fmt.Errorf("empty repeat type")
	}

	date, err := time.Parse(DateFormat, dateStr)
	if err != nil {
		return "", fmt.Errorf("invalid date format: %w", err)
	}

	parts := strings.Fields(repeatType)
	repeat := parts[0]

	switch repeat {
	case "daily":
		return handleDaily(now, date, parts)
	case "monthly":
		return handleMonthly(now, date, parts)
	case "even":
		return handleEvenOdd(now, date, true)
	case "odd":
		return handleEvenOdd(now, date, false)
	case "specific":
		return handleSpecificDates(now, parts)
	default:
		return "", fmt.Errorf("unsupported repeat type: %s", repeat)
	}
}

func handleEvenOdd(now time.Time, date time.Time, even bool) (string, error) {
	start := now
	if date.After(now) {
		start = date
	}

	next := start.AddDate(0, 0, 1)

	// Проверяем ближайшие дни
	for i := 0; i < 10; i++ {
		isEven := next.Day()%2 == 0
		if (even && isEven) || (!even && !isEven) {
			return next.Format(DateFormat), nil
		}
		next = next.AddDate(0, 0, 1)
	}
	return "", fmt.Errorf("could not calculate even/odd date")
}

func handleSpecificDates(now time.Time, parts []string) (string, error) {
	if len(parts) < 2 {
		return "", fmt.Errorf("dates repeat type requires a list of dates")
	}

	dates := strings.Split(parts[1], ",")
	var validDates []time.Time

	for _, date := range dates {
		d, err := time.Parse(DateFormat, strings.TrimSpace(date))
		if err != nil {
			continue
		}
		if d.After(now) {
			validDates = append(validDates, d)
		}
	}

	if len(validDates) == 0 {
		return "", fmt.Errorf("no future dates found in the list")
	}

	sort.Slice(validDates, func(i, j int) bool {
		return validDates[i].Before(validDates[j])
	})

	return validDates[0].Format(DateFormat), nil
}

func handleMonthly(now time.Time, date time.Time, parts []string) (string, error) {
	if len(parts) < 2 {
		return "", fmt.Errorf("missing days for monthle repeat type")
	}

	days := strings.Split(parts[1], ",")
	daysMap := make(map[int]bool)
	for _, day := range days {
		d, err := strconv.Atoi(strings.TrimSpace(day))
		if err != nil {
			return "", fmt.Errorf("invalid day: %s", day)
		}
		daysMap[d] = true
	}

	monthsMap := make(map[int]bool)
	if len(parts) >= 3 {
		months := strings.Split(parts[2], ",")
		for _, month := range months {
			m, err := strconv.Atoi(strings.TrimSpace(month))
			if err != nil || m < 1 || m > 12 {
				return "", fmt.Errorf("invalid month: %s", month)
			}
			monthsMap[m] = true
		}
	}

	start := now
	if date.After(now) {
		start = date
	}

	for i := 1; i <= 365; i++ {
		checkDate := start.AddDate(0, 0, i)

		if !checkDate.After(date) {
			continue
		}

		if len(monthsMap) > 0 && !monthsMap[int(checkDate.Month())] {
			continue
		}

		if isDayValid(checkDate, daysMap) {
			return checkDate.Format(DateFormat), nil
		}
	}

	return "", fmt.Errorf("no suitable date found within 2 years")
}

func isDayValid(t time.Time, daysMap map[int]bool) bool {
	if daysMap[t.Day()] {
		return true
	}

	lastDayOfMonth := time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, t.Location()).Day()

	if daysMap[-1] && t.Day() == lastDayOfMonth {
		return true
	}
	if daysMap[-2] && t.Day() == lastDayOfMonth-1 {
		return true
	}

	return false
}

func handleDaily(now time.Time, date time.Time, parts []string) (string, error) {
	if len(parts) < 2 {
		return "", fmt.Errorf("missing days interval for daily repeat type")
	}

	days, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", fmt.Errorf("invalid days interval: %s", parts[1])
	}

	if days < 1 || days > 400 {
		return "", fmt.Errorf("days interval must be between 1 and 400, got %d", days)
	}

	next := date

	next = next.AddDate(0, 0, days)

	for !next.After(now) {
		next = next.AddDate(0, 0, days)
	}

	return next.Format(DateFormat), nil
}

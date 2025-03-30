package datetools

import (
	"regexp"
	"strconv"
	"time"
)

func ParseRelativeDate(dateString string) (time.Time, bool) {
	patterns := map[string]string{
		"days":   `(\d+)\s+days?\s+ago`,
		"weeks":  `(\d+)\s+weeks?\s+ago`,
		"months": `(\d+)\s+months?\s+ago`,
	}

	for unit, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(dateString)

		if len(matches) > 1 {
			value, err := strconv.Atoi(matches[1])
			if err != nil {
				continue
			}

			today := time.Now()
			var pastDate time.Time

			switch unit {
			case "days":
				pastDate = today.AddDate(0, 0, -value)
			case "weeks":
				pastDate = today.AddDate(0, 0, -value*7)
			case "months":
				pastDate = today.AddDate(0, -value, 0)
			}

			return pastDate, true
		}
	}

	return time.Time{}, false
}

func DaysDifference(date1, date2 time.Time) int16 {
	// Narmalizing input (removing time)
	normalizedDate1 := time.Date(date1.Year(), date1.Month(), date1.Day(), 0, 0, 0, 0, time.UTC)
	normalizedDate2 := time.Date(date2.Year(), date2.Month(), date2.Day(), 0, 0, 0, 0, time.UTC)

	// Calculate difference between normalized dates
	diff := normalizedDate2.Sub(normalizedDate1)

	// Convert difference into days
	days := int16(diff.Hours() / 24)

	return days
}

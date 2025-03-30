package datetools

import (
	"regexp"
	"strconv"
	"time"
)

// ParseRelativeDate parst einen relativen Datumsstring wie "Updated 2 days ago"
// und gibt ein time.Time-Objekt zurück.
// Verarbeitet Tage, Wochen und Monate.
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

	return time.Time{}, false // Kein passendes Muster gefunden
}

// DaysDifference berechnet die Differenz in Tagen zwischen zwei Zeitpunkten,
// wobei die Uhrzeitkomponente ignoriert wird.
// Es werden nur die Datumsteile (Jahr, Monat, Tag) berücksichtigt.
func DaysDifference(date1, date2 time.Time) int16 {
	// Normalisiere die Zeitpunkte, um nur das Datum (Jahr, Monat, Tag) zu behalten
	// und die Uhrzeitkomponente zu eliminieren
	normalizedDate1 := time.Date(date1.Year(), date1.Month(), date1.Day(), 0, 0, 0, 0, time.UTC)
	normalizedDate2 := time.Date(date2.Year(), date2.Month(), date2.Day(), 0, 0, 0, 0, time.UTC)

	// Berechne die Differenz zwischen den normalisierten Zeitpunkten
	diff := normalizedDate2.Sub(normalizedDate1)

	// Konvertiere die Differenz in Tage
	days := int16(diff.Hours() / 24)

	return days
}

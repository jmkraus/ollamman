package datetools

import (
	"testing"
	"time"
)

func TestParseRelativeDate(t *testing.T) {
	// Fixierte Zeit für Tests
	fixedTime := time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)
	now = func() time.Time {
		return fixedTime
	}

	tests := []struct {
		input    string
		expected time.Time
		valid    bool
	}{
		{"3 days ago", fixedTime.AddDate(0, 0, -3), true},
		{"2 weeks ago", fixedTime.AddDate(0, 0, -14), true},
		{"1 month ago", fixedTime.AddDate(0, -1, 0), true},
		{"yesterday", fixedTime.AddDate(0, 0, -1), true},
		{"invalid input", time.Time{}, false},
	}

	for _, test := range tests {
		result, valid := ParseRelativeDate(test.input)
		if valid != test.valid || !result.Equal(test.expected) {
			t.Errorf("ParseRelativeDate(%q) = %v, %v; want %v, %v", test.input, result, valid, test.expected, test.valid)
		}
	}
}

func TestDaysDifference(t *testing.T) {
	// Fixierte Zeiten für Tests
	date1 := time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(2023, 10, 10, 0, 0, 0, 0, time.UTC)
	date3 := time.Date(2023, 9, 30, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		date1    time.Time
		date2    time.Time
		expected int16
	}{
		{date1, date2, 9},  // 9 Tage Unterschied
		{date2, date1, -9}, // -9 Tage Unterschied
		{date1, date1, 0},  // 0 Tage Unterschied
		{date1, date3, -1}, // -1 Tag Unterschied
	}

	for _, test := range tests {
		result := DaysDifference(test.date1, test.date2)
		if result != test.expected {
			t.Errorf("DaysDifference(%v, %v) = %v; want %v", test.date1, test.date2, result, test.expected)
		}
	}
}

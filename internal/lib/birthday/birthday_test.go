package birthday_test

import (
	"github.com/k6mil6/birthday-notificator/internal/lib/birthday"
	"testing"
	"time"
)

type birthdayTestCase struct {
	birthday time.Time
	expected time.Time
}

var birthdayTestCases = []birthdayTestCase{
	{
		birthday: time.Date(2002, time.January, 2, 0, 0, 0, 0, time.Local),
		expected: time.Date(2025, time.January, 2, 0, 0, 0, 0, time.Local),
	},
	{
		birthday: time.Date(2004, time.September, 2, 0, 0, 0, 0, time.Local),
		expected: time.Date(2024, time.September, 2, 0, 0, 0, 0, time.Local),
	},
}

func TestCalculateNextBirthday(t *testing.T) {
	for _, tc := range birthdayTestCases {
		t.Run(tc.birthday.String(), func(t *testing.T) {
			result := birthday.CalculateNextBirthday(tc.birthday)
			expectedUTC := tc.expected.UTC()
			if result.Year() != expectedUTC.Year() || result.Month() != expectedUTC.Month() || result.Day() != expectedUTC.Day() {
				t.Errorf("expected %v, got %v", expectedUTC, result)
			}
		})
	}
}

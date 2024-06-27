package offset_test

import (
	"github.com/k6mil6/birthday-notificator/internal/lib/notification/offset"
	"testing"
	"time"
)

type offsetTestCase struct {
	offset   offset.Offset
	expected time.Duration
}

var offsetTestCases = []offsetTestCase{
	{
		offset:   offset.Offset{Unit: "day", Value: 1},
		expected: 24 * time.Hour,
	},
	{
		offset:   offset.Offset{Unit: "hour", Value: 1},
		expected: time.Hour,
	},
	{
		offset:   offset.Offset{Unit: "minute", Value: 1},
		expected: time.Minute,
	},
	{
		offset:   offset.Offset{Unit: "second", Value: 1},
		expected: time.Second,
	},
	{
		offset:   offset.Offset{Unit: "day", Value: 0},
		expected: 0,
	},
}

func TestConvertToTimeDuration(t *testing.T) {
	for _, tc := range offsetTestCases {
		t.Run(tc.offset.Unit, func(t *testing.T) {
			result := offset.ConvertToTimeDuration(tc.offset)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

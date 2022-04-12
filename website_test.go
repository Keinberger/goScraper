package scraper

import (
	"testing"
	"time"

	"github.com/influxdata/influxdb/pkg/testing/assert"
)

func TestMinUntilUpdate(t *testing.T) {
	t.Run("12:30 | 18-90", func(t *testing.T) {
		w := Website{
			HourStartScraping: 18,
			HourStopScraping:  19,
		}
		assert.Equal(t, w.MinUntil(time.Date(2020, 10, 3, 12, 30, 0, 0, time.UTC)), time.Duration(330)*time.Minute)
	})

	t.Run("12:30 | 8-9", func(t *testing.T) {
		w := Website{
			HourStartScraping: 8,
			HourStopScraping:  9,
		}
		assert.Equal(t, w.MinUntil(time.Date(2020, 10, 3, 12, 30, 0, 0, time.UTC)), time.Duration(1170)*time.Minute)
	})

	t.Run("22:30 | 8-9", func(t *testing.T) {
		w := Website{
			HourStartScraping: 8,
			HourStopScraping:  9,
		}
		assert.Equal(t, w.MinUntil(time.Date(2020, 10, 3, 22, 30, 0, 0, time.UTC)), time.Duration(570)*time.Minute)
	})

	t.Run("5:45 | 12-13", func(t *testing.T) {
		w := Website{
			HourStartScraping: 12,
			HourStopScraping:  13,
		}
		assert.Equal(t, w.MinUntil(time.Date(2020, 10, 3, 5, 45, 0, 0, time.UTC)), time.Duration(375)*time.Minute)
	})
}

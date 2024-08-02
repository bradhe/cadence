package cadence

import (
	"time"

	"github.com/robfig/cron"
)

var cronParser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)

// Next calculates the next event in the series based on when the last event
// happened.
func Next(pattern string, last time.Time) (time.Time, error) {
	if spec, err := cronParser.Parse(pattern); err != nil {
		return time.Time{}, err
	} else {
		return spec.Next(last), nil
	}
}

package cadence

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/robfig/cron"
)

var cronParser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)

type interval string

const (
	second = interval("second")
	minute = interval("minute")
	hour   = interval("hour")
	day    = interval("day")
	week   = interval("week")
	month  = interval("month")
	year   = interval("year")
)

type englishPattern struct {
	Number   int
	Interval interval
}

func (e englishPattern) ToCrontab() string {
	switch e.Interval {
	case second:
		return fmt.Sprintf("*/%d * * * * *", e.Number)
	case minute:
		return fmt.Sprintf("* */%d * * * *", e.Number)
	case hour:
		return fmt.Sprintf("0 0 */%d * * *", e.Number)
	case week:
		return fmt.Sprintf("* * * */%d * *", (e.Number * 7))
	case month:
		return fmt.Sprintf("* * * * */%d *", e.Number)
	case year:
		return fmt.Sprintf("* * * * */%d *", (e.Number * 12))
	default:
		return "* * * * * *"
	}
}

func parseEnglishPattern(pattern string) (*englishPattern, error) {
	orig := pattern
	pattern = strings.TrimSpace(pattern)
	pattern = strings.ToLower(pattern)

	if !strings.HasPrefix(pattern, "every") {
		return nil, fmt.Errorf("invalid prefix for pattern `%s`", orig)
	}

	// there should be a number as the next character.
	pattern = strings.TrimPrefix(pattern, "every")
	pattern = strings.TrimSpace(pattern)

	var bld strings.Builder

	for _, r := range pattern {
		if unicode.IsDigit(r) {
			bld.WriteRune(r)
		} else {
			break
		}
	}

	i, _ := strconv.Atoi(bld.String())

	var num int

	if i == 0 {
		num = 1
	} else {
		num = i
	}

	pattern = strings.TrimLeftFunc(pattern, unicode.IsDigit)
	pattern = strings.TrimSpace(pattern)

	suffixes := map[string]interval{
		"second":  second,
		"seconds": second,
		"minute":  minute,
		"minutes": minute,
		"hour":    hour,
		"hours":   hour,
		"day":     day,
		"days":    day,
		"week":    week,
		"weeks":   week,
		"month":   month,
		"months":  month,
		"year":    year,
		"years":   year,
	}

	if i, ok := suffixes[pattern]; ok {
		return &englishPattern{
			Number:   num,
			Interval: i,
		}, nil
	} else {
		return nil, fmt.Errorf("couldn't define interval for patten `%s`", orig)
	}
}

func isValidEnglishPattern(pattern string) bool {
	if _, err := parseEnglishPattern(pattern); err != nil {
		return false
	} else {
		return true
	}
}

// Next uses the supplied pattern to determine when the next occurance of the
// event should be.
func Next(pattern string, last time.Time) (time.Time, error) {
	if spec, err := parseEnglishPattern(pattern); err == nil {
		if spec, err := cronParser.Parse(spec.ToCrontab()); err != nil {
			return time.Time{}, err
		} else {
			return spec.Next(last), nil
		}
	}

	if spec, err := cronParser.Parse(pattern); err != nil {
		return time.Time{}, err
	} else {
		return spec.Next(last), nil
	}
}

// IsValid will tell you if the pattern can be parsed by cadence.
func IsValid(pattern string) bool {
	if _, err := Next(pattern, time.Now()); err != nil {
		return false
	} else {
		return true
	}
}

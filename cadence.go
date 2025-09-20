package cadence

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/robfig/cron"
)

var cronParser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)

type timeUnit string

const (
	second = timeUnit("second")
	minute = timeUnit("minute")
	hour   = timeUnit("hour")
	day    = timeUnit("day")
	week   = timeUnit("week")
	month  = timeUnit("month")
	year   = timeUnit("year")
)

type interval struct {
	Number   int
	TimeUnit timeUnit
}

func (e interval) ToCrontab() string {
	switch e.TimeUnit {
	case second:
		return fmt.Sprintf("*/%d * * * * *", e.Number)
	case minute:
		return fmt.Sprintf("0 */%d * * * *", e.Number)
	case hour:
		return fmt.Sprintf("0 0 */%d * * *", e.Number)
	case day:
		return fmt.Sprintf("0 0 0 */%d * *", e.Number)
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

func parseEnglishPattern(pattern string) (*interval, error) {
	pattern = strings.ToLower(pattern) // case insensitive
	tokens := strings.Fields(pattern)
	var num int
	var tUnit timeUnit
	var timeUnitToken string

	// structure sanity check. There should be 3 tokens (or 2 if a "1" is implied)
	// in the form of "every <number> <time unit>".
	// We need to parse the number if it exists and retrieve the time unit.
	if len(tokens) == 3 {
		if n, err := strconv.Atoi(tokens[1]); err != nil {
			return nil, fmt.Errorf("invalid pattern: `%s` cannot be parsed as a number", tokens[1])
		} else {
			num = n
			timeUnitToken = tokens[2]
		}
	} else if len(tokens) == 2 {
		// there is no number so 1 is implied. Time unit should be the second token
		num = 1
		timeUnitToken = tokens[1]
	} else {
		return nil, fmt.Errorf("invalid number of tokens for pattern `%s`, expecting 2 or 3 tokens, got %d", pattern, len(tokens))
	}

	if tokens[0] != "every" {
		return nil, fmt.Errorf("invalid prefix for pattern `%s`, all patterns must start with \"every\"", pattern)
	}

	timeUnits := map[string]timeUnit{
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

	if unit, ok := timeUnits[timeUnitToken]; !ok {
		return nil, fmt.Errorf("invalid time unit `%s` for pattern `%s`", timeUnitToken, pattern)
	} else if num == 1 && isPlural(timeUnitToken) {
		// user supplied sth like "every (1) months", probable bug on their side, reject it.
		return nil, fmt.Errorf("invalid number and time unit combination. Number is 1 and time unit is in plural")
	} else if num > 1 && !isPlural(timeUnitToken) {
		// the inverse of before
		return nil, fmt.Errorf("invalid number and time unit combination. Number is > 1 and time unit is in singular")
	} else {
		tUnit = unit
	}

	return &interval{
		Number:   num,
		TimeUnit: tUnit,
	}, nil
}

func isPlural(s string) bool {
	return strings.HasSuffix(s, "s")
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
	// we lop off the second because that's the smallest interval we can align
	// to as per the semantics of the library.
	last = last.Truncate(time.Second)

	if spec, err := parseEnglishPattern(pattern); err == nil {
		if spec, err := cronParser.Parse(spec.ToCrontab()); err != nil {
			return time.Time{}, err
		} else {
			return spec.Next(last), nil
		}
	} else {
		if spec, err := cronParser.Parse(pattern); err != nil {
			return time.Time{}, err
		} else {
			next := spec.Next(last)

			return next, nil
		}
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

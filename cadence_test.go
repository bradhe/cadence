package cadence

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNext(t *testing.T) {
	t.Run("patterns with seconds", func(t *testing.T) {
		t.Run("every second", func(t *testing.T) {
			start := time.Now()
			next, err := Next("* * * * * *", start)
			assert.NoError(t, err)
			assert.NotEmpty(t, next)
			assert.Equal(t, start.Truncate(time.Second).Add(time.Second).Unix(), next.Unix())
		})

		t.Run("every 5 seconds", func(t *testing.T) {
			start := time.Now().Truncate(time.Nanosecond)
			next, err := Next("*/5 * * * * *", start)
			assert.NoError(t, err)
			assert.NotEmpty(t, next)
			assert.Equal(t, start.Truncate(5*time.Second).Add(5*time.Second), next)
		})

		t.Run("on the dot", func(t *testing.T) {
			start := MustParse(t, "2021-01-01T12:00:05Z")
			next, err := Next("*/5 * * * * *", start)
			assert.NoError(t, err)
			assert.NotEmpty(t, next)
			assert.Equal(t, MustParse(t, "2021-01-01T12:00:10Z"), next)
		})

		t.Run("every 5th second", func(t *testing.T) {
			start := MustParse(t, "2021-01-01T12:00:05Z")
			next, err := Next("5 * * * * *", start)
			assert.NoError(t, err)
			assert.NotEmpty(t, next)
			assert.Equal(t, MustParse(t, "2021-01-01T12:01:05Z"), next)
		})

		t.Run("every 5th second on Tuesday", func(t *testing.T) {
			next, err := Next("5 * * * * 2", time.Now().Add(-1*time.Second))
			assert.NoError(t, err)
			assert.NotEmpty(t, next)
		})

		t.Run("every 5th second on Tuesday", func(t *testing.T) {
			next, err := Next("*/5 * * */21 * 2", time.Now().Add(-1*time.Second))
			assert.NoError(t, err)
			assert.NotEmpty(t, next)
		})
	})

	t.Run("human readable", func(t *testing.T) {
		_, err := Next("every 1 hour", time.Now())
		assert.NoError(t, err)
	})

	t.Run("while inside an interval", func(t *testing.T) {
		start, err := time.Parse(time.RFC3339, "2021-01-01T12:00:01.001Z")
		require.NoError(t, err)

		next, err := Next("0 */1 * * *", start)
		require.NoError(t, err)
		assert.Equal(t, "2021-01-01 12:01:00 +0000 UTC", next.String())
	})

	t.Run("inside a 5 minute interval", func(t *testing.T) {
		start := MustParse(t, "2025-01-22T16:22:23.717803Z")
		next, err := Next("every 5 minutes", start)
		assert.NoError(t, err)
		assert.NotEmpty(t, next)
		assert.Equal(t, MustParse(t, "2025-01-22T16:25:00Z"), next)
	})

	t.Run("inside a 1 hour interval", func(t *testing.T) {
		start := MustParse(t, "2025-01-22T16:22:23.717803Z")
		next, err := Next("every 1 hour", start)
		assert.NoError(t, err)
		assert.NotEmpty(t, next)
		assert.Equal(t, MustParse(t, "2025-01-22T17:00:00Z"), next)
	})
}

func TestParseEnglishPattern(t *testing.T) {
	t.Run("general case", func(t *testing.T) {
		spec, err := parseEnglishPattern("every 2 years")
		require.NoError(t, err)
		assert.Equal(t, 2, spec.Number)
		assert.Equal(t, year, spec.TimeUnit)
	})

	t.Run("pattern too short", func(t *testing.T) {
		_, err := parseEnglishPattern("every ")
		require.Error(t, err)
	})

	t.Run("pattern too long", func(t *testing.T) {
		_, err := parseEnglishPattern("every 5 years now!")
		require.Error(t, err)
	})

	t.Run("pattern with an implicit 1", func(t *testing.T) {
		spec, err := parseEnglishPattern("every month")
		require.NoError(t, err)
		assert.Equal(t, 1, spec.Number)
		assert.Equal(t, month, spec.TimeUnit)
	})

	t.Run("pattern with 1 and time unit in plural", func(t *testing.T) {
		_, err := parseEnglishPattern("every 1 weeks")
		require.Error(t, err)
	})

	t.Run("pattern with > 1 and time unit in singular", func(t *testing.T) {
		_, err := parseEnglishPattern("every 2 month")
		require.Error(t, err)
	})

	t.Run("pattern with invalid number", func(t *testing.T) {
		_, err := parseEnglishPattern("every something months")
		require.Error(t, err)
	})

	t.Run("pattern with malformed number", func(t *testing.T) {
		_, err := parseEnglishPattern("every 2s months")
		require.Error(t, err)
	})

	t.Run("pattern not starting with every", func(t *testing.T) {
		_, err := parseEnglishPattern("1 year every")
		require.Error(t, err)
	})

}

func MustParse(t *testing.T, str string) time.Time {
	t.Helper()

	val, err := time.Parse(time.RFC3339, str)
	require.NoError(t, err)
	return val
}

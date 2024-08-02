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

		t.Run("every 5th second", func(t *testing.T) {
			start := time.Now()
			next, err := Next("5 * * * * *", start)
			assert.NoError(t, err)
			assert.NotEmpty(t, next)
			assert.Equal(t, start.Truncate(time.Minute).Add(time.Minute).Add(5*time.Second), next)
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
}

func TestParseEnglishPattern(t *testing.T) {
	spec, err := parseEnglishPattern("every 1 hour")
	require.NoError(t, err)
	assert.Equal(t, 1, spec.Number)
	assert.Equal(t, hour, spec.Interval)
}

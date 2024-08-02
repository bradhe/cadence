package cadence

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNext(t *testing.T) {
	t.Run("patterns with seconds", func(t *testing.T) {
		t.Run("every second", func(t *testing.T) {
			start := time.Now()
			next, err := Next("* * * * * *", start)
			assert.NoError(t, err)
			assert.NotEmpty(t, next)
			assert.Equal(t, start.Truncate(time.Second).Add(1*time.Second), next)
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
}

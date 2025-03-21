package ktime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTime(t *testing.T) {
	time.Sleep(time.Millisecond)
	now := time.Now()
	assert.True(t, now.After(initTime))
	time.Sleep(time.Millisecond)
	assert.True(t, time.Since(initTime) > 0)
}

func TestRelativeTime(t *testing.T) {
	time.Sleep(time.Millisecond)
	now := Now()
	assert.True(t, now > 0)
	time.Sleep(time.Millisecond)
	assert.True(t, Since(now) > 0)
}

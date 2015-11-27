package lib

import (
	"time"

	"github.com/cenkalti/backoff"
)

// LimitedConstantBackOff is backing off by constant interval but also has a max elapsed time
type LimitedConstantBackOff struct {
	Interval       time.Duration
	MaxElapsedTime time.Duration
	startTime      time.Time
}

// Reset resets the statTime
func (b *LimitedConstantBackOff) Reset() {
	b.startTime = time.Now()
}

// NextBackOff returns Stop if MaxElapsedTime is achieved, other returns internal.
func (b *LimitedConstantBackOff) NextBackOff() time.Duration {
	if b.GetElapsedTime() > b.MaxElapsedTime {
		return backoff.Stop
	}
	return b.Interval
}

// GetElapsedTime returns the elapsed time since startTime
func (b *LimitedConstantBackOff) GetElapsedTime() time.Duration {
	return time.Now().Sub(b.startTime)
}

// NewLimitedConstantBackOff creates LimitedExponentialBackOff
func NewLimitedConstantBackOff(interval, maxElapsedTime time.Duration) *LimitedConstantBackOff {
	return &LimitedConstantBackOff{
		Interval:       interval,
		MaxElapsedTime: maxElapsedTime,
		startTime:      time.Now(),
	}
}

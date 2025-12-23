package main

import (
	"math/rand"
	"sync"
	"time"
)

type Throttler struct {
	mu sync.Mutex

	minGlobal time.Duration
	maxGlobal time.Duration

	domainCooldown time.Duration
	lastByDomain    map[string]time.Time
}

func NewThrottler(minGlobal, maxGlobal, domainCooldown time.Duration) *Throttler {
	return &Throttler{
		minGlobal:       minGlobal,
		maxGlobal:       maxGlobal,
		domainCooldown:  domainCooldown,
		lastByDomain:    map[string]time.Time{},
	}
}

func (t *Throttler) SleepBefore(domain string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if last, ok := t.lastByDomain[domain]; ok {
		next := last.Add(t.domainCooldown)
		now := time.Now()
		if now.Before(next) {
			time.Sleep(next.Sub(now))
		}
	}

	if t.maxGlobal <= t.minGlobal {
		time.Sleep(t.minGlobal)
	} else {
		d := t.minGlobal + time.Duration(rand.Int63n(int64(t.maxGlobal-t.minGlobal)))
		time.Sleep(d)
	}

	t.lastByDomain[domain] = time.Now()
}

package cof

import (
	"maps"
	"sync"
	"time"
)

const (
	OFF                  = 0
	defaultCleanInterval = 1 * time.Minute
	defaultTTL           = 15 * time.Minute
)

type C[T any] struct {
	cache map[string]item[T]
	options
	sync.RWMutex
	stop chan bool
}

type item[T any] struct {
	value     T
	expiresOn int64
}

func (i *item[T]) isExpiredOn(timestamp int64) bool {
	return i.expiresOn > OFF && i.expiresOn <= timestamp
}

type options struct {
	cleanInterval time.Duration
	ttl           time.Duration
}

type option func(f *options)

func CleanInterval(ci time.Duration) option {
	return func(o *options) {
		o.cleanInterval = ci
	}
}

func TTL(ttl time.Duration) option {
	return func(o *options) {
		o.ttl = ttl
	}
}

func Init[T any](opts ...option) (*C[T], error) {
	o := options{
		cleanInterval: defaultCleanInterval,
		ttl:           defaultTTL,
	}

	for _, opt := range opts {
		opt(&o)
	}

	c := &C[T]{
		cache:   make(map[string]item[T]),
		stop:    make(chan bool, 1),
		options: o,
	}

	go c.cleaner()

	return c, nil
}

func (c *C[T]) Put(k string, v T) {
	c.Lock()
	defer c.Unlock()

	c.cache[k] = item[T]{
		value:     v,
		expiresOn: time.Now().Add(c.ttl).Unix(),
	}
}

func (c *C[T]) Pop(k string) (T, bool) {
	c.Lock()
	defer c.Unlock()

	v, ok := c.cache[k]
	if ok {
		delete(c.cache, k)
	}

	return v.value, ok
}

func (c *C[T]) Get(k string) (T, bool) {
	c.RLock()
	defer c.RUnlock()

	v, ok := c.cache[k]
	return v.value, ok
}

func (c *C[T]) Stop() {
	c.Lock()
	defer c.Unlock()

	select {
	case c.stop <- true:
	default:
	}

	clear(c.cache)
}

func (c *C[T]) cleaner() {
	if c.cleanInterval <= OFF {
		return
	}

	for {
		select {
		case <-c.stop:
			return
		case <-time.NewTicker(c.cleanInterval).C:
			c.cleanup()
		}
	}
}

func (c *C[T]) cleanup() {
	c.Lock()
	defer c.Unlock()

	now := time.Now().Unix()
	maps.DeleteFunc(c.cache, func(k string, v item[T]) bool {
		return v.isExpiredOn(now)
	})
}

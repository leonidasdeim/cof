package cof

import (
	"sync"
	"time"
)

const (
	defaultCleanInterval = 1 * time.Minute
)

type C[T any] struct {
	cache map[string]item[T]
	options
	sync.RWMutex
}

type item[T any] struct {
	value     T
	expiresOn int64
}

func (i *item[T]) isExpiredOn(timestamp int64) bool {
	return i.expiresOn > 0 && i.expiresOn <= timestamp
}

type options struct {
	cleanInterval time.Duration
}

type option func(f *options)

func CleanInterval(ci time.Duration) option {
	return func(o *options) {
		o.cleanInterval = ci
	}
}

func Init[T any](opts ...option) (*C[T], error) {
	o := options{
		cleanInterval: defaultCleanInterval,
	}

	for _, opt := range opts {
		opt(&o)
	}

	c := &C[T]{
		cache:   make(map[string]item[T]),
		options: o,
	}

	return c, nil
}

func (c *C[T]) cleaner() {
	for {

	}
}

package ratelimiter

import (
	"fmt"
	"sync"
	"time"
)

type Bucket struct {
	sync.Mutex
	Since    time.Time
	Attempts int
}

type BucketOption func(*Bucket)

func WithAttempts(cnt int) BucketOption {
	return func(bucket *Bucket) {
		bucket.Attempts = cnt
	}
}

func WithSince(since time.Time) BucketOption {
	return func(bucket *Bucket) {
		bucket.Since = since
	}
}

func NewBucket(opts ...BucketOption) *Bucket {
	var (
		defaultAttempts = 0
		defaultSince    = time.Now()
	)

	b := &Bucket{
		Mutex:    sync.Mutex{},
		Since:    defaultSince,
		Attempts: defaultAttempts,
	}

	for _, opt := range opts {
		opt(b)
	}

	return b
}

func (b *Bucket) Reset() {
	b.Lock()
	defer b.Unlock()
	b.Since = time.Now()
	b.Attempts = 0
}

func (b *Bucket) Incr() {
	b.Lock()
	defer b.Unlock()
	b.Attempts++
}

func (b *Bucket) IsMinuteSpent() bool {
	b.Lock()
	defer b.Unlock()
	return time.Since(b.Since) >= time.Minute
}

type Bunch struct {
	sync.RWMutex
	buckets map[string]*Bucket
}

func NewBunch() Bunch {
	return Bunch{
		RWMutex: sync.RWMutex{},
		buckets: make(map[string]*Bucket),
	}
}

func (b *Bunch) Get(key string) (*Bucket, error) {
	if bucket, ok := b.buckets[key]; ok {
		return bucket, nil
	}
	return nil, fmt.Errorf("there is no such key in the bunch")
}

func (b *Bunch) Set(key string, bucket *Bucket) {
	b.Lock()
	defer b.Unlock()
	b.buckets[key] = bucket
}

func (b *Bunch) Delete(key string) error {
	b.Lock()
	defer b.Unlock()
	if _, ok := b.buckets[key]; ok {
		delete(b.buckets, key)
		return nil
	}
	return fmt.Errorf("bucket for this key '%s' doesn't exists", key)
}

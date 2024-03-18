package limiter

import (
	"context"
	"github.com/Baraulia/anti_bruteforce_service/internal/models"
	"sync"
	"time"
)

type Limiter struct {
	sync.Mutex
	frequency int
	buckets   map[string]*models.Bucket
}

func NewLimiter(frequency int) *Limiter {
	return &Limiter{
		buckets:   make(map[string]*models.Bucket),
		frequency: frequency,
	}
}

func (l *Limiter) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(time.Duration(l.frequency) * time.Second)
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			l.clear()
		}
	}()
}

func (l *Limiter) clear() {
	l.Lock()
	defer l.Unlock()
	for key, value := range l.buckets {
		if time.Now().Sub(value.LastUpdate) > time.Duration(l.frequency)*time.Second {
			delete(l.buckets, key)
		}
	}
}

func (l *Limiter) CheckLimit(_ context.Context, key string) (int, error) {
	l.Lock()
	defer l.Unlock()
	value, ok := l.buckets[key]
	if !ok {
		l.buckets[key] = &models.Bucket{
			CurrentCount: 1,
			LastUpdate:   time.Now(),
		}
		return 1, nil
	} else {
		diffTime := int64(time.Now().Sub(value.LastUpdate).Seconds())
		switch diffTime >= int64(value.CurrentCount) {
		case true:
			value.CurrentCount = 1
			value.LastUpdate = time.Now()
			return 1, nil
		default:
			value.CurrentCount += 1
			value.CurrentCount -= int(diffTime)
			value.LastUpdate = time.Now()
			return value.CurrentCount, nil
		}
	}
}

func (l *Limiter) ClearBuckets(_ context.Context, ip, login string) error {
	l.Lock()
	defer l.Unlock()
	delete(l.buckets, ip)
	delete(l.buckets, login)

	return nil
}

func (l *Limiter) ClearAllBuckets(_ context.Context) error {
	l.buckets = make(map[string]*models.Bucket)

	return nil
}

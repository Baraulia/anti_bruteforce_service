package app

import "context"

//go:generate mockgen -source=limiterInterface.go -destination=mocks/limiter/limiter_mock.go -package=mocklimiter
type Limiter interface {
	CheckLimit(ctx context.Context, key string) (int, error)
	ClearBuckets(ctx context.Context, ip, login string) error
	ClearAllBuckets(ctx context.Context) error
	Start(ctx context.Context)
}

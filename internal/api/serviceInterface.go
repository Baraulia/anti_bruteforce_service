package api

import (
	"context"

	//nolint:depguard
	"github.com/Baraulia/anti_bruteforce_service/internal/models"
)

//go:generate mockgen -source=serviceInterface.go -destination=mocks/service_mock.go -package=mockservice
type ApplicationInterface interface {
	Check(ctx context.Context, data models.Data) (bool, error)
	AddToBlackList(ctx context.Context, ip string) error
	AddToWhiteList(ctx context.Context, ip string) error
	RemoveFromBlackList(ctx context.Context, ip string) error
	RemoveFromWhiteList(ctx context.Context, ip string) error
	ClearBuckets(ctx context.Context, data models.Data) error
	ClearAllBuckets(ctx context.Context) error
}

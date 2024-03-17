package app

import "context"

//go:generate mockgen -source=storageInterface.go -destination=mocks/storage/storage_mock.go -package=mockstorage
type Storage interface {
	CheckIPInWhiteList(ctx context.Context, ip string) (bool, error)
	CheckIPInBlackList(ctx context.Context, ip string) (bool, error)
	AddToBlackList(ctx context.Context, ip string) error
	AddToWhiteList(ctx context.Context, ip string) error
	RemoveFromBlackList(ctx context.Context, ip string) error
	RemoveFromWhiteList(ctx context.Context, ip string) error
}

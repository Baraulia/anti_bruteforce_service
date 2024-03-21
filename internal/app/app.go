package app

import (
	"context"

	"github.com/Baraulia/anti_bruteforce_service/internal/models"
)

type App struct {
	logger                Logger
	storage               Storage
	limiter               Limiter
	loginLimitAttempts    int
	passwordLimitAttempts int
	ipLimitAttempts       int
}

type Logger interface {
	Debug(msg string, fields map[string]interface{})
	Info(msg string, fields map[string]interface{})
	Warn(msg string, fields map[string]interface{})
	Error(msg string, fields map[string]interface{})
	Fatal(msg string, fields map[string]interface{})
}

func New(
	logger Logger,
	storage Storage,
	limiter Limiter,
	loginLimitAttempts, passwordLimitAttempts, ipLimitAttempts int,
) *App {
	return &App{
		logger:                logger,
		storage:               storage,
		limiter:               limiter,
		loginLimitAttempts:    loginLimitAttempts,
		passwordLimitAttempts: passwordLimitAttempts,
		ipLimitAttempts:       ipLimitAttempts,
	}
}

func (a *App) Check(ctx context.Context, data models.Data) (bool, error) {
	exists, err := a.storage.CheckIPInWhiteList(ctx, data.IP)
	if err != nil {
		return false, err
	} else if exists {
		return true, nil
	}

	exists, err = a.storage.CheckIPInBlackList(ctx, data.IP)
	if err != nil || exists {
		return false, err
	}

	count, err := a.limiter.CheckLimit(ctx, data.IP)
	if err != nil {
		return false, err
	}

	if count > a.ipLimitAttempts {
		a.logger.Info("the number of authorization attempts from this IP has been exceeded",
			map[string]interface{}{"count attempts": count})

		return false, nil
	}

	count, err = a.limiter.CheckLimit(ctx, data.Login)
	if err != nil {
		return false, err
	}

	if count > a.loginLimitAttempts {
		a.logger.Info("the number of authorization attempts from this login has been exceeded",
			map[string]interface{}{"count attempts": count})

		return false, nil
	}

	count, err = a.limiter.CheckLimit(ctx, data.Password)
	if err != nil {
		return false, err
	}

	if count > a.passwordLimitAttempts {
		a.logger.Info("the number of authorization attempts from this password has been exceeded",
			map[string]interface{}{"count attempts": count})

		return false, nil
	}

	return true, nil
}

func (a *App) ClearBuckets(ctx context.Context, data models.Data) error {
	return a.limiter.ClearBuckets(ctx, data.IP, data.Login)
}

func (a *App) AddToBlackList(ctx context.Context, ip string) error {
	return a.storage.AddToBlackList(ctx, ip)
}

func (a *App) AddToWhiteList(ctx context.Context, ip string) error {
	return a.storage.AddToWhiteList(ctx, ip)
}

func (a *App) RemoveFromBlackList(ctx context.Context, ip string) error {
	return a.storage.RemoveFromBlackList(ctx, ip)
}

func (a *App) RemoveFromWhiteList(ctx context.Context, ip string) error {
	return a.storage.RemoveFromWhiteList(ctx, ip)
}

func (a *App) ClearAllBuckets(ctx context.Context) error {
	return a.limiter.ClearAllBuckets(ctx)
}

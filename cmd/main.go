package main

import (
	"context"
	"flag"
	"github.com/Baraulia/anti_bruteforce_service/internal/api/http/handlers"
	"github.com/Baraulia/anti_bruteforce_service/internal/app"
	"github.com/Baraulia/anti_bruteforce_service/internal/storage/limiter"
	sqlstorage "github.com/Baraulia/anti_bruteforce_service/internal/storage/sql"
	"github.com/Baraulia/anti_bruteforce_service/pkg/logger"
	internalhttp "github.com/Baraulia/anti_bruteforce_service/pkg/server/http"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	configFile string
)

func init() {
	flag.StringVar(&configFile, "config", "./configs/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config, err := NewConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	logg, err := logger.GetLogger(config.Logger.Level, true)
	if err != nil {
		log.Fatal(err)
	}
	defer logg.Close()

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	storage := sqlstorage.NewPostgresStorage(sqlstorage.PgConfig{
		Host:     config.SQL.Host,
		Username: config.SQL.Username,
		Password: config.SQL.Password,
		Port:     config.SQL.Port,
		Database: config.SQL.Database,
	}, logg, true)

	lim := limiter.NewLimiter(config.App.Frequency)

	application := app.New(logg, storage, lim, config.App.LoginLimitAttempts, config.App.PasswordLimitAttempts, config.App.IpLimitAttempts)

	handler := handlers.NewHandler(logg, application)
	server := internalhttp.NewServer(logg, config.HTTPServer.Host, config.HTTPServer.Port, handler.InitRoutes())

	logg.Info("application is running...", nil)
	if err := server.Start(); err != nil {
		logg.Error("failed to start http server: "+err.Error(), nil)
		cancel()
		os.Exit(1) //nolint:gocritic
	}

	logg.Info("limiter is running...", nil)
	lim.Start(ctx)

	<-ctx.Done()

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	storage.Close()

	if err = server.Stop(ctx); err != nil {
		logg.Error("failed to stop http server: "+err.Error(), nil)
	}

}

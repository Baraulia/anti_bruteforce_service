package logger

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	logger *zap.SugaredLogger
	file   *os.File
}

type KeyLogger string

var (
	logger *ZapLogger
	once   sync.Once
)

func GetLogger(level string, saveLogsInFile bool) (*ZapLogger, error) {
	var err error
	once.Do(func() {
		var file *os.File
		level = strings.ToUpper(level)
		var zapLevel zapcore.Level
		switch level {
		case "DEBUG":
			zapLevel = zapcore.DebugLevel
		case "INFO":
			zapLevel = zapcore.InfoLevel
		case "WARN":
			zapLevel = zapcore.WarnLevel
		case "ERROR":
			zapLevel = zapcore.ErrorLevel
		case "PANIC":
			zapLevel = zapcore.PanicLevel
		case "FATAL":
			zapLevel = zapcore.FatalLevel
		default:
			logger = nil
			err = fmt.Errorf("unsupported level of logger: %s", level)
			return
		}

		var cores []zapcore.Core
		cores = append(cores, zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.Lock(os.Stdout),
			zapLevel,
		))

		if saveLogsInFile {
			err = os.MkdirAll("logs", 0o777)
			if err != nil {
				logger = nil
				err = fmt.Errorf("could not create directory %w", err)
				return
			}

			file, err = os.OpenFile("logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o777)
			if err != nil {
				logger = nil
				err = fmt.Errorf("could not open file %w", err)
				return
			}

			cores = append(cores, zapcore.NewCore(
				zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
				zapcore.Lock(file),
				zapLevel,
			))
		}

		logger = &ZapLogger{zap.New(zapcore.NewTee(cores...)).Sugar(), file}
	})

	return logger, err
}

func ContextWithLogger(ctx context.Context, logger *ZapLogger) context.Context {
	return context.WithValue(ctx, KeyLogger("logger"), logger)
}

func GetLoggerFromContext(ctx context.Context) (*ZapLogger, error) {
	if l, ok := ctx.Value(KeyLogger("logger")).(*ZapLogger); ok {
		return l, nil
	}

	return GetLogger("info", false)
}

func (l *ZapLogger) Close() {
	if err := l.file.Close(); err != nil {
		log.Fatalf("error while closing file with logs: %v", err)
	}
}

func (l *ZapLogger) Debug(msg string, fields map[string]interface{}) {
	l.logger.Debug(msg, zap.Any("args", fields))
}

func (l *ZapLogger) Info(msg string, fields map[string]interface{}) {
	l.logger.Infow(msg, zap.Any("args", fields))
}

func (l *ZapLogger) Warn(msg string, fields map[string]interface{}) {
	l.logger.Warnw(msg, zap.Any("args", fields))
}

func (l *ZapLogger) Error(msg string, fields map[string]interface{}) {
	l.logger.Errorw(msg, zap.Any("args", fields))
}

func (l *ZapLogger) Fatal(msg string, fields map[string]interface{}) {
	l.logger.Fatalw(msg, zap.Any("args", fields))
}

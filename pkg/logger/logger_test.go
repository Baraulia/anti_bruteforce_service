package logger

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		level         string
		expectedError bool
	}{
		{"DEBUG", false},
		{"INFO", false},
		{"WARN", false},
		{"ERROR", false},
		{"PANIC", false},
		{"FATAL", false},
		{"invalid", true},
	}

	for _, test := range tests {
		t.Run(test.level, func(t *testing.T) {
			once = sync.Once{}
			logg, err := GetLogger(test.level, false)
			if test.expectedError {
				require.Error(t, err)
				require.Nil(t, logg)
			} else {
				require.Nil(t, err)
				require.NotNil(t, logg)
			}
		})
	}
}

func TestLoggerWithCustomContext(t *testing.T) {
	once = sync.Once{}
	logg, err := GetLogger("debug", false)
	require.NoError(t, err)

	ctx := ContextWithLogger(context.Background(), logger)
	require.NotNil(t, ctx.Value(KeyLogger("logger")), "Expected logger in context")

	retrievedLogger, _ := GetLoggerFromContext(ctx)
	require.NotNil(t, retrievedLogger, "Expected logger retrieved from context")
	require.Equal(t, logg, retrievedLogger, "Retrieved logger does not match expected logger")
}

func TestLoggerFromContextWithoutLogger(t *testing.T) {
	once = sync.Once{}
	ctx := context.Background()
	logg, err := GetLoggerFromContext(ctx)
	require.NoError(t, err)
	require.NotNil(t, logg, "Expected a new logger when not present in context")
}

func TestLogOutputToFileAndStdout(t *testing.T) {
	once = sync.Once{}
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = w

	logg, err := GetLogger("debug", false)
	require.NoError(t, err)
	defer func() {
		os.Stdout = oldStdout
	}()

	logg.Info("Test logger message", nil)
	err = w.Close()
	require.NoError(t, err)

	got := make([]byte, 100)
	_, err = r.Read(got)
	require.NoError(t, err)

	require.Contains(t, string(got), "Test logger message", "Expected log message not found in file output")
}

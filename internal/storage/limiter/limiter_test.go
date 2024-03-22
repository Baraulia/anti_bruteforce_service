package limiter

//nolint:depguard
import (
	"context"
	"testing"
	"time"

	"github.com/Baraulia/anti_bruteforce_service/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClearBuckets(t *testing.T) {
	frequency := 5
	limiter := NewLimiter(frequency)
	ip := "127.0.0.1/25"
	login := "Test"
	limiter.buckets[ip] = &models.Bucket{
		CurrentCount: 1,
		LastUpdate:   time.Now(),
	}

	limiter.buckets[login] = &models.Bucket{
		CurrentCount: 1,
		LastUpdate:   time.Now(),
	}

	err := limiter.ClearBuckets(context.Background(), ip, login)
	require.NoError(t, err)

	if _, ok := limiter.buckets[ip]; ok {
		t.Errorf("Expected bucket for IP to be deleted, but it still exists")
	}
	if _, ok := limiter.buckets[login]; ok {
		t.Errorf("Expected bucket for login to be deleted, but it still exists")
	}
}

func TestClearAllBuckets(t *testing.T) {
	frequency := 5
	limiter := NewLimiter(frequency)
	ip := "127.0.0.1"
	login := "Test"
	limiter.buckets[ip] = &models.Bucket{
		CurrentCount: 1,
		LastUpdate:   time.Now(),
	}
	limiter.buckets[login] = &models.Bucket{
		CurrentCount: 1,
		LastUpdate:   time.Now(),
	}

	err := limiter.ClearAllBuckets(context.Background())
	require.NoError(t, err)

	if len(limiter.buckets) != 0 {
		t.Errorf("Expected all buckets to be deleted, but there are still buckets in the limiter")
	}
}

func TestCheckLimit(t *testing.T) {
	key := "test"
	count := 10
	frequency := 5

	limiter := NewLimiter(frequency)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	limiter.Start(ctx)

	for i := 0; i < count; i++ {
		currentCount, err := limiter.CheckLimit(context.Background(), key)
		require.NoError(t, err)

		assert.Equal(t, i+1, currentCount)
	}

	time.Sleep(time.Duration(frequency+1) * time.Second)

	currentCount, err := limiter.CheckLimit(context.Background(), key)
	require.NoError(t, err)

	assert.Equal(t, 1, currentCount)
}

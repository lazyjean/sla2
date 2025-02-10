package redis

import (
	"context"
	"testing"
	"time"

	"github.com/lazyjean/sla2/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRedis(t *testing.T) *RedisCache {
	cfg := &config.RedisConfig{
		Host:            "localhost",
		Port:            "6379",
		Password:        "",
		DB:              1, // 使用不同的数据库避免影响生产
		MaxRetries:      3,
		MinRetryBackoff: time.Millisecond * 100,
		MaxRetryBackoff: time.Second * 2,
		PoolSize:        10,
		MinIdleConns:    2,
		MaxConnAge:      time.Minute * 30,
	}

	cache, err := NewRedisCache(cfg)
	require.NoError(t, err)
	return cache
}

func TestRedisCache_SetGet(t *testing.T) {
	cache := setupTestRedis(t)
	defer cache.Close()

	ctx := context.Background()
	key := "test_key"
	value := "test_value"
	expiration := time.Minute

	// Test Set
	err := cache.Set(ctx, key, value, expiration)
	assert.NoError(t, err)

	// Test Get
	got, err := cache.Get(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, value, got)
}

func TestRedisCache_Delete(t *testing.T) {
	cache := setupTestRedis(t)
	defer cache.Close()

	ctx := context.Background()
	key := "test_delete_key"
	value := "test_value"

	// First set a value
	err := cache.Set(ctx, key, value, time.Minute)
	require.NoError(t, err)

	// Test Delete
	err = cache.Delete(ctx, key)
	assert.NoError(t, err)

	// Verify it's deleted
	_, err = cache.Get(ctx, key)
	assert.Error(t, err)
}

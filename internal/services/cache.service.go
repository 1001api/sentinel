package services

import (
	"context"
	"time"

	"github.com/gofiber/storage/redis/v2"
)

type CacheService interface {
	GetCache(key string) ([]byte, error)
	SetCache(key string, value []byte, exp time.Duration) error
	InvalidateCaches(keys []string) error
}

type CacheServiceImpl struct {
	RedisCon *redis.Storage
}

func InitCacheService(redisCon *redis.Storage) CacheServiceImpl {
	return CacheServiceImpl{
		RedisCon: redisCon,
	}
}

func (s *CacheServiceImpl) GetCache(key string) ([]byte, error) {
	return s.RedisCon.Get(key)
}

func (s *CacheServiceImpl) SetCache(key string, value []byte, exp time.Duration) error {
	return s.RedisCon.Set(key, value, exp)
}

func (s *CacheServiceImpl) InvalidateCaches(keys []string) error {
	if len(keys) > 0 {
		s.RedisCon.Conn().Unlink(context.Background(), keys...)
	}
	return nil
}

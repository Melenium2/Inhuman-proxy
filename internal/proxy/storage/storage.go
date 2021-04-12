package storage

import (
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

// TODO
//		health check all proxies and remove if proxy is unavailable or
//		mute proxy if it has delay
type ProxyStorage interface {
	GetRandom()
	Get(code string)
	Set(code string, proxy string)
}

type RedisStorage struct {
	rdb *redis.Client
	log *zap.SugaredLogger
}

func New(client *redis.Client, logger *zap.SugaredLogger) *RedisStorage {
	return &RedisStorage{
		rdb: client,
		log: logger,
	}
}

func (r RedisStorage) GetRandom() {
	panic("implement me")
}

func (r RedisStorage) Get(code string) {
	panic("implement me")
}

func (r RedisStorage) Set(code string, proxy string) {
}

func (r RedisStorage) healthCheck() error {
	return nil
}

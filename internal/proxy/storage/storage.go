package storage

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"math/rand"
	"time"
)

const (
	Free    = "1"
	Blocked = "0"
)

// TODO
//		health check all proxies and remove if proxy is unavailable or
//		mute proxy if it has delay
type ProxyStorage interface {
	GetRandom(ctx context.Context) (string, error)
	Get(ctx context.Context, code string) (map[string]string, error)
	Set(ctx context.Context, code string, proxy string) error
}

type RedisStorage struct {
	rdb *redis.Client
	log *zap.SugaredLogger
}

func New(client *redis.Client, logger *zap.SugaredLogger) *RedisStorage {
	rand.Seed(time.Now().UnixNano())

	return &RedisStorage{
		rdb: client,
		log: logger,
	}
}

// GetRandom returns first Free random proxy from store
func (r RedisStorage) GetRandom(ctx context.Context) (string, error) {
	cmd := r.rdb.Keys(ctx, "*")
	keys, err := cmd.Result()
	if err != nil {
		return "", err
	}

	kvs := make(map[string]string)
	for i := 0; i < len(keys); i++ {
		kv, err := r.Get(ctx, keys[i])
		if err != nil {
			return "", err
		}
		for k, v := range kv {
			if v == Free {
				kvs[k] = v
			}
		}
	}

	if len(kvs) > 0 {
		randN := rand.Intn(len(kvs))
		i := 0
		for k := range kvs {
			if i == randN {
				return k, nil
			}
			i++
		}
	}

	return "", nil
}

// Get get all proxies from store
func (r RedisStorage) Get(ctx context.Context, code string) (map[string]string, error) {
	cmd := r.rdb.HGetAll(ctx, code)
	res, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	return res, err
}

// Set new proxy address to proxy list by country code
func (r RedisStorage) Set(ctx context.Context, code string, proxy ...string) error {
	proxies := make(map[string]interface{})
	for _, v := range proxy {
		proxies[v] = Free
	}

	cmd := r.rdb.HSet(ctx, code, proxies)
	n, err := cmd.Result()
	if err != nil {
		return err
	}
	if int(n) != len(proxy) {
		return ErrOnPush(len(proxy), int(n))
	}
	return nil
}

// block sets Blocked status to proxy address
func (r RedisStorage) block(ctx context.Context, address string) error {
	return r.changeBlockStatus(ctx, address, Blocked)
}

// unblock sets Free status to proxy address
func (r RedisStorage) unblock(ctx context.Context, address string) error {
	return r.changeBlockStatus(ctx, address, Free)
}

// changeBlockStatus changes proxy address with given address
func (r RedisStorage) changeBlockStatus(ctx context.Context, address string, status string) error {
	cmd := r.rdb.Keys(ctx, "*")
	keys, err := cmd.Result()
	if err != nil {
		return err
	}

	for _, k := range keys {
		cmd := r.rdb.HExists(ctx, k, address)
		isExist, err := cmd.Result()
		if err != nil {
			return err
		}
		if isExist {
			cmd := r.rdb.HSet(ctx, k, address, status)
			if err := cmd.Err(); err != nil {
				return err
			}
		}
	}

	return nil
}

// delete removes proxy address from redis storage
func (r RedisStorage) delete(ctx context.Context, address string) error {
	cmd := r.rdb.Keys(ctx, "*")
	keys, err := cmd.Result()
	if err != nil {
		return err
	}

	var lastKey string
	for _, k := range keys {
		cmd := r.rdb.HDel(ctx, k, address)
		n, err := cmd.Result()
		if err != nil {
			return err
		}
		if n > 0 {
			lastKey = k
			break
		}
	}

	lenCmd := r.rdb.HLen(ctx, lastKey)
	nKey, err := lenCmd.Result()
	if err != nil {
		return err
	}
	if nKey > 0 {
		delCmd := r.rdb.Del(ctx, lastKey)
		err := delCmd.Err()
		if err != nil {
			return err
		}
	}

	return nil
}

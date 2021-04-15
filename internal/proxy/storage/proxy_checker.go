package storage

import (
	"context"
	"github.com/Melenium2/inhuman-reverse-proxy/internal/config"
	"github.com/Melenium2/inhuman-reverse-proxy/internal/proxy/health"
	"go.uber.org/zap"
	"time"
)

type CheckStatus uint8

const (
	Remove CheckStatus = iota
	Block
	Freed
)

type ProxyChecker struct {
	storage *RedisStorage
	cfg     config.CheckerConfig
	log     *zap.SugaredLogger
}

func NewChecker(storage *RedisStorage, logger *zap.SugaredLogger, config config.CheckerConfig) *ProxyChecker {
	return &ProxyChecker{
		storage: storage,
		cfg:     config,
		log:     logger,
	}
}

func (pc *ProxyChecker) Check() {
	for {
		ctx := context.Background()
		keys, err := pc.storage.keys(ctx)
		if err != nil {
			pc.log.Error(err)
			continue
		}

		kvs := make(map[string]string)
		for _, key := range keys {
			kv, err := pc.storage.Get(ctx, key)
			if err != nil {
				pc.log.Error(err)
			}
			for k, v := range kv {
				kvs[k] = v
			}
		}

		for address, status := range kvs {
			_ = pc.check(ctx, address, status)
		}

		time.Sleep(pc.cfg.Interval)
	}
}

func (pc *ProxyChecker) check(ctx context.Context, address, status string) CheckStatus {
	if err := health.Check(address, pc.cfg.MaxTimeout); err != nil {
		if err == health.ErrProxyUnreachable {
			if err := pc.storage.delete(ctx, address); err != nil {
				pc.log.Error(err)
			}

			pc.log.Infof("proxy %s removed from storage, becuase it is unreacheble", address)

			return Remove
		}
		if err := pc.storage.block(ctx, address); err != nil {
			pc.log.Error(err)
		}

		pc.log.Infof("proxy %s has timeout more then %s, it is blocked", address, pc.cfg.MaxTimeout)

		return Block
	}

	if status == Blocked {
		if err := pc.storage.unblock(ctx, address); err != nil {
			pc.log.Error(err)
		}
		pc.log.Infof("proxy %s has freedom", address)
	}

	return Freed
}

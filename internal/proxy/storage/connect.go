package storage

import (
	"context"
	"github.com/Melenium2/inhuman-reverse-proxy/internal/config"
	"github.com/go-redis/redis/v8"
)

func Connect(config config.StorageConfig) (*redis.Client, error) {
	rdclient := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		DB:       config.Database,
		Username: config.Username,
		Password: config.Password,
	})

	cmd := rdclient.Ping(context.Background())

	return rdclient, cmd.Err()
}

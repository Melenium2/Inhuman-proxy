package storage_test

import (
	"context"
	"github.com/Melenium2/inhuman-reverse-proxy/internal/config"
	"github.com/Melenium2/inhuman-reverse-proxy/internal/proxy/storage"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"testing"
)

var (
	cfg = config.StorageConfig{
		Addr:     "192.168.99.100:6379",
		Username: "",
		Password: "123456",
		Database: 0,
	}
)

type RedisStorageSuite struct {
	suite.Suite
	code    string
	rdb     *redis.Client
	storage *storage.RedisStorage
	clean   func() error
}

func (s *RedisStorageSuite) SetupTest() {
	var err error

	s.code = "ru"

	s.rdb, err = storage.Connect(cfg)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), s.rdb)

	s.storage = storage.New(s.rdb, zap.NewNop().Sugar())

	s.clean = func() error {
		cmd := s.rdb.FlushDB(context.Background())
		return cmd.Err()
	}
}

func (s *RedisStorageSuite) TestGet_ShouldReturnAllKeysFromRedis() {
	expected := map[string]interface{}{
		"key1": "1",
		"key2": "1",
		"key3": "1",
		"key4": "1",
	}

	cmd := s.rdb.HSet(context.Background(), s.code, expected)
	s.Require().NoError(cmd.Err())

	res, err := s.storage.Get(context.Background(), s.code)
	s.Require().NoError(err)

	for k := range expected {
		s.Assert().Equal(expected[k].(string), res[k])
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(RedisStorageSuite))
}

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
	defer s.clean() //nolint:errcheck

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

func (s *RedisStorageSuite) TestGet_ShouldReturnEmptyResultFromRedis() {
	res, err := s.storage.Get(context.Background(), s.code)
	s.Require().NoError(err)
	s.Assert().Empty(res)
}

func (s *RedisStorageSuite) TestSet_ShouldSetNewValueToRedis() {
	defer s.clean() //nolint:errcheck

	expected := []string{"1", "2", "3"}

	err := s.storage.Set(context.Background(), s.code, expected...)
	s.Require().NoError(err)

	cmd := s.rdb.HGetAll(context.Background(), s.code)
	res, err := cmd.Result()
	s.Require().NoError(err)

	s.Assert().Equal(len(expected), len(res))

	for _, k := range expected {
		value, ok := res[k]
		s.Assert().True(ok)
		s.Assert().Equal(storage.Free, value)
	}
}

func (s *RedisStorageSuite) TestSet_ShouldSetNewValueThenAppendNewValueToTheSameKey() {
	defer s.clean() //nolint:errcheck

	expected1 := []string{"1", "2", "3"}
	expected2 := []string{"5", "6"}

	err := s.storage.Set(context.Background(), s.code, expected1...)
	s.Require().NoError(err)
	err = s.storage.Set(context.Background(), s.code, expected2...)
	s.Require().NoError(err)

	cmd := s.rdb.HGetAll(context.Background(), s.code)
	res, err := cmd.Result()
	s.Require().NoError(err)

	s.Assert().Equal(len(expected1)+len(expected2), len(res))

	for _, k := range append(expected1, expected2...) {
		value, ok := res[k]
		s.Assert().True(ok)
		s.Assert().Equal(storage.Free, value)
	}
}

func (s *RedisStorageSuite) TestGetRandom_ShouldReturnRandomProxyFromStorage() {
	defer s.clean() //nolint:errcheck

	code1 := "ru"
	expected1 := []string{"p1", "1", "p2", "1"}
	code2 := "en"
	expected2 := []string{"p3", "1"}

	s.rdb.HSet(context.Background(), code1, expected1)
	s.rdb.HSet(context.Background(), code2, expected2)

	res, err := s.storage.GetRandom(context.Background())
	s.Require().NoError(err)

	for i, item := range append(expected1, expected2...) {
		if i&1 == 1 {
			continue
		}
		if item == res {
			s.T().Logf("find, %s res, %s", item, res)
			break
		}
	}
}

func (s *RedisStorageSuite) TestGetRandom_ShouldReturnExpectedValueIfOtherValuesBlocked() {
	var tt = []struct {
		name      string
		expected  string
		code1     string
		code2     string
		expected1 []string
		expected2 []string
	}{
		{
			name:      "should return p2",
			expected:  "p2",
			code1:     "ru",
			code2:     "en",
			expected1: []string{"p1", "0", "p2", "1"},
			expected2: []string{"p3", "0"},
		},
		{
			name:      "should return empty line",
			expected:  "",
			code1:     "ru",
			code2:     "en",
			expected1: []string{"p1", "0", "p2", "0"},
			expected2: []string{"p3", "0"},
		},
	}

	for _, test := range tt {
		s.T().Run(test.name, func(t *testing.T) {
			defer s.clean() //nolint:errcheck

			s.rdb.HSet(context.Background(), test.code1, test.expected1)
			s.rdb.HSet(context.Background(), test.code2, test.expected2)

			res, err := s.storage.GetRandom(context.Background())
			s.Require().NoError(err)

			s.Assert().Equal(test.expected, res)
		})
	}
}

func (s *RedisStorageSuite) TestChangeBlockStatus_ShouldChangeBlockStatusOfProxy() {
	var tt = []struct {
		name           string
		key            string
		example        []string
		expectedStatus string
		expectedError  bool
	}{
		{
			name:           "should change status to blocked in proxy 'p1'",
			key:            "p1",
			expectedStatus: "0",
			example:        []string{"p1", "1", "p2", "1", "p3", "1"},
		},
		{
			name:           "should change status to free in proxy 'p2'",
			key:            "p2",
			expectedStatus: "1",
			example:        []string{"p1", "1", "p2", "0", "p3", "1"},
		},
		{
			name:           "should return nil if store is empty",
			key:            "",
			expectedStatus: "",
			example:        []string{},
			expectedError:  true,
		},
	}

	for _, test := range tt {
		s.T().Run(test.name, func(t *testing.T) {
			defer s.clean() //nolint:errcheck

			s.rdb.HSet(context.Background(), s.code, test.example)

			err := s.storage.ChangeBlockStatus(context.Background(), test.key, test.expectedStatus)
			s.Require().NoError(err)

			cmd := s.rdb.HGet(context.Background(), s.code, test.key)
			res, err := cmd.Result()
			s.Require().Equal(test.expectedError, err != nil)

			s.Assert().Equal(res, test.expectedStatus)
		})
	}
}

func (s *RedisStorageSuite) TestDelete_ShouldRemoveProxyFromStorage() {
	var tt = []struct {
		name          string
		proxy         string
		example       []string
		expectedError bool
	}{
		{
			name:          "should remove proxy from storage",
			proxy:         "p3",
			example:       []string{"p1", "1", "p2", "0", "p3", "1"},
			expectedError: false,
		},
		{
			name:          "storage is empty, should do nothing",
			proxy:         "",
			example:       []string{},
			expectedError: false,
		},
	}

	for _, test := range tt {
		s.T().Run(test.name, func(t *testing.T) {
			defer s.clean() //nolint:errcheck

			s.rdb.HSet(context.Background(), s.code, test.example)

			err := s.storage.Delete(context.Background(), test.proxy)
			s.Require().Equal(test.expectedError, err != nil)
		})
	}
}

func (s *RedisStorageSuite) TestDelete_ShouldRemoveKeyIfAfterRemovingItemLengthEqualsZero() {
	defer s.clean() //nolint:errcheck

	s.rdb.HSet(context.Background(), s.code, []string{"p1", "1"})

	err := s.storage.Delete(context.Background(), "p1")
	s.Require().NoError(err)

	cmd := s.rdb.Get(context.Background(), s.code)
	_, err = cmd.Result()
	s.Assert().Equal(redis.Nil, err)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(RedisStorageSuite))
}

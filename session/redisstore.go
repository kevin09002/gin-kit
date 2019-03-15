package session

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/kevinma2010/gkits/redis"
)

type redisStore struct {
	rds *redis.Client
}

func NewRedisStore(cfg *rds.Config) *redisStore {
	return &redisStore{
		rds: rds.NewClient(cfg).RedisClient,
	}
}

func NewRedisStoreWithClient(client *redis.Client) *redisStore {
	return &redisStore{
		rds: client,
	}
}

func (s *redisStore) GetSession(sessionID, key string) ([]byte, bool, error) {
	r, err := s.rds.Get(sessionID + ":" + key).Result()
	if err != nil {
		if err != redis.Nil {
			return nil, false, err
		}
		return nil, false, nil
	}

	if r == "" {
		return nil, false, nil
	}
	return []byte(r), true, nil
}

func (s *redisStore) SetSession(sessionID, key string, data []byte) error {
	return s.rds.Set(sessionID+":"+key, data, time.Duration(SessionTimeout)*time.Minute).Err()
}

func (s *redisStore) ClearSession(sessionID, key string) error {
	return s.rds.Del(sessionID + ":" + key).Err()
}

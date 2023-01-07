package repositories

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	e "github.com/natanaelb27/oceana-backend/pkg/errors"
)

type RedisRepository interface {
	Get(c *gin.Context, key string) (result string, err error)
	Set(c *gin.Context, key string, value string, expire time.Duration) error
	Delete(c *gin.Context, key string) (err error)
}

type redisRepository struct {
	redisClient *redis.Client
}
type RedisRepositoryConfig struct {
	RedisClient *redis.Client
}

func NewRedisRepository(c *RedisRepositoryConfig) RedisRepository {
	return &redisRepository{
		redisClient: c.RedisClient,
	}
}

func (r *redisRepository) Get(c *gin.Context, key string) (result string, err error) {
	result, err = r.redisClient.Get(c, key).Result()
	if err != nil {
		if err == redis.Nil {
			err = e.ErrDataNotFound
			return
		}
		err = e.ErrCacheOperationFailed
		return
	}

	return
}

func (r *redisRepository) Set(c *gin.Context, key string, value string, expire time.Duration) (err error) {
	err = r.redisClient.Set(c, key, value, expire).Err()
	if err != nil {
		err = e.ErrCacheOperationFailed
		return
	}

	return nil
}

func (r *redisRepository) Delete(c *gin.Context, key string) (err error) {
	err = r.redisClient.Del(c, key).Err()
	if err != nil {
		err = e.ErrCacheOperationFailed
		return
	}

	return nil
}

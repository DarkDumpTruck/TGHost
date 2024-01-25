package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

var ErrRedisNotAvailable = fmt.Errorf("redis not available")

func InitCache(redisHost string, redisPassword string, redisDB int) (*redis.Client, error) {
	CacheClient := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPassword,
		DB:       redisDB,
	})
	resp, err := CacheClient.Ping(context.Background()).Result()
	if err != nil {
		return nil, errors.Wrap(err, "redis ping failed")
	}
	if resp != "PONG" {
		return nil, errors.New("redis failed to pong")
	}
	return CacheClient, nil
}

func GetByteCache(cacheClient *redis.Client, key string) ([]byte, error) {
	if cacheClient == nil {
		return nil, ErrRedisNotAvailable
	}
	return cacheClient.Get(context.Background(), key).Bytes()
}

func SetByteCache(cacheClient *redis.Client, key string, value []byte, ttl time.Duration) error {
	if cacheClient == nil {
		return ErrRedisNotAvailable
	}
	return cacheClient.Set(context.Background(), key, value, ttl).Err()
}

func GetCache(cacheClient *redis.Client, key string) (string, error) {
	if cacheClient == nil {
		return "", ErrRedisNotAvailable
	}
	return cacheClient.Get(context.Background(), key).Result()
}

func GetCacheWithCtx(ctx context.Context, cacheClient *redis.Client, key string) (string, error) {
	if cacheClient == nil {
		return "", ErrRedisNotAvailable
	}
	return cacheClient.Get(ctx, key).Result()
}

func GetCaches(cacheClient *redis.Client, keys []string) ([]string, error) {
	if cacheClient == nil {
		return nil, ErrRedisNotAvailable
	}
	result, err := cacheClient.MGet(context.Background(), keys...).Result()
	if err != nil {
		return nil, err
	}
	var ret []string
	for _, val := range result {
		if val == nil {
			ret = append(ret, "")
		} else {
			ret = append(ret, val.(string))
		}
	}
	return ret, nil
}

func GetCachesWithCtx(ctx context.Context, cacheClient *redis.Client, keys []string) ([]string, error) {
	if cacheClient == nil {
		return nil, ErrRedisNotAvailable
	}
	result, err := cacheClient.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}
	var ret []string
	for _, val := range result {
		if val == nil {
			ret = append(ret, "")
		} else {
			ret = append(ret, val.(string))
		}
	}
	return ret, nil
}

func SetCache(cacheClient *redis.Client, key string, value string, ttl time.Duration) error {
	if cacheClient == nil {
		return ErrRedisNotAvailable
	}
	return cacheClient.Set(context.Background(), key, value, ttl).Err()
}

func SetCacheWithCtx(ctx context.Context, cacheClient *redis.Client, key string, value string, ttl time.Duration) error {
	if cacheClient == nil {
		return ErrRedisNotAvailable
	}
	return cacheClient.Set(ctx, key, value, ttl).Err()
}

func ClearCache(cacheClient *redis.Client, key string) error {
	if cacheClient == nil {
		return ErrRedisNotAvailable
	}
	return cacheClient.Del(context.Background(), key).Err()
}

func ClearCacheByPrefix(cacheClient *redis.Client, prefix string) error {
	if cacheClient == nil {
		return ErrRedisNotAvailable
	}
	var cursor uint64
	for {
		var keys []string
		var err error
		keys, cursor, err = cacheClient.Scan(context.Background(), cursor, prefix+"*", 10).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			err = cacheClient.Del(context.Background(), keys...).Err()
			if err != nil {
				return err
			}
		}
		if cursor == 0 {
			break
		}
	}
	return nil
}

func GetCachesByPrefix(cacheClient *redis.Client, prefix string) (map[string]string, error) {
	if cacheClient == nil {
		return nil, ErrRedisNotAvailable
	}
	var cursor uint64
	result := make(map[string]string)
	for {
		var keys []string
		var err error
		keys, cursor, err = cacheClient.Scan(context.Background(), cursor, prefix+"*", 10).Result()
		if err != nil {
			return nil, err
		}
		if len(keys) > 0 {
			for _, key := range keys {
				val, err := cacheClient.Get(context.Background(), key).Result()
				if err != nil {
					return nil, err
				}
				result[key] = val
			}
		}
		if cursor == 0 {
			break
		}
	}
	return result, nil
}

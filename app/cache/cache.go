package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"gin-web/app/vars"
	"gin-web/libraty/redis"
	"time"

	"gorm.io/gorm"
)

type RedisCache struct {
	*redis.DBRedis
	*gorm.DB
	debug  bool
	expire int
}

func New(db *gorm.DB) *RedisCache {
	return &RedisCache{
		DBRedis: vars.DBRedis,
		DB:      db,
		expire:  vars.YmlConfig.GetInt("Redis.ExpireTime"),
		debug:   vars.YmlConfig.GetBool("Debug"),
	}
}

func (rc *RedisCache) Query(cacheKey string, v interface{}, fn func(db *gorm.DB, v interface{}) error) error {
	if rc.debug {
		return fn(rc.DB, v)
	}
	if val := rc.GetString(cacheKey); val != "" {
		return json.Unmarshal([]byte(val), v)
	} else {
		if err := fn(rc.DB, v); err != nil {
			return err
		} else {
			return rc.setCache(cacheKey, v)
		}
	}
}

func (rc *RedisCache) QueryRow(cacheKey string, v interface{}, id interface{}, fn func(db *gorm.DB, v interface{}, id interface{}) error) error {
	if rc.debug {
		return fn(rc.DB, v, id)
	}
	key := fmt.Sprintf("%s%v", cacheKey, id)
	if val := rc.GetString(key); val != "" {
		return json.Unmarshal([]byte(val), v)
	} else {
		if err := fn(rc.DB, v, id); err != nil {
			return err
		} else {
			return rc.setCache(key, v)
		}
	}
}

func (rc *RedisCache) DelQueryRowCache(cacheKey string, id interface{}) error {
	if rc.debug {
		return nil
	}
	key := fmt.Sprintf("%s%v", cacheKey, id)
	return rc.ExecRedis("del", key)
}

func (rc *RedisCache) SetRow(cacheKey string, v interface{}, id interface{}, fn func(db *gorm.DB, v interface{}, id interface{}) error) error {
	if err := fn(rc.DB, v, id); err != nil {
		return err
	}
	if rc.debug {
		return nil
	}
	key := fmt.Sprintf("%s%v", cacheKey, id)
	return rc.setCache(key, v)
}

func (rc *RedisCache) setCache(cacheKey string, v interface{}) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}
	if !rc.SetString(cacheKey, string(bytes), time.Duration(rc.expire)*time.Second) {
		return errors.New("设置缓存失败")
	}
	return nil
}

// SetExpire 设置缓存时长 单位秒
func (rc *RedisCache) SetExpire(expire int) *RedisCache {
	rc.expire = expire
	return rc
}

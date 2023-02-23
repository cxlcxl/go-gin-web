package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type DBRedis struct {
	*redis.Client
	keyPrefix string
	ctx       context.Context
}

func GetRedisInstance(host, password, prefix string, db int) (*DBRedis, error) {
	ctx := context.Background()
	r := redis.NewClient(&redis.Options{Addr: host, Password: password, DB: db})
	if _, err := r.Ping(ctx).Result(); err != nil {
		return nil, err
	}
	if prefix[len(prefix)-1:] != ":" {
		prefix = prefix + ":"
	}
	return &DBRedis{Client: r, keyPrefix: prefix, ctx: ctx}, nil
}

func (dr *DBRedis) HasKey(key string) bool {
	val, err := dr.Client.Exists(dr.ctx, dr.keyPrefix+key).Result()
	if err != nil {
		return false
	}
	return val > 0
}

func (dr *DBRedis) GetString(key string) string {
	val, err := dr.Client.Get(dr.ctx, dr.keyPrefix+key).Result()
	if err != nil {
		return ""
	}
	return val
}

func (dr *DBRedis) SetString(key, val string, expire time.Duration) bool {
	_, err := dr.Client.Set(dr.ctx, dr.keyPrefix+key, val, expire*time.Second).Result()
	if err != nil {
		return false
	}
	return true
}

func (dr *DBRedis) SetNX(key string, val interface{}, expire time.Duration) (bool, error) {
	return dr.Client.SetNX(dr.ctx, dr.keyPrefix+key, val, expire).Result()
}

// SMembers set 操作
func (dr *DBRedis) SMembers(key string) ([]string, error) {
	return dr.Client.SMembers(dr.ctx, dr.keyPrefix+key).Result()
}

func (dr *DBRedis) SAdd(key string, v ...interface{}) (err error) {
	_, err = dr.Client.SAdd(dr.ctx, dr.keyPrefix+key, v...).Result()
	return
}

func (dr *DBRedis) SRem(key string, v ...interface{}) (err error) {
	_, err = dr.Client.SRem(dr.ctx, dr.keyPrefix+key, v...).Result()
	return
}

func (dr *DBRedis) ExpireTime(key string, t int64) (err error) {
	_, err = dr.Client.Expire(dr.ctx, dr.keyPrefix+key, time.Second*time.Duration(t)).Result()
	return
}

func (dr *DBRedis) ExecRedis(args ...interface{}) (err error) {
	_, err = dr.Client.Do(dr.ctx, args...).Result()
	return
}

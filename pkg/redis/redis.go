package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
	// "github.com/gomodule/redigo/redis"

	"qiu/blog/pkg/setting"
)

// var RedisConn *redis.Pool
var ctx = context.Background()
var rdb *redis.Client

func Setup() error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     setting.RedisSetting.Host,
		Password: setting.RedisSetting.Password, // no password set
		DB:       0,                             // use default DB
	})
	if rdb == nil {
		panic("Redis Setup Fail!")
	}
	return nil
}

//redis.KeepTTL 10*time.Second
func SetBytes(key string, value interface{}, etime time.Duration) {
	data, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	if err := rdb.Set(ctx, key, data, etime).Err(); err != nil {
		panic(err)
	}
}
func Set(key string, value interface{}, etime time.Duration) {
	if err := rdb.Set(ctx, key, value, etime).Err(); err != nil {
		panic(err)
	}
}
func Get(key string) string {
	ret, err := rdb.Get(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	return ret
}

func Del(key string) {
	if err := rdb.Del(ctx, key).Err(); err != nil {
		panic(err)
	}
}

func GetBytes(key string) []byte {
	ret, err := rdb.Get(ctx, key).Bytes()
	if err != nil {
		panic(err)
	}
	return ret
}
func Exists(keys ...string) int64 {
	ret, err := rdb.Exists(ctx, keys...).Result()
	if err != nil {
		panic(err)
	}
	return ret
}
func HashSet(key string, data map[string]interface{}) {
	if err := rdb.HSet(ctx, key, data).Err(); err != nil {
		panic(err)
	}
}
func HashIncr(key string, field string, v int64) {
	if err := rdb.HIncrBy(ctx, key, field, v).Err(); err != nil {
		panic(err)
	}
}
func HashGetAll(key string) map[string]string {
	ret, err := rdb.HGetAll(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	return ret
}

func HashGet(key string, field string) string {
	ret, err := rdb.HGet(ctx, key, field).Result()
	if err != nil {
		return ""
	}
	return ret
}

func HashDel(key string, field string) {
	if err := rdb.HDel(ctx, key, field).Err(); err != nil {
		panic(err)
	}
}
func HashMGet(key string, fields ...string) []interface{} {
	ret, err := rdb.HMGet(ctx, key, fields...).Result()
	if err != nil {
		panic(err)
	}
	return ret
}

func HashKeyExist(key string, field string) bool {
	ret, err := rdb.HExists(ctx, key, field).Result()
	if err != nil {
		panic(err)
	}
	return ret
}

func SetBit(key string, offset int64, value int) {
	if err := rdb.SetBit(ctx, key, offset, value).Err(); err != nil {
		panic(err)
	}
}

func GetBit(key string, offset int64) int64 {
	ret, err := rdb.GetBit(ctx, key, offset).Result()
	if err != nil {
		panic(err)
	}
	return ret
}

func BitCount(key string) int64 {
	ret, err := rdb.BitCount(ctx, key, &redis.BitCount{}).Result()
	if err != nil {
		panic(err)
	}
	return ret
}

func SAdd(key string, value interface{}) {
	if err := rdb.SAdd(ctx, key, value).Err(); err != nil {
		panic(err)
	}
}

func SDEL(key string, value interface{}) {
	if err := rdb.SRem(ctx, key, value).Err(); err != nil {
		panic(err)
	}
}

func SExist(key string, memeber interface{}) bool {
	ret, err := rdb.SIsMember(ctx, key, memeber).Result()
	if err != nil {
		panic(err)
	}
	return ret
}
func SGET(key string) []string {
	ret, err := rdb.SMembers(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	return ret
}

func ZAdd(key string, score float64, value interface{}) {
	z := redis.Z{
		Score:  score,
		Member: value,
	}
	if err := rdb.ZAdd(ctx, key, z).Err(); err != nil {
		panic(err)
	}
}

func ZAddByte(key string, score float64, value interface{}) {
	data, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	z := redis.Z{
		Score:  score,
		Member: data,
	}
	if err := rdb.ZAdd(ctx, key, z).Err(); err != nil {
		panic(err)
	}
}
func ZRandMember(key string, cnt int) []string {
	ret, err := rdb.ZRandMember(ctx, key, cnt).Result()
	if err != nil {
		panic(err)
	}
	return ret
}

func ZRange(key string, start int64, stop int64) []string {
	ret, err := rdb.ZRange(ctx, key, start, stop).Result()
	if err != nil {
		panic(err)
	}
	return ret
}

func ZRevRange(key string, start int64, stop int64) []string {
	ret, err := rdb.ZRevRange(ctx, key, start, stop).Result()
	if err != nil {
		panic(err)
	}
	return ret
}

func ZRank(key string, member string) int64 {
	return rdb.ZRank(ctx, key, member).Val()
}

func ZRem(key string, member interface{}) {
	if err := rdb.ZRem(ctx, key, member).Err(); err != nil {
		panic(err)
	}

}

func ZCard(key string) int64 {
	ret, err := rdb.ZCard(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	return ret
}

func Expire(key string, timestamp time.Duration) {
	if err := rdb.Expire(ctx, key, timestamp).Err(); err != nil {
		panic(err)
	}
}
func ScanSetByPattern(pattern string) map[string][]string {
	iter := rdb.Scan(ctx, 0, pattern, 0).Iterator()
	ret := make(map[string][]string)
	var err error
	for iter.Next(ctx) {
		key := iter.Val()
		fmt.Println("Scan keys", key)
		ret[key], err = rdb.SMembers(ctx, iter.Val()).Result()
		if err != nil {
			panic(err)
		}
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}
	return ret
}

func ScanHashByPattern(pattern string) map[string]interface{} {
	iter := rdb.Scan(ctx, 0, pattern, 0).Iterator()
	ret := make(map[string]interface{})
	for iter.Next(ctx) {
		key := iter.Val()
		fmt.Println("Scan keys", key)
		ret[key] = HashGetAll(key)
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}
	return ret
}

func Publish(channel string, message interface{}) {
	err := rdb.Publish(ctx, channel, message).Err()
	if err != nil {
		panic(err)
	}
}

func Subscribe(channel string) *redis.PubSub {
	return rdb.Subscribe(ctx, channel)
}

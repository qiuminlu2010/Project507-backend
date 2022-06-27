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
	// RedisConn = &redis.Pool{
	// 	MaxIdle:     setting.RedisSetting.MaxIdle,
	// 	MaxActive:   setting.RedisSetting.MaxActive,
	// 	IdleTimeout: setting.RedisSetting.IdleTimeout,
	// 	Dial: func() (redis.Conn, error) {
	// 		c, err := redis.Dial("tcp", setting.RedisSetting.Host)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		if setting.RedisSetting.Password != "" {
	// 			if _, err := c.Do("AUTH", setting.RedisSetting.Password); err != nil {
	// 				c.Close()
	// 				return nil, err
	// 			}
	// 		}
	// 		return c, err
	// 	},
	// 	TestOnBorrow: func(c redis.Conn, t time.Time) error {
	// 		_, err := c.Do("PING")
	// 		return err
	// 	},
	// }

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
		panic(err)
	}
	return ret
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
	// l,err := rdb.StrLen(ctx, key).Result()
	// if err != nil {
	//     panic(err)
	// }
	// sit := redis.BitCount{}
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

// func SetString(key string, value string, time int) error {
// 	conn := RedisConn.Get()
// 	defer conn.Close()

// 	_, err := conn.Do("SET", key, value)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = conn.Do("EXPIRE", key, time)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func Set(key string, data interface{}, time int) error {
// 	conn := RedisConn.Get()
// 	defer conn.Close()

// 	value, err := json.Marshal(data)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = conn.Do("SET", key, value)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = conn.Do("EXPIRE", key, time)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func Exists(key string) bool {
// 	conn := RedisConn.Get()
// 	defer conn.Close()

// 	exists, err := redis.Bool(conn.Do("EXISTS", key))
// 	if err != nil {
// 		return false
// 	}

// 	return exists
// }

// func Get(key string) ([]byte, error) {
// 	conn := RedisConn.Get()
// 	defer conn.Close()

// 	reply, err := redis.Bytes(conn.Do("GET", key))
// 	if err != nil {
// 		return nil, err
// 	}

// 	return reply, nil
// }

// func Delete(key string) (bool, error) {
// 	conn := RedisConn.Get()
// 	defer conn.Close()

// 	return redis.Bool(conn.Do("DEL", key))
// }

// func LikeDeletes(key string) error {
// 	conn := RedisConn.Get()
// 	defer conn.Close()

// 	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
// 	if err != nil {
// 		return err
// 	}

// 	for _, key := range keys {
// 		_, err = Delete(key)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

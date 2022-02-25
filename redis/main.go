package redis

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var ctx = context.Background()
var rdb *redis.Client

func Init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ENDPOINT"),
		Password: os.Getenv("REDIS_DB_PASSWORD"),
		DB:       0, // use default DB
	})
}

func GetClient() (*redis.Client, context.Context) {
	return rdb, ctx
}

// Set cache for any struct value ex: project/category
func SetCache(key string, val interface{}) {
	expired := time.Second * 15
	b, marshalErr := json.Marshal(val)
	if marshalErr != nil {
		zap.S().Fatal("Cache: Marshal binary failed: " + marshalErr.Error())
	}
	err := rdb.Set(ctx, key, b, expired).Err()
	if err != nil {
		zap.S().Fatal("Set Cache error: " + err.Error())
	}
}

// get cache for any struct value ex: project/category
func GetCache(key string, dest interface{}) {
	val, err := rdb.Get(ctx, key).Result()
	val, resultErr := CheckResult(val, err)
	if resultErr != nil && resultErr.Error() != "Key does not exist" {
		zap.S().Fatal("Get Cache error: " + resultErr.Error())
	}
	err = json.Unmarshal([]byte(val), dest)
	if err != nil {
		zap.S().Fatal("Get Cache Unmarshal error: " + err.Error())
	}
}

// set cache for any int value ex: views
func SetCacheInt(key string, val int) {
	err := rdb.Set(ctx, key, val, 0).Err()
	if err != nil {
		zap.S().Fatal("Set View cache error: ", err.Error())
	}
}

// get cache for any int value ex: views
func GetCacheInt(key string) int {
	count, err := rdb.Get(ctx, key).Result()
	count, err = CheckResult(count, err)
	if err != nil && err.Error() != "Key does not exist" {
		zap.S().Fatal("Get View cache error: ", err.Error())
	}

	countInt, err := strconv.Atoi(count)
	if err != nil {
		zap.S().Fatal("Set View cache string to int error: ", err.Error())
	}
	return countInt
}

func DeleteCache(key string) {
	err := rdb.Del(ctx, key).Err()
	if err != nil {
		zap.S().Fatal("Delete cache failed, key=", key, ", error: ", err.Error())
	}
}

func FindAllCacheByPrefix(prefix string) ([]string, uint64) {
	var keys []string
	var err error
	var cursor uint64
	keys, cursor, err = rdb.Scan(ctx, cursor, prefix+"*", 0).Result()
	if err != nil {
		zap.S().Fatal("find all prefix=", prefix, " cache error:", err.Error())
	}
	return keys, cursor
}

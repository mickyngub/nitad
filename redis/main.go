package redis

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/birdglove2/nitad-backend/errors"
	"github.com/go-redis/redis/v8"
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
func SetCache(key string, val interface{}) errors.CustomError {
	log.Println("setting cached", key)
	expired := time.Second * 15
	b, marshalErr := json.Marshal(val)
	if marshalErr != nil {
		log.Println("Cache: Marshal binary failed: " + marshalErr.Error())
		return nil
	}
	err := rdb.Set(ctx, key, b, expired).Err()
	if err != nil {
		log.Println("Set cache error: ", err.Error())
		return nil
	}
	return nil
}

// get cache for any struct value ex: project/category
func GetCache(key string, dest interface{}) errors.CustomError {
	val, err := rdb.Get(ctx, key).Result()
	val, resultErr := CheckResult(val, err)
	if resultErr != nil {
		log.Println("Cache: get failed: " + resultErr.Error())
		return nil
	}
	err = json.Unmarshal([]byte(val), dest)
	if err != nil {
		log.Println("Cache: get failed: " + err.Error())
		return nil

	}
	return nil
}

// set cache for any int value ex: views
func SetCacheInt(key string, val int) {
	err := rdb.Set(ctx, key, val, 0).Err()
	if err != nil {
		log.Println("set view cache error: ", err.Error())
	}
	log.Println(key, " incrementing view... ", val)

}

// get cache for any int value ex: views
func GetCacheInt(key string) int {
	count, err := rdb.Get(ctx, key).Result()
	if err != nil && err.Error() != "Key does not exist" {
		log.Println("get view cache error: ", err.Error())
	}

	countInt, err := strconv.Atoi(count)
	if err != nil {
		log.Println("views: string to int error: ", err.Error())
	}
	return countInt
}

func DeleteCache(key string) {
	err := rdb.Del(ctx, key).Err()
	if err != nil {
		log.Println("Delete ", key, " failed!", err.Error())
	}
	log.Println("Deleting... ", key)
}

func FindAllCacheByPrefix(prefix string) ([]string, uint64) {
	var keys []string
	var err error
	var cursor uint64
	keys, cursor, err = rdb.Scan(ctx, cursor, prefix+"*", 0).Result()
	if err != nil {
		log.Println("find all prefix=", prefix, " cache error:", err.Error())
	}
	return keys, cursor
}

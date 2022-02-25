package redis

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"
)

// type Storage struct {
// 	db *redis.Client
// }

var store *Storage
var DefaultCacheExpireTime time.Duration = time.Second * 15

func GetStore() *Storage {
	return store
}

func Init() *Storage {
	port, _ := strconv.Atoi(os.Getenv("REDIS_PORT"))

	store = New(Config{
		Host:      os.Getenv("REDIS_ENDPOINT"),
		Port:      port,
		Username:  "",
		Password:  os.Getenv("REDIS_DB_PASSWORD"),
		URL:       "",
		Database:  0,
		Reset:     false,
		TLSConfig: nil,
	})
	return store
}

func SetCache(key string, val interface{}) {
	b, marshalErr := json.Marshal(val)
	if marshalErr != nil {
		zap.S().Warn("Cache: Marshal binary failed: " + marshalErr.Error())
	}
	err := store.Set(key, b, DefaultCacheExpireTime)
	if err != nil {
		zap.S().Warn("Set Cache error: " + err.Error())
	}
}

// set cache for any int value ex: views
func SetCacheInt(key string, val int) {
	err := store.Set(key, []byte(fmt.Sprint(val)), 0)
	if err != nil {
		zap.S().Warn("Set View cache error: ", err.Error())
	}
}

// get cache for any int value ex: views
func GetCacheInt(key string) int {
	bytes, err := store.Get(key)
	if err != nil {
		zap.S().Warn("Get View cache error: ", err.Error())
		return 0
	}

	count, err := strconv.Atoi(string(bytes))
	if err != nil && count != 0 {
		zap.S().Warn("Set View cache string to int error: ", err.Error())
		return 0
	}
	return count
}

func DeleteCache(key string) {
	err := store.Delete(key)
	if err != nil {
		zap.S().Warn("Delete cache failed, key=", key, ", error: ", err.Error())
	}
}

// func m() {
// 	redis.NewClient()
// }

// func FindAllCacheByPrefix(prefix string) ([]string, uint64) {
// 	var keys []string
// 	var err error
// 	var cursor uint64
// 	return store.Scan(prefix)
// 	keys, cursor, err = store.db.Scan(ctx, cursor, prefix+"*", 0).Result()
// 	if err != nil {
// 		zap.S().Warn("find all prefix=", prefix, " cache error:", err.Error())
// 	}
// 	return keys, cursor
// }

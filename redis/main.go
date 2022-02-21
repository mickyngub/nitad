package redis

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/birdglove2/nitad-backend/errors"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

var rdb *redis.Client

func SetCache(key string, value interface{}) errors.CustomError {
	log.Println("setting cached", key)
	expired := time.Hour
	b, marshalErr := json.Marshal(value)
	if marshalErr != nil {
		return errors.NewCacheError("Marshal binary failed: " + marshalErr.Error())
	}
	err := rdb.Set(ctx, key, b, expired).Err()
	if err != nil {
		return errors.NewCacheError("Set cache error: " + err.Error())
	}
	return nil
}

func GetCache(key string, dest interface{}) errors.CustomError {
	val, err := rdb.Get(ctx, key).Result()
	val, resultErr := CheckResult(val, err)
	if resultErr != nil {
		return resultErr
	}
	err = json.Unmarshal([]byte(val), dest)
	if err != nil {
		return errors.NewCacheError(err.Error())
	}
	return nil
}

func Init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     ":6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// subcate := subcategory.Subcategory{
	// 	Title: "cache subcategory test",
	// 	Image: "None",
	// }

	// log.Println("test 1")

	// sbyte, _ := MarshalBinary(subcate)
	// err := SetCache(subcate.ID.Hex(), sbyte)
	// if err != nil {
	// 	log.Println("error setting cache", err.Error())
	// }

	// log.Println("test 2")

	// val, err := GetCache(subcate.ID.Hex())
	// if err != nil {
	// 	log.Println("error getting cache", err.Error())
	// }
	// log.Println("test 3", val)

}

func CheckResult(val string, err error) (string, errors.CustomError) {
	switch {
	case err == redis.Nil:
		return "", errors.NewCacheError("Key does not exist")
	case err != nil:
		return "", errors.NewCacheError("Get failed" + err.Error())
	case val == "":
		return "", errors.NewCacheError("Value is empty")
	}
	return val, nil
}

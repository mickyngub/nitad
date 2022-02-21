package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/birdglove2/nitad-backend/errors"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

var rdb *redis.Client

func SetCache(key string, value interface{}) errors.CustomError {
	expired := time.Hour
	b, marshalErr := MarshalBinary(value)
	if marshalErr != nil {
		return errors.NewCacheError("Marshal binary failed: " + marshalErr.Error())
	}
	err := rdb.Set(ctx, key, b, expired).Err()
	if err != nil {
		return errors.NewCacheError("Set cache error: " + err.Error())
	}
	return nil
}

func GetCache(key string) (interface{}, errors.CustomError) {
	val, err := rdb.Get(ctx, key).Result()
	val, resultErr := CheckResult(val, err)
	if err != nil {
		return nil, resultErr
	}
	return val, nil
}

func MarshalBinary(i interface{}) ([]byte, error) {
	return json.Marshal(i)
}

func Init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     ":6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

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

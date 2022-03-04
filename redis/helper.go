package redis

import (
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/go-redis/redis/v8"
)

func CheckResult(b []byte, err error) (string, errors.CustomError) {
	val := string(b)
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

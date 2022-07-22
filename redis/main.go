package main

import (
	"context"
	"encoding/json"
	goRedis "github.com/go-redis/redis/v8"
	redisTrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/go-redis/redis.v8"
	"io/ioutil"
	"time"
)

const (
	operationTimeout = 100 * time.Millisecond
)

type redis struct {
	writeClient goRedis.UniversalClient
	readClient  goRedis.UniversalClient
}

func ProvideRedisCache() *redis {
	redis := new(redis)

	writeOptions := &goRedis.Options{
		Addr:         "localhost:6379",
		Password:     "",
		ReadTimeout:  operationTimeout,
		WriteTimeout: operationTimeout,
		MaxRetries:   1,
	}

	readOptions := &goRedis.Options{
		Addr:         "localhost:6379",
		Password:     "",
		ReadTimeout:  operationTimeout,
		WriteTimeout: operationTimeout,
		MaxRetries:   1,
	}

	redis.writeClient, redis.readClient = redisTrace.NewClient(writeOptions), redisTrace.NewClient(readOptions)

	return redis
}

func (r *redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.writeClient.Set(ctx, key, value, expiration).Err()
}

func main() {
	redis := ProvideRedisCache()

	expectedBody, _ := ioutil.ReadFile("test.json")

	var expectedResponse Model
	_ = json.Unmarshal(expectedBody, &expectedResponse)

	key1 := expectedResponse.Key1
	b, _ := json.Marshal(key1)
	_ = redis.Set(context.Background(), "key-1", string(b), 5*time.Minute)
	key2 := expectedResponse.Key2
	b, _ = json.Marshal(key2)
	_ = redis.Set(context.Background(), "key-2", string(b), 5*time.Minute)
}


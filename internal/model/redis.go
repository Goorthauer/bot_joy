package model

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	ErrRedisNotFound = errors.New("redis not found")
)

type RedisClient struct {
	client *redis.Client
}

func NewRedis() (*RedisClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Hour)
	defer cancel()
	r := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   2,
	})
	if err := r.Set(ctx, "key", "value", 0).Err(); err != nil {
		return nil, ErrRedisNotFound
	}
	return &RedisClient{client: r}, nil
}

func (rc *RedisClient) AddMemory(ctx context.Context, column, key string) {
	var count int
	countStr, err := rc.getMemory(ctx, key, column)
	if err != nil {
		fmt.Println("AddMemory get :", err)
	}
	if count, err = strconv.Atoi(countStr); err == nil {
		count++
	}

	if _, err = rc.client.HMSet(ctx, key, column, count).Result(); err != nil {
		fmt.Println("AddMemory HMSET :", err)
	}
}

func (rc *RedisClient) getMemory(ctx context.Context, key, field string) (string, error) {
	res, err := rc.client.HGet(ctx, key, field).Result()
	if err == nil {
		return fmt.Sprintf("%v", res), nil
	}
	rc.client.HDel(ctx, key, field)
	return "", err
}

func (rc *RedisClient) GetAll(ctx context.Context, key string) (string, error) {
	var outString string
	strCmd := rc.client.HGetAll(ctx, key)
	fmt.Println(strCmd)
	out, err := rc.client.HGetAll(ctx, key).Result()
	if err == nil {
		for i, v := range out {
			outString += fmt.Sprintf("%v: %v\n", i, v)
		}
		return outString, nil
	}
	return outString, err
}

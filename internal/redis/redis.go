package redis

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func Connect() (*redis.Client, error) {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	password := os.Getenv("REDIS_PASS")

	dbStr := os.Getenv("REDIS_DB")
	db, err := strconv.Atoi(dbStr)
	if err != nil {
		db = 0
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx := context.Background()
	status, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("cannot connect to Redis: %w", err)
	}

	if status != "PONG" {
		return nil, fmt.Errorf("connect Redis failed, take: %s", status)
	}

	return rdb, nil
}

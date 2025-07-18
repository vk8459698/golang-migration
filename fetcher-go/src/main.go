package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var ctx = context.Background()

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		data := fetchLiquidityPools()
		err := rdb.Set(ctx, "pools", data, 0).Err()
		if err != nil {
			fmt.Println("Redis error:", err)
		} else {
			fmt.Println("Updated pools in Redis.")
		}
	}
}

func fetchLiquidityPools() string {
	return `{"example_pool": {"token0": "ETH", "token1": "USDC"}}`
}

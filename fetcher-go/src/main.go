package main

import (
	"fetcher-go/src/common"
	"fmt"
	"time"
)

func initApp() {
	common.InitRedis()
}

func main() {
	initApp()

	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		data := fetchLiquidityPools()
		err := common.RedisClient.Set(common.Ctx, "pools", data, 0).Err()
		if err != nil {
			fmt.Println("[Main] Redis error:", err)
		} else {
			fmt.Println("[Main] Updated pools in Redis.")
		}
	}
}

func fetchLiquidityPools() string {
	return `{"example_pool": {"token0": "ETH", "token1": "USDC"}}`
}

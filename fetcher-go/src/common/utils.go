package common

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/redis/go-redis/v9"
	"os"
	"strings"
)

var Ctx = context.Background()
var RedisClient *redis.Client
var EvmChainIDs = []int{1, 10, 56, 137, 8453, 42161}

func InitRedis() {
	host := getEnv("REDIS_HOST", "localhost")
	port := getEnv("REDIS_PORT", "6379")

	RedisClient = redis.NewClient(&redis.Options{
		Addr: host + ":" + port,
	})

	err := checkRedisConnection()
	if err != nil {
		fmt.Println("[Utils] Redis not connected:", err)
	} else {
		fmt.Println("[Utils] Connected to Redis at", host+":"+port)
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func checkRedisConnection() error {
	return RedisClient.Ping(Ctx).Err()
}

func SortTokens(tokenA, tokenB string) [2]string {
	if strings.ToLower(tokenA) < strings.ToLower(tokenB) {
		return [2]string{tokenA, tokenB}
	}
	return [2]string{tokenB, tokenA}
}

func Checksum(addr string) string {
	return common.HexToAddress(addr).Hex()
}

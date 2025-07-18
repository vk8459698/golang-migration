package ledger

import (
	"encoding/json"
	"fetcher-go/src/common"
	"fetcher-go/src/dex/uniswap"
	"fmt"
	"strconv"
	"sync"
	"time"
)

type Ledger struct {
	quoteModelInterval time.Duration
}

func NewLedger() *Ledger {
	fmt.Println("[Ledger] Starting Ledger initialization")

	l := &Ledger{
		quoteModelInterval: 5 * time.Second,
	}

	l.startQuoteModel()

	fmt.Println("[Ledger] Completed Ledger initialization")
	return l
}

func (l *Ledger) startQuoteModel() {
	go func() {
		for {
			start := time.Now()
			fmt.Println("[Ledger] Starting pool data update cycle")

			err := l.updatePoolData()
			if err != nil {
				fmt.Println("[Ledger] Pool data update failed:", err)
			} else {
				fmt.Println("[Ledger] Completed pool data update cycle")
			}

			fmt.Printf("[Ledger] Pool data update duration: %v\n", time.Since(start))
			time.Sleep(l.quoteModelInterval)
		}
	}()
}

func (l *Ledger) updatePoolData() error {
	fmt.Println("[Ledger] Fetching and processing all chains and updating pools")

	ctx := common.Ctx
	redisClient := common.RedisClient
	chainConfigs := common.ChainConfigurations

	var wg sync.WaitGroup

	for chainId, chainConfig := range chainConfigs {
		chainIdStr := strconv.Itoa(chainId)
		chainGraph := map[string]map[string]bool{}
		var mu sync.Mutex

		dexWg := sync.WaitGroup{}

		for dexName, dexConfig := range chainConfig {
			dexWg.Add(1)

			go func(dexName string, dexConfig *common.DexConfiguration) {
				defer dexWg.Done()

				fmt.Printf("[Ledger] Fetching and storing data for dex: %s on chain: %s\n", dexName, chainIdStr)

				dexId := dexConfig.DexId

				switch dexId {
				case "UNISWAP_V3", "PANCAKESWAP_V3", "SUSHISWAP_V3":
					pools, err := uniswap.FetchPools(chainId, dexName, dexConfig)
					if err != nil {
						fmt.Printf("[Ledger] Error fetching pools for %s on chain: %s: %v\n", dexName, chainIdStr, err)
						return
					}

					for _, p := range pools {
						poolMap, ok := p.(map[string]interface{})
						if !ok {
							continue
						}

						id := common.Checksum(poolMap["id"].(string))
						if id == "0xe79d4ef25f12384f5a541d1e3485df69dd81447f" || id == "0xee9bf1d1e23784067bd7b0b3496f865038b766eb" {
							continue
						}

						token0 := common.Checksum(poolMap["token0"].(map[string]any)["id"].(string))
						token1 := common.Checksum(poolMap["token1"].(map[string]any)["id"].(string))
						symbol0 := poolMap["token0"].(map[string]interface{})["symbol"]
						symbol1 := poolMap["token1"].(map[string]interface{})["symbol"]

						tokens := common.SortTokens(token0, token1)
						tokenPairKey := fmt.Sprintf("%s:%s", tokens[0], tokens[1])

						info := map[string]interface{}{
							"poolId":    id,
							"dexId":     dexId,
							"token0":    token0,
							"token1":    token1,
							"symbol0":   symbol0,
							"symbol1":   symbol1,
							"liquidity": poolMap["liquidity"],
							"sqrtPrice": poolMap["sqrtPrice"],
							"fee":       poolMap["feeTier"],
							"tick":      poolMap["tick"],
							"ticks":     poolMap["ticks"],
							"gasPrice":  "0",
						}

						jsonInfo, _ := json.Marshal(info)
						redisClient.HSet(ctx, fmt.Sprintf("poolData:%s:%s", chainIdStr, tokenPairKey), id, jsonInfo)
						redisClient.HSet(ctx, fmt.Sprintf("tokenSymbol:%s", chainIdStr), token0, symbol0)
						redisClient.HSet(ctx, fmt.Sprintf("tokenSymbol:%s", chainIdStr), token1, symbol1)

						mu.Lock()
						if _, exists := chainGraph[token0]; !exists {
							chainGraph[token0] = map[string]bool{}
						}
						if _, exists := chainGraph[token1]; !exists {
							chainGraph[token1] = map[string]bool{}
						}
						chainGraph[token0][token1] = true
						chainGraph[token1][token0] = true
						mu.Unlock()
					}
				default:
					return
				}

				fmt.Printf("[Ledger] Fetched and stored data for dex: %s on chain: %s\n", dexName, chainIdStr)
			}(dexName, dexConfig)
		}

		wg.Add(1)
		go func(chainId int, chainIdStr string) {
			defer wg.Done()
			dexWg.Wait()

			fmt.Printf("[Ledger] Storing graph data in Redis for chain %s\n", chainIdStr)
			redisClient.Del(ctx, fmt.Sprintf("Graph:%s", chainIdStr))

			mu.Lock()
			for token, conns := range chainGraph {
				arr := []string{}
				for c := range conns {
					arr = append(arr, c)
				}
				j, _ := json.Marshal(arr)
				redisClient.HSet(ctx, fmt.Sprintf("Graph:%s", chainIdStr), token, j)
			}
			mu.Unlock()

			fmt.Printf("[Ledger] Stored graph data in Redis for chain %s\n", chainIdStr)
		}(chainId, chainIdStr)
	}

	wg.Wait()

	fmt.Println("[Ledger] Fetched and processed all chains and updated pools")
	return nil
}

package ledger

import (
	"encoding/json"
	"fetcher-go/src/common"
	"fetcher-go/src/dex/camelot"
	"fetcher-go/src/dex/maverick"
	"fetcher-go/src/dex/ramses"
	"fetcher-go/src/dex/ramseslegacy"
	"fetcher-go/src/dex/swaprv2"
	"fetcher-go/src/dex/traderjoe"
	"fetcher-go/src/dex/uniswap"
	"fetcher-go/src/dex/hyperion"
	"fetcher-go/src/dex/thala"
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
	ctx := common.Ctx
	redisClient := common.RedisClient
	chainConfigs := common.ChainConfigurations

	var wg sync.WaitGroup
	wg.Add(len(chainConfigs))

	for chainId, chainConfig := range chainConfigs {
		go func(chainId int, chainConfig map[string]*common.DexConfiguration) {
			defer wg.Done()

			chainIdStr := strconv.Itoa(chainId)
			chainGraph := make(map[string]map[string]bool)
			var mu sync.Mutex

			dexWg := sync.WaitGroup{}
			dexWg.Add(len(chainConfig))

			for dexName, dexConfig := range chainConfig {
				go func(dexName string, dexConfig *common.DexConfiguration) {
					defer dexWg.Done()

					fmt.Printf("[Ledger] Fetching data for dex: %s on chain: %s\n", dexName, chainIdStr)
					dexId := dexConfig.DexId

					switch dexId {
					case "UNISWAP_V3", "PANCAKESWAP_V3", "SUSHISWAP_V3":
						pools, err := uniswap.FetchPools(chainId, dexName, dexConfig)
						if err != nil {
							fmt.Printf("[Ledger] Error fetching pools for %s: %v\n", dexName, err)
							return
						}

						for _, p := range pools {
							poolMap, ok := p.(map[string]any)
							if !ok {
								continue
							}

							id := common.Checksum(poolMap["id"].(string))
							if id == "0xe79d4ef25f12384f5a541d1e3485df69dd81447f" || id == "0xee9bf1d1e23784067bd7b0b3496f865038b766eb" {
								continue
							}

							token0 := common.Checksum(poolMap["token0"].(map[string]any)["id"].(string))
							token1 := common.Checksum(poolMap["token1"].(map[string]any)["id"].(string))
							symbol0 := poolMap["token0"].(map[string]any)["symbol"]
							symbol1 := poolMap["token1"].(map[string]any)["symbol"]

							decimals0, err := strconv.Atoi(poolMap["token0"].(map[string]any)["decimals"].(string))
							if err != nil {
								continue
							}

							decimals1, err := strconv.Atoi(poolMap["token1"].(map[string]any)["decimals"].(string))
							if err != nil {
								continue
							}

							tickCurrent, err := strconv.Atoi(poolMap["tick"].(string))
							if err != nil {
								continue
							}

							fee, err := strconv.Atoi(poolMap["feeTier"].(string))
							if err != nil {
								continue
							}

							tickSpacing, ok := uniswap.TICK_SPACINGS[fee]
							if !ok {
								continue
							}

							tokens := common.SortTokens(token0, token1)
							tokenPairKey := fmt.Sprintf("%s:%s", tokens[0], tokens[1])

							info := map[string]any{
								"poolId":      id,
								"dexId":       dexId,
								"token0":      token0,
								"token1":      token1,
								"symbol0":     symbol0,
								"symbol1":     symbol1,
								"decimals0":   decimals0,
								"decimals1":   decimals1,
								"liquidity":   poolMap["liquidity"],
								"sqrtPrice":   poolMap["sqrtPrice"],
								"tickSpacing": tickSpacing,
								"tickCurrent": tickCurrent,
								"fee":         fee,
								"tick":        poolMap["tick"],
								"ticks":       poolMap["ticks"],
								"gasPrice":    "0",
							}

							jsonInfo, _ := json.Marshal(info)
							redisClient.HSet(ctx, fmt.Sprintf("poolData:%s:%s", chainIdStr, tokenPairKey), id, jsonInfo)
							redisClient.HSet(ctx, fmt.Sprintf("tokenSymbol:%s", chainIdStr), token0, symbol0)
							redisClient.HSet(ctx, fmt.Sprintf("tokenSymbol:%s", chainIdStr), token1, symbol1)

							mu.Lock()
							if _, exists := chainGraph[token0]; !exists {
								chainGraph[token0] = make(map[string]bool)
							}
							if _, exists := chainGraph[token1]; !exists {
								chainGraph[token1] = make(map[string]bool)
							}
							chainGraph[token0][token1] = true
							chainGraph[token1][token0] = true
							mu.Unlock()
						}

					case "CAMELOT_V3":
						pools, err := camelot.FetchPools(chainId, dexName, dexConfig)
						if err != nil {
							fmt.Printf("[Ledger] Error fetching pools for %s: %v\n", dexName, err)
							return
						}

						for _, p := range pools {
							poolMap, ok := p.(map[string]any)
							if !ok {
								continue
							}

							id := common.Checksum(poolMap["id"].(string))
							if id == "0xe79d4ef25f12384f5a541d1e3485df69dd81447f" || id == "0xee9bf1d1e23784067bd7b0b3496f865038b766eb" {
								continue
							}

							token0 := common.Checksum(poolMap["token0"].(map[string]any)["id"].(string))
							token1 := common.Checksum(poolMap["token1"].(map[string]any)["id"].(string))
							symbol0 := poolMap["token0"].(map[string]any)["symbol"]
							symbol1 := poolMap["token1"].(map[string]any)["symbol"]

							decimals0, err := strconv.Atoi(poolMap["token0"].(map[string]any)["decimals"].(string))
							if err != nil {
								continue
							}

							decimals1, err := strconv.Atoi(poolMap["token1"].(map[string]any)["decimals"].(string))
							if err != nil {
								continue
							}

							tickCurrent, err := strconv.Atoi(poolMap["tick"].(string))
							if err != nil {
								continue
							}

							tickSpacing, err := strconv.Atoi(poolMap["tickSpacing"].(string))
							if err != nil {
								continue
							}

							tokens := common.SortTokens(token0, token1)
							tokenPairKey := fmt.Sprintf("%s:%s", tokens[0], tokens[1])

							info := map[string]any{
								"poolId":      id,
								"dexId":       dexId,
								"token0":      token0,
								"token1":      token1,
								"symbol0":     symbol0,
								"symbol1":     symbol1,
								"decimals0":   decimals0,
								"decimals1":   decimals1,
								"liquidity":   poolMap["liquidity"],
								"sqrtPrice":   poolMap["sqrtPrice"],
								"tickSpacing": tickSpacing,
								"tickCurrent": tickCurrent,
								"ticks":       poolMap["ticks"],
								"fee":         0,
								"gasPrice":    "0",
							}

							jsonInfo, _ := json.Marshal(info)
							redisClient.HSet(ctx, fmt.Sprintf("poolData:%s:%s", chainIdStr, tokenPairKey), id, jsonInfo)
							redisClient.HSet(ctx, fmt.Sprintf("tokenSymbol:%s", chainIdStr), token0, symbol0)
							redisClient.HSet(ctx, fmt.Sprintf("tokenSymbol:%s", chainIdStr), token1, symbol1)

							mu.Lock()
							if _, exists := chainGraph[token0]; !exists {
								chainGraph[token0] = make(map[string]bool)
							}
							if _, exists := chainGraph[token1]; !exists {
								chainGraph[token1] = make(map[string]bool)
							}
							chainGraph[token0][token1] = true
							chainGraph[token1][token0] = true
							mu.Unlock()
						}

					case "RAMSES_CL":
						pools, err := ramses.FetchPools(chainId, dexName, dexConfig)
						if err != nil {
							fmt.Printf("[Ledger] Error fetching pools for %s: %v\n", dexName, err)
							return
						}

						for _, p := range pools {
							poolMap, ok := p.(map[string]any)
							if !ok {
								continue
							}

							id := common.Checksum(poolMap["id"].(string))
							if id == "0xe79d4ef25f12384f5a541d1e3485df69dd81447f" || id == "0xee9bf1d1e23784067bd7b0b3496f865038b766eb" {
								continue
							}

							token0 := common.Checksum(poolMap["token0"].(map[string]any)["id"].(string))
							token1 := common.Checksum(poolMap["token1"].(map[string]any)["id"].(string))
							symbol0 := poolMap["token0"].(map[string]any)["symbol"]
							symbol1 := poolMap["token1"].(map[string]any)["symbol"]

							decimals0, err := strconv.Atoi(poolMap["token0"].(map[string]any)["decimals"].(string))
							if err != nil {
								continue
							}

							decimals1, err := strconv.Atoi(poolMap["token1"].(map[string]any)["decimals"].(string))
							if err != nil {
								continue
							}

							tickCurrent, err := strconv.Atoi(poolMap["tick"].(string))
							if err != nil {
								continue
							}

							fee, err := strconv.Atoi(poolMap["feeTier"].(string))
							if err != nil {
								continue
							}

							tickSpacing, err := strconv.Atoi(poolMap["tickSpacing"].(string))
							if err != nil {
								continue
							}

							tokens := common.SortTokens(token0, token1)
							tokenPairKey := fmt.Sprintf("%s:%s", tokens[0], tokens[1])

							info := map[string]any{
								"poolId":      id,
								"dexId":       dexId,
								"token0":      token0,
								"token1":      token1,
								"symbol0":     symbol0,
								"symbol1":     symbol1,
								"decimals0":   decimals0,
								"decimals1":   decimals1,
								"liquidity":   poolMap["liquidity"],
								"sqrtPrice":   poolMap["sqrtPrice"],
								"tickSpacing": tickSpacing,
								"tickCurrent": tickCurrent,
								"fee":         fee,
								"tick":        poolMap["tick"],
								"ticks":       poolMap["ticks"],
								"gasPrice":    "0",
							}

							jsonInfo, _ := json.Marshal(info)
							redisClient.HSet(ctx, fmt.Sprintf("poolData:%s:%s", chainIdStr, tokenPairKey), id, jsonInfo)
							redisClient.HSet(ctx, fmt.Sprintf("tokenSymbol:%s", chainIdStr), token0, symbol0)
							redisClient.HSet(ctx, fmt.Sprintf("tokenSymbol:%s", chainIdStr), token1, symbol1)

							mu.Lock()
							if _, exists := chainGraph[token0]; !exists {
								chainGraph[token0] = make(map[string]bool)
							}
							if _, exists := chainGraph[token1]; !exists {
								chainGraph[token1] = make(map[string]bool)
							}
							chainGraph[token0][token1] = true
							chainGraph[token1][token0] = true
							mu.Unlock()
						}

					case "RAMSES_LEGACY":
						pools, err := ramseslegacy.FetchPools(chainId, dexName, dexConfig)
						if err != nil {
							fmt.Printf("[Ledger] Error fetching pools for %s: %v\n", dexName, err)
							return
						}

						for _, p := range pools {
							poolMap, ok := p.(map[string]any)
							if !ok {
								continue
							}

							id := common.Checksum(poolMap["id"].(string))
							if id == "0xe79d4ef25f12384f5a541d1e3485df69dd81447f" || id == "0xee9bf1d1e23784067bd7b0b3496f865038b766eb" {
								continue
							}

							token0 := common.Checksum(poolMap["token0"].(map[string]any)["id"].(string))
							token1 := common.Checksum(poolMap["token1"].(map[string]any)["id"].(string))
							symbol0 := poolMap["token0"].(map[string]any)["symbol"]
							symbol1 := poolMap["token1"].(map[string]any)["symbol"]

							decimals0, err := strconv.Atoi(poolMap["token0"].(map[string]any)["decimals"].(string))
							if err != nil {
								continue
							}

							decimals1, err := strconv.Atoi(poolMap["token1"].(map[string]any)["decimals"].(string))
							if err != nil {
								continue
							}

							stable := false
							if stableVal, exists := poolMap["stable"]; exists {
								stable, _ = stableVal.(bool)
							}

							reserve0 := "0"
							if reserve0Val, exists := poolMap["reserve0"]; exists {
								reserve0, _ = reserve0Val.(string)
							}

							reserve1 := "0"
							if reserve1Val, exists := poolMap["reserve1"]; exists {
								reserve1, _ = reserve1Val.(string)
							}

							tokens := common.SortTokens(token0, token1)
							tokenPairKey := fmt.Sprintf("%s:%s", tokens[0], tokens[1])

							info := map[string]any{
								"poolId":    id,
								"dexId":     dexId,
								"token0":    token0,
								"token1":    token1,
								"symbol0":   symbol0,
								"symbol1":   symbol1,
								"decimals0": decimals0,
								"decimals1": decimals1,
								"stable":    stable,
								"reserve0":  reserve0,
								"reserve1":  reserve1,
								"gasPrice":  "0",
							}

							jsonInfo, _ := json.Marshal(info)
							redisClient.HSet(ctx, fmt.Sprintf("poolData:%s:%s", chainIdStr, tokenPairKey), id, jsonInfo)
							redisClient.HSet(ctx, fmt.Sprintf("tokenSymbol:%s", chainIdStr), token0, symbol0)
							redisClient.HSet(ctx, fmt.Sprintf("tokenSymbol:%s", chainIdStr), token1, symbol1)

							mu.Lock()
							if _, exists := chainGraph[token0]; !exists {
								chainGraph[token0] = make(map[string]bool)
							}
							if _, exists := chainGraph[token1]; !exists {
								chainGraph[token1] = make(map[string]bool)
							}
							chainGraph[token0][token1] = true
							chainGraph[token1][token0] = true
							mu.Unlock()
						}

					case "SWAPR_V2":
						pools, err := swaprv2.FetchPools(chainId, dexName, dexConfig)
						if err != nil {
							fmt.Printf("[Ledger] Error fetching pools for %s: %v\n", dexName, err)
							return
						}

						for _, p := range pools {
							poolMap, ok := p.(map[string]any)
							if !ok {
								continue
							}

							id := common.Checksum(poolMap["id"].(string))
							if id == "0xe79d4ef25f12384f5a541d1e3485df69dd81447f" || id == "0xee9bf1d1e23784067bd7b0b3496f865038b766eb" {
								continue
							}

							token0 := common.Checksum(poolMap["token0"].(map[string]any)["id"].(string))
							token1 := common.Checksum(poolMap["token1"].(map[string]any)["id"].(string))
							symbol0 := poolMap["token0"].(map[string]any)["symbol"]
							symbol1 := poolMap["token1"].(map[string]any)["symbol"]

							decimals0, err := strconv.Atoi(poolMap["token0"].(map[string]any)["decimals"].(string))
							if err != nil {
								continue
							}

							decimals1, err := strconv.Atoi(poolMap["token1"].(map[string]any)["decimals"].(string))
							if err != nil {
								continue
							}

							reserve0 := "0"
							if reserve0Val, exists := poolMap["reserve0"]; exists {
								reserve0, _ = reserve0Val.(string)
							}

							reserve1 := "0"
							if reserve1Val, exists := poolMap["reserve1"]; exists {
								reserve1, _ = reserve1Val.(string)
							}

							tokens := common.SortTokens(token0, token1)
							tokenPairKey := fmt.Sprintf("%s:%s", tokens[0], tokens[1])

							info := map[string]any{
								"poolId":    id,
								"dexId":     dexId,
								"token0":    token0,
								"token1":    token1,
								"symbol0":   symbol0,
								"symbol1":   symbol1,
								"decimals0": decimals0,
								"decimals1": decimals1,
								"reserve0":  reserve0,
								"reserve1":  reserve1,
								"fee":       0,
								"gasPrice":  "0",
							}

							jsonInfo, _ := json.Marshal(info)
							redisClient.HSet(ctx, fmt.Sprintf("poolData:%s:%s", chainIdStr, tokenPairKey), id, jsonInfo)
							redisClient.HSet(ctx, fmt.Sprintf("tokenSymbol:%s", chainIdStr), token0, symbol0)
							redisClient.HSet(ctx, fmt.Sprintf("tokenSymbol:%s", chainIdStr), token1, symbol1)

							mu.Lock()
							if _, exists := chainGraph[token0]; !exists {
								chainGraph[token0] = make(map[string]bool)
							}
							if _, exists := chainGraph[token1]; !exists {
								chainGraph[token1] = make(map[string]bool)
							}
							chainGraph[token0][token1] = true
							chainGraph[token1][token0] = true
							mu.Unlock()
						}

					case "MAVERICK_V2":
						pools, err := maverick.FetchPools(chainId, dexName, dexConfig)
						if err != nil {
							fmt.Printf("[Ledger] Error fetching pools for %s: %v\n", dexName, err)
							return
						}

						for _, p := range pools {
							poolMap, ok := p.(*maverick.PoolDetails)
							if !ok {
								continue
							}

							id := common.Checksum(poolMap.ID)
							if id == "0xe79d4ef25f12384f5a541d1e3485df69dd81447f" || id == "0xee9bf1d1e23784067bd7b0b3496f865038b766eb" {
								continue
							}

							token0 := common.Checksum(poolMap.TokenA)
							token1 := common.Checksum(poolMap.TokenB)
							// symbol0 := "UNKNOWN" // Commented out as in the original JS
							// symbol1 := "UNKNOWN" // Commented out as in the original JS
							reserve0 := poolMap.State.ReserveA
							reserve1 := poolMap.State.ReserveB
							feeAIn := poolMap.FeeAIn
							feeBIn := poolMap.FeeBIn
							protocolFeeA := poolMap.ProtocolFeeA
							protocolFeeB := poolMap.ProtocolFeeB
							sqrtPrice := poolMap.SqrtPrice
							liquidity := poolMap.Liquidity
							sqrtLowerPrice := poolMap.SqrtLowerPrice
							sqrtUpperPrice := poolMap.SqrtUpperPrice

							tokens := common.SortTokens(token0, token1)
							tokenPairKey := fmt.Sprintf("%s:%s", tokens[0], tokens[1])

							info := map[string]any{
								"poolId":         id,
								"dexId":          dexId,
								"token0":         token0,
								"token1":         token1,
								"feeAIn":         feeAIn,
								"feeBIn":         feeBIn,
								"reserve0":       reserve0,
								"reserve1":       reserve1,
								"protocolFeeA":   protocolFeeA,
								"protocolFeeB":   protocolFeeB,
								"sqrtPrice":      sqrtPrice,
								"liquidity":      liquidity,
								"sqrtLowerPrice": sqrtLowerPrice,
								"sqrtUpperPrice": sqrtUpperPrice,
								"gasPrice":       "0",
							}

							jsonInfo, _ := json.Marshal(info)
							redisClient.HSet(ctx, fmt.Sprintf("poolData:%s:%s", chainIdStr, tokenPairKey), id, jsonInfo)
							// redisClient.HSet(ctx, fmt.Sprintf("tokenSymbol:%s", chainIdStr), token0, symbol0) // Commented out as in the original JS
							// redisClient.HSet(ctx, fmt.Sprintf("tokenSymbol:%s", chainIdStr), token1, symbol1) // Commented out as in the original JS

							mu.Lock()
							if _, exists := chainGraph[token0]; !exists {
								chainGraph[token0] = make(map[string]bool)
							}
							if _, exists := chainGraph[token1]; !exists {
								chainGraph[token1] = make(map[string]bool)
							}
							chainGraph[token0][token1] = true
							chainGraph[token1][token0] = true
							mu.Unlock()
						}
					case "HYPERION":
						pools, err := hyperion.FetchPools(chainId, dexName, dexConfig)
						if err != nil {
							fmt.Printf("[Ledger] Error fetching pools for %s: %v\n", dexName, err)
							return
						}

						for _, p := range pools {
							poolMap, ok := p.(map[string]any)
							if !ok {
								continue
							}

							id := common.Checksum(poolMap["poolId"].(string))
							if id == "0xe79d4ef25f12384f5a541d1e3485df69dd81447f" || id == "0xee9bf1d1e23784067bd7b0b3496f865038b766eb" {
								continue
							}

							token0 := common.Checksum(poolMap["token0"].(string))
							token1 := common.Checksum(poolMap["token1"].(string))
							symbol0 := poolMap["symbol0"]
							symbol1 := poolMap["symbol1"]

							decimals0, ok := poolMap["decimals0"].(int)
							if !ok {
								continue
							}

							decimals1, ok := poolMap["decimals1"].(int)
							if !ok {
								continue
							}

							tickCurrent, ok := poolMap["tickCurrent"].(int)
							if !ok {
								continue
							}

							fee, ok := poolMap["fee"].(int)
							if !ok {
								continue
							}

							tokens := common.SortTokens(token0, token1)
							tokenPairKey := fmt.Sprintf("%s:%s", tokens[0], tokens[1])

							info := map[string]any{
								"poolId":     id,
								"dexId":      dexId,
								"token0":     token0,
								"token1":     token1,
								"symbol0":    symbol0,
								"symbol1":    symbol1,
								"decimals0":  decimals0,
								"decimals1":  decimals1,
								"gasPrice":   "0",
								"tokenType0": poolMap["tokenType0"],
								"tokenType1": poolMap["tokenType1"],
								"fee":        fee,
								"ticks":      poolMap["ticks"],
								"sqrtPrice":  poolMap["sqrtPrice"],
								"tickCurrent": tickCurrent,
								"liquidity":  poolMap["liquidity"],
							}

							jsonInfo, _ := json.Marshal(info)
							redisClient.HSet(ctx, fmt.Sprintf("poolData:%s:%s", chainIdStr, tokenPairKey), id, jsonInfo)
							redisClient.HSet(ctx, fmt.Sprintf("tokenSymbol:%s", chainIdStr), token0, symbol0)
							redisClient.HSet(ctx, fmt.Sprintf("tokenSymbol:%s", chainIdStr), token1, symbol1)

							mu.Lock()
							if _, exists := chainGraph[token0]; !exists {
								chainGraph[token0] = make(map[string]bool)
							}
							if _, exists := chainGraph[token1]; !exists {
								chainGraph[token1] = make(map[string]bool)
							}
							chainGraph[token0][token1] = true
							chainGraph[token1][token0] = true
							mu.Unlock()
						}
					case "THALA_V2":
						pools, err := thala.FetchPools(chainId, dexName, dexConfig)
						if err != nil {
							fmt.Printf("[Ledger] Error fetching pools for %s: %v\n", dexName, err)
							return
						}

						for _, p := range pools {
							poolMap, ok := p.(map[string]any)
							if !ok {
								continue
							}

							id := common.Checksum(poolMap["poolId"].(string))
							if id == "0x7845c59627bf2ecd0a8d4e2e83f0008546868442c3027060f042398578213164" {
								continue
							}

							token0 := common.Checksum(poolMap["token0"].(string))
							token1 := common.Checksum(poolMap["token1"].(string))
							symbol0 := poolMap["symbol0"]
							symbol1 := poolMap["symbol1"]

							decimals0, ok := poolMap["decimals0"].(int)
							if !ok {
								continue
							}

							decimals1, ok := poolMap["decimals1"].(int)
							if !ok {
								continue
							}

							balances := poolMap["balances"]
							ampFactor := poolMap["ampFactor"]
							index0 := poolMap["index0"]
							index1 := poolMap["index1"]
							weights := poolMap["weights"]
							poolType := poolMap["poolType"]
							fee := poolMap["fee"]
							tokenType0 := poolMap["tokenType0"]
							tokenType1 := poolMap["tokenType1"]
							dexIdNew := poolMap["dexId"].(string)

							tokens := common.SortTokens(token0, token1)
							tokenPairKey := fmt.Sprintf("%s:%s", tokens[0], tokens[1])

							info := map[string]any{
								"poolId":     id,
								"dexId":      dexIdNew,
								"token0":     token0,
								"token1":     token1,
								"symbol0":    symbol0,
								"symbol1":    symbol1,
								"decimals0":  decimals0,
								"decimals1":  decimals1,
								"tokenType0": tokenType0,
								"tokenType1": tokenType1,
								"balances":   balances,
								"ampFactor":  ampFactor,
								"index0":     index0,
								"index1":     index1,
								"weights":    weights,
								"poolType":   poolType,
								"fee":        fee,
								"gasPrice":   "0",
							}

							jsonInfo, _ := json.Marshal(info)
							redisClient.HSet(ctx, fmt.Sprintf("poolData:%s:%s", chainIdStr, tokenPairKey), id, jsonInfo)
							redisClient.HSet(ctx, fmt.Sprintf("tokenSymbol:%s", chainIdStr), token0, symbol0)
							redisClient.HSet(ctx, fmt.Sprintf("tokenSymbol:%s", chainIdStr), token1, symbol1)

							mu.Lock()
							if _, exists := chainGraph[token0]; !exists {
								chainGraph[token0] = make(map[string]bool)
							}
							if _, exists := chainGraph[token1]; !exists {
								chainGraph[token1] = make(map[string]bool)
							}
							chainGraph[token0][token1] = true
							chainGraph[token1][token0] = true
							mu.Unlock()
						}
					case "TRADERJOE":
						pools, err := traderjoe.FetchPools(chainId, dexName, dexConfig)
						if err != nil {
							fmt.Printf("[Ledger] Error fetching pools for %s: %v\n", dexName, err)
							return
						}

						for _, p := range pools {
							poolMap, ok := p.(map[string]any)
							if !ok {
								continue
							}

							id := common.Checksum(poolMap["id"].(string))
							if id == "0xe79d4ef25f12384f5a541d1e3485df69dd81447f" || id == "0xee9bf1d1e23784067bd7b0b3496f865038b766eb" {
								continue
							}

							token0 := common.Checksum(poolMap["token0"].(map[string]any)["id"].(string))
							token1 := common.Checksum(poolMap["token1"].(map[string]any)["id"].(string))
							symbol0 := poolMap["token0"].(map[string]any)["symbol"]
							symbol1 := poolMap["token1"].(map[string]any)["symbol"]

							decimals0, err := strconv.Atoi(poolMap["token0"].(map[string]any)["decimals"].(string))
							if err != nil {
								continue
							}

							decimals1, err := strconv.Atoi(poolMap["token1"].(map[string]any)["decimals"].(string))
							if err != nil {
								continue
							}

							reserve0 := "0"
							if reserve0Val, exists := poolMap["reserve0"]; exists {
								reserve0, _ = reserve0Val.(string)
							}

							reserve1 := "0"
							if reserve1Val, exists := poolMap["reserve1"]; exists {
								reserve1, _ = reserve1Val.(string)
							}

							// Skip pools with zero reserves (as done in the JavaScript version)
							if reserve0 == "0" || reserve1 == "0" {
								continue
							}

							tokens := common.SortTokens(token0, token1)
							tokenPairKey := fmt.Sprintf("%s:%s", tokens[0], tokens[1])

							info := map[string]any{
								"poolId":    id,
								"dexId":     dexId,
								"token0":    token0,
								"token1":    token1,
								"symbol0":   symbol0,
								"symbol1":   symbol1,
								"decimals0": decimals0,
								"decimals1": decimals1,
								"reserve0":  reserve0,
								"reserve1":  reserve1,
								"fee":       0,
								"gasPrice":  "0",
							}

							jsonInfo, _ := json.Marshal(info)
							redisClient.HSet(ctx, fmt.Sprintf("poolData:%s:%s", chainIdStr, tokenPairKey), id, jsonInfo)
							redisClient.HSet(ctx, fmt.Sprintf("tokenSymbol:%s", chainIdStr), token0, symbol0)
							redisClient.HSet(ctx, fmt.Sprintf("tokenSymbol:%s", chainIdStr), token1, symbol1)

							mu.Lock()
							if _, exists := chainGraph[token0]; !exists {
								chainGraph[token0] = make(map[string]bool)
							}
							if _, exists := chainGraph[token1]; !exists {
								chainGraph[token1] = make(map[string]bool)
							}
							chainGraph[token0][token1] = true
							chainGraph[token1][token0] = true
							mu.Unlock()
						}

					default:
					}
					fmt.Printf("[Ledger] Fetched data for dex: %s on chain: %s\n", dexName, chainIdStr)
				}(dexName, dexConfig)
			}

			dexWg.Wait()

			fmt.Printf("[Ledger] Storing graph data in Redis for chain %s\n", chainIdStr)
			redisClient.Del(ctx, fmt.Sprintf("Graph:%s", chainIdStr))

			mu.Lock()
			for token, conns := range chainGraph {
				arr := make([]string, 0, len(conns))
				for c := range conns {
					arr = append(arr, c)
				}
				j, _ := json.Marshal(arr)
				redisClient.HSet(ctx, fmt.Sprintf("Graph:%s", chainIdStr), token, j)
			}
			mu.Unlock()
			fmt.Printf("[Ledger] Stored graph data in Redis for chain %s\n", chainIdStr)

		}(chainId, chainConfig)
	}

	wg.Wait()
	return nil
}

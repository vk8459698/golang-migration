package thala

import (
	"bytes"
	"context"
	"encoding/json"
	"fetcher-go/src/common"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/machinebox/graphql"
)

const (
	VIEW_API_URL      = "https://api.mainnet.aptoslabs.com/v1/view"
	POOL_MODULE_ADDR  = "0x007730cd28ee1cdc9e999336cbc430f99e7c44397c0aa77516f6f23a78559bb5"
	LENS_MODULE_ADDR  = "0xff1ac437457a839f7d07212d789b85dd77b3df00f59613fcba02388464bfcacb"
	POOL_TYPE_STABLE  = 100
	POOL_TYPE_WEIGHTED = 101
	EPSILON           = 0.000001
	MAX_LOOP_LIMIT    = 100
)

var APTOS_API_KEYS = []string{
	"aptoslabs_e3nyJCWpqE7_88xNQBDzVMcDzAXHTVzegHbQvZgjfCHrS",
	"aptoslabs_gmnkiejfwjC_JhoLSyHF66LAebPTZTVFkYKKK4fF3oJgC",
	"aptoslabs_5JtKHM5kBEW_KxYYzvxPnFCXh4mGJN5iLdwY5ur79yBax",
	"aptoslabs_CNnjSqzUu4_MLukRhm6U3hze5CUrW8PCtZqcsDmiYAWz",
	"aptoslabs_RvxLfaBYjpE_JXUedZha9QF8X87npNLWFodr6QLNpwmqX",
	"aptoslabs_945frd1sp5n_NVBS2xGq9w9jp3aBUK4mczX9DyxjC4rUa",
	"aptoslabs_8E2vZUaH9h6_JJa4J6STvcWa4jmP9JdFrBRayzxpdbYZ7",
	"aptoslabs_71qfkdRPB1d_A5zfUDzk1ibsySPqwTgyMCwAoVY2nKGQs",
	"aptoslabs_AgDp1EQLNU4_5pSKocSy8zYEYTyPcqtNJS4xX9QNvbNQf",
	"aptoslabs_RXTt6QDwVHA_JvbLAqvaVGcdVK3xppZwy9hNabJPx2UvW",
}

var (
	currentApiKeyIndex int
	lock               sync.Mutex
)

type PoolInfo struct {
	Pool struct {
		Inner string `json:"inner"`
	} `json:"pool"`
	Balances        []string `json:"balances"`
	AmpFactorOpt    *struct {
		Vec []string `json:"vec"`
	} `json:"amp_factor_opt"`
	SwapFeeBps      json.Number `json:"swap_fee_bps"`
	PoolType        json.Number `json:"pool_type"`
	WeightsOpt      *struct {
		Vec [][]string `json:"vec"`
	} `json:"weights_opt"`
	AssetsMetadata  []struct {
		Inner string `json:"inner"`
	} `json:"assets_metadata"`
}

type FungibleAssetMetadata struct {
	Symbol    string      `json:"symbol"`
	Name      string      `json:"name"`
	Decimals  json.Number `json:"decimals"`  // Changed from string to json.Number
	AssetType string      `json:"asset_type"`
}

type GraphQLResponse struct {
	FungibleAssetMetadata []FungibleAssetMetadata `json:"fungible_asset_metadata"`
}

func getNextApiKey() string {
	lock.Lock()
	defer lock.Unlock()
	
	key := APTOS_API_KEYS[currentApiKeyIndex]
	currentApiKeyIndex = (currentApiKeyIndex + 1) % len(APTOS_API_KEYS)
	return key
}

func normalizeAptosAddress(address string) string {
	normalized := strings.ToLower(address)
	normalized = strings.TrimPrefix(normalized, "0x")
	
	for len(normalized) < 64 {
		normalized = "0" + normalized
	}
	
	return "0x" + normalized
}

func fetchPoolAddresses() ([]string, error) {
	apiKey := getNextApiKey()
	
	payload := map[string]interface{}{
		"function":       fmt.Sprintf("%s::pool::pools", POOL_MODULE_ADDR),
		"type_arguments": []string{},
		"arguments":      []interface{}{},
	}
	
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %v", err)
	}
	
	req, err := http.NewRequest("POST", VIEW_API_URL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("[UtilsThala] Failed to fetch pool addresses: %s\n", string(body))
		return []string{}, nil
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}
	
	var result [][]map[string]string
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}
	
	if len(result) == 0 || len(result[0]) == 0 {
		return []string{}, nil
	}
	
	addresses := make([]string, len(result[0]))
	for i, obj := range result[0] {
		addresses[i] = obj["inner"]
	}
	
	return addresses, nil
}

func fetchPoolInfos(poolAddresses []string) ([]PoolInfo, error) {
	apiKey := getNextApiKey()
	
	wrappedAddresses := make([]map[string]string, len(poolAddresses))
	for i, addr := range poolAddresses {
		wrappedAddresses[i] = map[string]string{"inner": addr}
	}
	
	payload := map[string]interface{}{
		"function":       fmt.Sprintf("%s::lens::get_pools_info", LENS_MODULE_ADDR),
		"type_arguments": []string{},
		"arguments":      []interface{}{wrappedAddresses},
	}
	
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %v", err)
	}
	
	req, err := http.NewRequest("POST", VIEW_API_URL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("[UtilsThala] Failed to fetch pool infos: %s\n", string(body))
		return []PoolInfo{}, nil
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}
	
	var result [][]PoolInfo
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}
	
	if len(result) == 0 {
		return []PoolInfo{}, nil
	}
	
	return result[0], nil
}

func chunkArray(arr []string, chunkSize int) [][]string {
	var chunks [][]string
	for i := 0; i < len(arr); i += chunkSize {
		end := i + chunkSize
		if end > len(arr) {
			end = len(arr)
		}
		chunks = append(chunks, arr[i:end])
	}
	return chunks
}

func fetchAllThalaPools() ([]PoolInfo, error) {
	poolAddresses, err := fetchPoolAddresses()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pool addresses: %v", err)
	}
	
	if len(poolAddresses) == 0 {
		fmt.Println("[UtilsThala] No pool addresses found.")
		return []PoolInfo{}, nil
	}
	
	chunks := chunkArray(poolAddresses, 10)
	
	type result struct {
		pools []PoolInfo
		err   error
	}
	
	results := make(chan result, len(chunks))
	var wg sync.WaitGroup
	
	for _, batch := range chunks {
		wg.Add(1)
		go func(batch []string) {
			defer wg.Done()
			pools, err := fetchPoolInfos(batch)
			results <- result{pools: pools, err: err}
		}(batch)
	}
	
	wg.Wait()
	close(results)
	
	var allPools []PoolInfo
	for res := range results {
		if res.err == nil && len(res.pools) > 0 {
			allPools = append(allPools, res.pools...)
		}
	}
	
	return allPools, nil
}

func getY(xp []float64, x float64, a float64, i int, j int) (float64, error) {
	d, err := getD(xp, a)
	if err != nil {
		return 0, err
	}
	
	n := len(xp)
	ann := a * float64(n)
	
	c := d
	s := 0.0
	
	for k := 0; k < n; k++ {
		if k == j {
			continue
		}
		
		var x_k float64
		if k == i {
			x_k = x
		} else {
			x_k = xp[k]
		}
		
		s += x_k
		c = (c * d) / (x_k * float64(n))
	}
	
	c = (c * d) / (ann * float64(n))
	b := s + d/ann
	
	y := d
	for k := 0; k < MAX_LOOP_LIMIT; k++ {
		prev_y := y
		y = (y*y + c) / (2*y + b - d)
		if math.Abs(y-prev_y) < EPSILON {
			return y, nil
		}
	}
	
	return 0, fmt.Errorf("not converged in getY, xp: %v, x: %f, a: %f, i: %d, j: %d", xp, x, a, i, j)
}

func getD(xp []float64, a float64) (float64, error) {
	n := len(xp)
	
	s := 0.0
	for _, val := range xp {
		s += val
	}
	
	if s == 0 {
		return 0, nil
	}
	
	var prev_d float64
	d := s
	ann := a * float64(n)
	
	for i := 0; i < MAX_LOOP_LIMIT; i++ {
		dp := d
		
		for j := 0; j < n; j++ {
			dp = (dp * d) / (xp[j] * float64(n))
		}
		
		prev_d = d
		d = ((ann*s + float64(n)*dp) * d) / ((ann-1)*d + (float64(n)+1)*dp)
		if math.Abs(prev_d-d) < EPSILON {
			return d, nil
		}
	}
	
	return 0, fmt.Errorf("not converged in getD, xp: %v, a: %f", xp, a)
}

func getThalaAmountOutStable(pool map[string]interface{}, amountIn string, indexIn int, indexOut int) (string, error) {
	feeBpsFloat, _ := pool["fee"].(float64)
	feeBps := int(feeBpsFloat)
	amp, _ := pool["ampFactor"].(float64)
	balancesInterface, _ := pool["balances"].([]interface{})
	
	balances := make([]float64, len(balancesInterface))
	for i, b := range balancesInterface {
		if balanceStr, ok := b.(string); ok {
			if val, err := strconv.ParseFloat(balanceStr, 64); err == nil {
				balances[i] = val
			}
		}
	}
	
	amountInFloat, err := strconv.ParseFloat(amountIn, 64)
	if err != nil {
		return "0", err
	}
	
	amountInFloat = amountInFloat * (1 - float64(feeBps)/10000)
	
	newX := balances[indexIn] + amountInFloat
	newY, err := getY(balances, newX, amp, indexIn, indexOut)
	if err != nil {
		return "0", err
	}
	
	amountOut := balances[indexOut] - newY
	return strconv.FormatInt(int64(math.Round(amountOut)), 10), nil
}

func getThalaAmountOutWeighted(pool map[string]interface{}, amountIn string, indexIn int, indexOut int) (string, error) {
	balancesInterface, _ := pool["balances"].([]interface{})
	weightsInterface, _ := pool["weights"].([]interface{})
	
	var bI, bO float64
	if balanceStr, ok := balancesInterface[indexIn].(string); ok {
		bI, _ = strconv.ParseFloat(balanceStr, 64)
	}
	if balanceStr, ok := balancesInterface[indexOut].(string); ok {
		bO, _ = strconv.ParseFloat(balanceStr, 64)
	}
	
	var wI, wO float64
	if weightFloat, ok := weightsInterface[indexIn].(float64); ok {
		wI = weightFloat
	}
	if weightFloat, ok := weightsInterface[indexOut].(float64); ok {
		wO = weightFloat
	}
	
	aI, err := strconv.ParseFloat(amountIn, 64)
	if err != nil {
		return "0", err
	}
	
	feeBpsFloat, _ := pool["fee"].(float64)
	sF := feeBpsFloat / 10000
	
	denom := bI + aI*(1-sF)
	amountOut := bO * (1 - math.Pow(bI/denom, wI/wO))
	
	return strconv.FormatInt(int64(math.Round(amountOut)), 10), nil
}

func getThalaAmountOut(pool map[string]interface{}, amountIn string, indexIn int, indexOut int) (string, error) {
	poolTypeFloat, _ := pool["poolType"].(float64)
	poolType := int(poolTypeFloat)
	
	if poolType == POOL_TYPE_STABLE {
		return getThalaAmountOutStable(pool, amountIn, indexIn, indexOut)
	}
	
	if poolType == POOL_TYPE_WEIGHTED {
		return getThalaAmountOutWeighted(pool, amountIn, indexIn, indexOut)
	}
	
	return "0", fmt.Errorf("invalid thala swap pool type")
}

func FetchPools(chainId int, dexName string, dexConfig *common.DexConfiguration) ([]any, error) {
	allPoolsData, err := fetchAllThalaPools()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch thala pools: %v", err)
	}
	
	query := `query GetFungibleAssetInfo($in: [String!]) {
		fungible_asset_metadata(
			where: { asset_type: { _in: $in } }
		) {
			symbol
			name
			decimals
			asset_type
		}
	}`
	
	tokenTypeSet := make(map[string]bool)
	for _, pool := range allPoolsData {
		for _, tokenMeta := range pool.AssetsMetadata {
			tokenTypeSet[normalizeAptosAddress(tokenMeta.Inner)] = true
		}
	}
	
	tokenTypes := make([]string, 0, len(tokenTypeSet))
	for tokenType := range tokenTypeSet {
		tokenTypes = append(tokenTypes, tokenType)
	}
	
	client := graphql.NewClient("https://api.mainnet.aptoslabs.com/v1/graphql")
	req := graphql.NewRequest(query)
	req.Var("in", tokenTypes)
	
	var response GraphQLResponse
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = client.Run(ctx, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch fungible asset metadata: %v", err)
	}
	
	metadataMap := make(map[string]FungibleAssetMetadata)
	for _, token := range response.FungibleAssetMetadata {
		key := normalizeAptosAddress(token.AssetType)
		metadataMap[key] = token
	}
	
	var results []any
	
	for _, pool := range allPoolsData {
		id := normalizeAptosAddress(pool.Pool.Inner)
		
		if id == "0x7845c59627bf2ecd0a8d4e2e83f0008546868442c3027060f042398578213164" {
			continue
		}
		
		balances := pool.Balances
		var ampFactor float64
		if pool.AmpFactorOpt != nil && len(pool.AmpFactorOpt.Vec) > 0 {
			if val, err := strconv.ParseFloat(pool.AmpFactorOpt.Vec[0], 64); err == nil {
				ampFactor = val
			}
		}
		
		fee := 0
		if feeStr := string(pool.SwapFeeBps); feeStr != "" {
			if val, err := strconv.Atoi(feeStr); err == nil {
				fee = val
			}
		}
		
		poolType := 0
		if typeStr := string(pool.PoolType); typeStr != "" {
			if val, err := strconv.Atoi(typeStr); err == nil {
				poolType = val
			}
		}
		
		var weights []float64
		if pool.WeightsOpt != nil && len(pool.WeightsOpt.Vec) > 0 {
			weights = make([]float64, len(pool.WeightsOpt.Vec[0]))
			for i, w := range pool.WeightsOpt.Vec[0] {
				if val, err := strconv.ParseFloat(w, 64); err == nil {
					weights[i] = val
				}
			}
		}
		
		if poolType != 100 && poolType != 101 {
			continue
		}
		
		var dexIdNew string
		if poolType == 100 {
			dexIdNew = "THALA_V2_STABLE"
		}
		if poolType == 101 {
			dexIdNew = "THALA_V2_WEIGHTED"
		}
		
		assets := pool.AssetsMetadata
		
		for i := 0; i < len(assets); i++ {
			token0 := normalizeAptosAddress(assets[i].Inner)
			
			meta0, exists := metadataMap[token0]
			if !exists {
				continue
			}
			
			symbol0 := meta0.Symbol
			decimals0, err := strconv.Atoi(string(meta0.Decimals))  // Convert json.Number to string first
			if err != nil {
				continue
			}
			
			var assetType0 *string
			if token0 == "0x000000000000000000000000000000000000000000000000000000000000000a" {
				aptosType := "0x1::aptos_coin::AptosCoin"
				assetType0 = &aptosType
			}
			
			for j := i + 1; j < len(assets); j++ {
				token1 := normalizeAptosAddress(assets[j].Inner)
				
				meta1, exists := metadataMap[token1]
				if !exists {
					continue
				}
				
				symbol1 := meta1.Symbol
				decimals1, err := strconv.Atoi(string(meta1.Decimals))  // Convert json.Number to string first
				if err != nil {
					continue
				}
				
				var assetType1 *string
				if token1 == "0x000000000000000000000000000000000000000000000000000000000000000a" {
					aptosType := "0x1::aptos_coin::AptosCoin"
					assetType1 = &aptosType
				}
				
				poolInfo := map[string]any{
					"poolId":     id,
					"dexId":      dexIdNew,
					"token0":     token0,
					"token1":     token1,
					"symbol0":    symbol0,
					"symbol1":    symbol1,
					"decimals0":  decimals0,
					"decimals1":  decimals1,
					"tokenType0": assetType0,
					"tokenType1": assetType1,
					"balances":   balances,
					"ampFactor":  ampFactor,
					"index0":     i,
					"index1":     j,
					"weights":    weights,
					"poolType":   poolType,
					"fee":        fee,
					"gasPrice":   "0",
				}
				
				results = append(results, poolInfo)
			}
		}
	}
	
	return results, nil
}

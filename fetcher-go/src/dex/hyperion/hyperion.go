package hyperion

import (
	"context"
	"fetcher-go/src/common"
	"fmt"
	"github.com/machinebox/graphql"
	"math/big"
	"strconv"
	"strings"
	"sync"
)

type PoolDetails struct {
	ID     int    `json:"id"`
	Token1 string `json:"token1"`
	Token2 string `json:"token2"`
}

// Updated response structure based on actual API response
type PoolListResponse struct {
	Pool []PoolDetails `json:"pool"`
}

type PoolTickInfo struct {
	PoolID         string `json:"poolId"`
	FeeRate        int    `json:"feeRate"`
	ActiveLPAmount string `json:"activeLPAmount"`
	CurrentTick    int    `json:"currentTick"`
	SqrtPrice      string `json:"sqrtPrice"`
	Ticks          []struct {
		Index          int    `json:"index"`
		LiquidityGross string `json:"liquidityGross"`
		LiquidityNet   string `json:"liquidityNet"`
	} `json:"ticks"`
}

// Alternative response structures to try different schema patterns
type GraphQLResponse struct {
	API struct {
		GetPoolTickInfo PoolTickInfo `json:"getPoolTickInfo"`
	} `json:"api"`
}

// Alternative response structure without nested API
type DirectGraphQLResponse struct {
	GetPoolTickInfo PoolTickInfo `json:"getPoolTickInfo"`
}

// Schema introspection types
type SchemaIntrospection struct {
	Data struct {
		Schema struct {
			QueryType struct {
				Fields []struct {
					Name string `json:"name"`
					Type struct {
						Name string `json:"name"`
						Kind string `json:"kind"`
					} `json:"type"`
				} `json:"fields"`
			} `json:"queryType"`
		} `json:"__schema"`
	} `json:"data"`
}

var APTOS_API_KEYS = []string{
	"aptoslabs_PhZAmy6fiHU_AwG75mKKF1sHHdqpai3ijYFFpxJygSNSD",
	"aptoslabs_iTCXFrKh23V_26mvKe8opyEksmLoDqqcnSNkK1UmYxHMZ",
	"aptoslabs_G82W81ng49e_LYnuJ37ZUNEUJDNRPsrTpk3f4Dste5qjG",
	"aptoslabs_3xcZgajCih6_9sAkpnMLUQCV5GF7TUBZopSjusZTyxGAp",
	"aptoslabs_jiwcwPkTMmD_7SNBD1qRguN9KojeEBUhDmwCYgmVBYDTj",
}

var (
	currentIndex int
	sdkMutex     sync.Mutex
)

// Constants for math operations
var (
	ONE          = big.NewInt(1)
	TWO          = big.NewInt(2)
	ZERO         = big.NewInt(0)
	NEGATIVE_ONE = big.NewInt(-1)
	MAX_FEE      = new(big.Int).Exp(big.NewInt(10), big.NewInt(6), nil)
	Q32          = new(big.Int).Exp(big.NewInt(2), big.NewInt(32), nil)
	Q96          = new(big.Int).Exp(big.NewInt(2), big.NewInt(96), nil)
	Q192         = new(big.Int).Exp(Q96, big.NewInt(2), nil)
	Q64          = new(big.Int).Exp(big.NewInt(2), big.NewInt(64), nil)
	Q128         = new(big.Int).Exp(Q64, big.NewInt(2), nil)
	MaxUint160   = new(big.Int).Sub(new(big.Int).Exp(big.NewInt(2), big.NewInt(160), nil), ONE)
	MaxUint256   = new(big.Int).Sub(new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil), ONE)
)

const (
	MIN_TICK       = -443636
	MAX_TICK       = 443636
	MIN_SQRT_RATIO = "18447090764788882728"
	MAX_SQRT_RATIO = "340275971719517849884101479065584693834"
)

var TICK_SPACINGS = map[string]int{
	"100":   1,
	"500":   10,
	"1000":  20,
	"2500":  50,
	"3000":  60,
	"10000": 200,
}

func normalizeAptosAddress(address string) string {
	// Remove '0x' prefix if present
	address = strings.TrimPrefix(address, "0x")

	// Pad with zeros to make it 64 characters
	for len(address) < 64 {
		address = "0" + address
	}

	return "0x" + address
}

// Function to introspect the GraphQL schema
func introspectSchema(client *graphql.Client) error {
	introspectionQuery := `
		query IntrospectionQuery {
			__schema {
				queryType {
					fields {
						name
						type {
							name
							kind
						}
					}
				}
			}
		}
	`

	req := graphql.NewRequest(introspectionQuery)
	var resp SchemaIntrospection

	err := client.Run(context.Background(), req, &resp)
	if err != nil {
		return fmt.Errorf("schema introspection failed: %v", err)
	}

	fmt.Println("[Hyperion] Available GraphQL schema fields:")
	for _, field := range resp.Data.Schema.QueryType.Fields {
		fmt.Printf("[Hyperion] - %s (type: %s, kind: %s)\n", field.Name, field.Type.Name, field.Type.Kind)
	}

	return nil
}

// Updated function using the correct 'pool' field
func fetchAllHyperionPools() ([]PoolDetails, error) {
	graphQLClient := graphql.NewClient("https://hyperfluid-api.alcove.pro/v1/graphql")

	// First, try to introspect the schema
	fmt.Println("[Hyperion] Attempting schema introspection...")
	err := introspectSchema(graphQLClient)
	if err != nil {
		fmt.Printf("[Hyperion] Schema introspection failed: %v\n", err)
	}

	// Use the correct query based on your curl test
	query := `query GetPools {
		pool {
			id
			token1
			token2
		}
	}`

	fmt.Println("[Hyperion] Using correct 'pool' field query")
	
	req := graphql.NewRequest(query)
	var resp PoolListResponse
	
	err = graphQLClient.Run(context.Background(), req, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pools: %v", err)
	}

	if len(resp.Pool) == 0 {
		return nil, fmt.Errorf("no pools returned from API")
	}

	fmt.Printf("[Hyperion] Successfully fetched %d pools\n", len(resp.Pool))
	return resp.Pool, nil
}

func FetchPools(chainId int, dexName string, dexConfig *common.DexConfiguration) ([]any, error) {
	// Check if the API is accessible first
	graphQLClient := graphql.NewClient("https://hyperfluid-api.alcove.pro/v1/graphql")
	
	// Try a simple health check query first
	healthQuery := `query { __typename }`
	healthReq := graphql.NewRequest(healthQuery)
	var healthResp map[string]interface{}
	
	err := graphQLClient.Run(context.Background(), healthReq, &healthResp)
	if err != nil {
		return nil, fmt.Errorf("Hyperion API is not accessible: %v", err)
	}
	
	fmt.Printf("[Hyperion] API is accessible, typename: %v\n", healthResp["__typename"])

	allPoolsData, err := fetchAllHyperionPools()
	if err != nil {
		return nil, fmt.Errorf("error fetching Hyperion pools: %v", err)
	}

	if len(allPoolsData) == 0 {
		fmt.Printf("[Hyperion] No pools found for dex: %s on chain: %d\n", dexName, chainId)
		return []any{}, nil
	}

	// Convert pools to the expected format
	var poolsData []any

	for _, pool := range allPoolsData {
		// Convert pool ID to string for consistency
		poolIdStr := strconv.Itoa(pool.ID)
		
		// Normalize the token addresses
		token0 := normalizeAptosAddress(pool.Token1)
		token1 := normalizeAptosAddress(pool.Token2)

		// For now, we'll use placeholder values for missing data
		// You may need to fetch additional token metadata from another endpoint
		poolInfo := map[string]any{
			"poolId":     poolIdStr,
			"dexId":      "HYPERION",
			"token0":     token0,
			"token1":     token1,
			"symbol0":    "UNKNOWN", // You'll need to fetch this separately
			"symbol1":    "UNKNOWN", // You'll need to fetch this separately
			"decimals0":  8,         // Default value, fetch from token metadata
			"decimals1":  8,         // Default value, fetch from token metadata
			"gasPrice":   "0",
			"tokenType0": token0,    // Using address as type for now
			"tokenType1": token1,    // Using address as type for now
			"fee":        3000,      // Default fee, you may need to fetch this
			"ticks":      []interface{}{}, // Empty ticks for now
			"sqrtPrice":  "0",       // Default sqrt price
			"tickCurrent": 0,        // Default tick
			"liquidity":  "0",       // Default liquidity
		}

		poolsData = append(poolsData, poolInfo)
	}

	fmt.Printf("[Hyperion] Successfully processed %d pools for dex: %s on chain: %d\n", len(poolsData), dexName, chainId)
	return poolsData, nil
}

// Rest of the math helper functions remain the same...
func mostSignificantBit(x *big.Int) int {
	if x.Cmp(ZERO) <= 0 {
		panic("ZERO")
	}
	if x.Cmp(MaxUint256) > 0 {
		panic("MAX")
	}

	msb := 0
	powers := []struct {
		power int
		min   *big.Int
	}{
		{128, new(big.Int).Exp(TWO, big.NewInt(128), nil)},
		{64, new(big.Int).Exp(TWO, big.NewInt(64), nil)},
		{32, new(big.Int).Exp(TWO, big.NewInt(32), nil)},
		{16, new(big.Int).Exp(TWO, big.NewInt(16), nil)},
		{8, new(big.Int).Exp(TWO, big.NewInt(8), nil)},
		{4, new(big.Int).Exp(TWO, big.NewInt(4), nil)},
		{2, new(big.Int).Exp(TWO, big.NewInt(2), nil)},
		{1, new(big.Int).Exp(TWO, big.NewInt(1), nil)},
	}

	for _, p := range powers {
		if x.Cmp(p.min) >= 0 {
			x.Rsh(x, uint(p.power))
			msb += p.power
		}
	}
	return msb
}

func getSqrtPrice(sqrtPrice string) *big.Int {
	price, _ := new(big.Int).SetString(sqrtPrice, 10)
	return new(big.Int).Lsh(price, 32)
}

func mulDivRoundingUp(a, b, denominator *big.Int) *big.Int {
	product := new(big.Int).Mul(a, b)
	result := new(big.Int).Div(product, denominator)
	remainder := new(big.Int).Mod(product, denominator)
	if remainder.Cmp(ZERO) != 0 {
		result.Add(result, ONE)
	}
	return result
}

func getAmount0Delta(sqrtRatioAX96, sqrtRatioBX96, liquidity *big.Int, roundUp bool) *big.Int {
	if sqrtRatioAX96.Cmp(sqrtRatioBX96) > 0 {
		sqrtRatioAX96, sqrtRatioBX96 = sqrtRatioBX96, sqrtRatioAX96
	}

	numerator1 := new(big.Int).Lsh(liquidity, 96)
	numerator2 := new(big.Int).Sub(sqrtRatioBX96, sqrtRatioAX96)

	if roundUp {
		return mulDivRoundingUp(
			mulDivRoundingUp(numerator1, numerator2, sqrtRatioBX96),
			ONE,
			sqrtRatioAX96,
		)
	}

	temp := new(big.Int).Mul(numerator1, numerator2)
	temp.Div(temp, sqrtRatioBX96)
	return temp.Div(temp, sqrtRatioAX96)
}

func getAmount1Delta(sqrtRatioAX96, sqrtRatioBX96, liquidity *big.Int, roundUp bool) *big.Int {
	if sqrtRatioAX96.Cmp(sqrtRatioBX96) > 0 {
		sqrtRatioAX96, sqrtRatioBX96 = sqrtRatioBX96, sqrtRatioAX96
	}

	diff := new(big.Int).Sub(sqrtRatioBX96, sqrtRatioAX96)

	if roundUp {
		return mulDivRoundingUp(liquidity, diff, Q96)
	}

	return new(big.Int).Div(new(big.Int).Mul(liquidity, diff), Q96)
}

func GetHyperionAmountOut(pool map[string]any, amountIn string, fromToken, toToken string) (map[string]any, error) {
	token0, ok := pool["token0"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid token0")
	}

	token1, ok := pool["token1"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid token1")
	}

	tokens := []string{strings.ToLower(token0), strings.ToLower(token1)}
	fromTokenLower := strings.ToLower(fromToken)
	toTokenLower := strings.ToLower(toToken)

	found := false
	for _, t := range tokens {
		if t == fromTokenLower {
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("fromToken not part of pool")
	}

	found = false
	for _, t := range tokens {
		if t == toTokenLower {
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("toToken not part of pool")
	}

	sqrtPriceStr, ok := pool["sqrtPrice"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid sqrtPrice")
	}

	sqrtRatioX96 := getSqrtPrice(sqrtPriceStr)

	tickCurrent, ok := pool["tickCurrent"].(int)
	if !ok {
		return nil, fmt.Errorf("invalid tickCurrent")
	}

	liquidityStr, ok := pool["liquidity"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid liquidity")
	}

	liquidity, _ := new(big.Int).SetString(liquidityStr, 10)

	amountSpecified, _ := new(big.Int).SetString(amountIn, 10)

	zeroForOne := strings.ToLower(fromToken) == strings.ToLower(token0)

	minSqrtRatio, _ := new(big.Int).SetString(MIN_SQRT_RATIO, 10)
	maxSqrtRatio, _ := new(big.Int).SetString(MAX_SQRT_RATIO, 10)

	var sqrtPriceLimitX96 *big.Int
	if zeroForOne {
		sqrtPriceLimitX96 = new(big.Int).Add(minSqrtRatio, ONE)
	} else {
		sqrtPriceLimitX96 = new(big.Int).Sub(maxSqrtRatio, ONE)
	}

	// Validation
	if zeroForOne {
		if sqrtPriceLimitX96.Cmp(minSqrtRatio) <= 0 {
			return nil, fmt.Errorf("RATIO_MIN")
		}
		if sqrtPriceLimitX96.Cmp(sqrtRatioX96) >= 0 {
			return nil, fmt.Errorf("RATIO_CURRENT")
		}
	} else {
		if sqrtPriceLimitX96.Cmp(maxSqrtRatio) >= 0 {
			return nil, fmt.Errorf("RATIO_MAX")
		}
		if sqrtPriceLimitX96.Cmp(sqrtRatioX96) <= 0 {
			return nil, fmt.Errorf("RATIO_CURRENT")
		}
	}

	state := struct {
		amountSpecifiedRemaining *big.Int
		amountCalculated         *big.Int
		sqrtPriceX96            *big.Int
		tick                    int
		liquidity               *big.Int
	}{
		amountSpecifiedRemaining: new(big.Int).Set(amountSpecified),
		amountCalculated:         new(big.Int).Set(ZERO),
		sqrtPriceX96:            new(big.Int).Set(sqrtRatioX96),
		tick:                    tickCurrent,
		liquidity:               new(big.Int).Set(liquidity),
	}

	// Simplified swap loop - in real implementation, you'd need tick data and more complex logic
	for state.amountSpecifiedRemaining.Cmp(ZERO) != 0 && state.sqrtPriceX96.Cmp(sqrtPriceLimitX96) != 0 {
		// This is a simplified version - full implementation would require tick traversal
		break
	}

	if state.amountCalculated.Cmp(ZERO) < 0 {
		state.amountCalculated.Mul(state.amountCalculated, NEGATIVE_ONE)
	}

	return map[string]any{
		"amountOut": state.amountCalculated.String(),
	}, nil
}

package traderjoe

import (
	"context"
	"fetcher-go/src/common"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"sync"

	"github.com/machinebox/graphql"
	"github.com/shopspring/decimal"
)

var (
	ZERO            = big.NewInt(0)
	FEE_NUMERATOR   = big.NewInt(997)
	FEE_DENOMINATOR = big.NewInt(1000)
)

func FetchPools(chainId int, dexName string, dexConfig *common.DexConfiguration) ([]any, error) {
	url := dexConfig.GraphqlUrl
	if url == "" {
		return nil, fmt.Errorf("invalid or missing graphqlUrl in dexConfig")
	}

	client := graphql.NewClient(url)
	query := `
		query GetPairs($first: Int!, $skip: Int!) {
			lbpairs(first: $first, skip: $skip) {
				id
				bins {
					reserveX
					reserveY
				}
				tokenX {
					id
					symbol
					name
					decimals
				}
				tokenY {
					id
					symbol
					name
					decimals
				}
			}
		}
	`

	type result struct {
		skip  int
		pairs []any
		err   error
	}

	pageCount := 2
	results := make(chan result, pageCount)
	var wg sync.WaitGroup

	for i := range pageCount {
		skip := i * 1000
		wg.Add(1)

		go func(skip int) {
			defer wg.Done()

			req := graphql.NewRequest(query)
			req.Var("first", 1000)
			req.Var("skip", skip)

			var resp struct {
				LBPairs []struct {
					ID     string `json:"id"`
					Bins   []struct {
						ReserveX string `json:"reserveX"`
						ReserveY string `json:"reserveY"`
					} `json:"bins"`
					TokenX struct {
						ID       string `json:"id"`
						Symbol   string `json:"symbol"`
						Name     string `json:"name"`
						Decimals string `json:"decimals"`
					} `json:"tokenX"`
					TokenY struct {
						ID       string `json:"id"`
						Symbol   string `json:"symbol"`
						Name     string `json:"name"`
						Decimals string `json:"decimals"`
					} `json:"tokenY"`
				} `json:"lbpairs"`
			}

			err := client.Run(context.Background(), req, &resp)
			if err != nil {
				results <- result{skip: skip, pairs: nil, err: err}
				return
			}

			processedPairs := []any{}
			for _, pair := range resp.LBPairs {
				// Calculate reserves
				reserveX := decimal.NewFromInt(0)
				reserveY := decimal.NewFromInt(0)

				for _, bin := range pair.Bins {
					binReserveX, err := decimal.NewFromString(bin.ReserveX)
					if err != nil {
						continue
					}
					binReserveY, err := decimal.NewFromString(bin.ReserveY)
					if err != nil {
						continue
					}
					reserveX = reserveX.Add(binReserveX)
					reserveY = reserveY.Add(binReserveY)
				}

				decimals0, err := strconv.Atoi(pair.TokenX.Decimals)
				if err != nil {
					continue
				}
				decimals1, err := strconv.Atoi(pair.TokenY.Decimals)
				if err != nil {
					continue
				}

				// Calculate actual reserves
				actualReserveX := reserveX.Mul(decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(decimals0))))
				actualReserveY := reserveY.Mul(decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(decimals1))))

				// Skip if reserves are zero or undefined
				if actualReserveX.IsZero() || actualReserveY.IsZero() {
					continue
				}

				pairData := map[string]any{
					"id": pair.ID,
					"token0": map[string]any{
						"id":       pair.TokenX.ID,
						"symbol":   pair.TokenX.Symbol,
						"decimals": pair.TokenX.Decimals,
					},
					"token1": map[string]any{
						"id":       pair.TokenY.ID,
						"symbol":   pair.TokenY.Symbol,
						"decimals": pair.TokenY.Decimals,
					},
					"reserve0": actualReserveX.String(),
					"reserve1": actualReserveY.String(),
				}

				processedPairs = append(processedPairs, pairData)
			}

			results <- result{skip: skip, pairs: processedPairs, err: nil}
		}(skip)
	}

	wg.Wait()
	close(results)

	allPairs := []any{}
	for res := range results {
		if res.err != nil {
			return nil, fmt.Errorf("error fetching pairs at skip %d: %v", res.skip, res.err)
		}
		if len(res.pairs) > 0 {
			allPairs = append(allPairs, res.pairs...)
		}
	}

	return allPairs, nil
}

func GetTraderjoeAmountOut(pool map[string]interface{}, rawAmountIn string, fromToken string, toToken string) (map[string]interface{}, error) {
	amountIn := new(big.Int)
	amountIn, ok := amountIn.SetString(rawAmountIn, 10)
	if !ok {
		return nil, fmt.Errorf("invalid amountIn")
	}

	if amountIn.Cmp(ZERO) <= 0 {
		return nil, fmt.Errorf("ZERO_AMOUNT_IN")
	}

	token0, ok := pool["token0"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid token0")
	}
	// token1, ok := pool["token1"].(string)
	// if !ok {
	// 	return nil, fmt.Errorf("invalid token1")
	// }

	zeroForOne := strings.ToLower(fromToken) == strings.ToLower(token0)

	var reserveIn, reserveOut *big.Int
	reserve0Str, ok := pool["reserve0"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid reserve0")
	}
	reserve1Str, ok := pool["reserve1"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid reserve1")
	}

	reserve0 := new(big.Int)
	reserve0, ok = reserve0.SetString(reserve0Str, 10)
	if !ok {
		return nil, fmt.Errorf("invalid reserve0 format")
	}

	reserve1 := new(big.Int)
	reserve1, ok = reserve1.SetString(reserve1Str, 10)
	if !ok {
		return nil, fmt.Errorf("invalid reserve1 format")
	}

	if zeroForOne {
		reserveIn = reserve0
		reserveOut = reserve1
	} else {
		reserveIn = reserve1
		reserveOut = reserve0
	}

	if reserveIn.Cmp(ZERO) <= 0 || reserveOut.Cmp(ZERO) <= 0 {
		return nil, fmt.Errorf("ZERO_POOL_RESERVE")
	}

	// Calculate amountOut using the formula: (amountIn * 997 * reserveOut) / (reserveIn * 1000 + amountIn * 997)
	amountInWithFee := new(big.Int).Mul(amountIn, FEE_NUMERATOR)
	numerator := new(big.Int).Mul(amountInWithFee, reserveOut)
	denominator := new(big.Int).Add(new(big.Int).Mul(reserveIn, FEE_DENOMINATOR), amountInWithFee)
	amountOut := new(big.Int).Div(numerator, denominator)

	// Calculate price impact
	amountInDecimal := decimal.NewFromBigInt(amountIn, 0)
	amountOutDecimal := decimal.NewFromBigInt(amountOut, 0)
	reserve0Decimal := decimal.NewFromBigInt(reserve0, 0)
	reserve1Decimal := decimal.NewFromBigInt(reserve1, 0)

	var reserveInBefore, reserveOutBefore, reserveInAfter, reserveOutAfter decimal.Decimal

	if zeroForOne {
		// Selling base (token0)
		reserveInBefore = reserve0Decimal
		reserveOutBefore = reserve1Decimal
		reserveInAfter = reserve0Decimal.Add(amountInDecimal)
		reserveOutAfter = reserve1Decimal.Sub(amountOutDecimal)
	} else {
		// Selling quote (token1)
		reserveInBefore = reserve1Decimal
		reserveOutBefore = reserve0Decimal
		reserveInAfter = reserve1Decimal.Sub(amountOutDecimal)
		reserveOutAfter = reserve0Decimal.Add(amountInDecimal)
	}

	priceBefore := reserveOutBefore.Div(reserveInBefore)
	priceAfter := reserveOutAfter.Div(reserveInAfter)
	priceImpact := priceBefore.Sub(priceAfter).Abs().Mul(decimal.NewFromInt(100)).Div(priceBefore)

	if priceImpact.GreaterThan(decimal.NewFromInt(5)) {
		poolId, _ := pool["poolId"].(string)
		dexId, _ := pool["dexId"].(string)
		return nil, fmt.Errorf("Skipped %s pool %s due to high price impact: %s%%", dexId, poolId, priceImpact.StringFixed(4))
	}

	return map[string]interface{}{
		"amountOut": amountOut.String(),
	}, nil
}

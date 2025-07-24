package ramseslegacy

import (
	"context"
	"fetcher-go/src/common"
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/machinebox/graphql"
	"github.com/shopspring/decimal"
)

var ONE_E18 = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)

func ferFee(isStable bool) *big.Int {
	if isStable {
		return big.NewInt(5)
	}
	return big.NewInt(30)
}

func _f(x0, y *big.Int) *big.Int {
	y2 := new(big.Int).Div(new(big.Int).Mul(y, y), ONE_E18)
	y3 := new(big.Int).Div(new(big.Int).Mul(y2, y), ONE_E18)
	term1 := new(big.Int).Div(new(big.Int).Mul(x0, y3), ONE_E18)

	x02 := new(big.Int).Div(new(big.Int).Mul(x0, x0), ONE_E18)
	x03 := new(big.Int).Div(new(big.Int).Mul(x02, x0), ONE_E18)
	term2 := new(big.Int).Div(new(big.Int).Mul(x03, y), ONE_E18)

	return new(big.Int).Add(term1, term2)
}

func _d(x0, y *big.Int) *big.Int {
	y2 := new(big.Int).Div(new(big.Int).Mul(y, y), ONE_E18)
	term1 := new(big.Int).Div(new(big.Int).Mul(new(big.Int).Mul(big.NewInt(3), x0), y2), ONE_E18)

	x02 := new(big.Int).Div(new(big.Int).Mul(x0, x0), ONE_E18)
	x03 := new(big.Int).Div(new(big.Int).Mul(x02, x0), ONE_E18)

	return new(big.Int).Add(term1, x03)
}

func _get_y(x0, xy, y *big.Int) *big.Int {
	for i := 0; i < 255; i++ {
		y_prev := new(big.Int).Set(y)
		k := _f(x0, y)
		
		// Calculate absolute difference
		var absDiff *big.Int
		if k.Cmp(xy) >= 0 {
			absDiff = new(big.Int).Sub(k, xy)
		} else {
			absDiff = new(big.Int).Sub(xy, k)
		}
		
		dy := new(big.Int).Div(
			new(big.Int).Mul(absDiff, ONE_E18),
			_d(x0, y),
		)
		
		if k.Cmp(xy) < 0 {
			y = new(big.Int).Add(y, dy)
		} else {
			y = new(big.Int).Sub(y, dy)
		}

		var delta *big.Int
		if y.Cmp(y_prev) < 0 {
			delta = new(big.Int).Sub(y_prev, y)
		} else {
			delta = new(big.Int).Sub(y, y_prev)
		}

		if delta.Cmp(big.NewInt(2)) < 0 {
			return y
		}
	}
	return y
}

func _k(x, y *big.Int, stable bool, decimals0, decimals1 *big.Int) *big.Int {
	if stable {
		_x := new(big.Int).Div(new(big.Int).Mul(x, ONE_E18), decimals0)
		_y := new(big.Int).Div(new(big.Int).Mul(y, ONE_E18), decimals1)

		_a := new(big.Int).Div(new(big.Int).Mul(_x, _y), ONE_E18)
		_b := new(big.Int).Add(
			new(big.Int).Div(new(big.Int).Mul(_x, _x), ONE_E18),
			new(big.Int).Div(new(big.Int).Mul(_y, _y), ONE_E18),
		)

		return new(big.Int).Div(new(big.Int).Mul(_a, _b), ONE_E18)
	} else {
		return new(big.Int).Mul(x, y)
	}
}

func _getAmountOut(amountIn *big.Int, tokenIn string, _reserve0, _reserve1 *big.Int, stable bool, token0 string, decimals0, decimals1 *big.Int) *big.Int {
	if stable {
		xy := _k(_reserve0, _reserve1, stable, decimals0, decimals1)
		_reserve0 = new(big.Int).Div(new(big.Int).Mul(_reserve0, ONE_E18), decimals0)
		_reserve1 = new(big.Int).Div(new(big.Int).Mul(_reserve1, ONE_E18), decimals1)

		var reserveA, reserveB *big.Int
		var amountInNormalized *big.Int
		
		if strings.EqualFold(tokenIn, token0) {
			reserveA = _reserve0
			reserveB = _reserve1
			amountInNormalized = new(big.Int).Div(new(big.Int).Mul(amountIn, ONE_E18), decimals0)
		} else {
			reserveA = _reserve1
			reserveB = _reserve0
			amountInNormalized = new(big.Int).Div(new(big.Int).Mul(amountIn, ONE_E18), decimals1)
		}

		y := new(big.Int).Sub(reserveB, _get_y(new(big.Int).Add(amountInNormalized, reserveA), xy, reserveB))
		
		var targetDecimals *big.Int
		if strings.EqualFold(tokenIn, token0) {
			targetDecimals = decimals1
		} else {
			targetDecimals = decimals0
		}
		
		return new(big.Int).Div(
			new(big.Int).Mul(y, targetDecimals),
			ONE_E18,
		)
	} else {
		var reserveA, reserveB *big.Int
		if strings.EqualFold(tokenIn, token0) {
			reserveA = _reserve0
			reserveB = _reserve1
		} else {
			reserveA = _reserve1
			reserveB = _reserve0
		}

		return new(big.Int).Div(
			new(big.Int).Mul(amountIn, reserveB),
			new(big.Int).Add(reserveA, amountIn),
		)
	}
}

type Pool struct {
	IsStable  bool   `json:"isStable"`
	Reserve0  string `json:"reserve0"`
	Reserve1  string `json:"reserve1"`
	Decimals0 string `json:"decimals0"`
	Decimals1 string `json:"decimals1"`
	Token0    string `json:"token0"`
	DexId     string `json:"dexId"`
	PoolId    string `json:"poolId"`
}

type AmountOutResult struct {
	AmountOut string `json:"amountOut"`
}

func GetRamsesLegacyAmountOut(pool *Pool, amountIn string, tokenIn string) (*AmountOutResult, error) {
	stable := pool.IsStable
	
	reserve0, ok := new(big.Int).SetString(pool.Reserve0, 10)
	if !ok {
		return nil, fmt.Errorf("invalid reserve0: %s", pool.Reserve0)
	}
	
	reserve1, ok := new(big.Int).SetString(pool.Reserve1, 10)
	if !ok {
		return nil, fmt.Errorf("invalid reserve1: %s", pool.Reserve1)
	}
	
	decimals0, ok := new(big.Int).SetString(pool.Decimals0, 10)
	if !ok {
		return nil, fmt.Errorf("invalid decimals0: %s", pool.Decimals0)
	}
	
	decimals1, ok := new(big.Int).SetString(pool.Decimals1, 10)
	if !ok {
		return nil, fmt.Errorf("invalid decimals1: %s", pool.Decimals1)
	}
	
	token0 := pool.Token0

	amountInBI, ok := new(big.Int).SetString(amountIn, 10)
	if !ok {
		return nil, fmt.Errorf("invalid amountIn: %s", amountIn)
	}
	
	fee := ferFee(stable)
	feeAmount := new(big.Int).Div(new(big.Int).Mul(amountInBI, fee), big.NewInt(10000))
	amountAfterFee := new(big.Int).Sub(amountInBI, feeAmount)

	amountOut := _getAmountOut(amountAfterFee, tokenIn, reserve0, reserve1, stable, token0, decimals0, decimals1)

	zeroForOne := strings.EqualFold(tokenIn, pool.Token0)

	amountInDecimal, err := decimal.NewFromString(amountIn)
	if err != nil {
		return nil, fmt.Errorf("invalid amountIn decimal: %v", err)
	}
	
	amountOutDecimal, err := decimal.NewFromString(amountOut.String())
	if err != nil {
		return nil, fmt.Errorf("invalid amountOut decimal: %v", err)
	}
	
	reserve0Decimal, err := decimal.NewFromString(pool.Reserve0)
	if err != nil {
		return nil, fmt.Errorf("invalid reserve0 decimal: %v", err)
	}
	
	reserve1Decimal, err := decimal.NewFromString(pool.Reserve1)
	if err != nil {
		return nil, fmt.Errorf("invalid reserve1 decimal: %v", err)
	}

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
		return nil, fmt.Errorf("skipped %s pool %s due to high price impact: %s%%", pool.DexId, pool.PoolId, priceImpact.StringFixed(4))
	}

	return &AmountOutResult{
		AmountOut: amountOut.String(),
	}, nil
}

func FetchPools(chainId int, dexName string, dexConfig *common.DexConfiguration) ([]any, error) {
	url := dexConfig.GraphqlUrl
	if url == "" {
		return nil, fmt.Errorf("invalid or missing graphqlUrl in dexConfig")
	}

	client := graphql.NewClient(url)
	query := `
		query getPools($skip: Int) {
			pools(
				first: 1000,
				skip: $skip,
				orderBy: totalValueLockedUSD,
				orderDirection: desc,
				where: { liquidity_gt: 0, totalValueLockedUSD_gt: 10000 }
			) {
				id
				stable
				reserve0
				reserve1
				token0 { 
					id 
					symbol 
					decimals 
				}
				token1 { 
					id 
					symbol 
					decimals 
				}
			}
		}
	`

	type result struct {
		skip  int
		pools []any
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
			req.Var("skip", skip)

			var resp struct {
				Pools []any `json:"pools"`
			}

			err := client.Run(context.Background(), req, &resp)
			results <- result{skip: skip, pools: resp.Pools, err: err}
		}(skip)
	}

	wg.Wait()
	close(results)

	allPools := []any{}
	for res := range results {
		if res.err != nil {
			return nil, fmt.Errorf("error fetching pools at skip %d: %v", res.skip, res.err)
		}
		if len(res.pools) > 0 {
			allPools = append(allPools, res.pools...)
		}
	}

	return allPools, nil
}
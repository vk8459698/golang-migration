package uniswap

import (
	"context"
	"fetcher-go/src/common"
	"fmt"
	"github.com/machinebox/graphql"
	"sync"
)

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
				where: { liquidity_gt: 0, sqrtPrice_gt: 0, totalValueLockedUSD_gt: 10000 }
			) {
				id
				tick
				liquidity
				sqrtPrice
				feeTier
				token0 { id symbol decimals }
				token1 { id symbol decimals }
				ticks {
					tickIdx
					liquidityNet
					liquidityGross
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

	for i := 0; i < pageCount; i++ {
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

package uniswap

import (
	"context"
	"fetcher-go/src/common"
	"fmt"
	"github.com/machinebox/graphql"
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

	skip := 0
	allPools := []interface{}{}
	for range 2 {
		req := graphql.NewRequest(query)
		req.Var("skip", skip)

		var resp struct {
			Pools []any `json:"pools"`
		}

		err := client.Run(context.Background(), req, &resp)
		if err != nil {
			return nil, fmt.Errorf("error fetching pools: %v", err)
		}

		if len(resp.Pools) == 0 {
			break
		}

		allPools = append(allPools, resp.Pools...)

		if len(resp.Pools) < 1000 {
			break
		}

		skip += 1000
	}

	return allPools, nil
}

package swaprv2

import (
	"context"
	"fetcher-go/src/common"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"sync"

	"github.com/machinebox/graphql"
)

func FetchPools(chainId int, dexName string, dexConfig *common.DexConfiguration) ([]any, error) {
	url := dexConfig.GraphqlUrl
	if url == "" {
		return nil, fmt.Errorf("invalid or missing graphqlUrl in dexConfig")
	}

	client := graphql.NewClient(url)
	query := `
		query getPairs($skip: Int) {
			pairs(
				first: 1000,
				skip: $skip,
				where: { reserve0_gt: "0", reserve1_gt: "0" }
			) {
				id
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
				reserve0
				reserve1
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
			req.Var("skip", skip)

			var resp struct {
				Pairs []any `json:"pairs"`
			}

			err := client.Run(context.Background(), req, &resp)
			results <- result{skip: skip, pairs: resp.Pairs, err: err}
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

	// Process pairs to convert reserves using the original logic
	processedPairs := []any{}
	for _, p := range allPairs {
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("[Ledger] Recovered from panic while processing Swapr V2 pool data for dex %s on chain %d: %v", dexName, chainId, r)
				}
			}()

			pairMap, ok := p.(map[string]any)
			if !ok {
				return
			}

			token0Map, ok := pairMap["token0"].(map[string]any)
			if !ok {
				return
			}

			token1Map, ok := pairMap["token1"].(map[string]any)
			if !ok {
				return
			}

			decimals0Str, ok := token0Map["decimals"].(string)
			if !ok {
				return
			}

			decimals1Str, ok := token1Map["decimals"].(string)
			if !ok {
				return
			}

			reserve0Str, ok := pairMap["reserve0"].(string)
			if !ok {
				return
			}

			reserve1Str, ok := pairMap["reserve1"].(string)
			if !ok {
				return
			}

			decimals0, err := strconv.Atoi(decimals0Str)
			if err != nil {
				log.Printf("[Ledger] Error parsing decimals0 for Swapr V2 pool %s: %v", pairMap["id"], err)
				return
			}

			decimals1, err := strconv.Atoi(decimals1Str)
			if err != nil {
				log.Printf("[Ledger] Error parsing decimals1 for Swapr V2 pool %s: %v", pairMap["id"], err)
				return
			}

			reserve0Raw, _, err := new(big.Float).Parse(reserve0Str, 10)
			if err != nil {
				log.Printf("[Ledger] Error parsing reserve0 for Swapr V2 pool %s: %v", pairMap["id"], err)
				return
			}

			reserve1Raw, _, err := new(big.Float).Parse(reserve1Str, 10)
			if err != nil {
				log.Printf("[Ledger] Error parsing reserve1 for Swapr V2 pool %s: %v", pairMap["id"], err)
				return
			}

			exp0 := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals0)), nil))
			exp1 := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals1)), nil))

			reserve0 := new(big.Float).Mul(reserve0Raw, exp0).Text('f', 0)
			reserve1 := new(big.Float).Mul(reserve1Raw, exp1).Text('f', 0)

			// Update the reserves in the pair map
			pairMap["reserve0"] = reserve0
			pairMap["reserve1"] = reserve1

			processedPairs = append(processedPairs, pairMap)
		}()
	}

	return processedPairs, nil
}
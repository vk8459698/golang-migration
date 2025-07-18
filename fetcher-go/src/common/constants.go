package common

type TokenInfo struct {
	ID       string
	Symbol   string
	Decimals int
}

type PoolInfo struct {
	Token0 TokenInfo
	Token1 TokenInfo
	Fee    int
}

type ProtocolParameters struct {
	MaxTxSize  int
	MaxValSize int
}

type DexConfiguration struct {
	DexId               string
	ChainId             int
	ChainName           string
	NativeSymbol        string
	NativePrecision     int
	RpcUrl              string
	GraphqlUrl          string
	BlockfrostProjectId string
	ProtocolParameters  *ProtocolParameters
	BaseTokens          map[string]TokenInfo
	QuoteTokens         map[string]TokenInfo
	Pools               map[string]PoolInfo
}

var UniswapConfigurations = map[int]*DexConfiguration{
	1: {
		DexId:           "UNISWAP_V3",
		ChainId:         1,
		ChainName:       "Ethereum",
		NativeSymbol:    "ETH",
		NativePrecision: 18,
		RpcUrl: func() string {
			url, _ := getRpcUrl("uniswap", "1")
			return url
		}(),
		GraphqlUrl: func() string {
			url, _ := getGraphqlUrl("uniswap", "1")
			return url
		}(),
	},
	10: {
		DexId:           "UNISWAP_V3",
		ChainId:         10,
		ChainName:       "Optimism",
		NativeSymbol:    "ETH",
		NativePrecision: 18,
		RpcUrl: func() string {
			url, _ := getRpcUrl("uniswap", "10")
			return url
		}(),
		GraphqlUrl: func() string {
			url, _ := getGraphqlUrl("uniswap", "10")
			return url
		}(),
	},
	56: {
		DexId:           "UNISWAP_V3",
		ChainId:         56,
		ChainName:       "BNB",
		NativeSymbol:    "BNB",
		NativePrecision: 18,
		RpcUrl: func() string {
			url, _ := getRpcUrl("uniswap", "56")
			return url
		}(),
		GraphqlUrl: func() string {
			url, _ := getGraphqlUrl("uniswap", "56")
			return url
		}(),
	},
	137: {
		DexId:           "UNISWAP_V3",
		ChainId:         137,
		ChainName:       "Polygon",
		NativeSymbol:    "POL",
		NativePrecision: 18,
		RpcUrl: func() string {
			url, _ := getRpcUrl("uniswap", "137")
			return url
		}(),
		GraphqlUrl: func() string {
			url, _ := getGraphqlUrl("uniswap", "137")
			return url
		}(),
	},
	42161: {
		DexId:           "UNISWAP_V3",
		ChainId:         42161,
		ChainName:       "Arbitrum",
		NativeSymbol:    "ETH",
		NativePrecision: 18,
		RpcUrl: func() string {
			url, _ := getRpcUrl("uniswap", "42161")
			return url
		}(),
		GraphqlUrl: func() string {
			url, _ := getGraphqlUrl("uniswap", "42161")
			return url
		}(),
	},
}

var PancakeswapConfigurations = map[int]*DexConfiguration{
	1: {
		DexId:           "PANCAKESWAP_V3",
		ChainId:         1,
		ChainName:       "Ethereum",
		NativeSymbol:    "ETH",
		NativePrecision: 18,
		RpcUrl: func() string {
			url, _ := getRpcUrl("pancakeswap", "1")
			return url
		}(),
		GraphqlUrl: func() string {
			url, _ := getGraphqlUrl("pancakeswap", "1")
			return url
		}(),
	},
	56: {
		DexId:           "PANCAKESWAP_V3",
		ChainId:         56,
		ChainName:       "BNB",
		NativeSymbol:    "BNB",
		NativePrecision: 18,
		RpcUrl: func() string {
			url, _ := getRpcUrl("pancakeswap", "56")
			return url
		}(),
		GraphqlUrl: func() string {
			url, _ := getGraphqlUrl("pancakeswap", "56")
			return url
		}(),
	},
	8453: {
		DexId:           "PANCAKESWAP_V3",
		ChainId:         8453,
		ChainName:       "Base",
		NativeSymbol:    "ETH",
		NativePrecision: 18,
		RpcUrl: func() string {
			url, _ := getRpcUrl("pancakeswap", "8453")
			return url
		}(),
		GraphqlUrl: func() string {
			url, _ := getGraphqlUrl("pancakeswap", "8453")
			return url
		}(),
	},
	42161: {
		DexId:           "PANCAKESWAP_V3",
		ChainId:         42161,
		ChainName:       "Arbitrum",
		NativeSymbol:    "ETH",
		NativePrecision: 18,
		RpcUrl: func() string {
			url, _ := getRpcUrl("pancakeswap", "42161")
			return url
		}(),
		GraphqlUrl: func() string {
			url, _ := getGraphqlUrl("pancakeswap", "42161")
			return url
		}(),
	},
}

var SushiswapConfigurations = map[int]*DexConfiguration{
	1: {
		DexId:           "SUSHISWAP_V3",
		ChainId:         1,
		ChainName:       "Ethereum",
		NativeSymbol:    "ETH",
		NativePrecision: 18,
		RpcUrl: func() string {
			url, _ := getRpcUrl("sushiswap", "1")
			return url
		}(),
		GraphqlUrl: func() string {
			url, _ := getGraphqlUrl("sushiswap", "1")
			return url
		}(),
	},
	10: {
		DexId:           "SUSHISWAP_V3",
		ChainId:         10,
		ChainName:       "Optimism",
		NativeSymbol:    "ETH",
		NativePrecision: 18,
		RpcUrl: func() string {
			url, _ := getRpcUrl("sushiswap", "10")
			return url
		}(),
		GraphqlUrl: func() string {
			url, _ := getGraphqlUrl("sushiswap", "10")
			return url
		}(),
	},
	56: {
		DexId:           "SUSHISWAP_V3",
		ChainId:         56,
		ChainName:       "BNB",
		NativeSymbol:    "BNB",
		NativePrecision: 18,
		RpcUrl: func() string {
			url, _ := getRpcUrl("sushiswap", "56")
			return url
		}(),
		GraphqlUrl: func() string {
			url, _ := getGraphqlUrl("sushiswap", "56")
			return url
		}(),
	},
	137: {
		DexId:           "SUSHISWAP_V3",
		ChainId:         137,
		ChainName:       "Polygon",
		NativeSymbol:    "POL",
		NativePrecision: 18,
		RpcUrl: func() string {
			url, _ := getRpcUrl("sushiswap", "137")
			return url
		}(),
		GraphqlUrl: func() string {
			url, _ := getGraphqlUrl("sushiswap", "137")
			return url
		}(),
	},
	8453: {
		DexId:           "SUSHISWAP_V3",
		ChainId:         8453,
		ChainName:       "Base",
		NativeSymbol:    "ETH",
		NativePrecision: 18,
		RpcUrl: func() string {
			url, _ := getRpcUrl("sushiswap", "8453")
			return url
		}(),
		GraphqlUrl: func() string {
			url, _ := getGraphqlUrl("sushiswap", "8453")
			return url
		}(),
	},
	42161: {
		DexId:           "SUSHISWAP_V3",
		ChainId:         42161,
		ChainName:       "Arbitrum",
		NativeSymbol:    "ETH",
		NativePrecision: 18,
		RpcUrl: func() string {
			url, _ := getRpcUrl("sushiswap", "42161")
			return url
		}(),
		GraphqlUrl: func() string {
			url, _ := getGraphqlUrl("sushiswap", "42161")
			return url
		}(),
	},
}

var CamelotConfigurations = map[int]*DexConfiguration{
	42161: {
		DexId:           "CAMELOT_V3",
		ChainId:         42161,
		ChainName:       "Arbitrum",
		NativeSymbol:    "ETH",
		NativePrecision: 18,
		RpcUrl: func() string {
			url, _ := getRpcUrl("camelot", "42161")
			return url
		}(),
		GraphqlUrl: func() string {
			url, _ := getGraphqlUrl("camelot", "42161")
			return url
		}(),
	},
}

var RaydiumConfigurations = map[int]*DexConfiguration{
	101: {
		DexId:           "RAYDIUM_V4",
		ChainId:         101,
		ChainName:       "Solana",
		NativeSymbol:    "SOL",
		NativePrecision: 9,
		RpcUrl:          "https://api.mainnet-beta.solana.com",
	},
}

var MeteoraConfigurations = map[int]*DexConfiguration{
	101: {
		DexId:           "METEORA",
		ChainId:         101,
		ChainName:       "Solana",
		NativeSymbol:    "SOL",
		NativePrecision: 9,
		RpcUrl:          "https://api.mainnet-beta.solana.com",
	},
}

var OrcaConfigurations = map[int]*DexConfiguration{
	101: {
		DexId:           "ORCA_WHIRL",
		ChainId:         101,
		ChainName:       "Solana",
		NativeSymbol:    "SOL",
		NativePrecision: 9,
		RpcUrl:          "https://api.mainnet-beta.solana.com",
	},
}

var LifintyConfigurations = map[int]*DexConfiguration{
	101: {
		DexId:           "LIFINITY",
		ChainId:         101,
		ChainName:       "Solana",
		NativeSymbol:    "SOL",
		NativePrecision: 9,
		RpcUrl:          "https://api.mainnet-beta.solana.com",
	},
}

var MinswapConfigurations = map[int]*DexConfiguration{
	921: {
		DexId:               "MINSWAP",
		ChainId:             921,
		ChainName:           "Cardano",
		NativeSymbol:        "ADA",
		NativePrecision:     6,
		RpcUrl:              "https://cardano-mainnet.blockfrost.io/api/v0/",
		BlockfrostProjectId: "mainnetAuG3MT1aYz6iVKfutPMss65Gu0eJV8hA",
		ProtocolParameters: &ProtocolParameters{
			MaxTxSize:  16384,
			MaxValSize: 5000,
		},
	},
}

var DodoConfigurations = map[int]*DexConfiguration{
	1: {
		DexId:           "DODO_V2",
		ChainId:         1,
		ChainName:       "Ethereum",
		NativeSymbol:    "ETH",
		NativePrecision: 18,
		RpcUrl: func() string {
			url, _ := getRpcUrl("dodo", "1")
			return url
		}(),
		GraphqlUrl: func() string {
			url, _ := getGraphqlUrl("dodo", "1")
			return url
		}(),
	},
	10: {
		DexId:           "DODO_V2",
		ChainId:         10,
		ChainName:       "Optimism",
		NativeSymbol:    "ETH",
		NativePrecision: 18,
		RpcUrl: func() string {
			url, _ := getRpcUrl("dodo", "10")
			return url
		}(),
		GraphqlUrl: func() string {
			url, _ := getGraphqlUrl("dodo", "10")
			return url
		}(),
	},
	56: {
		DexId:           "DODO_V2",
		ChainId:         56,
		ChainName:       "BNB",
		NativeSymbol:    "BNB",
		NativePrecision: 18,
		RpcUrl: func() string {
			url, _ := getRpcUrl("dodo", "56")
			return url
		}(),
		GraphqlUrl: func() string {
			url, _ := getGraphqlUrl("dodo", "56")
			return url
		}(),
	},
	137: {
		DexId:           "DODO_V2",
		ChainId:         137,
		ChainName:       "Polygon",
		NativeSymbol:    "POL",
		NativePrecision: 18,
		RpcUrl: func() string {
			url, _ := getRpcUrl("dodo", "137")
			return url
		}(),
		GraphqlUrl: func() string {
			url, _ := getGraphqlUrl("dodo", "137")
			return url
		}(),
	},
	8453: {
		DexId:           "DODO_V2",
		ChainId:         8453,
		ChainName:       "Base",
		NativeSymbol:    "ETH",
		NativePrecision: 18,
		RpcUrl: func() string {
			url, _ := getRpcUrl("dodo", "8453")
			return url
		}(),
		GraphqlUrl: func() string {
			url, _ := getGraphqlUrl("dodo", "8453")
			return url
		}(),
	},
	42161: {
		DexId:           "DODO_V2",
		ChainId:         42161,
		ChainName:       "Arbitrum",
		NativeSymbol:    "ETH",
		NativePrecision: 18,
		RpcUrl: func() string {
			url, _ := getRpcUrl("dodo", "42161")
			return url
		}(),
		GraphqlUrl: func() string {
			url, _ := getGraphqlUrl("dodo", "42161")
			return url
		}(),
	},
}

var CurveConfigurations = map[int]*DexConfiguration{
	1: {
		DexId:           "CURVE",
		ChainId:         1,
		ChainName:       "Ethereum",
		NativeSymbol:    "ETH",
		NativePrecision: 18,
		RpcUrl: func() string {
			url, _ := getRpcUrl("curve", "1")
			return url
		}(),
		GraphqlUrl: func() string {
			url, _ := getGraphqlUrl("curve", "1")
			return url
		}(),
	},
}

var MaverickConfigurations = map[int]*DexConfiguration{
	1: {
		DexId:           "MAVERICK_V2",
		ChainId:         1,
		ChainName:       "Ethereum",
		NativeSymbol:    "ETH",
		NativePrecision: 18,
		RpcUrl: func() string {
			url, _ := getRpcUrl("maverick", "1")
			return url
		}(),
	},
	56: {
		DexId:           "MAVERICK_V2",
		ChainId:         56,
		ChainName:       "BNB",
		NativeSymbol:    "BNB",
		NativePrecision: 18,
		RpcUrl: func() string {
			url, _ := getRpcUrl("maverick", "56")
			return url
		}(),
	},
	8453: {
		DexId:           "MAVERICK_V2",
		ChainId:         8453,
		ChainName:       "Base",
		NativeSymbol:    "ETH",
		NativePrecision: 18,
		RpcUrl: func() string {
			url, _ := getRpcUrl("maverick", "8453")
			return url
		}(),
	},
	42161: {
		DexId:           "MAVERICK_V2",
		ChainId:         42161,
		ChainName:       "Arbitrum",
		NativeSymbol:    "ETH",
		NativePrecision: 18,
		RpcUrl: func() string {
			url, _ := getRpcUrl("maverick", "42161")
			return url
		}(),
	},
}

var HyperionConfigurations = map[int]*DexConfiguration{
	1400: {
		DexId:           "HYPERION",
		ChainId:         1400,
		ChainName:       "Aptos",
		NativeSymbol:    "APT",
		NativePrecision: 8,
	},
}

var ThalaConfigurations = map[int]*DexConfiguration{
	1400: {
		DexId:           "THALA_V2",
		ChainId:         1400,
		ChainName:       "Aptos",
		NativeSymbol:    "APT",
		NativePrecision: 8,
	},
}

var CellanaConfigurations = map[int]*DexConfiguration{
	1400: {
		DexId:           "CELLANA",
		ChainId:         1400,
		ChainName:       "Aptos",
		NativeSymbol:    "APT",
		NativePrecision: 8,
	},
}

var TraderJoeConfigurations = map[int]*DexConfiguration{
	56: {
		DexId:           "TRADERJOE",
		ChainId:         56,
		ChainName:       "BNB",
		NativeSymbol:    "BNB",
		NativePrecision: 18,
		GraphqlUrl: func() string {
			url, _ := getGraphqlUrl("traderJoe", "56")
			return url
		}(),
	},
	42161: {
		DexId:           "TRADERJOE",
		ChainId:         42161,
		ChainName:       "Arbitrum",
		NativeSymbol:    "ETH",
		NativePrecision: 18,
		GraphqlUrl: func() string {
			url, _ := getGraphqlUrl("traderJoe", "42161")
			return url
		}(),
	},
}

var RamsesClConfiguration = map[int]*DexConfiguration{
	42161: {
		DexId:           "RAMSES_CL",
		ChainId:         42161,
		ChainName:       "Arbitrum",
		NativeSymbol:    "ETH",
		NativePrecision: 18,
		GraphqlUrl: func() string {
			url, _ := getGraphqlUrl("ramsesCl", "42161")
			return url
		}(),
	},
}

var RamsesLegacyConfiguration = map[int]*DexConfiguration{
	42161: {
		DexId:           "RAMSES_LEGACY",
		ChainId:         42161,
		ChainName:       "Arbitrum",
		NativeSymbol:    "ETH",
		NativePrecision: 18,
		GraphqlUrl: func() string {
			url, _ := getGraphqlUrl("ramsesLegacy", "42161")
			return url
		}(),
	},
}

var SwaprConfiguration = map[int]*DexConfiguration{
	1: {
		DexId:           "SWAPR_V2",
		ChainId:         1,
		ChainName:       "Ethereum",
		NativeSymbol:    "ETH",
		NativePrecision: 18,
		GraphqlUrl: func() string {
			url, _ := getGraphqlUrl("swapr", "1")
			return url
		}(),
	},
	42161: {
		DexId:           "SWAPR_V2",
		ChainId:         42161,
		ChainName:       "Arbitrum",
		NativeSymbol:    "ETH",
		NativePrecision: 18,
		GraphqlUrl: func() string {
			url, _ := getGraphqlUrl("swapr", "42161")
			return url
		}(),
	},
}

var WoofiConfigurations = map[int]*DexConfiguration{
	10: {
		DexId:           "WOOFI_V2",
		ChainId:         10,
		ChainName:       "Optimism",
		NativeSymbol:    "ETH",
		NativePrecision: 18,
		BaseTokens:      map[string]TokenInfo{},
		QuoteTokens:     map[string]TokenInfo{},
	},
	56: {
		DexId:           "WOOFI_V2",
		ChainId:         56,
		ChainName:       "BNB",
		NativeSymbol:    "BNB",
		NativePrecision: 18,
		BaseTokens:      map[string]TokenInfo{},
		QuoteTokens:     map[string]TokenInfo{},
	},
	8453: {
		DexId:           "WOOFI_V2",
		ChainId:         8453,
		ChainName:       "Base",
		NativeSymbol:    "ETH",
		NativePrecision: 18,
		BaseTokens:      map[string]TokenInfo{},
		QuoteTokens:     map[string]TokenInfo{},
	},
	42161: {
		DexId:           "WOOFI_V2",
		ChainId:         42161,
		ChainName:       "Arbitrum",
		NativeSymbol:    "ETH",
		NativePrecision: 18,
		BaseTokens: map[string]TokenInfo{
			"0x82aF49447D8a07e3bd95BD0d56f35241523fBab1": {
				ID:       "0x82aF49447D8a07e3bd95BD0d56f35241523fBab1",
				Symbol:   "WETH",
				Decimals: 18,
			},
			"0x2f2a2543B76A4166549F7aaB2e75Bef0aefC5B0f": {
				ID:       "0x2f2a2543B76A4166549F7aaB2e75Bef0aefC5B0f",
				Symbol:   "WBTC",
				Decimals: 8,
			},
			"0xFF970A61A04b1cA14834A43f5dE4533eBDDB5CC8": {
				ID:       "0xFF970A61A04b1cA14834A43f5dE4533eBDDB5CC8",
				Symbol:   "USDC.e",
				Decimals: 6,
			},
			"0x912CE59144191C1204E64559FE8253a0e49E6548": {
				ID:       "0x912CE59144191C1204E64559FE8253a0e49E6548",
				Symbol:   "ARB",
				Decimals: 18,
			},
			"0xFd086bC7CD5C481DCC9C85ebE478A1C0b69FCbb9": {
				ID:       "0xFd086bC7CD5C481DCC9C85ebE478A1C0b69FCbb9",
				Symbol:   "USDT",
				Decimals: 6,
			},
		},
		QuoteTokens: map[string]TokenInfo{
			"0xaf88d065e77c8cC2239327C5EDb3A432268e5831": {
				ID:       "0xaf88d065e77c8cC2239327C5EDb3A432268e5831",
				Symbol:   "USDC",
				Decimals: 6,
			},
		},
	},
}

var FluidConfigurations = map[int]*DexConfiguration{
	1: {
		DexId:           "FLUID_V2",
		ChainId:         1,
		ChainName:       "Ethereum",
		NativeSymbol:    "ETH",
		NativePrecision: 18,
		Pools:           map[string]PoolInfo{},
	},
	137: {
		DexId:           "FLUID_V2",
		ChainId:         137,
		ChainName:       "Polygon",
		NativeSymbol:    "POL",
		NativePrecision: 18,
		Pools:           map[string]PoolInfo{},
	},
	8453: {
		DexId:           "FLUID_V2",
		ChainId:         8453,
		ChainName:       "Base",
		NativeSymbol:    "ETH",
		NativePrecision: 18,
		Pools:           map[string]PoolInfo{},
	},
	42161: {
		DexId:           "FLUID_V2",
		ChainId:         42161,
		ChainName:       "Arbitrum",
		NativeSymbol:    "ETH",
		NativePrecision: 18,
		Pools: map[string]PoolInfo{
			"0x0B1a513ee24972DAEf112bC777a5610d4325C9e7": {
				Token0: TokenInfo{
					ID:       "0x5979D7b546E38E414F7E9822514be443A4800529",
					Symbol:   "wstETH",
					Decimals: 18,
				},
				Token1: TokenInfo{
					ID:       "0x82aF49447D8a07e3bd95BD0d56f35241523fBab1",
					Symbol:   "WETH",
					Decimals: 18,
				},
				Fee: 100,
			},
			"0x667701e51B4D1Ca244F17C78F7aB8744B4C99F9B": {
				Token0: TokenInfo{
					ID:       "0x35751007a407ca6FEFfE80b3cB397736D2cf4dbe",
					Symbol:   "weETH",
					Decimals: 18,
				},
				Token1: TokenInfo{
					ID:       "0x82aF49447D8a07e3bd95BD0d56f35241523fBab1",
					Symbol:   "WETH",
					Decimals: 18,
				},
				Fee: 100,
			},
			"0x3C0441B42195F4aD6aa9a0978E06096ea616CDa7": {
				Token0: TokenInfo{
					ID:       "0xaf88d065e77c8cC2239327C5EDb3A432268e5831",
					Symbol:   "USDC",
					Decimals: 6,
				},
				Token1: TokenInfo{
					ID:       "0xFd086bC7CD5C481DCC9C85ebE478A1C0b69FCbb9",
					Symbol:   "USDT",
					Decimals: 6,
				},
				Fee: 10,
			},
			"0xdE632C3a214D5f14C1d8ddF0b92F8BCd188fee45": {
				Token0: TokenInfo{
					ID:       "0xaf88d065e77c8cC2239327C5EDb3A432268e5831",
					Symbol:   "USDC",
					Decimals: 6,
				},
				Token1: TokenInfo{
					ID:       "0x82aF49447D8a07e3bd95BD0d56f35241523fBab1",
					Symbol:   "ETH",
					Decimals: 18,
				},
				Fee: 200,
			},
			"0x2886a01a0645390872a9eb99dAe1283664b0c524": {
				Token0: TokenInfo{
					ID:       "0x61E030A56D33e8260FdD81f03B162A79Fe3449Cd",
					Symbol:   "FLUID",
					Decimals: 18,
				},
				Token1: TokenInfo{
					ID:       "0x82aF49447D8a07e3bd95BD0d56f35241523fBab1",
					Symbol:   "ETH",
					Decimals: 18,
				},
				Fee: 5000,
			},
		},
	},
}

var BalancerConfigurations = map[int]*DexConfiguration{
	1: {
		DexId:           "BALANCER",
		ChainId:         1,
		ChainName:       "Ethereum",
		NativeSymbol:    "ETH",
		NativePrecision: 18,
		GraphqlUrl: func() string {
			url, _ := getGraphqlUrl("balancer", "1")
			return url
		}(),
	},
	10: {
		DexId:           "BALANCER",
		ChainId:         10,
		ChainName:       "Optimism",
		NativeSymbol:    "ETH",
		NativePrecision: 18,
		GraphqlUrl: func() string {
			url, _ := getGraphqlUrl("balancer", "10")
			return url
		}(),
	},
	137: {
		DexId:           "BALANCER",
		ChainId:         137,
		ChainName:       "Polygon",
		NativeSymbol:    "POL",
		NativePrecision: 18,
		GraphqlUrl: func() string {
			url, _ := getGraphqlUrl("balancer", "137")
			return url
		}(),
	},
	8453: {
		DexId:           "BALANCER",
		ChainId:         8453,
		ChainName:       "Base",
		NativeSymbol:    "ETH",
		NativePrecision: 18,
		GraphqlUrl: func() string {
			url, _ := getGraphqlUrl("balancer", "8453")
			return url
		}(),
	},
	42161: {
		DexId:           "BALANCER",
		ChainId:         42161,
		ChainName:       "Arbitrum",
		NativeSymbol:    "ETH",
		NativePrecision: 18,
		GraphqlUrl: func() string {
			url, _ := getGraphqlUrl("balancer", "42161")
			return url
		}(),
	},
}

var DexConfigurations = map[string]map[int]*DexConfiguration{
	"uniswap":      UniswapConfigurations,
	"pancakeswap":  PancakeswapConfigurations,
	"sushiswap":    SushiswapConfigurations,
	"camelot":      CamelotConfigurations,
	"dodo":         DodoConfigurations,
	"curve":        CurveConfigurations,
	"ramsesCl":     RamsesClConfiguration,
	"ramsesLegacy": RamsesLegacyConfiguration,
	"swapr":        SwaprConfiguration,
	"woofi":        WoofiConfigurations,
	"fluid":        FluidConfigurations,
	"balancer":     BalancerConfigurations,
	"maverick":     MaverickConfigurations,
	"traderjoe":    TraderJoeConfigurations,
	"raydium":      RaydiumConfigurations,
	"orca":         OrcaConfigurations,
	"meteora":      MeteoraConfigurations,
	"lifinity":     LifintyConfigurations,
	"minswap":      MinswapConfigurations,
	"hyperion":     HyperionConfigurations,
	"thala":        ThalaConfigurations,
	"cellana":      CellanaConfigurations,
}

var ChainConfigurations = make(map[int]map[string]*DexConfiguration)

func init() {
	for dexName, dexConfigs := range DexConfigurations {
		for chainID, config := range dexConfigs {
			if ChainConfigurations[chainID] == nil {
				ChainConfigurations[chainID] = make(map[string]*DexConfiguration)
			}
			ChainConfigurations[chainID][dexName] = config
		}
	}
}

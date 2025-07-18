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
	dexId               string
	chainId             int
	chainName           string
	nativeSymbol        string
	nativePrecision     int
	rpcUrl              string
	graphqlUrl          string
	blockfrostProjectId string
	protocolParameters  *ProtocolParameters
	baseTokens          map[string]TokenInfo
	quoteTokens         map[string]TokenInfo
	pools               map[string]PoolInfo
}

var uniswapConfigurations = map[int]*DexConfiguration{
	1: {
		dexId:           "UNISWAP_V3",
		chainId:         1,
		chainName:       "Ethereum",
		nativeSymbol:    "ETH",
		nativePrecision: 18,
		rpcUrl: func() string {
			url, _ := getRpcUrl("uniswap", "1")
			return url
		}(),
		graphqlUrl: func() string {
			url, _ := getGraphqlUrl("uniswap", "1")
			return url
		}(),
	},
	10: {
		dexId:           "UNISWAP_V3",
		chainId:         10,
		chainName:       "Optimism",
		nativeSymbol:    "ETH",
		nativePrecision: 18,
		rpcUrl: func() string {
			url, _ := getRpcUrl("uniswap", "10")
			return url
		}(),
		graphqlUrl: func() string {
			url, _ := getGraphqlUrl("uniswap", "10")
			return url
		}(),
	},
	56: {
		dexId:           "UNISWAP_V3",
		chainId:         56,
		chainName:       "BNB",
		nativeSymbol:    "BNB",
		nativePrecision: 18,
		rpcUrl: func() string {
			url, _ := getRpcUrl("uniswap", "56")
			return url
		}(),
		graphqlUrl: func() string {
			url, _ := getGraphqlUrl("uniswap", "56")
			return url
		}(),
	},
	137: {
		dexId:           "UNISWAP_V3",
		chainId:         137,
		chainName:       "Polygon",
		nativeSymbol:    "POL",
		nativePrecision: 18,
		rpcUrl: func() string {
			url, _ := getRpcUrl("uniswap", "137")
			return url
		}(),
		graphqlUrl: func() string {
			url, _ := getGraphqlUrl("uniswap", "137")
			return url
		}(),
	},
	42161: {
		dexId:           "UNISWAP_V3",
		chainId:         42161,
		chainName:       "Arbitrum",
		nativeSymbol:    "ETH",
		nativePrecision: 18,
		rpcUrl: func() string {
			url, _ := getRpcUrl("uniswap", "42161")
			return url
		}(),
		graphqlUrl: func() string {
			url, _ := getGraphqlUrl("uniswap", "42161")
			return url
		}(),
	},
}

var pancakeswapConfigurations = map[int]*DexConfiguration{
	1: {
		dexId:           "PANCAKESWAP_V3",
		chainId:         1,
		chainName:       "Ethereum",
		nativeSymbol:    "ETH",
		nativePrecision: 18,
		rpcUrl: func() string {
			url, _ := getRpcUrl("pancakeswap", "1")
			return url
		}(),
		graphqlUrl: func() string {
			url, _ := getGraphqlUrl("pancakeswap", "1")
			return url
		}(),
	},
	56: {
		dexId:           "PANCAKESWAP_V3",
		chainId:         56,
		chainName:       "BNB",
		nativeSymbol:    "BNB",
		nativePrecision: 18,
		rpcUrl: func() string {
			url, _ := getRpcUrl("pancakeswap", "56")
			return url
		}(),
		graphqlUrl: func() string {
			url, _ := getGraphqlUrl("pancakeswap", "56")
			return url
		}(),
	},
	8453: {
		dexId:           "PANCAKESWAP_V3",
		chainId:         8453,
		chainName:       "Base",
		nativeSymbol:    "ETH",
		nativePrecision: 18,
		rpcUrl: func() string {
			url, _ := getRpcUrl("pancakeswap", "8453")
			return url
		}(),
		graphqlUrl: func() string {
			url, _ := getGraphqlUrl("pancakeswap", "8453")
			return url
		}(),
	},
	42161: {
		dexId:           "PANCAKESWAP_V3",
		chainId:         42161,
		chainName:       "Arbitrum",
		nativeSymbol:    "ETH",
		nativePrecision: 18,
		rpcUrl: func() string {
			url, _ := getRpcUrl("pancakeswap", "42161")
			return url
		}(),
		graphqlUrl: func() string {
			url, _ := getGraphqlUrl("pancakeswap", "42161")
			return url
		}(),
	},
}

var sushiswapConfigurations = map[int]*DexConfiguration{
	1: {
		dexId:           "SUSHISWAP_V3",
		chainId:         1,
		chainName:       "Ethereum",
		nativeSymbol:    "ETH",
		nativePrecision: 18,
		rpcUrl: func() string {
			url, _ := getRpcUrl("sushiswap", "1")
			return url
		}(),
		graphqlUrl: func() string {
			url, _ := getGraphqlUrl("sushiswap", "1")
			return url
		}(),
	},
	10: {
		dexId:           "SUSHISWAP_V3",
		chainId:         10,
		chainName:       "Optimism",
		nativeSymbol:    "ETH",
		nativePrecision: 18,
		rpcUrl: func() string {
			url, _ := getRpcUrl("sushiswap", "10")
			return url
		}(),
		graphqlUrl: func() string {
			url, _ := getGraphqlUrl("sushiswap", "10")
			return url
		}(),
	},
	56: {
		dexId:           "SUSHISWAP_V3",
		chainId:         56,
		chainName:       "BNB",
		nativeSymbol:    "BNB",
		nativePrecision: 18,
		rpcUrl: func() string {
			url, _ := getRpcUrl("sushiswap", "56")
			return url
		}(),
		graphqlUrl: func() string {
			url, _ := getGraphqlUrl("sushiswap", "56")
			return url
		}(),
	},
	137: {
		dexId:           "SUSHISWAP_V3",
		chainId:         137,
		chainName:       "Polygon",
		nativeSymbol:    "POL",
		nativePrecision: 18,
		rpcUrl: func() string {
			url, _ := getRpcUrl("sushiswap", "137")
			return url
		}(),
		graphqlUrl: func() string {
			url, _ := getGraphqlUrl("sushiswap", "137")
			return url
		}(),
	},
	8453: {
		dexId:           "SUSHISWAP_V3",
		chainId:         8453,
		chainName:       "Base",
		nativeSymbol:    "ETH",
		nativePrecision: 18,
		rpcUrl: func() string {
			url, _ := getRpcUrl("sushiswap", "8453")
			return url
		}(),
		graphqlUrl: func() string {
			url, _ := getGraphqlUrl("sushiswap", "8453")
			return url
		}(),
	},
	42161: {
		dexId:           "SUSHISWAP_V3",
		chainId:         42161,
		chainName:       "Arbitrum",
		nativeSymbol:    "ETH",
		nativePrecision: 18,
		rpcUrl: func() string {
			url, _ := getRpcUrl("sushiswap", "42161")
			return url
		}(),
		graphqlUrl: func() string {
			url, _ := getGraphqlUrl("sushiswap", "42161")
			return url
		}(),
	},
}

var camelotConfigurations = map[int]*DexConfiguration{
	42161: {
		dexId:           "CAMELOT_V3",
		chainId:         42161,
		chainName:       "Arbitrum",
		nativeSymbol:    "ETH",
		nativePrecision: 18,
		rpcUrl: func() string {
			url, _ := getRpcUrl("camelot", "42161")
			return url
		}(),
		graphqlUrl: func() string {
			url, _ := getGraphqlUrl("camelot", "42161")
			return url
		}(),
	},
}

var raydiumConfigurations = map[int]*DexConfiguration{
	101: {
		dexId:           "RAYDIUM_V4",
		chainId:         101,
		chainName:       "Solana",
		nativeSymbol:    "SOL",
		nativePrecision: 9,
		rpcUrl:          "https://api.mainnet-beta.solana.com",
	},
}

var meteoraConfigurations = map[int]*DexConfiguration{
	101: {
		dexId:           "METEORA",
		chainId:         101,
		chainName:       "Solana",
		nativeSymbol:    "SOL",
		nativePrecision: 9,
		rpcUrl:          "https://api.mainnet-beta.solana.com",
	},
}

var orcaConfigurations = map[int]*DexConfiguration{
	101: {
		dexId:           "ORCA_WHIRL",
		chainId:         101,
		chainName:       "Solana",
		nativeSymbol:    "SOL",
		nativePrecision: 9,
		rpcUrl:          "https://api.mainnet-beta.solana.com",
	},
}

var lifinityConfigurations = map[int]*DexConfiguration{
	101: {
		dexId:           "LIFINITY",
		chainId:         101,
		chainName:       "Solana",
		nativeSymbol:    "SOL",
		nativePrecision: 9,
		rpcUrl:          "https://api.mainnet-beta.solana.com",
	},
}

var minswapConfigurations = map[int]*DexConfiguration{
	921: {
		dexId:               "MINSWAP",
		chainId:             921,
		chainName:           "Cardano",
		nativeSymbol:        "ADA",
		nativePrecision:     6,
		rpcUrl:              "https://cardano-mainnet.blockfrost.io/api/v0/",
		blockfrostProjectId: "mainnetAuG3MT1aYz6iVKfutPMss65Gu0eJV8hA",
		protocolParameters: &ProtocolParameters{
			MaxTxSize:  16384,
			MaxValSize: 5000,
		},
	},
}

var dodoConfigurations = map[int]*DexConfiguration{
	1: {
		dexId:           "DODO_V2",
		chainId:         1,
		chainName:       "Ethereum",
		nativeSymbol:    "ETH",
		nativePrecision: 18,
		rpcUrl: func() string {
			url, _ := getRpcUrl("dodo", "1")
			return url
		}(),
		graphqlUrl: func() string {
			url, _ := getGraphqlUrl("dodo", "1")
			return url
		}(),
	},
	10: {
		dexId:           "DODO_V2",
		chainId:         10,
		chainName:       "Optimism",
		nativeSymbol:    "ETH",
		nativePrecision: 18,
		rpcUrl: func() string {
			url, _ := getRpcUrl("dodo", "10")
			return url
		}(),
		graphqlUrl: func() string {
			url, _ := getGraphqlUrl("dodo", "10")
			return url
		}(),
	},
	56: {
		dexId:           "DODO_V2",
		chainId:         56,
		chainName:       "BNB",
		nativeSymbol:    "BNB",
		nativePrecision: 18,
		rpcUrl: func() string {
			url, _ := getRpcUrl("dodo", "56")
			return url
		}(),
		graphqlUrl: func() string {
			url, _ := getGraphqlUrl("dodo", "56")
			return url
		}(),
	},
	137: {
		dexId:           "DODO_V2",
		chainId:         137,
		chainName:       "Polygon",
		nativeSymbol:    "POL",
		nativePrecision: 18,
		rpcUrl: func() string {
			url, _ := getRpcUrl("dodo", "137")
			return url
		}(),
		graphqlUrl: func() string {
			url, _ := getGraphqlUrl("dodo", "137")
			return url
		}(),
	},
	8453: {
		dexId:           "DODO_V2",
		chainId:         8453,
		chainName:       "Base",
		nativeSymbol:    "ETH",
		nativePrecision: 18,
		rpcUrl: func() string {
			url, _ := getRpcUrl("dodo", "8453")
			return url
		}(),
		graphqlUrl: func() string {
			url, _ := getGraphqlUrl("dodo", "8453")
			return url
		}(),
	},
	42161: {
		dexId:           "DODO_V2",
		chainId:         42161,
		chainName:       "Arbitrum",
		nativeSymbol:    "ETH",
		nativePrecision: 18,
		rpcUrl: func() string {
			url, _ := getRpcUrl("dodo", "42161")
			return url
		}(),
		graphqlUrl: func() string {
			url, _ := getGraphqlUrl("dodo", "42161")
			return url
		}(),
	},
}

var curveConfigurations = map[int]*DexConfiguration{
	1: {
		dexId:           "CURVE",
		chainId:         1,
		chainName:       "Ethereum",
		nativeSymbol:    "ETH",
		nativePrecision: 18,
		rpcUrl: func() string {
			url, _ := getRpcUrl("curve", "1")
			return url
		}(),
		graphqlUrl: func() string {
			url, _ := getGraphqlUrl("curve", "1")
			return url
		}(),
	},
}

var maverickConfigurations = map[int]*DexConfiguration{
	1: {
		dexId:           "MAVERICK_V2",
		chainId:         1,
		chainName:       "Ethereum",
		nativeSymbol:    "ETH",
		nativePrecision: 18,
		rpcUrl: func() string {
			url, _ := getRpcUrl("maverick", "1")
			return url
		}(),
	},
	56: {
		dexId:           "MAVERICK_V2",
		chainId:         56,
		chainName:       "BNB",
		nativeSymbol:    "BNB",
		nativePrecision: 18,
		rpcUrl: func() string {
			url, _ := getRpcUrl("maverick", "56")
			return url
		}(),
	},
	8453: {
		dexId:           "MAVERICK_V2",
		chainId:         8453,
		chainName:       "Base",
		nativeSymbol:    "ETH",
		nativePrecision: 18,
		rpcUrl: func() string {
			url, _ := getRpcUrl("maverick", "8453")
			return url
		}(),
	},
	42161: {
		dexId:           "MAVERICK_V2",
		chainId:         42161,
		chainName:       "Arbitrum",
		nativeSymbol:    "ETH",
		nativePrecision: 18,
		rpcUrl: func() string {
			url, _ := getRpcUrl("maverick", "42161")
			return url
		}(),
	},
}

var hyperionConfigurations = map[int]*DexConfiguration{
	1400: {
		dexId:           "HYPERION",
		chainId:         1400,
		chainName:       "Aptos",
		nativeSymbol:    "APT",
		nativePrecision: 8,
	},
}

var thalaConfigurations = map[int]*DexConfiguration{
	1400: {
		dexId:           "THALA_V2",
		chainId:         1400,
		chainName:       "Aptos",
		nativeSymbol:    "APT",
		nativePrecision: 8,
	},
}

var cellanaConfigurations = map[int]*DexConfiguration{
	1400: {
		dexId:           "CELLANA",
		chainId:         1400,
		chainName:       "Aptos",
		nativeSymbol:    "APT",
		nativePrecision: 8,
	},
}

var traderJoeConfigurations = map[int]*DexConfiguration{
	56: {
		dexId:           "TRADERJOE",
		chainId:         56,
		chainName:       "BNB",
		nativeSymbol:    "BNB",
		nativePrecision: 18,
		graphqlUrl: func() string {
			url, _ := getGraphqlUrl("traderJoe", "56")
			return url
		}(),
	},
	42161: {
		dexId:           "TRADERJOE",
		chainId:         42161,
		chainName:       "Arbitrum",
		nativeSymbol:    "ETH",
		nativePrecision: 18,
		graphqlUrl: func() string {
			url, _ := getGraphqlUrl("traderJoe", "42161")
			return url
		}(),
	},
}

var ramsesClConfiguration = map[int]*DexConfiguration{
	42161: {
		dexId:           "RAMSES_CL",
		chainId:         42161,
		chainName:       "Arbitrum",
		nativeSymbol:    "ETH",
		nativePrecision: 18,
		graphqlUrl: func() string {
			url, _ := getGraphqlUrl("ramsesCl", "42161")
			return url
		}(),
	},
}

var ramsesLegacyConfiguration = map[int]*DexConfiguration{
	42161: {
		dexId:           "RAMSES_LEGACY",
		chainId:         42161,
		chainName:       "Arbitrum",
		nativeSymbol:    "ETH",
		nativePrecision: 18,
		graphqlUrl: func() string {
			url, _ := getGraphqlUrl("ramsesLegacy", "42161")
			return url
		}(),
	},
}

var swaprConfiguration = map[int]*DexConfiguration{
	1: {
		dexId:           "SWAPR_V2",
		chainId:         1,
		chainName:       "Ethereum",
		nativeSymbol:    "ETH",
		nativePrecision: 18,
		graphqlUrl: func() string {
			url, _ := getGraphqlUrl("swapr", "1")
			return url
		}(),
	},
	42161: {
		dexId:           "SWAPR_V2",
		chainId:         42161,
		chainName:       "Arbitrum",
		nativeSymbol:    "ETH",
		nativePrecision: 18,
		graphqlUrl: func() string {
			url, _ := getGraphqlUrl("swapr", "42161")
			return url
		}(),
	},
}

var woofiConfigurations = map[int]*DexConfiguration{
	10: {
		dexId:           "WOOFI_V2",
		chainId:         10,
		chainName:       "Optimism",
		nativeSymbol:    "ETH",
		nativePrecision: 18,
		baseTokens:      map[string]TokenInfo{},
		quoteTokens:     map[string]TokenInfo{},
	},
	56: {
		dexId:           "WOOFI_V2",
		chainId:         56,
		chainName:       "BNB",
		nativeSymbol:    "BNB",
		nativePrecision: 18,
		baseTokens:      map[string]TokenInfo{},
		quoteTokens:     map[string]TokenInfo{},
	},
	8453: {
		dexId:           "WOOFI_V2",
		chainId:         8453,
		chainName:       "Base",
		nativeSymbol:    "ETH",
		nativePrecision: 18,
		baseTokens:      map[string]TokenInfo{},
		quoteTokens:     map[string]TokenInfo{},
	},
	42161: {
		dexId:           "WOOFI_V2",
		chainId:         42161,
		chainName:       "Arbitrum",
		nativeSymbol:    "ETH",
		nativePrecision: 18,
		baseTokens: map[string]TokenInfo{
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
		quoteTokens: map[string]TokenInfo{
			"0xaf88d065e77c8cC2239327C5EDb3A432268e5831": {
				ID:       "0xaf88d065e77c8cC2239327C5EDb3A432268e5831",
				Symbol:   "USDC",
				Decimals: 6,
			},
		},
	},
}

var fluidConfigurations = map[int]*DexConfiguration{
	1: {
		dexId:           "FLUID_V2",
		chainId:         1,
		chainName:       "Ethereum",
		nativeSymbol:    "ETH",
		nativePrecision: 18,
		pools:           map[string]PoolInfo{},
	},
	137: {
		dexId:           "FLUID_V2",
		chainId:         137,
		chainName:       "Polygon",
		nativeSymbol:    "POL",
		nativePrecision: 18,
		pools:           map[string]PoolInfo{},
	},
	8453: {
		dexId:           "FLUID_V2",
		chainId:         8453,
		chainName:       "Base",
		nativeSymbol:    "ETH",
		nativePrecision: 18,
		pools:           map[string]PoolInfo{},
	},
	42161: {
		dexId:           "FLUID_V2",
		chainId:         42161,
		chainName:       "Arbitrum",
		nativeSymbol:    "ETH",
		nativePrecision: 18,
		pools: map[string]PoolInfo{
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

var balancerConfigurations = map[int]*DexConfiguration{
	1: {
		dexId:           "BALANCER",
		chainId:         1,
		chainName:       "Ethereum",
		nativeSymbol:    "ETH",
		nativePrecision: 18,
		graphqlUrl: func() string {
			url, _ := getGraphqlUrl("balancer", "1")
			return url
		}(),
	},
	10: {
		dexId:           "BALANCER",
		chainId:         10,
		chainName:       "Optimism",
		nativeSymbol:    "ETH",
		nativePrecision: 18,
		graphqlUrl: func() string {
			url, _ := getGraphqlUrl("balancer", "10")
			return url
		}(),
	},
	137: {
		dexId:           "BALANCER",
		chainId:         137,
		chainName:       "Polygon",
		nativeSymbol:    "POL",
		nativePrecision: 18,
		graphqlUrl: func() string {
			url, _ := getGraphqlUrl("balancer", "137")
			return url
		}(),
	},
	8453: {
		dexId:           "BALANCER",
		chainId:         8453,
		chainName:       "Base",
		nativeSymbol:    "ETH",
		nativePrecision: 18,
		graphqlUrl: func() string {
			url, _ := getGraphqlUrl("balancer", "8453")
			return url
		}(),
	},
	42161: {
		dexId:           "BALANCER",
		chainId:         42161,
		chainName:       "Arbitrum",
		nativeSymbol:    "ETH",
		nativePrecision: 18,
		graphqlUrl: func() string {
			url, _ := getGraphqlUrl("balancer", "42161")
			return url
		}(),
	},
}

var DexConfigurations = map[string]map[int]*DexConfiguration{
	"uniswap":      uniswapConfigurations,
	"pancakeswap":  pancakeswapConfigurations,
	"sushiswap":    sushiswapConfigurations,
	"camelot":      camelotConfigurations,
	"dodo":         dodoConfigurations,
	"curve":        curveConfigurations,
	"ramsesCl":     ramsesClConfiguration,
	"ramsesLegacy": ramsesLegacyConfiguration,
	"swapr":        swaprConfiguration,
	"woofi":        woofiConfigurations,
	"fluid":        fluidConfigurations,

	// "balancer":  balancerConfigurations,
	// "maverick":  maverickConfigurations,
	// "traderjoe": traderJoeConfigurations,
	// "orca":      orcaConfigurations,

	"raydium":  raydiumConfigurations,
	"meteora":  meteoraConfigurations,
	"lifinity": lifinityConfigurations,

	// "minswap": minswapConfigurations,

	"thala":    thalaConfigurations,
	"cellana":  cellanaConfigurations,
	"hyperion": hyperionConfigurations,
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

package common

import "fmt"

type InfuraAccount struct {
	ApiKey string
}

type GraphqlAccount struct {
	ApiKey  string
	GraphId string
}

type Accounts struct {
	BaseUrls        map[string]string
	InfuraAccounts  map[string]InfuraAccount
	GraphqlAccounts map[string]GraphqlAccount
}

var Data = &Accounts{
	BaseUrls: map[string]string{
		"1":     "https://mainnet.infura.io/v3/",
		"10":    "https://optimism-mainnet.infura.io/v3/",
		"56":    "https://bsc-mainnet.infura.io/v3/",
		"137":   "https://polygon-mainnet.infura.io/v3/",
		"8453":  "https://base-mainnet.infura.io/v3/",
		"42161": "https://arbitrum-mainnet.infura.io/v3/",
	},

	InfuraAccounts: map[string]InfuraAccount{
		"uniswap_1":         {"034b2e82e87c4e3185c36a860393f7a3"},
		"uniswap_10":        {"50990608bb724acd9a76710ab6703a94"},
		"uniswap_56":        {"9a5af1db8626469f9edfa8a811f74451"},
		"uniswap_137":       {"f59d17c5aa644f3f856fc36fe98a549d"},
		"uniswap_42161":     {"45cbaa516e354e1b953126b0aedda069"},
		"uniswap_8453":      {"0dbbd46b18d74685a621d2878063ed45"},
		"pancakeswap_1":     {"3e34039c9795478f9dfed0e6193d840d"},
		"pancakeswap_56":    {"d5035e370598423db97f83bd11fc3ed8"},
		"pancakeswap_42161": {"69688f21099f4129a2e9e535608e4c42"},
		"pancakeswap_8453":  {"6f2c48ce2d6f4636b05f01581b16579e"},
		"sushiswap_1":       {"89c9979d7b45483b9771ac7b6dd33595"},
		"sushiswap_10":      {"d5bd99c7df0648f299ed0f127d428b0d"},
		"sushiswap_56":      {"397a0d6c9d284be38abadb305057e618"},
		"sushiswap_137":     {"9e647ea71d1741838fb0ce7b680a6dc9"},
		"sushiswap_42161":   {"9113bbc69d0d45d8ba282a356e819286"},
		"sushiswap_8453":    {"bd99b67d409c4095bfd6caa1e7dde8b8"},
		"camelot_42161":     {"1e6f6a02dbec485da79a1f2602746a9d"},
		"dodo_1":            {"f742f5dfd20e411dbd4c4727b49ac149"},
		"dodo_10":           {"e74f5aac1e9e4668a67afe7182b707f8"},
		"dodo_56":           {"7c9227d97c264e34b451baafcaf0969c"},
		"dodo_137":          {"a91f031436184db494fa3577e6040af1"},
		"dodo_42161":        {"de1dab6ca3c1440aa23ad2b05c8c0d1d"},
		"dodo_8453":         {"d7e97ea786d043b1be79c87c19d955ff"},
		"curve_1":           {"be938b2b2187463fb3baa42bfba038ff"},
		"maverick_1":        {"006d6ba09d4741409b71ef18736b809c"},
		"maverick_56":       {"dc0362a8f883491c847cdeda4ac7e339"},
		"maverick_8453":     {"8d8d7f5b7b26411d93b72cb06c481719"},
		"maverick_42161":    {"3adf2006cb6d443ea9414104f8d3a17e"},
	},

	// Uniswap GeaphQL:     https://thegraph.com/explorer/profile/0xddaaed8b88ac0ccfdbfabdceba1c619391760f7f?view=Subgraphs
	// Pancakeswap GraphQL: https://developer.pancakeswap.finance/apis/subgraph
	// Sushiswap GraphQL:   https://thegraph.com/explorer/profile/0x23defc2ca207e7fbd84ae43b00048fb5cb4db5b2?view=Subgraphs
	// Camelot GraphQL:     https://docs.camelot.exchange/contracts/arbitrum/one-mainnet
	GraphqlAccounts: map[string]GraphqlAccount{
		"uniswap_1":          {"7afa1d2efbe70d5017ae09fa1a890332", "5zvR82QoaXYFyDEKLZ9t6v9adgnptxYpKpSbxtgVENFV"},
		"uniswap_10":         {"522e0a3096f10bf29bc4fe2d9aae0461", "Cghf4LfVqPiFw6fp6Y5X5Ubc8UpmUhSfJL82zwiBFLaj"},
		"uniswap_56":         {"057081ab34b3ca087eae988c58890e6a", "F85MNzUGYqgSHSHRGgeVMNsdnW1KtZSVgFULumXRZTw2"},
		"uniswap_137":        {"a0f78d521ae7f8119c1c7bafe9dffeba", "3hCPRGf4z88VC5rsBKU5AA9FBBq5nF3jbKJG7VZCbhjm"},
		"uniswap_8453":       {"3d1016f520496972e4e5606754e8911c", "43Hwfi3dJSoGpyas9VwNoDAv55yjgGrPpNSmbQZArzMG"},
		"uniswap_42161":      {"ef3c198785d19fba398ef5682d766d26", "FbCGRftH4a3yZugY7TnbYgPJVEv2LvMT6oF1fxPe9aJM"},
		"pancakeswap_1":      {"48e4566da15aaead28e540ea7f8fd7c3", "CJYGNhb7RvnhfBDjqpRnD3oxgyhibzc7fkAMa38YV3oS"},
		"pancakeswap_56":     {"e6db7eaf02393e457eeced4fff525d60", "Hv1GncLY5docZoGtXjo4kwbTvxm3MAhVZqBZE4sUT9eZ"},
		"pancakeswap_8453":   {"3068f4cbd2232dd0b539336bd85f8245", "BHWNsedAHtmTCzXxCCDfhPmm6iN9rxUhoRHdHKyujic3"},
		"pancakeswap_42161":  {"855a39f13d16ef4efb4758ddf8898d3e", "251MHFNN1rwjErXD2efWMpNS73SANZN8Ua192zw6iXve"},
		"sushiswap_1":        {"e04859ef837b4d9b73d65be2070bbd11", "5nnoU1nUFeWqtXgbpC54L9PWdpgo7Y9HYinR3uTMsfzs"},
		"sushiswap_10":       {"2665cb979ba0361a965454e820e34428", "Dr3FkshPgTMMDwxckz3oZdwLxaPcbzZuAbE92i6arYtJ"},
		"sushiswap_56":       {"287bae314381c945ae14695dd5137da6", "FiJDXMFCBv88GP17g2TtPh8BcA8jZozn5WRW7hCN7cUT"},
		"sushiswap_137":      {"d0195452217346bb1852cc620e56f7f8", "CqLnQY1d6DLcBYu7aZvGmt17LoNdTe4fDYnGbE2EgotR"},
		"sushiswap_8453":     {"57c2307487e1aaab986cde55ecf42a78", "Cz4Snpih41NNNPZcbj1gd3fYXPwFr5q92iWMoZjCarEb"},
		"sushiswap_42161":    {"0627101bb46562a66074f1b1fa26bfc7", "96EYD64NqmnFxMELu2QLWB95gqCmA9N96ssYsZfFiYHg"},
		"camelot_42161":      {"9d7db82aaa32347cc55e6e9983a3f82c", "3utanEBA9nqMjPnuQP1vMCCys6enSM3EawBpKTVwnUw2"},
		"dodo_1":             {"f0c3ff0191d0c272b6816ebe61f1687c", "Ef7uxVMFPgMkwKdAYEKsvcLmFjM5GSgTTXS4VpaT1nf7"},
		"dodo_10":            {"163ee3970e8b5b59ad4057a4f3c647b2", "GscH44twEB19DR8ad1oM7zZzP7wVBkvjNQmknHtigPqR"},
		"dodo_56":            {"e75fb04365a1b8aa02f5614a561810fd", "2TVoLeQK4uSASrpoHkQga4wkdpDxoUAerajauzhuSNFq"},
		"dodo_137":           {"6e4b57d434b009b3df3e6df8b03bf9fa", "6Dr9u3dGRnShJVTTHrTSTMyhDkiL32RmuTxM63TFuFkN"},
		"dodo_8453":          {"ef333bf5588495656d0283dff4d1cf31", "DraGctRHSrQDG2UsLDz882MibtVCYeCedGWG6UvpQ1yV"},
		"dodo_42161":         {"6014a0f42456d875dd1fe1fc80308a74", "5G1mhjXywQyo482AtNrmhXi6sg75f51QCAFxydP4iEoY"},
		"curve_1":            {"1bb03e4e5b69b9206f99268a941bddfb", "3wabEWN3YikX16VkKZzLdeaKHYTaB2bzyrJxNa3YbWq2"},
		"traderJoe_56":       {"404472c3c705a42e26144a8a8594ba43", "8hVhyg8FK4Zf4urb5AYFzjpJBbp23Mv5LM2FA4jk3PFr"},
		"traderJoe_42161":    {"1baec0942307d05e945fba3cba3c6b9c", "J4WH2VYKFzMjCMLWNvTgS597Bz7PjGFyProcEbid68Ns"},
		"ramsesCl_42161":     {"9f3d23607a6153c6c5955d2bf97d51be", "ATQTt3wRTgXy4canCh6t1yeczAz4ZuEkFQL2mrLXEMyQ"},
		"ramsesLegacy_42161": {"bfe959b26418af01f5ebe2d14f6ead4e", "6McF7JVH7zLWw9TVexwAiBjuMmY1A87CJhBpDx58BBVP"},
		"swapr_1":            {"e874bf3d432e668759c13798eb4854e5", "DQApa5vhVyx1sajkrF8zEFmLJTYyyMyw8WdiYt5hw9Fn"},
		"swapr_42161":        {"6dcaf28e1c5d95b61f6b286df25ab537", "H2EYoeTL5qDNeUeFecReCBvEq6BPCx4EnEDoAv7UTyL5"},
		"balancer_1":         {"d9c4d56893e2f9c247dffa5d637450ff", "C4ayEZP2yTXRAB8vSaTrgN4m9anTe9Mdm2ViyiAuV9TV"},
		"balancer_10":        {"a42e691df43d5511a0d147010a700c98", "FsmdxmvBJLGjUQPxKMRtcWKzuCNpomKuMTbSbtRtggZ7"},
		"balancer_137":       {"f734c7f62fef73ad3797b4fe8fca07a5", "H9oPAbXnobBRq1cB3HDmbZ1E8MWQyJYQjT1QDJMrdbNp"},
		"balancer_8453":      {"bb023fc70b068ede74196cdd8fea6a42", "E7XyutxXVLrp8njmjF16Hh38PCJuHm12RRyMt5ma4ctX"},
		"balancer_42161":     {"88233de0b8f7697ba19235b2322c5bbe", "98cQDy6tufTJtshDCuhh9z2kWXsQWBHVh2bqnLHsGAeS"},
	},
}

func getRpcUrl(dexId, chainId string) (string, error) {
	key := dexId + "_" + chainId
	baseUrl, exists := Data.BaseUrls[chainId]
	if !exists {
		return "", fmt.Errorf("baseUrl not found for chainId %s", chainId)
	}

	account, exists := Data.InfuraAccounts[key]
	if !exists {
		return "", fmt.Errorf("infura account not found for key %s", key)
	}

	return baseUrl + account.ApiKey, nil
}

func getGraphqlUrl(dexId, chainId string) (string, error) {
	key := dexId + "_" + chainId
	account, exists := Data.GraphqlAccounts[key]
	if !exists {
		return "", fmt.Errorf("graphql account not found for key %s", key)
	}

	return fmt.Sprintf(
		"https://gateway.thegraph.com/api/%s/subgraphs/id/%s",
		account.ApiKey,
		account.GraphId,
	), nil
}

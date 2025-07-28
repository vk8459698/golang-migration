package maverick

import (
	//"context"
	//"encoding/json"
	"fetcher-go/src/common"
	"fmt"
	"math/big"
	"strings"
	"sync"

	//"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var ONE = big.NewInt(0).Exp(big.NewInt(10), big.NewInt(18), nil) // 1e18
var MAX_TICK = big.NewInt(322378)

var CONTRACTS = map[int]map[string]string{
	1: {
		"WETH":                             "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
		"MaverickV2Factory":                "0x0A7e848Aca42d879EF06507Fca0E7b33A0a63c1e",
		"MaverickV2PoolLens":               "0x6A9EB38DE5D349Fe751E0aDb4c0D9D391f94cc8D",
		"MaverickV2Quoter":                 "0xb40AfdB85a07f37aE217E7D6462e609900dD8D7A",
		"MaverickV2Router":                 "0x62e31802c6145A2D5E842EeD8efe01fC224422fA",
		"MaverickV2Position":               "0x116193c58B40D50687c0433B2aa0cC4AE00bC32c",
		"MaverickV2BoostedPositionFactory": "0xd94C8f6D13Cf480FfAC686712C63471D1596cc29",
		"MaverickV2BoostedPositionLens":    "0x12DD145927CECF616cbD196789c89C2573A53244",
		"MaverickV2IncentiveMatcherFactory": "0x924Dd05c2325829fa4063CAbE1456273084009d7",
		"MaverickV2VotingEscrowFactory":     "0x451d47fd6207781dc053551edFD98De8d5EB4Cda",
		"MaverickV2RewardFactory":           "0x63EF1a657cc53747689B201aa07A76E9ef22f8Fe",
		"MaverickV2RewardRouter":            "0xc0C3BC532690af8922a2f260c6e1dEb6CFaB45A0",
		"MaverickV2VotingEscrowLens":        "0x102f936B0fc2E74dC34E45B601FaBaA522f381F0",
		"MaverickToken":                     "0x7448c7456a97769F6cD04F1E83A4a23cCdC46aBD",
		"LegacyMaverickVe":                  "0x4949Ac21d5b2A0cCd303C20425eeb29DCcba66D8",
		"MaverickVeV2":                      "0xC6addB3327A7D4b3b604227f82A6259Ca7112053",
		"MaverickTokenIncentiveMatcher":     "0x9172a390Cb35a15a890293f59EA5aF250b234D55",
	},
	11155111: {
		"WETH":                             "0x7b79995e5f793A07Bc00c21412e50Ecae098E7f9",
		"MaverickV2Factory":                "0x0A7e848Aca42d879EF06507Fca0E7b33A0a63c1e",
		"MaverickV2PoolLens":               "0x6A9EB38DE5D349Fe751E0aDb4c0D9D391f94cc8D",
		"MaverickV2Quoter":                 "0xb40AfdB85a07f37aE217E7D6462e609900dD8D7A",
		"MaverickV2Router":                 "0x4563d58D072C3198A66EAfCf3333024330dE9104",
		"MaverickV2Position":               "0x116193c58B40D50687c0433B2aa0cC4AE00bC32c",
		"MaverickV2BoostedPositionFactory": "0xd94C8f6D13Cf480FfAC686712C63471D1596cc29",
		"MaverickV2BoostedPositionLens":    "0x12DD145927CECF616cbD196789c89C2573A53244",
		"MaverickV2IncentiveMatcherFactory": "0x11C0F55102790f84A6F132d8B25FDFe1c96d0992",
		"MaverickV2VotingEscrowFactory":     "0x51E4AE1BA70D657eEF8e31a2Cb6a8b9AA61aB84e",
		"MaverickV2RewardFactory":           "0x873b272D7493Da5860E9c513cB805Ff3287D8470",
		"MaverickV2RewardRouter":            "0x0d17027A98F1396EC2A250d99Dc349e8cf93abb1",
		"MaverickV2VotingEscrowLens":        "0x102f936B0fc2E74dC34E45B601FaBaA522f381F0",
	},
	8453: {
		"WETH":                             "0x4200000000000000000000000000000000000006",
		"MaverickV2Factory":                "0x0A7e848Aca42d879EF06507Fca0E7b33A0a63c1e",
		"MaverickV2PoolLens":               "0x6A9EB38DE5D349Fe751E0aDb4c0D9D391f94cc8D",
		"MaverickV2Quoter":                 "0xb40AfdB85a07f37aE217E7D6462e609900dD8D7A",
		"MaverickV2Router":                 "0x5eDEd0d7E76C563FF081Ca01D9d12D6B404Df527",
		"MaverickV2Position":               "0x116193c58B40D50687c0433B2aa0cC4AE00bC32c",
		"MaverickV2BoostedPositionFactory": "0xd94C8f6D13Cf480FfAC686712C63471D1596cc29",
		"MaverickV2BoostedPositionLens":    "0x12DD145927CECF616cbD196789c89C2573A53244",
		"MaverickV2IncentiveMatcherFactory": "0xa476bb7DfCDD4E59dDaA6Ea9311A24cF28561544",
		"MaverickV2VotingEscrowFactory":     "0x1dE8C03c2D5DD021bd456bc4bB4F0ecD85f99443",
		"MaverickV2RewardFactory":           "0x1cdC67950a68256c5157987bBF700e94595807F8",
		"MaverickV2RewardRouter":            "0xE7c73727c1b67A2fA47E63DCBaa4859777aeF392",
		"MaverickV2VotingEscrowLens":        "0x102f936B0fc2E74dC34E45B601FaBaA522f381F0",
		"MaverickToken":                     "0x64b88c73A5DfA78D1713fE1b4c69a22d7E0faAa7",
		"MaverickVeV2":                      "0x05b1b801191B41a21B9C0bFd4c4ef8952eb28cd9",
		"MaverickTokenIncentiveMatcher":     "0xc84bDDC0C45FEeFB0F59e1c48332E4d47e29D112",
	},
	42161: {
		"WETH":                             "0x82aF49447D8a07e3bd95BD0d56f35241523fBab1",
		"MaverickV2Factory":                "0x0A7e848Aca42d879EF06507Fca0E7b33A0a63c1e",
		"MaverickV2PoolLens":               "0x6A9EB38DE5D349Fe751E0aDb4c0D9D391f94cc8D",
		"MaverickV2Quoter":                 "0xb40AfdB85a07f37aE217E7D6462e609900dD8D7A",
		"MaverickV2Router":                 "0x5c3b380e5Aeec389d1014Da3Eb372FA2C9e0fc76",
		"MaverickV2Position":               "0x116193c58B40D50687c0433B2aa0cC4AE00bC32c",
		"MaverickV2BoostedPositionFactory": "0xd94C8f6D13Cf480FfAC686712C63471D1596cc29",
		"MaverickV2BoostedPositionLens":    "0x12DD145927CECF616cbD196789c89C2573A53244",
		"MaverickV2IncentiveMatcherFactory": "0x11C0F55102790f84A6F132d8B25FDFe1c96d0992",
		"MaverickV2VotingEscrowFactory":     "0x51E4AE1BA70D657eEF8e31a2Cb6a8b9AA61aB84e",
		"MaverickV2RewardFactory":           "0x873b272D7493Da5860E9c513cB805Ff3287D8470",
		"MaverickV2RewardRouter":            "0x293A7D159C5AD1b36b784998DE5563fe36963460",
		"MaverickV2VotingEscrowLens":        "0x102f936B0fc2E74dC34E45B601FaBaA522f381F0",
		"MaverickToken":                     "0x7448c7456a97769F6cD04F1E83A4a23cCdC46aBD",
		"MaverickVeV2":                      "0xd5d8cB7569BB843c3b8FA98dBD5960d37E83eA8d",
		"MaverickTokenIncentiveMatcher":     "0xB1F334176AadC61F74afc6381210e8786CcEc37D",
	},
	324: {
		"WETH":                             "0x5AEa5775959fBC2557Cc8789bC1bf90A239D9a91",
		"MaverickV2Factory":                "0x7A6902af768a06bdfAb4F076552036bf68D1dc56",
		"MaverickV2PoolLens":               "0x9439280a7d04FCa28d12a4eB74c92173241d5b2F",
		"MaverickV2Quoter":                 "0x3e1c4b57c9d9624f2841f07C6328D3c25ca30C79",
		"MaverickV2Router":                 "0xad8262e847676E7eDdAFEe664c4fd492789260ba",
		"MaverickV2Position":               "0x4D93c58B348d99969257cec007cFb31B410b21A0",
		"MaverickV2BoostedPositionFactory": "0x270a03bfc3EA123c041d4A0c72D30202A514D845",
		"MaverickV2BoostedPositionLens":    "0xd32CE31CaC98CAC0631764B8286358c0606D87F9",
		"MaverickV2IncentiveMatcherFactory": "0x11244D8b724De7788f62667791e35284E191745F",
		"MaverickV2VotingEscrowFactory":     "0x521B444d5f9bb4B36CDd771f4D85cCd0B291FB92",
		"MaverickV2RewardFactory":           "0xc9e5F0832C96F8E2EEDe472C1B87621Cbb86D7e0",
		"MaverickV2RewardRouter":            "0x432e6791d35dc6c638f44E949A5c0228e4048244",
		"MaverickV2VotingEscrowLens":        "0x74E56528CDd2F831cc4ecc9414bCE9C4d540ceC7",
		"MaverickToken":                     "0x787c09494Ec8Bcb24DcAf8659E7d5D69979eE508",
		"LegacyMaverickVe":                  "0x7EDcB053d4598a145DdaF5260cf89A32263a2807",
		"MaverickVeV2":                      "0xe86151Af9cc43533add87921c381dA11c314DEBf",
		"MaverickTokenIncentiveMatcher":     "0x57FA162aCb48376455c5Ff4D45FE0d36E947D79b",
	},
	56: {
		"WETH":                             "0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c",
		"MaverickV2Factory":                "0x0A7e848Aca42d879EF06507Fca0E7b33A0a63c1e",
		"MaverickV2PoolLens":               "0x6A9EB38DE5D349Fe751E0aDb4c0D9D391f94cc8D",
		"MaverickV2Quoter":                 "0xb40AfdB85a07f37aE217E7D6462e609900dD8D7A",
		"MaverickV2Router":                 "0x374bFCc264678c67a582D067AD91f1951bC6b20f",
		"MaverickV2Position":               "0x116193c58B40D50687c0433B2aa0cC4AE00bC32c",
		"MaverickV2BoostedPositionFactory": "0xd94C8f6D13Cf480FfAC686712C63471D1596cc29",
		"MaverickV2BoostedPositionLens":    "0x12DD145927CECF616cbD196789c89C2573A53244",
		"MaverickV2IncentiveMatcherFactory": "0x53EEE0a9d1D301eA570329C298Af3f19d1D556c7",
		"MaverickV2VotingEscrowFactory":     "0x790d33B4271EDD0a611d91E971F2143D8a7DD936",
		"MaverickV2RewardFactory":           "0x443b1F86D45C1dDC60b355D5A8A931656aB25267",
		"MaverickV2RewardRouter":            "0x5DeB1bAe837374f988d8a30Cc0Fbccbc63892Bb3",
		"MaverickV2VotingEscrowLens":        "0x102f936B0fc2E74dC34E45B601FaBaA522f381F0",
		"MaverickToken":                     "0xd691d9a68C887BDF34DA8c36f63487333ACfD103",
		"LegacyMaverickVe":                  "0xE6108f1869d37E5076a56168C66A1607EdB10819",
		"MaverickVeV2":                      "0x675178AE86A75EE7D7Ef81e30a91E1798306094C",
		"MaverickTokenIncentiveMatcher":     "0x053D0eC15e60c7D8936Ab966A82BB62cCb7E3Ced",
	},
}

type PoolState struct {
	ReserveA string `json:"reserveA"`
	ReserveB string `json:"reserveB"`
}

type PoolDetails struct {
	ID             string    `json:"id"`
	TokenA         string    `json:"tokenA"`
	TokenB         string    `json:"tokenB"`
	FeeAIn         string    `json:"feeAIn"`
	FeeBIn         string    `json:"feeBIn"`
	State          PoolState `json:"state"`
	TokenAScale    string    `json:"tokenAScale"`
	TokenBScale    string    `json:"tokenBScale"`
	ProtocolFeeA   string    `json:"protocolFeeA"`
	ProtocolFeeB   string    `json:"protocolFeeB"`
	SqrtPrice      string    `json:"sqrtPrice"`
	Liquidity      string    `json:"liquidity"`
	SqrtLowerPrice string    `json:"sqrtLowerPrice"`
	SqrtUpperPrice string    `json:"sqrtUpperPrice"`
}

type TickSqrtPrices struct {
	SqrtLowerPrice *big.Int
	SqrtUpperPrice *big.Int
}

type TickSqrtPriceAndL struct {
	Liquidity *big.Int
	SqrtPrice *big.Int
}

func sqrt(x *big.Int) *big.Int {
	if x.Sign() == 0 {
		return big.NewInt(0)
	}

	y := new(big.Int).Set(x)
	z := big.NewInt(181)

	// Constants for comparison
	const_128_bit := new(big.Int)
	const_128_bit.SetString("0x10000000000000000000000000000000000", 0)

	const_64_bit := new(big.Int)
	const_64_bit.SetString("0x1000000000000000000", 0)

	const_32_bit := new(big.Int)
	const_32_bit.SetString("0x10000000000", 0)

	const_16_bit := new(big.Int)
	const_16_bit.SetString("0x1000000", 0)

	if y.Cmp(const_128_bit) >= 0 {
		y.Rsh(y, 128)
		z.Lsh(z, 64)
	}
	if y.Cmp(const_64_bit) >= 0 {
		y.Rsh(y, 64)
		z.Lsh(z, 32)
	}
	if y.Cmp(const_32_bit) >= 0 {
		y.Rsh(y, 32)
		z.Lsh(z, 16)
	}
	if y.Cmp(const_16_bit) >= 0 {
		y.Rsh(y, 16)
		z.Lsh(z, 8)
	}

	// z = z.mul(y.add(65536)).shr(18)
	temp := new(big.Int).Add(y, big.NewInt(65536))
	z.Mul(z, temp)
	z.Rsh(z, 18)

	// Newton's method iterations
	for i := 0; i < 7; i++ {
		temp := new(big.Int).Div(x, z)
		z.Add(z, temp)
		z.Rsh(z, 1)
	}

	// Final floor check
	zFloorCheck := new(big.Int).Div(x, z)
	if zFloorCheck.Cmp(z) < 0 {
		z.Sub(z, big.NewInt(1))
	}

	return z
}

func mulDivDown(a, b, c *big.Int) *big.Int {
	result := new(big.Int).Mul(a, b)
	return result.Div(result, c)
}

func mulDown(a, b *big.Int) *big.Int {
	return mulDivDown(a, b, ONE)
}

func divDown(a, b *big.Int) *big.Int {
	return mulDivDown(a, ONE, b)
}

func boundValue(value, limitA, limitB *big.Int) *big.Int {
	var lower, upper *big.Int
	if limitA.Cmp(limitB) < 0 {
		lower = limitA
		upper = limitB
	} else {
		lower = limitB
		upper = limitA
	}

	if value.Cmp(lower) < 0 {
		return new(big.Int).Set(lower)
	}
	if value.Cmp(upper) > 0 {
		return new(big.Int).Set(upper)
	}
	return new(big.Int).Set(value)
}

func subTickIndex(tickSpacing, _tick *big.Int) *big.Int {
	absTick := new(big.Int)
	if _tick.Sign() < 0 {
		absTick.Mul(_tick, big.NewInt(-1))
	} else {
		absTick.Set(_tick)
	}

	subTick := new(big.Int).Mul(absTick, tickSpacing)

	if subTick.Cmp(MAX_TICK) > 0 {
		panic(fmt.Sprintf("TickMaxExceeded: %s", _tick.String()))
	}

	return subTick
}

func tickSqrtPrice(tickSpacing, tick *big.Int) *big.Int {
	absTick := subTickIndex(tickSpacing, tick)

	var ratio *big.Int
	if new(big.Int).And(absTick, big.NewInt(0x1)).Sign() == 0 {
		ratio = new(big.Int)
		ratio.SetString("0x100000000000000000000000000000000", 0)
	} else {
		ratio = new(big.Int)
		ratio.SetString("0xfffcb933bd6fad9d3af5f0b9f25db4d6", 0)
	}

	// Apply bit-by-bit transformations
	constants := []string{
		"0xfff97272373d41fd789c8cb37ffcaa1c",
		"0xfff2e50f5f656ac9229c67059486f389",
		"0xffe5caca7e10e81259b3cddc7a064941",
		"0xffcb9843d60f67b19e8887e0bd251eb7",
		"0xff973b41fa98cd2e57b660be99eb2c4a",
		"0xff2ea16466c9838804e327cb417cafcb",
		"0xfe5dee046a99d51e2cc356c2f617dbe0",
		"0xfcbe86c7900aecf64236ab31f1f9dcb5",
		"0xf987a7253ac4d9194200696907cf2e37",
		"0xf3392b0822b88206f8abe8a3b44dd9be",
		"0xe7159475a2c578ef4f1d17b2b235d480",
		"0xd097f3bdfd254ee83bdd3f248e7e785e",
		"0xa9f746462d8f7dd10e744d913d033333",
		"0x70d869a156ddd32a39e257bc3f50aa9b",
		"0x31be135f97da6e09a19dc367e3b6da40",
		"0x9aa508b5b7e5a9780b0cc4e25d61a56",
		"0x5d6af8dedbcb3a6ccb7ce618d14225",
		"0x2216e584f630389b2052b8db590e",
	}

	for i, constStr := range constants {
		bit := big.NewInt(int64(1 << (i + 1)))
		if new(big.Int).And(absTick, bit).Sign() != 0 {
			temp := new(big.Int)
			temp.SetString(constStr, 0)
			ratio.Mul(ratio, temp)
			ratio.Rsh(ratio, 128)
		}
	}

	if tick.Sign() > 0 {
		maxUint256 := new(big.Int)
		maxUint256.SetString("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", 0)
		ratio = new(big.Int).Div(maxUint256, ratio)
	}

	result := new(big.Int).Mul(ratio, ONE)
	return result.Rsh(result, 128)
}

func getTickL(reserveA, reserveB, sqrtLowerTickPrice, sqrtUpperTickPrice *big.Int) *big.Int {
	diff := new(big.Int).Sub(sqrtUpperTickPrice, sqrtLowerTickPrice)

	precisionBump := 0

	// Check if both reserves are small enough for precision bump
	reserveA_shifted := new(big.Int).Rsh(reserveA, 78)
	reserveB_shifted := new(big.Int).Rsh(reserveB, 78)

	if reserveA_shifted.Sign() == 0 && reserveB_shifted.Sign() == 0 {
		precisionBump = 57
		reserveA = new(big.Int).Lsh(reserveA, uint(precisionBump))
		reserveB = new(big.Int).Lsh(reserveB, uint(precisionBump))
	}

	if reserveB.Sign() == 0 {
		result := divDown(reserveA, diff)
		return result.Rsh(result, uint(precisionBump))
	}

	if reserveA.Sign() == 0 {
		temp1 := mulDown(reserveB, sqrtLowerTickPrice)
		temp2 := mulDivDown(temp1, sqrtUpperTickPrice, diff)
		return temp2.Rsh(temp2, uint(precisionBump))
	}

	// Calculate b = (reserveA / sqrtUpperTickPrice + reserveB * sqrtLowerTickPrice) / 2
	term1 := divDown(reserveA, sqrtUpperTickPrice)
	term2 := mulDown(reserveB, sqrtLowerTickPrice)
	b := new(big.Int).Add(term1, term2)
	b.Rsh(b, 1)

	// Calculate bSquared = b * b / ONE
	bSquared := new(big.Int).Mul(b, b)
	bSquared.Div(bSquared, ONE)

	// Calculate extra = (reserveA * reserveB * diff) / sqrtUpperTickPrice
	temp := new(big.Int).Mul(reserveA, reserveB)
	temp.Mul(temp, diff)
	extra := new(big.Int).Div(temp, sqrtUpperTickPrice)

	// Calculate innerSqrt = sqrt(bSquared + extra)
	innerSqrt := new(big.Int).Add(bSquared, extra)
	sqrtInner := sqrt(innerSqrt)

	// Calculate term = sqrtInner * 1000000000
	term := new(big.Int).Mul(sqrtInner, big.NewInt(1000000000))

	// Calculate numerator = (b + term) * sqrtUpperTickPrice
	numerator := new(big.Int).Add(b, term)
	numerator.Mul(numerator, sqrtUpperTickPrice)

	// Calculate liquidity = numerator / diff >> precisionBump
	liquidity := new(big.Int).Div(numerator, diff)
	return liquidity.Rsh(liquidity, uint(precisionBump))
}

func getSqrtPrice(reserveA, reserveB, sqrtLowerTickPrice, sqrtUpperTickPrice, liquidity *big.Int) *big.Int {
	if reserveA.Sign() == 0 {
		return new(big.Int).Set(sqrtLowerTickPrice)
	}
	if reserveB.Sign() == 0 {
		return new(big.Int).Set(sqrtUpperTickPrice)
	}

	// numerator = reserveA + liquidity * sqrtLowerTickPrice
	temp1 := mulDown(liquidity, sqrtLowerTickPrice)
	numerator := new(big.Int).Add(reserveA, temp1)

	// denominator = reserveB + liquidity / sqrtUpperTickPrice
	temp2 := divDown(liquidity, sqrtUpperTickPrice)
	denominator := new(big.Int).Add(reserveB, temp2)

	// ratio = (numerator / denominator) * ONE
	ratio := divDown(numerator, denominator)
	ratio.Mul(ratio, ONE)

	// sqrtPrice = sqrt(ratio)
	sqrtPrice := sqrt(ratio)

	// Bound the result
	return boundValue(sqrtPrice, sqrtLowerTickPrice, sqrtUpperTickPrice)
}

func GetTickSqrtPrices(tickSpacing, tick *big.Int) TickSqrtPrices {
	sqrtLowerPrice := tickSqrtPrice(tickSpacing, tick)
	nextTick := new(big.Int).Add(tick, big.NewInt(1))
	sqrtUpperPrice := tickSqrtPrice(tickSpacing, nextTick)

	return TickSqrtPrices{
		SqrtLowerPrice: sqrtLowerPrice,
		SqrtUpperPrice: sqrtUpperPrice,
	}
}

func GetTickSqrtPriceAndL(reserveA, reserveB, sqrtLowerTickPrice, sqrtUpperTickPrice *big.Int) TickSqrtPriceAndL {
	liquidity := getTickL(reserveA, reserveB, sqrtLowerTickPrice, sqrtUpperTickPrice)
	sqrtPrice := getSqrtPrice(reserveA, reserveB, sqrtLowerTickPrice, sqrtUpperTickPrice, liquidity)

	return TickSqrtPriceAndL{
		Liquidity: liquidity,
		SqrtPrice: sqrtPrice,
	}
}

// ABI definitions would need to be loaded from JSON files in practice
var MaverickV2PoolABI = `[
	{"inputs":[],"name":"tokenA","outputs":[{"type":"address"}],"stateMutability":"view","type":"function"},
	{"inputs":[],"name":"tokenB","outputs":[{"type":"address"}],"stateMutability":"view","type":"function"},
	{"inputs":[{"type":"bool","name":"tokenAIn"}],"name":"fee","outputs":[{"type":"uint256"}],"stateMutability":"view","type":"function"},
	{"inputs":[],"name":"getState","outputs":[{"type":"uint128","name":"reserveA"},{"type":"uint128","name":"reserveB"},{"type":"int64","name":"lastTwaD8"},{"type":"int64","name":"lastLogPriceD8"},{"type":"uint40","name":"lastTimestamp"},{"type":"int32","name":"activeTick"},{"type":"bool","name":"isLocked"},{"type":"uint32","name":"binCounter"},{"type":"uint8","name":"protocolFeeRatioD3"}],"stateMutability":"view","type":"function"},
	{"inputs":[],"name":"tokenAScale","outputs":[{"type":"uint256"}],"stateMutability":"view","type":"function"},
	{"inputs":[],"name":"tokenBScale","outputs":[{"type":"uint256"}],"stateMutability":"view","type":"function"},
	{"inputs":[],"name":"protocolFeeA","outputs":[{"type":"uint128"}],"stateMutability":"view","type":"function"},
	{"inputs":[],"name":"protocolFeeB","outputs":[{"type":"uint128"}],"stateMutability":"view","type":"function"}
]`

var MaverickV2FactoryABI = `[
	{"inputs":[{"type":"uint256","name":"startIndex"},{"type":"uint256","name":"endIndex"}],"name":"lookup","outputs":[{"type":"address[]"}],"stateMutability":"view","type":"function"}
]`

var MaverickV2PoolLensABI = `[
	{"inputs":[{"type":"address","name":"pool"},{"type":"int32","name":"tick"}],"name":"getTickSqrtPriceAndL","outputs":[{"type":"uint256"},{"type":"uint256"}],"stateMutability":"view","type":"function"},
	{"inputs":[{"type":"address","name":"pool"},{"type":"uint8","name":"kind"}],"name":"getTicksAroundActiveWLiquidity","outputs":[{"type":"int32[]"},{"type":"uint256[]"},{"type":"uint256[]"},{"type":"uint256[]"},{"type":"uint256[]"}],"stateMutability":"view","type":"function"}
]`

// Multicall structure for batch calls
type Call struct {
	Target   ethcommon.Address
	CallData []byte
}

type CallResult struct {
	Success    bool
	ReturnData []byte
}

func fetchPoolDetails(rpcUrl string, chainId int, poolAddress string) (*PoolDetails, error) {
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to client: %v", err)
	}
	defer client.Close()

	poolAddr := ethcommon.HexToAddress(poolAddress)

	// Parse ABIs
	poolABI, err := abi.JSON(strings.NewReader(MaverickV2PoolABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse pool ABI: %v", err)
	}

	lensABI, err := abi.JSON(strings.NewReader(MaverickV2PoolLensABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse lens ABI: %v", err)
	}

	// Create contract instances
	poolContract := bind.NewBoundContract(poolAddr, poolABI, client, client, client)
	lensAddr := ethcommon.HexToAddress(CONTRACTS[chainId]["MaverickV2PoolLens"])
	lensContract := bind.NewBoundContract(lensAddr, lensABI, client, client, client)

	// Call tokenA
	var tokenAResult []interface{}
	err = poolContract.Call(&bind.CallOpts{}, &tokenAResult, "tokenA")
	if err != nil {
		return nil, fmt.Errorf("failed to call tokenA: %v", err)
	}
	tokenA := tokenAResult[0].(ethcommon.Address)

	// Call tokenB
	var tokenBResult []interface{}
	err = poolContract.Call(&bind.CallOpts{}, &tokenBResult, "tokenB")
	if err != nil {
		return nil, fmt.Errorf("failed to call tokenB: %v", err)
	}
	tokenB := tokenBResult[0].(ethcommon.Address)

	// Call fee for tokenA (true)
	var feeAInResult []interface{}
	err = poolContract.Call(&bind.CallOpts{}, &feeAInResult, "fee", true)
	if err != nil {
		return nil, fmt.Errorf("failed to call feeAIn: %v", err)
	}
	feeAIn := feeAInResult[0].(*big.Int)

	// Call fee for tokenB (false)
	var feeBInResult []interface{}
	err = poolContract.Call(&bind.CallOpts{}, &feeBInResult, "fee", false)
	if err != nil {
		return nil, fmt.Errorf("failed to call feeBIn: %v", err)
	}
	feeBIn := feeBInResult[0].(*big.Int)

	// Call getState
	var stateResult []interface{}
	err = poolContract.Call(&bind.CallOpts{}, &stateResult, "getState")
	if err != nil {
		return nil, fmt.Errorf("failed to call getState: %v", err)
	}

	// Extract state values
	reserveA := stateResult[0].(*big.Int)
	reserveB := stateResult[1].(*big.Int)
	activeTick := stateResult[5].(int32)

	// Call tokenAScale
	var tokenAScaleResult []interface{}
	err = poolContract.Call(&bind.CallOpts{}, &tokenAScaleResult, "tokenAScale")
	if err != nil {
		return nil, fmt.Errorf("failed to call tokenAScale: %v", err)
	}
	tokenAScale := tokenAScaleResult[0].(*big.Int)

	// Call tokenBScale
	var tokenBScaleResult []interface{}
	err = poolContract.Call(&bind.CallOpts{}, &tokenBScaleResult, "tokenBScale")
	if err != nil {
		return nil, fmt.Errorf("failed to call tokenBScale: %v", err)
	}
	tokenBScale := tokenBScaleResult[0].(*big.Int)

	// Call protocolFeeA
	var protocolFeeAResult []interface{}
	err = poolContract.Call(&bind.CallOpts{}, &protocolFeeAResult, "protocolFeeA")
	if err != nil {
		return nil, fmt.Errorf("failed to call protocolFeeA: %v", err)
	}
	protocolFeeA := protocolFeeAResult[0].(*big.Int)

	// Call protocolFeeB
	var protocolFeeBResult []interface{}
	err = poolContract.Call(&bind.CallOpts{}, &protocolFeeBResult, "protocolFeeB")
	if err != nil {
		return nil, fmt.Errorf("failed to call protocolFeeB: %v", err)
	}
	protocolFeeB := protocolFeeBResult[0].(*big.Int)

	// FIXED: Pass activeTick as int32 directly instead of converting to big.Int
	var sqrtPriceAndLResult []interface{}
	var ticksAroundActiveResult []interface{}
	
	// Try to call lens functions, but don't fail the entire pool fetch if they fail
	err = lensContract.Call(&bind.CallOpts{}, &sqrtPriceAndLResult, "getTickSqrtPriceAndL", poolAddr, activeTick)
	if err != nil {
		fmt.Printf("[UtilsMaverick] Warning: failed to call getTickSqrtPriceAndL for pool %s: %v\n", poolAddress, err)
		// Continue with nil results - defaults will be set below
	}

	err = lensContract.Call(&bind.CallOpts{}, &ticksAroundActiveResult, "getTicksAroundActiveWLiquidity", poolAddr, uint8(0))
	if err != nil {
		fmt.Printf("[UtilsMaverick] Warning: failed to call getTicksAroundActiveWLiquidity for pool %s: %v\n", poolAddress, err)
		// Continue with nil results - defaults will be set below
	}

	// Extract sqrt price and liquidity
	var sqrtPrice, liquidity *big.Int
	if len(sqrtPriceAndLResult) >= 2 {
		sqrtPrice = sqrtPriceAndLResult[0].(*big.Int)
		liquidity = sqrtPriceAndLResult[1].(*big.Int)
	}

	// Extract sqrt lower and upper prices
	var sqrtLowerPrice, sqrtUpperPrice *big.Int
	if len(ticksAroundActiveResult) >= 5 {
		if lowerPricesInterface, ok := ticksAroundActiveResult[3].([]interface{}); ok && len(lowerPricesInterface) > 0 {
			sqrtLowerPrice = lowerPricesInterface[0].(*big.Int)
		}
		if upperPricesInterface, ok := ticksAroundActiveResult[4].([]interface{}); ok && len(upperPricesInterface) > 0 {
			sqrtUpperPrice = upperPricesInterface[0].(*big.Int)
		}
	}

	// Set default values if nil
	if sqrtPrice == nil {
		sqrtPrice = big.NewInt(0)
	}
	if liquidity == nil {
		liquidity = big.NewInt(0)
	}
	if sqrtLowerPrice == nil {
		sqrtLowerPrice = big.NewInt(0)
	}
	if sqrtUpperPrice == nil {
		sqrtUpperPrice = big.NewInt(0)
	}

	return &PoolDetails{
		ID:     poolAddress,
		TokenA: tokenA.Hex(),
		TokenB: tokenB.Hex(),
		FeeAIn: feeAIn.String(),
		FeeBIn: feeBIn.String(),
		State: PoolState{
			ReserveA: reserveA.String(),
			ReserveB: reserveB.String(),
		},
		TokenAScale:    tokenAScale.String(),
		TokenBScale:    tokenBScale.String(),
		ProtocolFeeA:   protocolFeeA.String(),
		ProtocolFeeB:   protocolFeeB.String(),
		SqrtPrice:      sqrtPrice.String(),
		Liquidity:      liquidity.String(),
		SqrtLowerPrice: sqrtLowerPrice.String(),
		SqrtUpperPrice: sqrtUpperPrice.String(),
	}, nil
}

func FetchPools(chainId int, dexName string, dexConfig *common.DexConfiguration) ([]any, error) {
	rpcUrl := dexConfig.RpcUrl
	if rpcUrl == "" {
		return nil, fmt.Errorf("invalid or missing rpcUrl in dexConfig")
	}

	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to client: %v", err)
	}
	defer client.Close()

	// Check if contracts exist for this chain
	contracts, exists := CONTRACTS[chainId]
	if !exists {
		return nil, fmt.Errorf("no contracts configured for chain %d", chainId)
	}

	factoryAddr := ethcommon.HexToAddress(contracts["MaverickV2Factory"])
	factoryABI, err := abi.JSON(strings.NewReader(MaverickV2FactoryABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse factory ABI: %v", err)
	}

	factoryContract := bind.NewBoundContract(factoryAddr, factoryABI, client, client, client)

	// Get pool addresses
	var lookupResult []interface{}
	err = factoryContract.Call(&bind.CallOpts{}, &lookupResult, "lookup", big.NewInt(0), big.NewInt(1000))
	if err != nil {
		return nil, fmt.Errorf("failed to call lookup: %v", err)
	}

	poolAddresses := lookupResult[0].([]ethcommon.Address)

	// Process pools in parallel
	type result struct {
		index int
		pool  any
		err   error
	}

	results := make(chan result, len(poolAddresses))
	var wg sync.WaitGroup

	// Limit concurrent goroutines
	maxConcurrent := 10
	semaphore := make(chan struct{}, maxConcurrent)

	for i, poolAddr := range poolAddresses {
		wg.Add(1)
		go func(index int, poolAddress string) {
			defer wg.Done()

			semaphore <- struct{}{}        // Acquire semaphore
			defer func() { <-semaphore }() // Release semaphore

			details, err := fetchPoolDetails(rpcUrl, chainId, poolAddress)
			if err != nil {
				fmt.Printf("[UtilsMaverick] Failed to fetch details for pool %s on chain %d: %v\n", poolAddress, chainId, err)
				results <- result{index: index, pool: nil, err: err}
				return
			}

			results <- result{index: index, pool: details, err: nil}
		}(i, poolAddr.Hex())
	}

	wg.Wait()
	close(results)

	// Collect results
	allPools := []any{}
	for res := range results {
		if res.err == nil && res.pool != nil {
			allPools = append(allPools, res.pool)
		}
	}

	return allPools, nil
}

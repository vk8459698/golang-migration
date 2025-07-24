package ramses

import (
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"strings"
    "context"
	"fetcher-go/src/common"
	"sync"

	"github.com/machinebox/graphql"
	"github.com/shopspring/decimal"
)

var (
	MIN_TICK = -887272
	MAX_TICK = -MIN_TICK

	MIN_SQRT_RATIO, _ = big.NewInt(0).SetString("4295128739", 10)
	MAX_SQRT_RATIO, _ = big.NewInt(0).SetString("1461446703485210103287273052203988822378723970342", 10)

	ONE          = big.NewInt(1)
	TWO          = big.NewInt(2)
	ZERO         = big.NewInt(0)
	NEGATIVE_ONE = big.NewInt(-1)

	MAX_FEE = new(big.Int).Exp(big.NewInt(10), big.NewInt(6), nil)

	Q96  = new(big.Int).Exp(big.NewInt(2), big.NewInt(96), nil)
	Q32  = new(big.Int).Exp(big.NewInt(2), big.NewInt(32), nil)
	Q192 = new(big.Int).Exp(Q96, big.NewInt(2), nil)
	Q64  = new(big.Int).Exp(big.NewInt(2), big.NewInt(64), nil)
	Q128 = new(big.Int).Exp(Q64, big.NewInt(2), nil)

	MaxUint160 = new(big.Int).Sub(new(big.Int).Exp(big.NewInt(2), big.NewInt(160), nil), ONE)
	MaxUint256 = new(big.Int).Sub(new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil), ONE)
)

var POWERS_OF_2 = []struct {
	Power int64
	Value *big.Int
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

type Tick struct {
	TickIdx      string `json:"tickIdx"`
	LiquidityNet string `json:"liquidityNet"`
}

func mostSignificantBit(x *big.Int) int {
	if x.Cmp(ZERO) <= 0 {
		panic("ZERO")
	}
	if x.Cmp(MaxUint256) > 0 {
		panic("MAX")
	}

	msb := 0
	temp := new(big.Int).Set(x)
	for _, pow := range POWERS_OF_2 {
		if temp.Cmp(pow.Value) >= 0 {
			temp.Rsh(temp, uint(pow.Power))
			msb += int(pow.Power)
		}
	}
	return msb
}

func isBelowSmallest(ticks []Tick, tick int) bool {
	if len(ticks) == 0 {
		panic("LENGTH")
	}
	tickIdx, _ := strconv.Atoi(ticks[0].TickIdx)
	return tick < tickIdx
}

func binarySearch(ticks []Tick, tick int) int {
	if isBelowSmallest(ticks, tick) {
		panic("BELOW_SMALLEST")
	}

	l := 0
	r := len(ticks) - 1
	var i int

	for {
		i = (l + r) / 2
		tickIdx, _ := strconv.Atoi(ticks[i].TickIdx)
		
		var nextTickIdx int
		if i == len(ticks)-1 {
			nextTickIdx = tick + 1 // Make condition true
		} else {
			nextTickIdx, _ = strconv.Atoi(ticks[i+1].TickIdx)
		}

		if tickIdx <= tick && (i == len(ticks)-1 || nextTickIdx > tick) {
			return i
		}

		if tickIdx < tick {
			l = i + 1
		} else {
			r = i - 1
		}
	}
}

func nextInitializedTickFunction(ticks []Tick, tick int, lte bool) Tick {
	if lte {
		if isBelowSmallest(ticks, tick) {
			panic("BELOW_SMALLEST")
		}
		if isAtOrAboveLargest(ticks, tick) {
			return ticks[len(ticks)-1]
		}
		index := binarySearch(ticks, tick)
		return ticks[index]
	} else {
		if isAtOrAboveLargest(ticks, tick) {
			panic("AT_OR_ABOVE_LARGEST")
		}
		if isBelowSmallest(ticks, tick) {
			return ticks[0]
		}
		index := binarySearch(ticks, tick)
		return ticks[index+1]
	}
}

func isAtOrAboveLargest(ticks []Tick, tick int) bool {
	if len(ticks) == 0 {
		panic("LENGTH")
	}
	tickIdx, _ := strconv.Atoi(ticks[len(ticks)-1].TickIdx)
	return tick >= tickIdx
}

func nextInitializedTickWithinOneWord(ticks []Tick, tick int, lte bool, tickSpacing int) (int, bool) {
	compressed := tick / tickSpacing

	if lte {
		wordPos := compressed >> 8
		minimum := (wordPos << 8) * tickSpacing

		if isBelowSmallest(ticks, tick) {
			return minimum, false
		}

		nextTick := nextInitializedTickFunction(ticks, tick, lte)
		index, _ := strconv.Atoi(nextTick.TickIdx)
		nextInitializedTick := maximum(minimum, index)
		return nextInitializedTick, nextInitializedTick == index
	} else {
		wordPos := (compressed + 1) >> 8
		maximum := (((wordPos + 1) << 8) - 1) * tickSpacing

		if isAtOrAboveLargest(ticks, tick) {
			return maximum, false
		}

		nextTick := nextInitializedTickFunction(ticks, tick, lte)
		index, _ := strconv.Atoi(nextTick.TickIdx)
		nextInitializedTick := minimum(maximum, index)
		return nextInitializedTick, nextInitializedTick == index
	}
}

func maximum(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minimum(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func mulShift(val *big.Int, mulBy string) *big.Int {
	mulByBig := new(big.Int)
	mulByBig.SetString(mulBy, 0)
	
	result := new(big.Int).Mul(val, mulByBig)
	result.Rsh(result, 128)
	return result
}

func getSqrtRatioAtTick(tick int) *big.Int {
	if tick < MIN_TICK || tick > MAX_TICK {
		panic(fmt.Sprintf("TICK: invalid tick value (%d)", tick))
	}

	absTick := tick
	if tick < 0 {
		absTick = -tick
	}

	var ratio *big.Int
	if (absTick & 0x1) != 0 {
		ratio = new(big.Int)
		ratio.SetString("0xfffcb933bd6fad37aa2d162d1a594001", 0)
	} else {
		ratio = new(big.Int)
		ratio.SetString("0x100000000000000000000000000000000", 0)
	}

	if (absTick & 0x2) != 0 {
		ratio = mulShift(ratio, "0xfff97272373d413259a46990580e213a")
	}
	if (absTick & 0x4) != 0 {
		ratio = mulShift(ratio, "0xfff2e50f5f656932ef12357cf3c7fdcc")
	}
	if (absTick & 0x8) != 0 {
		ratio = mulShift(ratio, "0xffe5caca7e10e4e61c3624eaa0941cd0")
	}
	if (absTick & 0x10) != 0 {
		ratio = mulShift(ratio, "0xffcb9843d60f6159c9db58835c926644")
	}
	if (absTick & 0x20) != 0 {
		ratio = mulShift(ratio, "0xff973b41fa98c081472e6896dfb254c0")
	}
	if (absTick & 0x40) != 0 {
		ratio = mulShift(ratio, "0xff2ea16466c96a3843ec78b326b52861")
	}
	if (absTick & 0x80) != 0 {
		ratio = mulShift(ratio, "0xfe5dee046a99a2a811c461f1969c3053")
	}
	if (absTick & 0x100) != 0 {
		ratio = mulShift(ratio, "0xfcbe86c7900a88aedcffc83b479aa3a4")
	}
	if (absTick & 0x200) != 0 {
		ratio = mulShift(ratio, "0xf987a7253ac413176f2b074cf7815e54")
	}
	if (absTick & 0x400) != 0 {
		ratio = mulShift(ratio, "0xf3392b0822b70005940c7a398e4b70f3")
	}
	if (absTick & 0x800) != 0 {
		ratio = mulShift(ratio, "0xe7159475a2c29b7443b29c7fa6e889d9")
	}
	if (absTick & 0x1000) != 0 {
		ratio = mulShift(ratio, "0xd097f3bdfd2022b8845ad8f792aa5825")
	}
	if (absTick & 0x2000) != 0 {
		ratio = mulShift(ratio, "0xa9f746462d870fdf8a65dc1f90e061e5")
	}
	if (absTick & 0x4000) != 0 {
		ratio = mulShift(ratio, "0x70d869a156d2a1b890bb3df62baf32f7")
	}
	if (absTick & 0x8000) != 0 {
		ratio = mulShift(ratio, "0x31be135f97d08fd981231505542fcfa6")
	}
	if (absTick & 0x10000) != 0 {
		ratio = mulShift(ratio, "0x9aa508b5b7a84e1c677de54f3e99bc9")
	}
	if (absTick & 0x20000) != 0 {
		ratio = mulShift(ratio, "0x5d6af8dedb81196699c329225ee604")
	}
	if (absTick & 0x40000) != 0 {
		ratio = mulShift(ratio, "0x2216e584f5fa1ea926041bedfe98")
	}
	if (absTick & 0x80000) != 0 {
		ratio = mulShift(ratio, "0x48a170391f7dc42444e8fa2")
	}

	if tick > 0 {
		ratio = new(big.Int).Div(MaxUint256, ratio)
	}

	remainder := new(big.Int).Mod(ratio, Q32)
	if remainder.Cmp(ZERO) > 0 {
		return new(big.Int).Add(new(big.Int).Div(ratio, Q32), ONE)
	}
	return new(big.Int).Div(ratio, Q32)
}

func getTickAtSqrtRatio(sqrtRatioX96 *big.Int) int {
	if sqrtRatioX96.Cmp(MIN_SQRT_RATIO) < 0 || sqrtRatioX96.Cmp(MAX_SQRT_RATIO) >= 0 {
		panic("SQRT_RATIO")
	}

	sqrtRatioX128 := new(big.Int).Lsh(sqrtRatioX96, 32)
	msb := mostSignificantBit(sqrtRatioX128)

	var r *big.Int
	if msb >= 128 {
		r = new(big.Int).Rsh(sqrtRatioX128, uint(msb-127))
	} else {
		r = new(big.Int).Lsh(sqrtRatioX128, uint(127-msb))
	}

	log_2 := new(big.Int).Lsh(big.NewInt(int64(msb-128)), 64)

	for i := 0; i < 14; i++ {
		r = new(big.Int).Rsh(new(big.Int).Mul(r, r), 127)
		f := new(big.Int).Rsh(r, 128)
		log_2 = new(big.Int).Or(log_2, new(big.Int).Lsh(f, uint(63-i)))
		r = new(big.Int).Rsh(r, uint(f.Int64()))
	}

	// Fixed line:
	constant, _ := big.NewInt(0).SetString("255738958999603826347141", 10)
	log_sqrt10001 := new(big.Int).Mul(log_2, constant)

	constant1, _ := big.NewInt(0).SetString("3402992956809132418596140100660247210", 10)
	constant2, _ := big.NewInt(0).SetString("291339464771989622907027621153398088495", 10)

	tickLow := int(new(big.Int).Rsh(new(big.Int).Sub(log_sqrt10001, constant1), 128).Int64())
	tickHigh := int(new(big.Int).Rsh(new(big.Int).Add(log_sqrt10001, constant2), 128).Int64())

	if tickLow == tickHigh {
		return tickLow
	}

	if getSqrtRatioAtTick(tickHigh).Cmp(sqrtRatioX96) <= 0 {
		return tickHigh
	}
	return tickLow
}

func getTick(ticks []Tick, index int) Tick {
	tick := ticks[binarySearch(ticks, index)]
	tickIdx, _ := strconv.Atoi(tick.TickIdx)
	if tickIdx != index {
		panic("NOT_CONTAINED")
	}
	return tick
}

func addDelta(x, y *big.Int) *big.Int {
	if y.Cmp(ZERO) < 0 {
		return new(big.Int).Sub(x, new(big.Int).Mul(y, NEGATIVE_ONE))
	}
	return new(big.Int).Add(x, y)
}

func mulDivRoundingUp(a, b, denominator *big.Int) *big.Int {
	product := new(big.Int).Mul(a, b)
	result := new(big.Int).Div(product, denominator)
	if new(big.Int).Mod(product, denominator).Cmp(ZERO) != 0 {
		result = new(big.Int).Add(result, ONE)
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
	return new(big.Int).Div(
		new(big.Int).Div(new(big.Int).Mul(numerator1, numerator2), sqrtRatioBX96),
		sqrtRatioAX96,
	)
}

func getAmount1Delta(sqrtRatioAX96, sqrtRatioBX96, liquidity *big.Int, roundUp bool) *big.Int {
	if sqrtRatioAX96.Cmp(sqrtRatioBX96) > 0 {
		sqrtRatioAX96, sqrtRatioBX96 = sqrtRatioBX96, sqrtRatioAX96
	}

	if roundUp {
		return mulDivRoundingUp(
			liquidity,
			new(big.Int).Sub(sqrtRatioBX96, sqrtRatioAX96),
			Q96,
		)
	}
	return new(big.Int).Div(
		new(big.Int).Mul(liquidity, new(big.Int).Sub(sqrtRatioBX96, sqrtRatioAX96)),
		Q96,
	)
}

func multiplyIn256(x, y *big.Int) *big.Int {
	product := new(big.Int).Mul(x, y)
	return new(big.Int).And(product, MaxUint256)
}

func addIn256(x, y *big.Int) *big.Int {
	sum := new(big.Int).Add(x, y)
	return new(big.Int).And(sum, MaxUint256)
}

func getNextSqrtPriceFromAmount0RoundingUp(sqrtPX96, liquidity, amount *big.Int, add bool) *big.Int {
	if amount.Cmp(ZERO) == 0 {
		return sqrtPX96
	}
	numerator1 := new(big.Int).Lsh(liquidity, 96)

	if add {
		product := multiplyIn256(amount, sqrtPX96)
		if new(big.Int).Div(product, amount).Cmp(sqrtPX96) == 0 {
			denominator := addIn256(numerator1, product)
			if denominator.Cmp(numerator1) >= 0 {
				return mulDivRoundingUp(numerator1, sqrtPX96, denominator)
			}
		}

		return mulDivRoundingUp(
			numerator1,
			ONE,
			new(big.Int).Add(new(big.Int).Div(numerator1, sqrtPX96), amount),
		)
	} else {
		product := multiplyIn256(amount, sqrtPX96)

		if new(big.Int).Div(product, amount).Cmp(sqrtPX96) != 0 {
			panic("Division check failed")
		}
		if numerator1.Cmp(product) <= 0 {
			panic("Numerator1 not greater than product")
		}
		denominator := new(big.Int).Sub(numerator1, product)
		return mulDivRoundingUp(numerator1, sqrtPX96, denominator)
	}
}

func getNextSqrtPriceFromAmount1RoundingDown(sqrtPX96, liquidity, amount *big.Int, add bool) *big.Int {
	if add {
		var quotient *big.Int
		if amount.Cmp(MaxUint160) <= 0 {
			quotient = new(big.Int).Div(new(big.Int).Lsh(amount, 96), liquidity)
		} else {
			quotient = new(big.Int).Div(new(big.Int).Mul(amount, Q96), liquidity)
		}

		return new(big.Int).Add(sqrtPX96, quotient)
	} else {
		quotient := mulDivRoundingUp(amount, Q96, liquidity)

		if sqrtPX96.Cmp(quotient) <= 0 {
			panic("SqrtPX96 not greater than quotient")
		}
		return new(big.Int).Sub(sqrtPX96, quotient)
	}
}

func getNextSqrtPriceFromInput(sqrtPX96, liquidity, amountIn *big.Int, zeroForOne bool) *big.Int {
	if sqrtPX96.Cmp(ZERO) <= 0 {
		panic("SqrtPX96 must be greater than zero")
	}
	if liquidity.Cmp(ZERO) <= 0 {
		panic("Liquidity must be greater than zero")
	}

	if zeroForOne {
		return getNextSqrtPriceFromAmount0RoundingUp(sqrtPX96, liquidity, amountIn, true)
	}
	return getNextSqrtPriceFromAmount1RoundingDown(sqrtPX96, liquidity, amountIn, true)
}

func getNextSqrtPriceFromOutput(sqrtPX96, liquidity, amountOut *big.Int, zeroForOne bool) *big.Int {
	if sqrtPX96.Cmp(ZERO) <= 0 {
		panic("SqrtPX96 must be greater than zero")
	}
	if liquidity.Cmp(ZERO) <= 0 {
		panic("Liquidity must be greater than zero")
	}

	if zeroForOne {
		return getNextSqrtPriceFromAmount1RoundingDown(sqrtPX96, liquidity, amountOut, false)
	}
	return getNextSqrtPriceFromAmount0RoundingUp(sqrtPX96, liquidity, amountOut, false)
}

type SwapStepResult struct {
	SqrtRatioNextX96 *big.Int
	AmountIn         *big.Int
	AmountOut        *big.Int
	FeeAmount        *big.Int
}

func computeSwapStep(sqrtRatioCurrentX96, sqrtRatioTargetX96, liquidity, amountRemaining *big.Int, feePips int64) SwapStepResult {
	returnValues := SwapStepResult{}

	feePipsBig := big.NewInt(feePips)
	zeroForOne := sqrtRatioCurrentX96.Cmp(sqrtRatioTargetX96) >= 0
	exactIn := amountRemaining.Cmp(ZERO) >= 0

	if exactIn {
		amountRemainingLessFee := new(big.Int).Div(
			new(big.Int).Mul(amountRemaining, new(big.Int).Sub(MAX_FEE, feePipsBig)),
			MAX_FEE,
		)
		if zeroForOne {
			returnValues.AmountIn = getAmount0Delta(
				sqrtRatioTargetX96,
				sqrtRatioCurrentX96,
				liquidity,
				true,
			)
		} else {
			returnValues.AmountIn = getAmount1Delta(
				sqrtRatioCurrentX96,
				sqrtRatioTargetX96,
				liquidity,
				true,
			)
		}
		if amountRemainingLessFee.Cmp(returnValues.AmountIn) >= 0 {
			returnValues.SqrtRatioNextX96 = sqrtRatioTargetX96
		} else {
			returnValues.SqrtRatioNextX96 = getNextSqrtPriceFromInput(
				sqrtRatioCurrentX96,
				liquidity,
				amountRemainingLessFee,
				zeroForOne,
			)
		}
	} else {
		if zeroForOne {
			returnValues.AmountOut = getAmount1Delta(
				sqrtRatioTargetX96,
				sqrtRatioCurrentX96,
				liquidity,
				false,
			)
		} else {
			returnValues.AmountOut = getAmount0Delta(
				sqrtRatioCurrentX96,
				sqrtRatioTargetX96,
				liquidity,
				false,
			)
		}
		if new(big.Int).Mul(amountRemaining, NEGATIVE_ONE).Cmp(returnValues.AmountOut) >= 0 {
			returnValues.SqrtRatioNextX96 = sqrtRatioTargetX96
		} else {
			returnValues.SqrtRatioNextX96 = getNextSqrtPriceFromOutput(
				sqrtRatioCurrentX96,
				liquidity,
				new(big.Int).Mul(amountRemaining, NEGATIVE_ONE),
				zeroForOne,
			)
		}
	}

	max := sqrtRatioTargetX96.Cmp(returnValues.SqrtRatioNextX96) == 0

	if zeroForOne {
		if max && exactIn {
			// returnValues.AmountIn already set
		} else {
			returnValues.AmountIn = getAmount0Delta(
				returnValues.SqrtRatioNextX96,
				sqrtRatioCurrentX96,
				liquidity,
				true,
			)
		}
		if max && !exactIn {
			// returnValues.AmountOut already set
		} else {
			returnValues.AmountOut = getAmount1Delta(
				returnValues.SqrtRatioNextX96,
				sqrtRatioCurrentX96,
				liquidity,
				false,
			)
		}
	} else {
		if max && exactIn {
			// returnValues.AmountIn already set
		} else {
			returnValues.AmountIn = getAmount1Delta(
				sqrtRatioCurrentX96,
				returnValues.SqrtRatioNextX96,
				liquidity,
				true,
			)
		}
		if max && !exactIn {
			// returnValues.AmountOut already set
		} else {
			returnValues.AmountOut = getAmount0Delta(
				sqrtRatioCurrentX96,
				returnValues.SqrtRatioNextX96,
				liquidity,
				false,
			)
		}
	}

	if !exactIn && returnValues.AmountOut.Cmp(new(big.Int).Mul(amountRemaining, NEGATIVE_ONE)) > 0 {
		returnValues.AmountOut = new(big.Int).Mul(amountRemaining, NEGATIVE_ONE)
	}

	if exactIn && returnValues.SqrtRatioNextX96.Cmp(sqrtRatioTargetX96) != 0 {
		// we didn't reach the target, so take the remainder of the maximum input as fee
		returnValues.FeeAmount = new(big.Int).Sub(amountRemaining, returnValues.AmountIn)
	} else {
		returnValues.FeeAmount = mulDivRoundingUp(
			returnValues.AmountIn,
			feePipsBig,
			new(big.Int).Sub(MAX_FEE, feePipsBig),
		)
	}

	return returnValues
}

type Pool struct {
	PoolId      string `json:"poolId"`
	DexId       string `json:"dexId"`
	Token0      string `json:"token0"`
	Token1      string `json:"token1"`
	Symbol0     string `json:"symbol0"`
	Symbol1     string `json:"symbol1"`
	Decimals0   int    `json:"decimals0"`
	Decimals1   int    `json:"decimals1"`
	Liquidity   string `json:"liquidity"`
	SqrtPrice   string `json:"sqrtPrice"`
	TickSpacing int    `json:"tickSpacing"`
	TickCurrent int    `json:"tickCurrent"`
	Fee         int    `json:"fee"`
	Tick        string `json:"tick"`
	Ticks       []Tick `json:"ticks"`
	GasPrice    string `json:"gasPrice"`
}

type AmountOutResult struct {
	AmountOut string `json:"amountOut"`
}

func GetRamsesClAmountOut(pool Pool, amountIn string, fromToken, toToken string) (AmountOutResult, error) {
	token0 := strings.ToLower(pool.Token0)
	token1 := strings.ToLower(pool.Token1)

	tokens := []string{token0, token1}
	fromTokenLower := strings.ToLower(fromToken)
	toTokenLower := strings.ToLower(toToken)

	hasFromToken := false
	hasToToken := false
	for _, token := range tokens {
		if token == fromTokenLower {
			hasFromToken = true
		}
		if token == toTokenLower {
			hasToToken = true
		}
	}

	if !hasFromToken || !hasToToken {
		return AmountOutResult{}, errors.New("fromToken or toToken not part of pool")
	}

	fee := big.NewInt(int64(pool.Fee))
	sqrtRatioX96 := new(big.Int)
	sqrtRatioX96.SetString(pool.SqrtPrice, 10)
	tickCurrent := pool.TickCurrent
	liquidity := new(big.Int)
	liquidity.SetString(pool.Liquidity, 10)
	tickSpacing := pool.TickSpacing
	amountSpecified := new(big.Int)
	amountSpecified.SetString(amountIn, 10)

	ticks := pool.Ticks
	sort.Slice(ticks, func(i, j int) bool {
		tickI, _ := strconv.Atoi(ticks[i].TickIdx)
		tickJ, _ := strconv.Atoi(ticks[j].TickIdx)
		return tickI < tickJ
	})

	zeroForOne := fromTokenLower == token0

	var sqrtPriceLimitX96 *big.Int
	if zeroForOne {
		sqrtPriceLimitX96 = new(big.Int).Add(MIN_SQRT_RATIO, ONE)
	} else {
		sqrtPriceLimitX96 = new(big.Int).Sub(MAX_SQRT_RATIO, ONE)
	}

	if zeroForOne {
		if sqrtPriceLimitX96.Cmp(MIN_SQRT_RATIO) <= 0 {
			panic("RATIO_MIN")
		}
		if sqrtPriceLimitX96.Cmp(sqrtRatioX96) >= 0 {
			panic("RATIO_CURRENT")
		}
	} else {
		if sqrtPriceLimitX96.Cmp(MAX_SQRT_RATIO) >= 0 {
			panic("RATIO_MAX")
		}
		if sqrtPriceLimitX96.Cmp(sqrtRatioX96) <= 0 {
			panic("RATIO_CURRENT")
		}
	}

	exactInput := amountSpecified.Cmp(ZERO) >= 0

	state := struct {
		AmountSpecifiedRemaining *big.Int
		AmountCalculated         *big.Int
		SqrtPriceX96            *big.Int
		Tick                    int
		Liquidity               *big.Int
	}{
		AmountSpecifiedRemaining: new(big.Int).Set(amountSpecified),
		AmountCalculated:         new(big.Int).Set(ZERO),
		SqrtPriceX96:            new(big.Int).Set(sqrtRatioX96),
		Tick:                    tickCurrent,
		Liquidity:               new(big.Int).Set(liquidity),
	}

	for state.AmountSpecifiedRemaining.Cmp(ZERO) != 0 && state.SqrtPriceX96.Cmp(sqrtPriceLimitX96) != 0 {
		step := struct {
			SqrtPriceStartX96 *big.Int
			TickNext          int
			Initialized       bool
			SqrtPriceNextX96  *big.Int
			AmountIn          *big.Int
			AmountOut         *big.Int
			FeeAmount         *big.Int
		}{}

		step.SqrtPriceStartX96 = new(big.Int).Set(state.SqrtPriceX96)

		step.TickNext, step.Initialized = nextInitializedTickWithinOneWord(ticks, state.Tick, zeroForOne, tickSpacing)

		if step.TickNext < MIN_TICK {
			step.TickNext = MIN_TICK
		} else if step.TickNext > MAX_TICK {
			step.TickNext = MAX_TICK
		}

		step.SqrtPriceNextX96 = getSqrtRatioAtTick(step.TickNext)

		var targetPrice *big.Int
		if zeroForOne {
			if step.SqrtPriceNextX96.Cmp(sqrtPriceLimitX96) < 0 {
				targetPrice = sqrtPriceLimitX96
			} else {
				targetPrice = step.SqrtPriceNextX96
			}
		} else {
			if step.SqrtPriceNextX96.Cmp(sqrtPriceLimitX96) > 0 {
				targetPrice = sqrtPriceLimitX96
			} else {
				targetPrice = step.SqrtPriceNextX96
			}
		}

		swapResult := computeSwapStep(
			state.SqrtPriceX96,
			targetPrice,
			state.Liquidity,
			state.AmountSpecifiedRemaining,
			fee.Int64(),
		)

		state.SqrtPriceX96 = swapResult.SqrtRatioNextX96
		step.AmountIn = swapResult.AmountIn
		step.AmountOut = swapResult.AmountOut
		step.FeeAmount = swapResult.FeeAmount

		if exactInput {
			state.AmountSpecifiedRemaining = new(big.Int).Sub(
				state.AmountSpecifiedRemaining,
				new(big.Int).Add(step.AmountIn, step.FeeAmount),
			)
			state.AmountCalculated = new(big.Int).Sub(
				state.AmountCalculated,
				step.AmountOut,
			)
		} else {
			state.AmountSpecifiedRemaining = new(big.Int).Add(
				state.AmountSpecifiedRemaining,
				step.AmountOut,
			)
			state.AmountCalculated = new(big.Int).Add(
				state.AmountCalculated,
				new(big.Int).Add(step.AmountIn, step.FeeAmount),
			)
		}

		if state.SqrtPriceX96.Cmp(step.SqrtPriceNextX96) == 0 {
			// if the tick is initialized, run the tick transition
			if step.Initialized {
				tick := getTick(ticks, step.TickNext)
				liquidityNet := new(big.Int)
				liquidityNet.SetString(tick.LiquidityNet, 10)

				// if we're moving leftward, we interpret liquidityNet as the opposite sign
				// safe because liquidityNet cannot be type(int128).min
				if zeroForOne {
					liquidityNet = new(big.Int).Mul(liquidityNet, NEGATIVE_ONE)
				}

				state.Liquidity = addDelta(state.Liquidity, liquidityNet)
			}

			if zeroForOne {
				state.Tick = step.TickNext - 1
			} else {
				state.Tick = step.TickNext
			}
		} else if state.SqrtPriceX96.Cmp(step.SqrtPriceStartX96) != 0 {
			// updated comparison function
			// recompute unless we're on a lower tick boundary (i.e. already transitioned ticks), and haven't moved
			state.Tick = getTickAtSqrtRatio(state.SqrtPriceX96)
		}
	}

	if state.AmountCalculated.Cmp(ZERO) < 0 {
		state.AmountCalculated = new(big.Int).Mul(state.AmountCalculated, NEGATIVE_ONE)
	}

	sqrtPriceX96Before := new(big.Int)
	sqrtPriceX96Before.SetString(pool.SqrtPrice, 10)
	sqrtPriceX96After := new(big.Int).Set(state.SqrtPriceX96)

	before, err := decimal.NewFromString(sqrtPriceX96Before.String())
	if err != nil {
		return AmountOutResult{}, fmt.Errorf("invalid before price: %v", err)
	}
	after, err := decimal.NewFromString(sqrtPriceX96After.String())
	if err != nil {
		return AmountOutResult{}, fmt.Errorf("invalid after price: %v", err)
	}

	priceImpact := after.Sub(before).Div(before).Mul(decimal.NewFromInt(100))

	if priceImpact.LessThan(decimal.NewFromInt(-5)) || priceImpact.GreaterThan(decimal.NewFromInt(5)) {
		return AmountOutResult{}, fmt.Errorf("Skipped %s pool %s due to high price impact: %s%%", pool.DexId, pool.PoolId, priceImpact.StringFixed(4))
	}

	return AmountOutResult{
		AmountOut: state.AmountCalculated.String(),
	}, nil
}

var TICK_SPACINGS = map[int]int{
	100:   1,
	200:   4,
	300:   6,
	400:   8,
	500:   10,
	2500:  50,
	3000:  60,
	10000: 200,
}

// FetchPools function to be compatible with ledger.go pattern
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
				tickSpacing
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
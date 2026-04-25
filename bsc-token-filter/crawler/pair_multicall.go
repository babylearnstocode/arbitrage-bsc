package crawler

import (
	"math/big"

	"github.com/babylearnstocode/bsc-token-filter/eth"
	"github.com/ethereum/go-ethereum/common"
)

type PairData struct {
	Address      string
	Token0       string
	Token1       string
	Reserve0     string
	Reserve1     string
	LiquidityUSD float64
}

// build base token prices dynamically from USDT pairs
func BuildBasePriceMap(pairs []PairData) map[common.Address]float64 {
	priceMap := map[common.Address]float64{}

	usdt := common.HexToAddress("0x55d398326f99059fF775485246999027B3197955")
	priceMap[usdt] = 1

	for _, p := range pairs {
		t0 := common.HexToAddress(p.Token0)
		t1 := common.HexToAddress(p.Token1)

		r0, _ := new(big.Float).SetString(p.Reserve0)
		r1, _ := new(big.Float).SetString(p.Reserve1)

		f0, _ := r0.Float64()
		f1, _ := r1.Float64()

		// if token0 is USDT → price token1
		if t0 == usdt && f1 > 0 {
			priceMap[t1] = f0 / f1
		}

		// if token1 is USDT → price token0
		if t1 == usdt && f0 > 0 {
			priceMap[t0] = f1 / f0
		}
	}

	return priceMap
}

func BuildPairCalls(pairs []common.Address) []eth.Call3 {
	var calls []eth.Call3

	for _, p := range pairs {
		t0, _ := eth.PairABI().Pack("token0")
		t1, _ := eth.PairABI().Pack("token1")
		r, _ := eth.PairABI().Pack("getReserves")

		calls = append(calls,
			eth.Call3{Target: p, AllowFailure: true, CallData: t0},
			eth.Call3{Target: p, AllowFailure: true, CallData: t1},
			eth.Call3{Target: p, AllowFailure: true, CallData: r},
		)
	}

	return calls
}

func DecodePairResults(results []eth.Result, pairs []common.Address) []PairData {
	var raw []PairData

	if len(results) == 0 || len(pairs) == 0 {
		return raw
	}

	for i := 0; i < len(pairs); i++ {
		base := i * 3

		if base+2 >= len(results) {
			break
		}

		if !results[base].Success ||
			!results[base+1].Success ||
			!results[base+2].Success {
			continue
		}

		var token0 common.Address
		var token1 common.Address
		var reserves struct {
			Reserve0           *big.Int
			Reserve1           *big.Int
			BlockTimestampLast uint32
		}

		err0 := eth.PairABI().UnpackIntoInterface(&token0, "token0", results[base].ReturnData)
		err1 := eth.PairABI().UnpackIntoInterface(&token1, "token1", results[base+1].ReturnData)
		err2 := eth.PairABI().UnpackIntoInterface(&reserves, "getReserves", results[base+2].ReturnData)

		if err0 != nil || err1 != nil || err2 != nil {
			continue
		}

		if !(eth.BaseTokens[token0] || eth.BaseTokens[token1]) {
			continue
		}

		raw = append(raw, PairData{
			Address:  pairs[i].Hex(),
			Token0:   token0.Hex(),
			Token1:   token1.Hex(),
			Reserve0: reserves.Reserve0.String(),
			Reserve1: reserves.Reserve1.String(),
		})

	}

	// build dynamic price map
	priceMap := BuildBasePriceMap(raw)

	var out []PairData

	for _, p := range raw {
		t0 := common.HexToAddress(p.Token0)
		t1 := common.HexToAddress(p.Token1)

		r0, _ := new(big.Float).SetString(p.Reserve0)
		r1, _ := new(big.Float).SetString(p.Reserve1)

		f0, _ := r0.Float64()
		f1, _ := r1.Float64()

		price0, ok0 := priceMap[t0]
		price1, ok1 := priceMap[t1]

		var liquidity float64

		if ok0 && ok1 {
			liquidity = f0*price0 + f1*price1
		} else if ok0 {
			liquidity = 2 * f0 * price0
		} else if ok1 {
			liquidity = 2 * f1 * price1
		} else {
			continue
		}

		p.LiquidityUSD = liquidity
		out = append(out, p)
	}

	return out
}

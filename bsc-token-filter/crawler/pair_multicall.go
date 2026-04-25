package crawler

import (
	"fmt"
	"math/big"

	"github.com/babylearnstocode/bsc-token-filter/eth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type PairData struct {
	Address      string
	Token0       string
	Token1       string
	Reserve0     string
	Reserve1     string
	LiquidityUSD float64
}

// BuildBasePriceMap builds token prices from USDT pairs
func BuildBasePriceMap(pairs []PairData) map[common.Address]float64 {
	usdt := common.HexToAddress("0x55d398326f99059fF775485246999027B3197955")
	priceMap := map[common.Address]float64{usdt: 1}

	for _, p := range pairs {
		t0 := common.HexToAddress(p.Token0)
		t1 := common.HexToAddress(p.Token1)

		r0, _ := new(big.Float).SetString(p.Reserve0)
		r1, _ := new(big.Float).SetString(p.Reserve1)
		f0, _ := r0.Float64()
		f1, _ := r1.Float64()

		if t0 == usdt && f1 > 0 {
			priceMap[t1] = f0 / f1
		}
		if t1 == usdt && f0 > 0 {
			priceMap[t0] = f1 / f0
		}
	}

	return priceMap
}

// FetchFilteredPairs runs the two-round preliminary filter:
//
// Round 1: token0 + token1 for all pairs → keep only pairs touching BaseTokens
// Round 2: getReserves only for survivors → compute TVL
//
// This avoids fetching reserves for millions of pairs with unknown tokens.
func FetchFilteredPairs(client *ethclient.Client, pairs []common.Address) ([]PairData, error) {

	// ── round 1: fetch token0 + token1 ───────────────────────────────────────
	fmt.Printf("round 1: fetching token addresses for %d pairs\n", len(pairs))

	calls := make([]eth.Call3, 0, len(pairs)*2)
	t0Data, _ := eth.PairABI().Pack("token0")
	t1Data, _ := eth.PairABI().Pack("token1")

	for _, p := range pairs {
		calls = append(calls,
			eth.Call3{Target: p, AllowFailure: true, CallData: t0Data},
			eth.Call3{Target: p, AllowFailure: true, CallData: t1Data},
		)
	}

	results, err := eth.ExecuteMulticall(client, calls)
	if err != nil {
		return nil, fmt.Errorf("round 1 multicall: %w", err)
	}

	// filter to pairs where token0 or token1 is a base token
	type tokenPair struct {
		addr   common.Address
		token0 common.Address
		token1 common.Address
	}

	var survivors []tokenPair
	for i, pair := range pairs {
		base := i * 2
		if base+1 >= len(results) {
			break
		}
		if !results[base].Success || !results[base+1].Success {
			continue
		}

		var t0, t1 common.Address
		if err := eth.PairABI().UnpackIntoInterface(&t0, "token0", results[base].ReturnData); err != nil {
			continue
		}
		if err := eth.PairABI().UnpackIntoInterface(&t1, "token1", results[base+1].ReturnData); err != nil {
			continue
		}

		if !eth.BaseTokens[t0] && !eth.BaseTokens[t1] {
			continue
		}

		survivors = append(survivors, tokenPair{pair, t0, t1})
	}

	fmt.Printf("round 1 done: %d / %d pairs touch base tokens\n", len(survivors), len(pairs))

	// ── round 2: getReserves only for survivors ───────────────────────────────
	fmt.Printf("round 2: fetching reserves for %d pairs\n", len(survivors))

	rData, _ := eth.PairABI().Pack("getReserves")
	reserveCalls := make([]eth.Call3, len(survivors))
	for i, s := range survivors {
		reserveCalls[i] = eth.Call3{Target: s.addr, AllowFailure: true, CallData: rData}
	}

	reserveResults, err := eth.ExecuteMulticall(client, reserveCalls)
	if err != nil {
		return nil, fmt.Errorf("round 2 multicall: %w", err)
	}

	var raw []PairData
	for i, s := range survivors {
		if i >= len(reserveResults) || !reserveResults[i].Success {
			continue
		}

		var reserves struct {
			Reserve0           *big.Int
			Reserve1           *big.Int
			BlockTimestampLast uint32
		}
		if err := eth.PairABI().UnpackIntoInterface(&reserves, "getReserves", reserveResults[i].ReturnData); err != nil {
			continue
		}

		raw = append(raw, PairData{
			Address:  s.addr.Hex(),
			Token0:   s.token0.Hex(),
			Token1:   s.token1.Hex(),
			Reserve0: reserves.Reserve0.String(),
			Reserve1: reserves.Reserve1.String(),
		})
	}

	// ── compute TVL ───────────────────────────────────────────────────────────
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
		switch {
		case ok0 && ok1:
			liquidity = f0*price0 + f1*price1
		case ok0:
			liquidity = 2 * f0 * price0
		case ok1:
			liquidity = 2 * f1 * price1
		default:
			continue
		}

		p.LiquidityUSD = liquidity
		out = append(out, p)
	}

	fmt.Printf("round 2 done: %d pairs with computable TVL\n", len(out))
	return out, nil
}

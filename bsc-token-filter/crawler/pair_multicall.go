package crawler

import (
	"fmt"
	"math/big"

	"github.com/babylearnstocode/bsc-token-filter/eth"
	"github.com/ethereum/go-ethereum/common"
)

type PairData struct {
	Address  string
	Token0   string
	Token1   string
	Reserve0 string
	Reserve1 string
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
	var out []PairData

	for i := 0; i < len(pairs); i++ {
		base := i * 3

		if !results[base].Success {
			continue
		}

		var t0 []interface{}
		var t1 []interface{}
		var r []interface{}

		eth.PairABI().UnpackIntoInterface(&t0, "token0", results[base].ReturnData)
		eth.PairABI().UnpackIntoInterface(&t1, "token1", results[base+1].ReturnData)
		eth.PairABI().UnpackIntoInterface(&r, "getReserves", results[base+2].ReturnData)

		pairData := PairData{
			Address:  pairs[i].Hex(),
			Token0:   t0[0].(common.Address).Hex(),
			Token1:   t1[0].(common.Address).Hex(),
			Reserve0: r[0].(*big.Int).String(),
			Reserve1: r[1].(*big.Int).String(),
		}
		fmt.Println(pairData)
		out = append(out, pairData)
	}

	return out
}

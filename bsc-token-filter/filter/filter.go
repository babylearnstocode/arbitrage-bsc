package filter

import (
	"context"
	"strings"

	"github.com/babylearnstocode/bsc-token-filter/crawler"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const pauseABIJson = `[
	{"constant":true,"inputs":[],"name":"paused","outputs":[{"type":"bool"}],"stateMutability":"view","type":"function"},
	{"constant":true,"inputs":[],"name":"tradingEnabled","outputs":[{"type":"bool"}],"stateMutability":"view","type":"function"},
	{"constant":true,"inputs":[],"name":"isTrading","outputs":[{"type":"bool"}],"stateMutability":"view","type":"function"},
	{"constant":true,"inputs":[],"name":"swapEnabled","outputs":[{"type":"bool"}],"stateMutability":"view","type":"function"}
]`

var pauseABI abi.ABI

func init() {
	parsed, err := abi.JSON(strings.NewReader(pauseABIJson))
	if err != nil {
		panic(err)
	}
	pauseABI = parsed
}

func FilterHighLiquidity(pairs []crawler.PairData, minUSD float64) []crawler.PairData {
	var out []crawler.PairData

	for _, p := range pairs {
		if p.LiquidityUSD >= minUSD {
			out = append(out, p)
		}
	}

	return out
}

func callBool(ctx context.Context, client *ethclient.Client, token common.Address, method string) (bool, bool) {
	data, err := pauseABI.Pack(method)
	if err != nil {
		return false, false
	}
	res, err := client.CallContract(ctx, ethereum.CallMsg{To: &token, Data: data}, nil)
	if err != nil || len(res) == 0 {
		return false, false
	}
	out, err := pauseABI.Unpack(method, res)
	if err != nil || len(out) == 0 {
		return false, false
	}
	val, ok := out[0].(bool)
	return val, ok
}

func IsPaused(ctx context.Context, client *ethclient.Client, token common.Address) bool {
	if val, ok := callBool(ctx, client, token, "paused"); ok && val {
		return true
	}
	if val, ok := callBool(ctx, client, token, "tradingEnabled"); ok && !val {
		return true
	}
	if val, ok := callBool(ctx, client, token, "isTrading"); ok && !val {
		return true
	}
	if val, ok := callBool(ctx, client, token, "swapEnabled"); ok && !val {
		return true
	}
	return false
}

func FilterPausedPairs(ctx context.Context, client *ethclient.Client, pairs []crawler.PairData) []crawler.PairData {
	var out []crawler.PairData
	for _, p := range pairs {
		t0 := common.HexToAddress(p.Token0)
		t1 := common.HexToAddress(p.Token1)
		if IsPaused(ctx, client, t0) || IsPaused(ctx, client, t1) {
			continue
		}
		out = append(out, p)
	}
	return out
}

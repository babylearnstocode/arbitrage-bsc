package filter

import (
	"context"
	"strings"

	"github.com/babylearnstocode/bsc-token-filter/crawler"
	"github.com/babylearnstocode/bsc-token-filter/eth"
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

func FilterPausedTokens(ctx context.Context, client *ethclient.Client, tokens []common.Address) []common.Address {
	var out []common.Address
	for _, t := range tokens {

		if IsPaused(ctx, client, t) {
			continue
		}
		out = append(out, t)
	}
	return out
}

func ExtractUniqueTokens(pairs []crawler.PairData) (base []common.Address, unknown []common.Address) {
	seen := make(map[common.Address]bool)

	for _, p := range pairs {
		for _, raw := range []string{p.Token0, p.Token1} {
			addr := common.HexToAddress(raw)
			if seen[addr] {
				continue
			}
			seen[addr] = true

			if eth.BaseTokens[addr] {
				base = append(base, addr)
			} else {
				unknown = append(unknown, addr)
			}
		}
	}
	return base, unknown
}

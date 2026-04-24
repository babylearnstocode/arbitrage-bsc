package main

import (
	"fmt"

	"github.com/babylearnstocode/bsc-token-filter/config"
	"github.com/babylearnstocode/bsc-token-filter/crawler"
	"github.com/babylearnstocode/bsc-token-filter/eth"
	"github.com/babylearnstocode/bsc-token-filter/filter"
)

func main() {
	cfg := config.LoadConfig()

	localClient := eth.NewClient(cfg.IpcPath)

	defer localClient.Close()
	// Stage 0 load local cache

	// Stage 1 fetch all pairs from index 0
	fetchRes := crawler.Start(localClient, cfg)

	// Stage 2 multicall info,
	calls := crawler.BuildPairCalls(fetchRes)

	results, _ := eth.ExecuteMulticall(localClient, calls)

	decoded := crawler.DecodePairResults(results, fetchRes)

	// 1.1 filter liquidity >= 500k USD
	// usage
	const MinLiquidity = 5e23
	filtered := filter.FilterHighLiquidity(decoded, MinLiquidity)
	fmt.Println(filtered)
	// Stage 2 activity filter: >= 200 swap per day

	// Stage 3 safety filter

	// 3.1 eth_call simulate buy_sell, tax > 0.3% no pause()

	// stage 4 reachability filter
	// BFS from WBNB/BNB <= 2 hops

	// Rank top 50 pairs by EL

}

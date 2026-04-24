package main

import (
	"github.com/babylearnstocode/bsc-token-filter/config"
	"github.com/babylearnstocode/bsc-token-filter/crawler"
	"github.com/babylearnstocode/bsc-token-filter/eth"
)

func main() {
	cfg := config.LoadConfig()

	localClient := eth.NewClient(cfg.IpcPath)

	defer localClient.Close()
	// Stage 0 load local cache

	// Stage 1 fetch all pairs resume from cache index or 0
	fetchRes := crawler.Start(localClient, cfg)

	// Stage 2 multicall info,
	calls := crawler.BuildPairCalls(fetchRes)

	results, _ := eth.ExecuteMulticall(localClient, calls)

	crawler.DecodePairResults(results, fetchRes)

	// 1.1 calculate effective liquidity >= 500k USD

	// 1.2 save to cache

	// Stage 2 activity filter: >= 200 swap per day

	// Stage 3 safety filter

	// 3.1 eth_call simulate buy_sell, tax > 0.3% no pause()

	// stage 4 reachability filter
	// BFS from WBNB/BNB <= 2 hops

	// Rank top 50 pairs by EL

}

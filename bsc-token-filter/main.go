package main

import (
	"context"
	"fmt"

	"github.com/babylearnstocode/bsc-token-filter/config"
	"github.com/babylearnstocode/bsc-token-filter/crawler"
	"github.com/babylearnstocode/bsc-token-filter/eth"
	"github.com/babylearnstocode/bsc-token-filter/filter"
)

func main() {
	ctx := context.Background()
	cfg := config.LoadConfig()

	localClient := eth.NewClient(cfg.IpcPath)
	//infuraClient := eth.NewClient(cfg.ArchiveUrl)

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
	fmt.Println(len(filtered))
	// Stage 2 activity filter: >= 200 swap per day -- drop this

	// Stage 3 safety filter

	// 3.1 paused() check
	filtered = filter.FilterPausedPairs(ctx, localClient, filtered)
	fmt.Println(len(filtered))
	// transfer tax simulation, tax > 0.3%

	// sell simulation

	// max tx simulation

	// Rank top 50 pairs by EL

}

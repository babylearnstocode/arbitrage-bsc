package main

import (
	"context"
	"fmt"
	"log"

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
	pairs, err := crawler.FetchFilteredPairs(localClient, fetchRes)

	if err != nil {
		log.Fatal(err)
	}
	// 1.1 filter liquidity >= 500k USD
	// usage
	const MinLiquidity = 5e22
	filtered := filter.FilterHighLiquidity(pairs, MinLiquidity)
	fmt.Printf("filtered: %v", len(filtered))
	// Stage 2 activity filter: >= 200 swap per day -- drop this

	// Stage 3 safety filter

	// 3.1 paused() check
	_, unknownTokens := filter.ExtractUniqueTokens(filtered)

	filter.FilterPausedTokens(ctx, localClient, unknownTokens)

	// transfer tax simulation, tax > 0.3%, sell simulation
	safeUnknown := filter.FilterHoneypotTokens(ctx, unknownTokens, 0.3)
	fmt.Printf("safeUnknown: %v", len(safeUnknown))
	// Rank top 50 tokens by Vol
	res := filter.TopTokensByVolume(ctx, safeUnknown, 50)
	fmt.Println(len(res))

	// save

	filter.SaveTokens("data/data.json", res)
}

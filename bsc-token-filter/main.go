package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/babylearnstocode/bsc-token-filter/config"
	"github.com/babylearnstocode/bsc-token-filter/crawler"
	"github.com/babylearnstocode/bsc-token-filter/eth"
	"github.com/babylearnstocode/bsc-token-filter/filter"
	"github.com/babylearnstocode/bsc-token-filter/pool"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	tokenFile    = "data/data.json"
	MinLiquidity = 5e22
)

func main() {
	ctx := context.Background()
	cfg := config.LoadConfig()

	localClient := eth.NewClient(cfg.IpcPath)
	defer localClient.Close()

	// if token list already exists, skip crawling and go straight to pool discovery
	if _, err := os.Stat(tokenFile); err == nil {
		fmt.Println("token file found — loading and finding pools")
		runFindPools(ctx, localClient)
		return
	}

	fmt.Println("no token file — running full crawler pipeline")
	runCrawler(ctx, cfg, localClient)
}

func runCrawler(ctx context.Context, cfg *config.Config, localClient *ethclient.Client) {
	// Stage 1 — fetch all pairs
	fetchRes := crawler.Start(localClient, cfg)

	// Stage 2 — multicall info, filter by liquidity
	pairs, err := crawler.FetchFilteredPairs(localClient, fetchRes)
	if err != nil {
		log.Fatal(err)
	}

	filtered := filter.FilterHighLiquidity(pairs, MinLiquidity)
	fmt.Printf("after liquidity filter: %d pairs\n", len(filtered))

	// Stage 3 — safety filter
	_, unknownTokens := filter.ExtractUniqueTokens(filtered)

	filter.FilterPausedTokens(ctx, localClient, unknownTokens)

	safeUnknown := filter.FilterHoneypotTokens(ctx, unknownTokens, 0.3)
	fmt.Printf("safe unknown tokens: %d\n", len(safeUnknown))

	// Stage 4 — rank top 50 by volume and save
	res := filter.TopTokensByVolume(ctx, safeUnknown, 50)
	fmt.Printf("top tokens: %d\n", len(res))

	if err := filter.SaveTokens(tokenFile, res); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("saved to %s\n", tokenFile)
}

func runFindPools(ctx context.Context, localClient *ethclient.Client) {
	tokens, err := filter.LoadTokens(tokenFile)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("loaded %d tokens from %s\n", len(tokens), tokenFile)

	v2, v3, err := pool.FindAllPools(localClient, tokens)
	if err != nil {
		log.Fatal(err)
	}
	if err := filter.SavePools(v2, v3); err != nil {
		log.Fatal(err)
	}
}

package crawler

import (
	"fmt"
	"log"

	"github.com/babylearnstocode/bsc-token-filter/config"
	"github.com/babylearnstocode/bsc-token-filter/eth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func Start(client *ethclient.Client, cfg *config.Config) []common.Address {
	totalBig, err := eth.GetTotalPairs(client, cfg)
	if err != nil {
		log.Fatal(err)
	}
	total := totalBig.Int64()

	fmt.Println("Total pairs:", total)

	var pairs []common.Address

	for i := int64(0); i < 10; i++ {
		if i%1000 == 0 {
			fmt.Printf("Fetching %d / %d\n", i, total)
		}

		addr, err := eth.GetPair(client, cfg, i)
		if err != nil {
			continue
		}

		fmt.Printf("Index %v - addr: %v\n", i, addr.Hex())
		pairs = append(pairs, addr)
	}

	return pairs
}

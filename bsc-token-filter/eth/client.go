package eth

import (
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

func NewClient(clientString string) *ethclient.Client {
	client, err := ethclient.Dial(clientString)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

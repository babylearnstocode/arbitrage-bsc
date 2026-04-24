// eth/pair.go
package eth

import (
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

const pairABIJson = `[
  {
    "name":"token0",
    "outputs":[{"type":"address"}],
    "stateMutability":"view",
    "type":"function"
  },
  {
    "name":"token1",
    "outputs":[{"type":"address"}],
    "stateMutability":"view",
    "type":"function"
  },
  {
    "name":"getReserves",
    "outputs":[
      {"type":"uint112","name":"reserve0"},
      {"type":"uint112","name":"reserve1"},
      {"type":"uint32","name":"blockTimestampLast"}
    ],
    "stateMutability":"view",
    "type":"function"
  }
]`

var pairABI abi.ABI

func init() {
	var err error
	pairABI, err = abi.JSON(strings.NewReader(pairABIJson))
	if err != nil {
		log.Fatal("parse pair ABI:", err)
	}
}

func PairABI() abi.ABI {
	return pairABI
}

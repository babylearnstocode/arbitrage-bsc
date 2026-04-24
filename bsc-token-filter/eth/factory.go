package eth

import (
	"context"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/babylearnstocode/bsc-token-filter/config"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ── ABIs ─────────────────────────────────────────────────────────────────────

const factoryABIJson = `[
  {"constant":true,"inputs":[],"name":"allPairsLength","outputs":[{"type":"uint256"}],"stateMutability":"view","type":"function"},
  {"constant":true,"inputs":[{"type":"uint256"}],"name":"allPairs","outputs":[{"type":"address"}],"stateMutability":"view","type":"function"}
]`

var (
	factoryABI abi.ABI
)

func init() {
	var err error

	factoryABI, err = abi.JSON(strings.NewReader(factoryABIJson))
	if err != nil {
		log.Fatal("parse factoryABI:", err)
	}

}

func GetTotalPairs(client *ethclient.Client, cfg *config.Config) (*big.Int, error) {
	addr := common.HexToAddress(cfg.V2FactoryAddress)

	data, err := factoryABI.Pack("allPairsLength")
	if err != nil {
		return nil, err
	}

	res, err := client.CallContract(context.Background(), ethereum.CallMsg{
		To:   &addr,
		Data: data,
	}, nil)
	if err != nil {
		return nil, err
	}

	out, err := factoryABI.Unpack("allPairsLength", res)
	if err != nil {
		return nil, err
	}

	return out[0].(*big.Int), nil
}

func GetPair(client *ethclient.Client, cfg *config.Config, i int64) (common.Address, error) {
	addr := common.HexToAddress(cfg.V2FactoryAddress)

	data, err := factoryABI.Pack("allPairs", big.NewInt(i))
	if err != nil {
		return common.Address{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := client.CallContract(ctx, ethereum.CallMsg{
		To:   &addr,
		Data: data,
	}, nil)
	if err != nil {
		return common.Address{}, err
	}

	out, err := factoryABI.Unpack("allPairs", res)
	if err != nil {
		return common.Address{}, err
	}

	return out[0].(common.Address), nil
}

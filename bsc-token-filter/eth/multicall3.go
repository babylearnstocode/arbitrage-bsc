package eth

import (
	"context"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const Multicall3Address = "0xcA11bde05977b3631167028862bE2a173976CA11"

const multicallABIJson = `
[
  {
    "inputs": [
      {
        "components": [
          {"internalType":"address","name":"target","type":"address"},
          {"internalType":"bool","name":"allowFailure","type":"bool"},
          {"internalType":"bytes","name":"callData","type":"bytes"}
        ],
        "internalType":"struct Multicall3.Call3[]",
        "name":"calls",
        "type":"tuple[]"
      }
    ],
    "name":"aggregate3",
    "outputs":[
      {
        "components":[
          {"internalType":"bool","name":"success","type":"bool"},
          {"internalType":"bytes","name":"returnData","type":"bytes"}
        ],
        "internalType":"struct Multicall3.Result[]",
        "name":"returnData",
        "type":"tuple[]"
      }
    ],
    "stateMutability":"payable",
    "type":"function"
  }
]
`

var multicallABI abi.ABI

func InitMultiCall3ABI() {
	var err error
	multicallABI, err = abi.JSON(strings.NewReader(multicallABIJson))
	if err != nil {
		log.Fatal(err)
	}
}

type Call3 struct {
	Target       common.Address
	AllowFailure bool
	CallData     []byte
}

type Result struct {
	Success    bool
	ReturnData []byte
}

func ExecuteMulticall(
	client *ethclient.Client,
	calls []Call3,
) ([]Result, error) {

	addr := common.HexToAddress(Multicall3Address)

	data, err := multicallABI.Pack("aggregate3", calls)
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

	var out []Result
	err = multicallABI.UnpackIntoInterface(&out, "aggregate3", res)
	if err != nil {
		return nil, err
	}

	return out, nil
}

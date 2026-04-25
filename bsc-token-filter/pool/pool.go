package pool

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/babylearnstocode/bsc-token-filter/eth"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ── DEX registry ──────────────────────────────────────────────────────────────

type DEX struct {
	Name    string
	Factory common.Address
	Fee     uint32 // fixed fee in bps×100: 9975 = 0.25%, 9980 = 0.20%, 9970 = 0.30%
}

var V2DEXes = []DEX{
	{Name: "pancakeswap_v2", Factory: common.HexToAddress("0xca143ce32fe78f1f7019d7d551a6402fc5350c73"), Fee: 9975},
	{Name: "biswap", Factory: common.HexToAddress("0x858E3312ed3A876947EA49d572A7C42DE08af7EE"), Fee: 9990},
	{Name: "apeswap", Factory: common.HexToAddress("0x0841BD0B734E4F5853f0dD8d7Ea041c241fb0Da6"), Fee: 9980},
	{Name: "squadswap_v2", Factory: common.HexToAddress("0x918Adf1f2C03b244823Cd712E010B6e3CD653DbA"), Fee: 9975},
	{Name: "fstswap", Factory: common.HexToAddress("0x9A272d734c5a0d7d84E0a892e891a553e8066dce"), Fee: 9970},
}

// V3 fee tiers as *big.Int — uint24 in ABI maps to *big.Int in go-ethereum
var V3FeeTiers = []*big.Int{
	big.NewInt(100),
	big.NewInt(500),
	big.NewInt(2500),
	big.NewInt(10000),
}

var V3DEXes = []DEX{
	{Name: "pancakeswap_v3", Factory: common.HexToAddress("0x0BFbCF9fa4f9C56B0F40a671Ad40E0805A091865")},
	{Name: "uniswap_v3", Factory: common.HexToAddress("0xdB1d10011AD0Ff90774D0C6Bb92e5C5c8b4461F7")},
	{Name: "squadswap_v3", Factory: common.HexToAddress("0x10d8612D9D8269e322AB551C18a307cB4D6BC07B")},
}

// ── Pool ──────────────────────────────────────────────────────────────────────

type Pool struct {
	Address  common.Address
	Token0   common.Address
	Token1   common.Address
	DEX      string
	Fee      uint32 // 0 for V2
	Protocol string // "v2" or "v3"
}

// ── ABIs ──────────────────────────────────────────────────────────────────────

const getPairABIJson = `[{
	"name": "getPair",
	"type": "function",
	"stateMutability": "view",
	"inputs": [
		{"name": "tokenA", "type": "address"},
		{"name": "tokenB", "type": "address"}
	],
	"outputs": [{"name": "pair", "type": "address"}]
}]`

const getPoolABIJson = `[{
	"name": "getPool",
	"type": "function",
	"stateMutability": "view",
	"inputs": [
		{"name": "tokenA", "type": "address"},
		{"name": "tokenB", "type": "address"},
		{"name": "fee",    "type": "uint24"}
	],
	"outputs": [{"name": "pool", "type": "address"}]
}]`

// ── Helpers ───────────────────────────────────────────────────────────────────

type combo struct {
	TokenA common.Address
	TokenB common.Address
	DEX    string
	Fee    *big.Int
}

var zero = common.Address{}

func activeDEXes(list []DEX) []DEX {
	var out []DEX
	for _, d := range list {
		if d.Factory == zero {
			fmt.Printf("skipping %s — factory address not set\n", d.Name)
			continue
		}
		out = append(out, d)
	}
	return out
}

func decodeAddressResult(data []byte) (common.Address, bool) {
	if len(data) < 32 {
		return zero, false
	}
	addr := common.BytesToAddress(data[12:32])
	return addr, addr != zero
}

// ── V2 ────────────────────────────────────────────────────────────────────────

func FindV2Pools(client *ethclient.Client, tokens []common.Address) ([]Pool, error) {
	pairABI, err := abi.JSON(strings.NewReader(getPairABIJson))
	if err != nil {
		return nil, fmt.Errorf("parse getPair abi: %w", err)
	}

	dexes := activeDEXes(V2DEXes)
	if len(dexes) == 0 {
		return nil, nil
	}

	var combos []combo
	var calls []eth.Call3

	for _, dex := range dexes {
		for i := 0; i < len(tokens); i++ {
			for j := i + 1; j < len(tokens); j++ {
				calldata, err := pairABI.Pack("getPair", tokens[i], tokens[j])
				if err != nil {
					return nil, err
				}
				combos = append(combos, combo{TokenA: tokens[i], TokenB: tokens[j], DEX: dex.Name, Fee: big.NewInt(int64(dex.Fee))})
				calls = append(calls, eth.Call3{Target: dex.Factory, AllowFailure: true, CallData: calldata})
			}
		}
	}

	n := len(tokens)
	fmt.Printf("v2: %d DEXes × %d combinations = %d calls\n", len(dexes), n*(n-1)/2, len(calls))

	results, err := eth.ExecuteMulticall(client, calls)
	if err != nil {
		return nil, err
	}

	var pools []Pool
	for i, res := range results {
		if !res.Success {
			continue
		}
		addr, ok := decodeAddressResult(res.ReturnData)
		if !ok {
			continue
		}
		pools = append(pools, Pool{
			Address:  addr,
			Token0:   combos[i].TokenA,
			Token1:   combos[i].TokenB,
			DEX:      combos[i].DEX,
			Fee:      uint32(combos[i].Fee.Uint64()),
			Protocol: "v2",
		})
	}

	fmt.Printf("v2: found %d pools\n", len(pools))
	return pools, nil
}

// ── V3 ────────────────────────────────────────────────────────────────────────

func FindV3Pools(client *ethclient.Client, tokens []common.Address) ([]Pool, error) {
	poolABI, err := abi.JSON(strings.NewReader(getPoolABIJson))
	if err != nil {
		return nil, fmt.Errorf("parse getPool abi: %w", err)
	}

	dexes := activeDEXes(V3DEXes)
	if len(dexes) == 0 {
		return nil, nil
	}

	var combos []combo
	var calls []eth.Call3

	for _, dex := range dexes {
		for i := 0; i < len(tokens); i++ {
			for j := i + 1; j < len(tokens); j++ {
				for _, fee := range V3FeeTiers {
					calldata, err := poolABI.Pack("getPool", tokens[i], tokens[j], fee)
					if err != nil {
						return nil, err
					}
					combos = append(combos, combo{TokenA: tokens[i], TokenB: tokens[j], DEX: dex.Name, Fee: fee})
					calls = append(calls, eth.Call3{Target: dex.Factory, AllowFailure: true, CallData: calldata})
				}
			}
		}
	}

	n := len(tokens)
	fmt.Printf("v3: %d DEXes × %d combinations × %d fee tiers = %d calls\n",
		len(dexes), n*(n-1)/2, len(V3FeeTiers), len(calls))

	results, err := eth.ExecuteMulticall(client, calls)
	if err != nil {
		return nil, err
	}

	var pools []Pool
	for i, res := range results {
		if !res.Success {
			continue
		}
		addr, ok := decodeAddressResult(res.ReturnData)
		if !ok {
			continue
		}
		pools = append(pools, Pool{
			Address:  addr,
			Token0:   combos[i].TokenA,
			Token1:   combos[i].TokenB,
			DEX:      combos[i].DEX,
			Fee:      uint32(combos[i].Fee.Uint64()),
			Protocol: "v3",
		})
	}

	fmt.Printf("v3: found %d pools\n", len(pools))
	return pools, nil
}

// ── Combined ──────────────────────────────────────────────────────────────────

func FindAllPools(client *ethclient.Client, tokens []common.Address) ([]Pool, []Pool, error) {
	v2, err := FindV2Pools(client, tokens)
	if err != nil {
		return nil, nil, fmt.Errorf("v2: %w", err)
	}

	v3, err := FindV3Pools(client, tokens)
	if err != nil {
		fmt.Println(err)
		return v2, nil, fmt.Errorf("v3: %w", err)
	}

	all := append(v2, v3...)
	fmt.Printf("total: %d pools (%d v2, %d v3)\n", len(all), len(v2), len(v3))
	return v2, v3, nil
}

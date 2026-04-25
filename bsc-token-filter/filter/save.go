package filter

import (
	"encoding/json"
	"os"

	"github.com/babylearnstocode/bsc-token-filter/eth"
	"github.com/babylearnstocode/bsc-token-filter/pool"
	"github.com/ethereum/go-ethereum/common"
)

type poolJSON struct {
	Address  string `json:"address"`
	Token0   string `json:"token0"`
	Token1   string `json:"token1"`
	DEX      string `json:"dex"`
	Fee      uint32 `json:"fee,omitempty"`
	Protocol string `json:"protocol"`
}

func SaveTokens(path string, topTokens []TokenVolume) error {
	seen := make(map[common.Address]bool)
	var addrs []string

	for addr := range eth.BaseTokens {
		seen[addr] = true
		addrs = append(addrs, addr.Hex())
	}

	for _, t := range topTokens {
		if !seen[t.Address] {
			addrs = append(addrs, t.Address.Hex())
		}
	}

	data, err := json.MarshalIndent(addrs, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll("data", 0755); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func LoadTokens(path string) ([]common.Address, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var addrs []string
	if err := json.Unmarshal(data, &addrs); err != nil {
		return nil, err
	}
	tokens := make([]common.Address, len(addrs))
	for i, a := range addrs {
		tokens[i] = common.HexToAddress(a)
	}
	return tokens, nil
}

func savePools(path string, pools []pool.Pool) error {
	entries := make([]poolJSON, len(pools))
	for i, p := range pools {
		entries[i] = poolJSON{
			Address:  p.Address.Hex(),
			Token0:   p.Token0.Hex(),
			Token1:   p.Token1.Hex(),
			DEX:      p.DEX,
			Fee:      p.Fee,
			Protocol: p.Protocol,
		}
	}

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll("data", 0755); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func SavePools(v2 []pool.Pool, v3 []pool.Pool) error {
	if err := savePools("data/v2.json", v2); err != nil {
		return err
	}
	return savePools("data/v3.json", v3)
}

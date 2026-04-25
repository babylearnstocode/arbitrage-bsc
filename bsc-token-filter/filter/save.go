package filter

import (
	"encoding/json"
	"os"

	"github.com/babylearnstocode/bsc-token-filter/eth"
	"github.com/ethereum/go-ethereum/common"
)

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

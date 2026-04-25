package filter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

const dexscreenerTokenPairsURL = "https://api.dexscreener.com/token-pairs/v1/bsc"

type TokenVolume struct {
	Address   common.Address
	VolumeH24 float64
}

type dexPair struct {
	BaseToken struct {
		Address string `json:"address"`
	} `json:"baseToken"`
	QuoteToken struct {
		Address string `json:"address"`
	} `json:"quoteToken"`
	Volume struct {
		H24 float64 `json:"h24"`
	} `json:"volume"`
}

func fetchTokenVolume(ctx context.Context, token common.Address) (float64, error) {
	url := fmt.Sprintf("%s/%s", dexscreenerTokenPairsURL, token.Hex())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return 0, nil
	}
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("dexscreener returned %d for %s", resp.StatusCode, token.Hex())
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var pairs []dexPair
	if err := json.Unmarshal(body, &pairs); err != nil {
		return 0, fmt.Errorf("unmarshal: %w — body: %.200s", err, body)
	}

	tokenLow := strings.ToLower(token.Hex())
	var total float64
	for _, pair := range pairs {
		if strings.ToLower(pair.BaseToken.Address) == tokenLow ||
			strings.ToLower(pair.QuoteToken.Address) == tokenLow {
			total += pair.Volume.H24
		}
	}
	return total, nil
}

func TopTokensByVolume(ctx context.Context, tokens []common.Address, top int) []TokenVolume {
	var results []TokenVolume

	for _, token := range tokens {
		vol, err := fetchTokenVolume(ctx, token)
		if err != nil {
			fmt.Printf("volume error %s: %v\n", token.Hex(), err)
		}
		results = append(results, TokenVolume{Address: token, VolumeH24: vol})
		fmt.Printf("token %s - vol24h: $%.0f\n", token.Hex(), vol)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].VolumeH24 > results[j].VolumeH24
	})

	if top > len(results) {
		top = len(results)
	}
	return results[:top]
}

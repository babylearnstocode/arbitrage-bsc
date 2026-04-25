package filter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

const honeypotAPIURL = "https://api.honeypot.is/v2/IsHoneypot"

var httpClient = &http.Client{Timeout: 10 * time.Second}

type HoneypotResult struct {
	IsHoneypot bool
	BuyTax     float64
	SellTax    float64
}

type honeypotResponse struct {
	HoneypotResult *struct {
		IsHoneypot bool `json:"isHoneypot"`
	} `json:"honeypotResult"`
	SimulationResult *struct {
		BuyTax  float64 `json:"buyTax"`
		SellTax float64 `json:"sellTax"`
	} `json:"simulationResult"`
}

func CheckHoneypot(ctx context.Context, token common.Address) (HoneypotResult, error) {
	url := fmt.Sprintf("%s?address=%s&chainID=56", honeypotAPIURL, token.Hex())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return HoneypotResult{}, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return HoneypotResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return HoneypotResult{}, fmt.Errorf("honeypot.is returned %d for %s", resp.StatusCode, token.Hex())
	}

	var data honeypotResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return HoneypotResult{}, err
	}

	// API returns null honeypotResult when simulation could not run
	if data.HoneypotResult == nil {
		return HoneypotResult{IsHoneypot: false}, nil
	}

	result := HoneypotResult{
		IsHoneypot: data.HoneypotResult.IsHoneypot,
	}
	if data.SimulationResult != nil {
		result.BuyTax = data.SimulationResult.BuyTax
		result.SellTax = data.SimulationResult.SellTax
	}
	return result, nil
}

func FilterHoneypotTokens(ctx context.Context, tokens []common.Address, maxTax float64) []common.Address {
	var safe []common.Address
	for _, token := range tokens {
		result, err := CheckHoneypot(ctx, token)
		if err != nil {
			fmt.Printf("honeypot check error %s: %v\n", token.Hex(), err)
			continue
		}
		if result.IsHoneypot || result.BuyTax > maxTax || result.SellTax > maxTax {
			fmt.Printf("unsafe token %s - honeypot:%v buyTax:%.2f sellTax:%.2f\n",
				token.Hex(), result.IsHoneypot, result.BuyTax, result.SellTax)
			continue
		}
		safe = append(safe, token)
	}
	return safe
}

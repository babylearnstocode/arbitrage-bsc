package filter

import (
	"github.com/babylearnstocode/bsc-token-filter/crawler"
)

func FilterHighLiquidity(pairs []crawler.PairData, minUSD float64) []crawler.PairData {
	var out []crawler.PairData

	for _, p := range pairs {
		if p.LiquidityUSD >= minUSD {
			out = append(out, p)
		}
	}

	return out
}

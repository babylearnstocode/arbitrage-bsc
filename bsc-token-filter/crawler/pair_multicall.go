package crawler

import (
	"math/big"
	"strconv"

	"github.com/babylearnstocode/bsc-token-filter/eth"
	"github.com/ethereum/go-ethereum/common"
)

type PairData struct {
	Address      string
	Token0       string
	Token1       string
	Reserve0     string
	Reserve1     string
	LiquidityUSD float64
}

var BaseTokens = map[common.Address]bool{
	common.HexToAddress("0xbb4cdb9cbd36b01bd1cbaebf2de08d9173bc095c"): true, //WBNB
	common.HexToAddress("0x55d398326f99059fF775485246999027B3197955"): true, //USDT
	common.HexToAddress("0x8AC76a51cc950d9822D68b83fE1Ad97B32Cd580d"): true, //USDC
	common.HexToAddress("0xe9e7CEA3DedcA5984780Bafc599bD69ADd087D56"): true, //BUSD
	common.HexToAddress("0x0555E30da8f98308EdB960aa94C0Db47230d2B9c"): true, //WBTC
	common.HexToAddress("0x7130d2A12B9BCbFAe4f2634d864A1Ee1Ce3Ead9c"): true, //BTCB
	common.HexToAddress("0x2170Ed0880ac9A755fd29B2688956BD959F933F8"): true, //ETH
	common.HexToAddress("0x40af3827F39D0EAcBF4A168f8D4ee67c121D11c9"): true, //TUSD
	common.HexToAddress("0x1AF3F329e8BE154074D8769D1FFa4eE058B1DBc3"): true, //DAI
	common.HexToAddress("0xc5f0f7b66764F6ec8C8Dff7BA683102295E16409"): true, //FDUSD
	common.HexToAddress("0x5d3a1Ff2b6BAb83b63cd9AD0787074081a52ef34"): true, //USDE

	// Dex tokens
	common.HexToAddress("0x0E09FaBB73Bd3Ade0a17ECC321fD13a19e81cE82"): true, //CAKE
	common.HexToAddress("0xcF6BB5389c92Bdda8a3747Ddb454cB7a64626C63"): true, //XVS
	common.HexToAddress("0xBf5140A22578168FD562DCcF235E5D43A02ce9B1"): true, //UNI

	//Wrapped majors
	common.HexToAddress("0xbA2aE424d960c26247Dd6c32edC70B295c744C43"): true, //DOGE
	common.HexToAddress("0x3EE2200Efb3400fAbB9AacF31297cBdD1d435D47"): true, //ADA
	common.HexToAddress("0x1D2F0da169ceB9fC7B3144628dB156f3F6c60dBE"): true, //XRP
	common.HexToAddress("0xCC42724C6683B7E57334c4E856f4c9965ED682bD"): true, //MATIC
	common.HexToAddress("0x1CE0c2827e2eF14D5C4f29a091d735A204794041"): true, //AVAX
	common.HexToAddress("0x7083609fCE4d1d8Dc0C979AAb8c869Ea2C873402"): true, //DOT
	common.HexToAddress("0xF8A0BF9cF54Bb92F17374d9e9A321E6a111a51bD"): true, //LINK
	common.HexToAddress("0x4338665CBB7B2485A8855A139b75D5e34AB0DB94"): true, //LTC
	common.HexToAddress("0x570A5D26f7765Ecb712C0924E4De545B89fD43dF"): true, //SOL
	common.HexToAddress("0x76A797A59Ba2C17726896976B7B3747BfD1d220f"): true, //TON
	common.HexToAddress("0x04C0599Ae5A44757c0af6F9eC3b93da8976c150A"): true, //weETH
	common.HexToAddress("0x211Cc4DD073734dA055fbF44a2b4667d5E5fE5d2"): true, //sUSDe
	common.HexToAddress("0xb3b02E4A9Fb2bD28CC2ff97B0aB3F6B3Ec1eE9D2"): true, //USDf
	common.HexToAddress("0x45e51bc23d592eb2dba86da3985299f7895d66ba"): true, //USDD
	common.HexToAddress("0xce24439f2d9c6a2289f741120fe202248b666666"): true, //U
	common.HexToAddress("0x0eb3a705fc54725037cc9e008bdede697f62f335"): true, //ATOM
	common.HexToAddress("0x4aae823a6a0b376de6a78e74ecc5b079d38cbcf7"): true, //SOLVBTC

}

func toFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func bigToFloat(b *big.Int) float64 {
	f, _ := new(big.Float).SetInt(b).Float64()
	return f
}

// build base token prices dynamically from USDT pairs
func BuildBasePriceMap(pairs []PairData) map[common.Address]float64 {
	priceMap := map[common.Address]float64{}

	usdt := common.HexToAddress("0x55d398326f99059fF775485246999027B3197955")
	priceMap[usdt] = 1

	for _, p := range pairs {
		t0 := common.HexToAddress(p.Token0)
		t1 := common.HexToAddress(p.Token1)

		r0, _ := new(big.Float).SetString(p.Reserve0)
		r1, _ := new(big.Float).SetString(p.Reserve1)

		f0, _ := r0.Float64()
		f1, _ := r1.Float64()

		// if token0 is USDT → price token1
		if t0 == usdt && f1 > 0 {
			priceMap[t1] = f0 / f1
		}

		// if token1 is USDT → price token0
		if t1 == usdt && f0 > 0 {
			priceMap[t0] = f1 / f0
		}
	}

	return priceMap
}

func BuildPairCalls(pairs []common.Address) []eth.Call3 {
	var calls []eth.Call3

	for _, p := range pairs {
		t0, _ := eth.PairABI().Pack("token0")
		t1, _ := eth.PairABI().Pack("token1")
		r, _ := eth.PairABI().Pack("getReserves")

		calls = append(calls,
			eth.Call3{Target: p, AllowFailure: true, CallData: t0},
			eth.Call3{Target: p, AllowFailure: true, CallData: t1},
			eth.Call3{Target: p, AllowFailure: true, CallData: r},
		)
	}

	return calls
}

func DecodePairResults(results []eth.Result, pairs []common.Address) []PairData {
	var raw []PairData

	if len(results) == 0 || len(pairs) == 0 {
		return raw
	}

	for i := 0; i < len(pairs); i++ {
		base := i * 3

		if base+2 >= len(results) {
			break
		}

		if !results[base].Success ||
			!results[base+1].Success ||
			!results[base+2].Success {
			continue
		}

		var token0 common.Address
		var token1 common.Address
		var reserves struct {
			Reserve0           *big.Int
			Reserve1           *big.Int
			BlockTimestampLast uint32
		}

		err0 := eth.PairABI().UnpackIntoInterface(&token0, "token0", results[base].ReturnData)
		err1 := eth.PairABI().UnpackIntoInterface(&token1, "token1", results[base+1].ReturnData)
		err2 := eth.PairABI().UnpackIntoInterface(&reserves, "getReserves", results[base+2].ReturnData)

		if err0 != nil || err1 != nil || err2 != nil {
			continue
		}

		if !(BaseTokens[token0] || BaseTokens[token1]) {
			continue
		}

		raw = append(raw, PairData{
			Address:  pairs[i].Hex(),
			Token0:   token0.Hex(),
			Token1:   token1.Hex(),
			Reserve0: reserves.Reserve0.String(),
			Reserve1: reserves.Reserve1.String(),
		})

	}

	// build dynamic price map
	priceMap := BuildBasePriceMap(raw)

	var out []PairData

	for _, p := range raw {
		t0 := common.HexToAddress(p.Token0)
		t1 := common.HexToAddress(p.Token1)

		r0, _ := new(big.Float).SetString(p.Reserve0)
		r1, _ := new(big.Float).SetString(p.Reserve1)

		f0, _ := r0.Float64()
		f1, _ := r1.Float64()

		price0, ok0 := priceMap[t0]
		price1, ok1 := priceMap[t1]

		var liquidity float64

		if ok0 && ok1 {
			liquidity = f0*price0 + f1*price1
		} else if ok0 {
			liquidity = 2 * f0 * price0
		} else if ok1 {
			liquidity = 2 * f1 * price1
		} else {
			continue
		}

		p.LiquidityUSD = liquidity
		out = append(out, p)
	}

	return out
}

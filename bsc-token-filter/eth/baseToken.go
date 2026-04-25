package eth

import "github.com/ethereum/go-ethereum/common"

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

	//Trusted tokens
	common.HexToAddress("0x4B0F1812e5Df2A09796481Ff14017e6005508003"): true, // TWT
	common.HexToAddress("0x4BD17003473389A42DAF6a0a729f6Fdb328BbBd7"): true, // VAI

}

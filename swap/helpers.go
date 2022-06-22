package swap

import (
	"github.com/MinterTeam/explorer-sdk/helpers"
	"math/big"
)

func getVolumeInBip(price *big.Float, volume string) *big.Float {
	firstCoinBaseVolume := helpers.Pip2Bip(helpers.StringToBigInt(volume))
	return new(big.Float).Mul(firstCoinBaseVolume, price)
}

func computePrice(reserve1, reserve2 string) *big.Float {
	return new(big.Float).Quo(
		helpers.Pip2Bip(helpers.StringToBigInt(reserve1)),
		helpers.Pip2Bip(helpers.StringToBigInt(reserve2)),
	)
}

func str2bigint(string string) *big.Int {
	bInt, _ := new(big.Int).SetString(string, 10)
	return bInt
}

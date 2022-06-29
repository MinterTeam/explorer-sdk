package swap

import (
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"math/big"
	"reflect"
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

func inArray(needle interface{}, haystack interface{}) bool {
	switch reflect.TypeOf(haystack).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(haystack)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(needle, s.Index(i).Interface()) == true {
				return true
			}
		}
	}

	return false
}

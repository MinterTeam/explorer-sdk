package models

import (
	"encoding/json"
	"fmt"
)

type Stake struct {
	ID             uint       `json:"id"               bun:",pk"`
	OwnerAddressID uint       `json:"owner_address_id"`
	ValidatorID    uint       `json:"validator_id"`
	CoinID         uint       `json:"coin_id"          pg:",use_zero"`
	Value          string     `json:"value"            bun:"type:numeric(70)"`
	BipValue       string     `json:"bip_value"        bun:"type:numeric(70)"`
	IsKicked       bool       `json:"is_kicked"`
	Coin           *Coin      `json:"coins"            bun:"rel:belongs-to"`                          //Relation has one to Coins
	OwnerAddress   *Address   `json:"owner_address"    bun:"rel:belongs-to,join:owner_address_id=id"` //Relation has one to Addresses
	Validator      *Validator `json:"validator"        bun:"rel:belongs-to"`                          //Relation has one to Validators
}

func (s Stake) String() string {
	bytes, err := json.Marshal(&s)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(bytes)
}

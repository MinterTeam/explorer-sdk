package models

const SlashEvent = "minter/SlashEvent"

type Slash struct {
	ID          uint64     `json:"id"            bun:",pk"`
	CoinID      uint       `json:"coin_id"       pg:",use_zero"`
	BlockID     uint64     `json:"block_id"`
	AddressID   uint       `json:"address_id"`
	ValidatorID uint64     `json:"validator_id"`
	Amount      string     `json:"amount"        bun:"type:numeric(70)"`
	Coin        *Coin      `bun:"rel:belongs-to"` //Relation has one to Coins
	Block       *Block     `bun:"rel:belongs-to"` //Relation has one to Blocks
	Address     *Address   `bun:"rel:belongs-to"` //Relation has one to Addresses
	Validator   *Validator `bun:"rel:belongs-to"` //Relation has one to Validators
}

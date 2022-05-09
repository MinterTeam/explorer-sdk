package models

const RewardEvent = "minter/RewardEvent"

type Reward struct {
	BlockID     uint64     `json:"block"        bun:",pk"`
	AddressID   uint       `json:"address_id"   bun:",pk"`
	ValidatorID uint64     `json:"validator_id" bun:",pk"`
	Role        string     `json:"role"         bun:",pk"`
	Amount      string     `json:"amount"       bun:"type:numeric(70)"`
	Block       *Block     `bun:"rel:belongs-to"` //Relation has one to Blocks
	Address     *Address   `bun:"rel:belongs-to"` //Relation has one to Addresses
	Validator   *Validator `bun:"rel:belongs-to"` //Relation has one to Validators
}

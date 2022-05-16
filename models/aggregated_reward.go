package models

import "time"

type AggregatedReward struct {
	FromBlockID uint64     `json:"from_block_id"`
	ToBlockID   uint64     `json:"to_block_id"`
	AddressID   uint64     `json:"address_id"    bun:",pk"`
	ValidatorID uint64     `json:"validator_id"  bun:",pk"`
	Role        string     `json:"role"          bun:",pk"`
	Amount      string     `json:"amount"        bun:"type:numeric(70)"`
	TimeID      time.Time  `json:"time_id"`
	FromBlock   *Block     `pg:"rel:belongs-to,join:from_block_id=id"` //Relation has one to Blocks
	ToBlock     *Block     `pg:"rel:belongs-to,join:to_block_id=id"`   //Relation has one to Blocks
	Address     *Address   `pg:"rel:belongs-to"`                       //Relation has one to Addresses
	Validator   *Validator `pg:"rel:belongs-to"`                       //Relation has one to Validators
}

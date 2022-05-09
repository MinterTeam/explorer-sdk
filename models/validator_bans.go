package models

type ValidatorBan struct {
	ValidatorId uint       `json:"validator_id" bun:",pk"`
	BlockId     uint64     `json:"block_id"     bun:",pk"`
	ToBlockId   uint64     `json:"to_block_id"`
	Validator   *Validator `json:"validator"    bun:"rel:belongs-to"`
	Block       *Block     `json:"block"        bun:"rel:belongs-to"`
	ToBlock     *Block     `json:"to_block"     bun:"rel:belongs-to,join:to_block_id=id"`
}

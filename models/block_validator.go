package models

import "github.com/uptrace/bun"

type BlockValidator struct {
	bun.BaseModel `bun:"table:block_validator"`
	Signed        bool       `json:"signed"`
	BlockID       uint64     `json:"block_id"     bun:",pk"`
	Block         *Block     `json:"block"        bun:"rel:belongs-to,join:block_id=id"`
	ValidatorID   uint64     `json:"validator_id" bun:",pk"`
	Validator     *Validator `json:"validator"    bun:"rel:belongs-to,join:validator_id=id"`
}

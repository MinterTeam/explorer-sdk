package models

import "time"

type ValidatorPublicKeys struct {
	ID          uint       `json:"id"  bun:",pk"`
	Key         string     `json:"key" bun:"type:varchar(64)"`
	ValidatorId uint       `json:"validator_id"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdateAt    *time.Time `json:"update_at"`
	Validator   *Validator `json:"validator" bun:"rel:belongs-to,join:validator_id=id"`
}

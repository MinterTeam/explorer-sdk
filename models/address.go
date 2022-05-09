package models

type Address struct {
	ID                  uint                  `json:"id"                   bun:",pk"`
	Address             string                `json:"address"              bun:"type:varchar(64)"`
	Balances            []*Balance            `json:"balances"             bun:"rel:has-many"` //relation has many to Balances
	Rewards             []*Reward             `json:"rewards"              bun:"rel:has-many"` //relation has many to Rewards
	Slashes             []*Slash              `json:"slashes"              bun:"rel:has-many"` //relation has many to Slashes
	Transactions        []*Transaction        `json:"transactions"         bun:"rel:has-many"` //relation has many to Transactions
	InvalidTransactions []*InvalidTransaction `json:"invalid_transactions" bun:"rel:has-many"` //relation has many to InvalidTransactions
}

// GetAddress Return address with prefix
func (a *Address) GetAddress() string {
	return `Mx` + a.Address
}

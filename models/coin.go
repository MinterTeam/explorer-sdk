package models

import (
	"fmt"
	"time"
)

type CoinType byte

const (
	_ CoinType = iota
	CoinTypeBase
	CoinTypeToken
	CoinTypePoolToken
)

type Coin struct {
	ID               uint           `json:"id"         bun:"id,pk"`
	Type             CoinType       `json:"type"`
	Name             string         `json:"name"`
	Symbol           string         `json:"symbol"`
	Volume           string         `json:"volume"     bun:"type:numeric(70)"`
	Crr              uint           `json:"crr"`
	Reserve          string         `json:"reserve"    bun:"type:numeric(70)"`
	MaxSupply        string         `json:"max_supply" bun:"type:numeric(70)"`
	Version          uint           `json:"version"    pg:",use_zero"`
	OwnerAddressId   uint           `json:"owner_address_id"`
	CreatedAtBlockId uint           `json:"created_at_block_id"`
	Burnable         bool           `json:"burnable"`
	Mintable         bool           `json:"mintable"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        *time.Time     `json:"updated_at"`
	DeletedAt        *time.Time     `bun:",soft_delete"`
	OwnerAddress     *Address       `bun:"rel:belongs-to,join:owner_address_id=id"`
	CreatedAtBlock   *Block         `bun:"rel:belongs-to,join:created_at_block_id=id"`
	Contracts        *TokenContract `bun:"rel:has-one"`
}

// GetSymbol Return coin with version
func (c *Coin) GetSymbol() string {
	if c.Version == 0 {
		return c.Symbol
	}
	return fmt.Sprintf("%s-%d", c.Symbol, c.Version)
}

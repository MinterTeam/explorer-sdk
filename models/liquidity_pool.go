package models

import (
	"fmt"
	"time"
)

const LockedLiquidityVolume = 1000

type LiquidityPool struct {
	Id               uint64 `json:"id"                 bun:",pk"`
	TokenId          uint64 `json:"token_id"`
	FirstCoinId      uint64 `json:"first_coin_id"`
	SecondCoinId     uint64 `json:"second_coin_id" `
	FirstCoinVolume  string `json:"first_coin_volume"  bun:"type:numeric(100)"`
	SecondCoinVolume string `json:"second_coin_volume" bun:"type:numeric(100)"`
	Liquidity        string `json:"liquidity"`
	LiquidityBip     string `json:"liquidity_bip"`
	UpdatedAtBlockId uint64 `json:"updated_at_block_id"`
	FirstCoin        *Coin  `json:"first_coin"  bun:"rel:belongs-to,join:first_coin_id=id"`
	SecondCoin       *Coin  `json:"second_coin" bun:"rel:belongs-to,join:second_coin_id=id"`
	Token            *Coin  `json:"token"       bun:"rel:belongs-to,join:token_id=id"`
}

func (lp *LiquidityPool) GetTokenSymbol() string {
	return fmt.Sprintf("LP-%d", lp.Id)
}

type AddressLiquidityPool struct {
	LiquidityPoolId  uint64         `json:"liquidity_pool_id"  bun:",pk"`
	AddressId        uint64         `json:"address_id"         bun:",pk"`
	FirstCoinVolume  string         `json:"first_coin_volume"  bun:"type:numeric(100)"`
	SecondCoinVolume string         `json:"second_coin_volume" bun:"type:numeric(100)"`
	Liquidity        string         `json:"liquidity"`
	Address          *Address       `json:"address"            bun:"rel:belongs-to"`
	LiquidityPool    *LiquidityPool `json:"liquidity_pool"     bun:"rel:belongs-to,join:liquidity_pool_id=id"`
}

type TagLiquidityPool struct {
	PoolID   uint64 `json:"pool_id"`
	CoinIn   uint64 `json:"coin_in"`
	ValueIn  string `json:"value_in"`
	CoinOut  uint64 `json:"coin_out"`
	ValueOut string `json:"value_out"`
}

type LiquidityPoolSnapshot struct {
	BlockId          uint64    `json:"block_id"`
	LiquidityPoolId  uint64    `json:"liquidity_pool_id"`
	FirstCoinVolume  string    `json:"first_coin_volume"`
	SecondCoinVolume string    `json:"second_coin_volume"`
	Liquidity        string    `json:"liquidity"`
	LiquidityBip     string    `json:"liquidity_bip"`
	CreatedAt        time.Time `json:"created_at"`
}

type LiquidityPoolTrade struct {
	BlockId          uint64         `json:"block_id"`
	LiquidityPoolId  uint64         `json:"liquidity_pool_id"`
	TransactionId    uint64         `json:"transaction_id"`
	FirstCoinVolume  string         `json:"first_coin_volume"`
	SecondCoinVolume string         `json:"second_coin_volume"`
	CreatedAt        time.Time      `json:"created_at"`
	Block            *Block         `json:"block"          bun:"rel:belongs-to"`
	LiquidityPool    *LiquidityPool `json:"liquidity_pool" bun:"rel:belongs-to,join:liquidity_pool_id=id"`
	Transaction      *Transaction   `json:"transaction"    bun:"rel:belongs-to"`
}

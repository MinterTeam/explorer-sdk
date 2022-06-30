package models

import (
	"encoding/json"
	"github.com/uptrace/bun"
	"time"
)

type Transaction struct {
	ID                  uint64               `json:"id" bun:"id,pk,autoincrement"`
	FromAddressID       uint64               `json:"from_address_id"`
	Nonce               uint64               `json:"nonce"`
	GasPrice            uint64               `json:"gas_price"`
	Gas                 uint64               `json:"gas"`
	Commission          string               `json:"commission"`
	BlockID             uint64               `json:"block_id"`
	GasCoinID           uint64               `json:"gas_coin_id" pg:",use_zero"`
	CreatedAt           time.Time            `json:"created_at"`
	Type                uint8                `json:"type"`
	Hash                string               `json:"hash"`
	ServiceData         string               `json:"service_data"`
	Data                json.RawMessage      `json:"data"`
	IData               interface{}          `json:"-"            bun:"-"`
	Tags                map[string]string    `json:"tags"`
	Payload             []byte               `json:"payload"`
	RawTx               []byte               `json:"raw_tx"`
	CommissionPriceCoin interface{}          `json:"commission_price_coin" bun:"-"`
	Block               *Block               `json:"block"        bun:"rel:belongs-to"`                                       //Relation has one to Blocks
	FromAddress         *Address             `json:"from_address" bun:"rel:belongs-to,join:from_address_id=id"`               //Relation has one to Address
	GasCoin             *Coin                `json:"gas_coin"     bun:"rel:belongs-to,join:gas_coin_id=id"`                   //Relation has one to Coin
	Validators          []*Validator         `json:"validators"   bun:"m2m:transaction_validator,join:Transaction=Validator"` //Relation has many to Validators
	TxOutputs           []*TransactionOutput `json:"tx_outputs"   bun:"rel:has-many"`
	TxOutput            *TransactionOutput   `json:"tx_output"    bun:"rel:has-one"`
}

type TransactionValidator struct {
	bun.BaseModel `bun:"table:transaction_validator"`
	TransactionID uint64       `bun:",pk"`
	Transaction   *Transaction `bun:"rel:belongs-to,join:transaction_id=id"`
	ValidatorID   uint64       `bun:",pk"`
	Validator     *Validator   `bun:"rel:belongs-to,join:validator_id=id"`
}

type TransactionLiquidityPool struct {
	tableName       struct{} `pg:"transaction_liquidity_pool"`
	TransactionID   uint64
	LiquidityPoolID uint64
}

// GetFee Return fee for transaction
func (t *Transaction) GetFee() uint64 {
	return t.GasPrice * t.Gas
}

// GetHash Return transactions hash with prefix
func (t *Transaction) GetHash() string {
	return `Mt` + t.Hash
}

package models

type TransactionOutput struct {
	ID            uint64       `json:"id"`
	TransactionID uint64       `json:"transaction_id"`
	ToAddressID   uint64       `json:"to_address_id"`
	CoinID        uint         `json:"coin_id"     pg:",use_zero"`
	Value         string       `json:"value"       bun:"type:numeric(70)"`
	Coin          *Coin        `json:"coin"        bun:"rel:belongs-to"`                       //Relation has one to Coins
	ToAddress     *Address     `json:"to_address"  bun:"rel:belongs-to,join:to_address_id=id"` //Relation has one to Addresses
	Transaction   *Transaction `json:"transaction" bun:"rel:belongs-to"`                       //Relation has one to Transactions
}

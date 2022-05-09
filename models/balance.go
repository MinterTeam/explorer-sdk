package models

type Balance struct {
	AddressID uint     `json:"address_id"  bun:",pk"`
	CoinID    uint     `json:"coin_id"     bun:",pk"`
	Value     string   `json:"value"       bun:"type:numeric(70)"`
	Address   *Address `bun:"rel:belongs-to,join:address_id=id"` //Relation has one to Address
	Coin      *Coin    `bun:"rel:belongs-to,join:coin_id=id"`    //Relation has one to Coin
}

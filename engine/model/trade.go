package model

import (
	"github.com/shopspring/decimal"
)

type Trade struct {
	ID       int64  `json:"id" xorm:"'id' pk autoincr"`
	MarketID string `json:"marketID" xorm:"market_id"`
	Maker    string `json:"maker"`
	Taker    string `json:"taker"`
	// TakerSide       string          `json:"takerSide"`
	// MakerOrderID    string          `json:"makerOrderID" xorm:"maker_order_id"`
	// TakerOrderID    string          `json:"takerOrderID" xorm:"taker_order_id"`
	// Sequence        int             `json:"sequence"`
	Amount    decimal.Decimal `json:"amount"`
	Price     decimal.Decimal `json:"price"`
	UpdatedAt int64           `json:"updatedAt"`
	TakerSide string          `json:"takerSide"`
	HashUrl   string          `json:"hashUrl"`
	// Fee             decimal.Decimal `json:"fee"`
}

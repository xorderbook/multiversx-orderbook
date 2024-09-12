package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type MatchRecord struct {
	ID           int64  `json:"id" xorm:"'id' pk autoincr"`
	MarketID     string `json:"marketID" xorm:"market_id"`
	Maker        string `json:"maker"`
	Taker        string `json:"taker"`
	TakerSide    string `json:"takerSide"`
	MakerOrderID string `json:"makerOrderID" xorm:"maker_order_id"`
	TakerOrderID string `json:"takerOrderID" xorm:"taker_order_id"`
	// Sequence        int             `json:"sequence"`
	Amount        decimal.Decimal `json:"amount"`
	MatchedAmount string          `json:"matched_amount"`
	Price         decimal.Decimal `json:"price"`
	UpdatedAt     time.Time       `json:"updatedAt"`
	Fee           decimal.Decimal `json:"fee"`
	State         string          `json:"state"`
	HashUrl       string          `json:"hashUrl"`
}

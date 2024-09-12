package common

import (
	"encoding/json"

	"github.com/shopspring/decimal"
)

type EngineEvent struct {
	EventID int
	Data    string
}

type LimitOrder struct {
	Id              string
	Qty             decimal.Decimal
	Price           decimal.Decimal
	Side            int
	Creator         string
	ContractAddress string
	TimeStamp       int64
	Filled          decimal.Decimal
}

func (ev LimitOrder) ToString() string {
	b, _ := json.Marshal(ev)
	return string(b)
}

type UpdateLimitOrder struct {
	ContractAddress string
	OrderList       []string
}

func (ev UpdateLimitOrder) ToString() string {
	b, _ := json.Marshal(ev)
	return string(b)
}

type MatchOrdersExt struct {
	Caller       string
	Epoch        uint64
	OrderType    uint64
	OrderId      uint64
	OrderCreator string
}

package model

import (
	"github.com/shopspring/decimal"
)

type LimitOrder struct {
	Id              string          `xorm:"'id' pk"`
	Qty             decimal.Decimal `xorm:"'qty'"`
	Price           decimal.Decimal `xorm:"'price'"`
	Creator         string          `xorm:"'creator'"`
	ContractAddress string          `xorm:"'contract_address'"`
	TimeStamp       int64           `xorm:"'timestamp'"`
	Filled          decimal.Decimal `xorm:"'filled'"`
	State           string          `xorm:"'state'"`
	Market          string          `xorm:"'market'"`
	Side            int             `xorm:"'side'"`
}

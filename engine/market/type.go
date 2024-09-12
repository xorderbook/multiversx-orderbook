package market

import (
	"engine/orderbook"
)

type MatchOrder struct {
	TakerOrder     orderbook.Order
	MakerOrderList []orderbook.Order
}

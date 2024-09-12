package market

import (
	"context"
	"encoding/json"
	"engine/common"
	"engine/conf"
	"engine/db"
	"engine/launcher"
	"engine/model"
	"engine/orderbook"
	util "engine/util/redis"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

type Market struct {
	orderbook *orderbook.OrderBook
	conf      conf.Market
	sequence  uint64
	ch        chan interface{}
	rdb       *redis.Client

	orders map[string]*common.LimitOrder
}

func NewMarket(cfg conf.Market) *Market {
	m := new(Market)
	m.orderbook = orderbook.NewOrderBook()
	m.conf = cfg
	m.ch = make(chan interface{}, 1000)
	m.sequence = 0
	m.rdb = util.NewRedis("127.0.0.1:6379")
	m.orders = make(map[string]*common.LimitOrder)
	return m
}

func (m *Market) Push(order interface{}) {
	m.ch <- order
}

func (m *Market) run() {
	log.Infof("%v started", m.conf.Market)
	for {
		select {
		case payload := <-m.ch:
			switch s := payload.(type) {
			case *common.LimitOrder:
				m.processLimitOrder(s)
			case *common.UpdateLimitOrder:
				m.confirmMatch(s)
			}
		}
	}
}

func (m *Market) processLimitOrder(data interface{}) {
	order := data.(*common.LimitOrder)
	log.Infof("processLimitOrder market: %s, order id: %v, qty: %v, price: %v, side: %v", m.conf.Market, order.Id, order.Qty, order.Price, order.Side)
	m.orders[order.Id] = order
	var side orderbook.Side
	var takeSide string
	if order.Side == common.NEW_BUY_ORDER {
		side = orderbook.Buy
		takeSide = "buy"
	} else {
		side = orderbook.Sell
		takeSide = "sell"
	}

	err := db.SaveNewOrder(m.conf.Market, *order)
	if err != nil {
		return
	}

	done, partial, partialQuantityProcessed, err := m.orderbook.ProcessLimitOrder(
		order.Creator,
		side,
		order.Id,
		order.Qty,
		order.Price)

	if err != nil {
		log.Error(err)
	}
	log.Debug("processLimitOrder result...")
	log.Debugf(" done: %v", done)
	log.Debugf(" partial: %v", partial)
	log.Debugf(" partialQuantityProcessed: %v", partialQuantityProcessed)

	makerOrderList := make([]string, 0)
	makerOrderObjectList := make([]*orderbook.Order, 0)
	hasMatch := false
	// whole match
	if len(done) > 0 {
		hasMatch = true
		for _, o := range done {
			if o.ID() != order.Id {
				makerOrderList = append(makerOrderList, o.ID())
				makerOrderObjectList = append(makerOrderObjectList, o)
			}
			// log.Debugf(" update order done: %v", o.ID())
			// db.UpdateOrderDone(o.ID())
			db.UpdateOrderFilled(o.ID())
		}
	}
	// partial match
	if partial != nil {
		hasMatch = true
		if partial.ID() != order.Id {
			makerOrderList = append(makerOrderList, partial.ID())
			makerOrderObjectList = append(makerOrderObjectList, partial)
		}
		log.Debugf(" UpdateOrderPartialFilled: %v , %v", partial.ID(), partialQuantityProcessed)
		db.UpdateOrderPartialFilled(partial.ID(), partialQuantityProcessed)
	}
	if hasMatch {
		matchId := db.SaveMatchRecord(order.Id, makerOrderObjectList, partialQuantityProcessed, m.conf.Market, takeSide)
		launcher.SubmitBlockchainOrderMatch(matchId, m.conf.ExchangeAddress, order.Id, makerOrderList)
	}
	m.updateOrderbookView()
}

func (m *Market) confirmMatch(confirm *common.UpdateLimitOrder) {
	log.Debugf("... confirmMatch: %v, confirm: %v", m.conf.Market, confirm)
	pendingMatch, ok := db.GetPendingMatch(m.conf.Market)
	if !ok {
		return
	}
	log.Debugf("pending match: %#v", pendingMatch)

	if pendingMatch.State != "pending" {
		return
	}

	// get maker order list
	var makerOrderIds []string
	json.Unmarshal([]byte(pendingMatch.MakerOrderID), &makerOrderIds)
	log.Debugf("  makerOrderIds: %#v", makerOrderIds)

	// get maker match amount list
	var makerOrderMatchedAmount []string
	json.Unmarshal([]byte(pendingMatch.MatchedAmount), &makerOrderMatchedAmount)
	log.Debugf("  maker order matched amount: %#v", makerOrderMatchedAmount)

	// update maker order state
	db.UpdateMakerOrderState(makerOrderIds, makerOrderMatchedAmount)

	// update taker order state
	db.UpdateTakerOrderState(pendingMatch.TakerOrderID, makerOrderMatchedAmount)

	// create trade record
	dbTradeList := make([]model.Trade, 0)
	taker := pendingMatch.Taker
	for k, makerOrderId := range makerOrderIds {
		var trade model.Trade
		trade.MarketID = m.conf.Market
		trade.Taker = taker
		trade.Maker = db.GetOrderCreator(makerOrderId)
		trade.Amount = decimal.RequireFromString(makerOrderMatchedAmount[k])
		trade.UpdatedAt = time.Now().Unix()
		trade.TakerSide = pendingMatch.TakerSide
		trade.Price = db.GetOrderPrice(makerOrderId)
		trade.HashUrl = pendingMatch.HashUrl
		dbTradeList = append(dbTradeList, trade)
	}
	// insert db
	db.InsertTradeRecords(dbTradeList)
	// update match done
	db.UpdatePendingMatchFinish(pendingMatch.ID)
}

func (m *Market) getOrderbookView() string {
	data, err := m.orderbook.MarshalJSON()
	if err != nil {
		log.Error(err)
	}
	return string(data)
}

func (m *Market) updateOrderbookView() {
	ctx := context.Background()
	orderbook := m.getOrderbookView()
	log.Debugf("updateOrderbookView: %v, orderbook: %v", m.conf.Market, orderbook)
	err := m.rdb.Set(ctx, m.conf.Market, orderbook, 0).Err()
	if err != nil {
		log.Error("updateOrderbookView:", err)
		return
	}
}

// func (m *Market) GetOpenOrders(orderIds []string) []*common.LimitOrder {
// 	list := make([]*common.LimitOrder, 0)
// 	for _, orderId := range orderIds {
// 		if _, ok := m.orders[orderId]; ok {
// 			list = append(list, m.orders[orderId])
// 		}
// 	}
// 	for i := 0; i < 20; i++ {
// 		order := new(common.LimitOrder)
// 		order.Id = fmt.Sprintf("%v", i)
// 		order.Creator = "elgdxxx"
// 		// order.Filled = decimal.New(10, 0)
// 		// order.Price = decimal.NewFromFloat(12.2)
// 		// order.Qty = decimal.NewFromFloat(1.2)
// 		order.Side = 0
// 		order.TimeStamp = time.Now().Unix()

// 		list = append(list, order)
// 	}
// 	return list
// }

func (m *Market) UpdateOrder(orderIds []string) []*common.LimitOrder {
	list := make([]*common.LimitOrder, 0)
	for _, orderId := range orderIds {
		if _, ok := m.orders[orderId]; ok {
			list = append(list, m.orders[orderId])
		}
	}
	return list
}

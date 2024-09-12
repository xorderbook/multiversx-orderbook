package api

import (
	"context"
	"encoding/json"
	"engine/model"
	"engine/util/convert"
	util "engine/util/redis"
	"sort"

	db "engine/util/mysql"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

var rdb *redis.Client
var UNIT = decimal.RequireFromString("1000000000000000000")

func init() {
	rdb = util.NewRedis("127.0.0.1:6379")
}

func toEtherAmount(amount string) float64 {
	val := decimal.RequireFromString(amount).Div(UNIT)
	return convert.StrToUint64(val.String())
}

func OrderBook(c echo.Context) error {
	req := new(OrderbookReq)
	err := c.Bind(req)
	if err != nil {
		log.Errorf("OrderBook %v", err)
		return err
	}

	cmd := rdb.Get(context.Background(), req.MarketID)
	redisData, err := cmd.Result()

	if err == redis.Nil {
		return c.JSON(200, map[string]interface{}{"code": 0, "data": ""})
	} else if err != nil {
		log.Error(err)
		return c.JSON(200, map[string]interface{}{"code": 0, "data": ""})
	}

	view := new(OrderbookView)
	err = json.Unmarshal([]byte(redisData), view)
	if err != nil {
		log.Error(err)
		return c.JSON(200, map[string]interface{}{"code": 0, "data": ""})
	}

	askList := make([][]float64, 0)
	for _, v := range view.Asks.Prices {
		price := convert.StrToUint64(v.Price)
		volume := toEtherAmount(v.Volume)
		askList = append(askList, []float64{price, volume})
	}
	sortByPrice(askList, true)

	bidList := make([][]float64, 0)
	for _, v := range view.Bids.Prices {
		price := convert.StrToUint64(v.Price)
		volume := toEtherAmount(v.Volume)
		bidList = append(bidList, []float64{price, volume})
	}
	sortByPrice(bidList, false)

	resp := make(map[string]interface{})
	resp["bids"] = bidList
	resp["asks"] = askList
	data, _ := json.Marshal(resp)

	return c.JSON(200, map[string]interface{}{"code": 0, "data": string(data)})
}

func sortByPrice(askList [][]float64, asc bool) {
	if asc {
		sort.Slice(askList, func(i, j int) bool {
			return askList[i][0] < askList[j][0]
		})
	} else {
		sort.Slice(askList, func(i, j int) bool {
			return askList[i][0] > askList[j][0]
		})
	}
}

func OpenOrder(c echo.Context) error {
	req := new(GetOpenOrderReq)
	err := c.Bind(req)
	if err != nil {
		log.Errorf("OrderBook %v", err)
		return err
	}

	sess := db.GetSession()
	defer sess.Close()

	openOrders := make([]*model.LimitOrder, 0)
	sess.Where("market = ? and creator = ? and state != ?", req.MarketID, convert.HexAddress(req.Address), "done").Find(&openOrders)

	for _, o := range openOrders {
		o.Qty = o.Qty.Div(UNIT)
		o.Filled = o.Filled.Div(UNIT)
	}

	// // test code
	// for i := 0; i < 20; i++ {
	// 	order := new(model.LimitOrder)
	// 	order.Id = fmt.Sprintf("%v", i)
	// 	order.Creator = "elgdxxx"
	// 	// order.Filled = decimal.New(10, 0)
	// 	// order.Price = decimal.NewFromFloat(12.2)
	// 	// order.Qty = decimal.NewFromFloat(1.2)
	// 	// order.Side = 0
	// 	order.TimeStamp = time.Now().Unix()

	// 	openOrders = append(openOrders, order)
	// }

	return c.JSON(200, map[string]interface{}{"code": 0, "data": openOrders})
}

func MyTrade(c echo.Context) error {
	req := new(GetMyTradeReq)
	err := c.Bind(req)
	if err != nil {
		log.Errorf("OrderBook %v", err)
		return err
	}

	sess := db.GetSession()
	defer sess.Close()

	list := make([]model.Trade, 0)
	hexAddress := convert.HexAddress(req.Address)
	sess.Where("market_id = ? and (maker = ? or taker = ?)", req.MarketID, hexAddress, hexAddress).OrderBy("id desc").Find(&list)

	return c.JSON(200, map[string]interface{}{"code": 0, "data": list})
}

func TradeHistory(c echo.Context) error {
	req := new(GetTradeHistory)
	err := c.Bind(req)
	if err != nil {
		log.Errorf("OrderBook %v", err)
		return err
	}

	sess := db.GetSession()
	defer sess.Close()

	list := make([]model.Trade, 0)
	sess.Where("maker = ? or taker = ?", req.Address).OrderBy("id desc").Find(&list)

	return c.JSON(200, map[string]interface{}{"code": 0, "data": list})
}

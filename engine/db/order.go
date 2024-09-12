package db

import (
	"encoding/json"
	types "engine/common"
	"engine/model"
	"engine/orderbook"
	"engine/util/mysql"
	"time"

	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

func getOrderIds(makerOrderList []*orderbook.Order) string {
	var list []string
	for _, v := range makerOrderList {
		list = append(list, v.ID())
	}
	b, _ := json.Marshal(list)

	return string(b)
}

func SaveNewOrder(market string, order types.LimitOrder) error {
	sess := mysql.GetSession()
	defer sess.Close()

	dbOrder := new(model.LimitOrder)
	dbOrder.ContractAddress = order.ContractAddress
	dbOrder.Creator = order.Creator
	dbOrder.Id = order.Id
	dbOrder.Filled = decimal.Decimal{}
	dbOrder.Price = order.Price
	dbOrder.Qty = order.Qty
	dbOrder.State = "open"
	dbOrder.Market = market
	dbOrder.Side = order.Side
	dbOrder.TimeStamp = time.Now().Unix()
	rowsAffected, err := sess.Insert(dbOrder)
	if err != nil {
		log.Error(err)
		return err
	}
	if rowsAffected == 0 {
		log.Errorf("order inserted not affected: %#v", dbOrder)
	}

	return nil
}

func UpdateOrderDone(orderId string) {
	sess := mysql.GetSession()
	defer sess.Close()

	dbOrder := new(model.LimitOrder)
	has, dbOrder := getOrderFromDB(orderId)
	if !has {
		return
	}
	dbOrder.Filled = dbOrder.Qty
	dbOrder.State = "done"
	rowsAffected, err := sess.Where("id = ?", orderId).Cols("state", "filled").Update(dbOrder)
	if err != nil {
		log.Error(err)
		return
	}

	if rowsAffected == 0 {
		log.Warnf("order not updated (UpdateOrderDone) orderId: %s", orderId)
	}
}

func UpdateOrderFilled(orderId string) {
	log.Debugf("UpdateOrderFilled: %v", orderId)
	sess := mysql.GetSession()
	defer sess.Close()

	dbOrder := new(model.LimitOrder)
	has, dbOrder := getOrderFromDB(orderId)
	if !has {
		return
	}
	dbOrder.Filled = dbOrder.Qty
	dbOrder.State = "settling"
	rowsAffected, err := sess.Where("id = ?", orderId).Cols("state", "filled").Update(dbOrder)
	if err != nil {
		log.Error(err)
		return
	}

	if rowsAffected == 0 {
		log.Warnf("order not updated (UpdateOrderFilled) orderId: %s", orderId)
	}
}

func UpdateOrderPartialFilled(orderId string, filledQty decimal.Decimal) {
	log.Debugf("UpdateOrderPartialFilled: %v", orderId)
	sess := mysql.GetSession()
	defer sess.Close()

	has, dbOrder := getOrderFromDB(orderId)
	if !has {
		return
	}

	dbOrder.Filled.Add(filledQty)
	dbOrder.State = "settling"
	rowsAffected, err := sess.Cols("state", "filled").Update(dbOrder)
	if err != nil {
		log.Error(err)
		return
	}

	if rowsAffected == 0 {
		log.Warnf("order not updated [UpdateOrderPartialFilled] orderId: %s", orderId)
	}
}

func getOrderFromDB(orderId string) (bool, *model.LimitOrder) {
	engine := mysql.GetEngine()

	dbVal := make(map[string]string)
	dbOrder := new(model.LimitOrder)
	has, err := engine.Table("limit_order").Where("id = ?", orderId).Get(&dbVal)
	if err != nil {
		log.Error(err)
		return false, nil
	}

	if !has {
		log.Errorf("getOrderFromDB not found orderId: %s", orderId)
		return false, nil
	}
	// log.Infof("dbVal: %#v", dbVal)

	if dbVal["filled"] != "" {
		dbOrder.Filled = decimal.RequireFromString(dbVal["filled"])
	}
	if dbVal["qty"] != "" {
		dbOrder.Qty = decimal.RequireFromString(dbVal["qty"])
	}
	dbOrder.Id = dbVal["id"]
	dbOrder.Creator = dbVal["creator"]
	dbOrder.Price = decimal.RequireFromString(dbVal["price"])

	return true, dbOrder
}

func GetOrderCreator(orderId string) string {
	has, o := getOrderFromDB(orderId)
	if !has {
		log.Errorf("GetOrderCreator not found orderId: %s", orderId)
		return ""
	}
	return o.Creator
}

func GetOrderPrice(orderId string) decimal.Decimal {
	has, o := getOrderFromDB(orderId)
	if !has {
		log.Errorf("GetOrderPrice not found orderId: %s", orderId)
		return decimal.Decimal{}
	}
	return o.Price
}

func UpdateMakerOrderState(makerOrderIds []string, makerOrderMatchedAmount []string) {
	sess := mysql.GetSession()
	defer sess.Close()

	for _, orderId := range makerOrderIds {
		has, dbOrder := getOrderFromDB(orderId)
		if !has {
			log.Errorf("UpdateMakerOrderState not found orderId: %s", orderId)
			return
		}
		// dbOrder.Filled = dbOrder.Filled.Add(decimal.RequireFromString(makerOrderMatchedAmount[k]))
		if dbOrder.Filled.Equal(dbOrder.Qty) {
			dbOrder.State = "done"
		} else {
			dbOrder.State = "open"
		}
		log.Debugf("UpdateMakerOrderState: %#v , state:%v", orderId, dbOrder.State)
		rowsAffected, err := sess.ID(orderId).Cols("filled", "state").Update(dbOrder)
		if err != nil {
			log.Error(err)
			return
		}

		if rowsAffected == 0 {
			log.Errorf("order not updated[UpdateMakerOrder] orderId: %s", orderId)
		}
	}
}

func UpdateTakerOrderState(takerOrderId string, makerOrderMatchedAmount []string) {
	sess := mysql.GetSession()
	defer sess.Close()

	filledSum := decimal.Decimal{}
	for _, amt := range makerOrderMatchedAmount {
		filledSum = filledSum.Add(decimal.RequireFromString(amt))
	}
	has, dbOrder := getOrderFromDB(takerOrderId)
	if !has {
		log.Errorf("UpdateTakerOrderState not found orderId: %s", takerOrderId)
		return
	}

	// dbOrder.Filled = dbOrder.Filled.Add(filledSum)
	if dbOrder.Filled.Equal(dbOrder.Qty) {
		dbOrder.State = "done"
	} else {
		dbOrder.State = "open"
	}

	rowsAffected, err := sess.ID(takerOrderId).Cols("filled", "state").Update(dbOrder)
	if err != nil {
		log.Error(err)
		return
	}

	if rowsAffected == 0 {
		log.Errorf("order not updated[UpdateTakerOrderState] orderId: %s", takerOrderId)
	}
}

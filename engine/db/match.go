package db

import (
	"encoding/json"
	"engine/model"
	"engine/orderbook"
	"engine/util/mysql"

	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

func GetPendingMatch(market string) (model.MatchRecord, bool) {
	sess := mysql.GetSession()
	defer sess.Close()

	list := make([]model.MatchRecord, 0)
	err := sess.Where("market_id = ? and state = ?", market, "pending").OrderBy("id asc").Limit(1).Find(&list)
	if err != nil {
		log.Error()
	}

	if len(list) == 1 {
		return list[0], true
	}
	return model.MatchRecord{}, false
}

func getOrderMatchedAmount(makerOrderList []*orderbook.Order) string {
	var list []string
	for _, v := range makerOrderList {
		list = append(list, v.Quantity().String())
	}
	b, _ := json.Marshal(list)

	return string(b)
}

func SaveMatchRecord(takerOrderId string, makerOrderList []*orderbook.Order, partialQuantityProcessed decimal.Decimal, MarketID string, takerSide string) int64 {
	sess := mysql.GetSession()
	defer sess.Close()

	dbMatch := new(model.MatchRecord)
	dbMatch.Taker = GetOrderCreator(takerOrderId)
	dbMatch.TakerOrderID = takerOrderId
	dbMatch.MakerOrderID = getOrderIds(makerOrderList)
	dbMatch.MatchedAmount = getOrderMatchedAmount(makerOrderList)
	dbMatch.MarketID = MarketID
	dbMatch.TakerSide = takerSide
	dbMatch.State = "pending"

	log.Infof("SaveMatchRecord...... : %#v", dbMatch)

	rowsAffected, err := sess.Insert(dbMatch)
	if err != nil {
		log.Error(err)
		return 0
	}
	if rowsAffected == 0 {
		log.Errorf("order inserted not effected: %#v", dbMatch)
	}

	return dbMatch.ID
}

func UpdatePendingMatchFinish(id int64) {
	sess := mysql.GetSession()
	defer sess.Close()

	// Execute a SQL statement
	result, err := sess.Exec("UPDATE match_record SET state = ? WHERE id = ?", "done", id)
	if err != nil {
		log.Error(err)
	}

	// Check the number of affected rows
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error(err)
		return
	}

	if rowsAffected == 0 {
		log.Errorf("order not updated (UpdatePendingMatchFinish) id: %v", id)
	}
}

func UpdateMatchTxHash(id int64, hash, url string) {
	sess := mysql.GetSession()
	defer sess.Close()

	// Execute a SQL statement
	result, err := sess.Exec("UPDATE match_record SET hash_url = ? WHERE id = ?", url, id)
	if err != nil {
		log.Error(err)
	}

	// Check the number of affected rows
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error(err)
		return
	}

	if rowsAffected == 0 {
		log.Errorf("order not updated (UpdateMatchTxHash) id: %v", id)
	}
}

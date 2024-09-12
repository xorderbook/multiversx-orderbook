package db

import (
	"engine/model"
	"engine/util/mysql"

	log "github.com/sirupsen/logrus"
)

func InsertTradeRecords(tradeList []model.Trade) {
	sess := mysql.GetSession()
	defer sess.Close()

	for _, v := range tradeList {
		log.Debugf("InsertTradeRecords: %#v\n", v)
	}

	rowsAffected, err := sess.Insert(tradeList)
	if err != nil {
		log.Error()
		return
	}

	if rowsAffected == 0 {
		log.Warnf("InsertTradeRecords not affected: %#v", tradeList)
	}
}

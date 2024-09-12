package common

import (
	"encoding/json"
	"fmt"
)

type NotifierEvent struct {
	EventID int
	Data    string
}

type WaitOrderMatch struct {
	NumOfEventsNeedsToWait int
	ReceivedEvents         []MatchOrdersExt
	TxHash                 string
}

func NewWaitOrderMatch() *WaitOrderMatch {
	return &WaitOrderMatch{
		NumOfEventsNeedsToWait: 0,
		ReceivedEvents:         make([]MatchOrdersExt, 0),
		TxHash:                 "",
	}
}

func (e *WaitOrderMatch) AllEventReceived() bool {
	return e.NumOfEventsNeedsToWait == len(e.ReceivedEvents)
}

func (e *WaitOrderMatch) GetMatchedOrderIds() []string {
	var orderIds []string
	for _, v := range e.ReceivedEvents {
		orderIds = append(orderIds, fmt.Sprintf("%v", v.OrderId))
	}
	return orderIds
}

func (e *WaitOrderMatch) ToString() string {
	b, _ := json.Marshal(e)
	return string(b)
}

package websocket

import (
	"context"
	"encoding/json"
	"engine/notifier"
	"fmt"
	"time"

	util "engine/util/redis"

	"github.com/go-redis/redis/v8"

	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/data"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"

	types "engine/common"
)

const (
	// hostURL = "127.0.0.1:5000"
	// hostURL = "147.182.246.6:5000"
	wsPath = "/hub/ws"
)

type EventNotifier struct {
	exchangeAddressList []string
	rdb                 *redis.Client
	// order match event
	// key hash  v numbers of event needs to wait
	waitTxHashList map[string]*types.WaitOrderMatch
	hostURL        string
}

func NewEventNotifier(addressList []string, hostURL string) *EventNotifier {
	eventNotifier := new(EventNotifier)
	eventNotifier.exchangeAddressList = addressList
	eventNotifier.rdb = util.NewRedis("127.0.0.1:6379")
	eventNotifier.waitTxHashList = make(map[string]*types.WaitOrderMatch)
	eventNotifier.hostURL = hostURL
	return eventNotifier
}

func (e *EventNotifier) Start() {
	go e.run()
	go e.drainMsg()
}

func (e *EventNotifier) run() {
	ws, err := NewWSClient(e.hostURL)
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer ws.Close()

	subscribeEvent := &data.SubscribeEvent{
		SubscriptionEntries: []data.SubscriptionEntry{
			{
				EventType: common.BlockEvents,
			},
			{
				EventType: common.BlockScrs,
			},
			{
				EventType: common.BlockTxs,
			},
		},
	}

	ws.SendSubscribeMessage(subscribeEvent)

	for {
		m, err := ws.ReadMessage()
		if err != nil {
			log.Error(err.Error())
			continue
		}

		var reply data.WebSocketEvent
		err = json.Unmarshal(m, &reply)
		if err != nil {
			log.Error(err.Error())
			continue
		}

		switch reply.Type {
		case common.BlockEvents:
			var event data.BlockEventsWithOrder
			_ = json.Unmarshal(reply.Data, &event)
			// log.Info("BlockEvents", event)
			for _, event := range event.Events {
				// log.Infof(" %#v\n", event)

				if e.isInterestedEvent(event) {
					// log.Debugf("contract event: %#v", event)
					if event.Identifier == "createBuyOrder" {
						e.handleCreateOrder(event, "buy")
					} else if event.Identifier == "createSellOrder" {
						e.handleCreateOrder(event, "sell")
					} else if event.Identifier == "matchOrdersExt" {
						e.handleMatch(event)
					}
				}
			}

		case common.RevertBlockEvents:
			var event *data.RevertBlock
			_ = json.Unmarshal(reply.Data, &event)
			// log.Info("RevertBlockEvents", event)
		case common.FinalizedBlockEvents:
			var event *data.FinalizedBlock
			_ = json.Unmarshal(reply.Data, &event)
			// log.Info("FinalizedBlockEvents", event)
		case common.BlockTxs:
			var event *data.BlockTxs
			_ = json.Unmarshal(reply.Data, &event)
			// log.Info("BlockTxs", event)
		case common.BlockScrs:
			var event data.BlockScrs
			_ = json.Unmarshal(reply.Data, &event)
			// log.Info("BlockScrs", event)
		default:
			log.Error("invalid message type")
		}
	}
}

func (e *EventNotifier) isInterestedEvent(ev data.Event) bool {
	for _, addr := range e.exchangeAddressList {
		if ev.Address == addr {
			return true
		}
	}
	return false
}

func (e *EventNotifier) drainMsg() {
	ticker := time.NewTicker(200 * time.Millisecond)

	for {
		select {
		case <-ticker.C:
			// log.Infof("engine drainMsg")
			key := types.NOTIFIER_KEY
			value, err := e.rdb.BRPop(context.Background(), 1*time.Second, key).Result()
			if err == redis.Nil {
				continue
			}
			if err != nil {
				continue
			}

			notifyEvent := parseRedisMsg(value[1])
			if notifyEvent != nil {
				e.handleEvent(notifyEvent)
			}
		}
	}
}

func (e *EventNotifier) handleEvent(ev *types.NotifierEvent) {
	if ev.EventID == types.WAIT_MATCH_EVENT {
		waitOrderMatch := types.NewWaitOrderMatch()
		err := json.Unmarshal([]byte(ev.Data), waitOrderMatch)
		if err != nil {
			log.Error(err)
			return
		}
		e.waitTxHashList[waitOrderMatch.TxHash] = waitOrderMatch
	}
}

func (e *EventNotifier) handleCreateOrder(ev data.Event, side string) {
	log.Debugf("===> handleCreateOrder %#v", ev)
	id, creator, _, inputAmount, outputAmount := notifier.ParseOrderDataByte(ev.Data)
	log.Debugf(" id:%v, creator:%v, inputAmount:%v, outputAmount:%v, side:%v", id, creator, inputAmount, outputAmount, side)

	if side == "buy" {
		event := types.EngineEvent{}
		newOrder := types.LimitOrder{
			Id:              fmt.Sprintf("%v", id),
			Qty:             outputAmount,
			Price:           inputAmount.Div(outputAmount),
			Side:            types.NEW_BUY_ORDER,
			ContractAddress: ev.Address,
			Creator:         creator,
			TimeStamp:       time.Now().Unix(),
			Filled:          decimal.RequireFromString("0"),
		}
		event.EventID = types.NEW_BUY_ORDER
		event.Data = newOrder.ToString()
		e.pushEngineEvent(event)
	} else {
		event := types.EngineEvent{}
		newOrder := types.LimitOrder{
			Id:              fmt.Sprintf("%v", id),
			Qty:             inputAmount,
			Price:           outputAmount.Div(inputAmount),
			Side:            types.NEW_SELL_ORDER,
			ContractAddress: ev.Address,
			Creator:         creator,
			TimeStamp:       time.Now().Unix(),
			Filled:          decimal.RequireFromString("0"),
		}
		event.EventID = types.NEW_SELL_ORDER
		event.Data = newOrder.ToString()
		e.pushEngineEvent(event)
	}
}

func (e *EventNotifier) handleMatch(ev data.Event) {
	matchOrdersExtEvent := notifier.ParseMatchOrdersExtEventByte(ev.Topics)
	log.Debugf("caller:%v, epoch:%v, orderType:%v, orderId:%v, orderCreator:%v", matchOrdersExtEvent.Caller, matchOrdersExtEvent.Epoch, matchOrdersExtEvent.OrderType, matchOrdersExtEvent.OrderId, matchOrdersExtEvent.OrderCreator)
	waitBlockMatchEvent := e.waitTxHashList[ev.TxHash]

	waitBlockMatchEvent.ReceivedEvents = append(waitBlockMatchEvent.ReceivedEvents, matchOrdersExtEvent)
	log.Debugf("waitBlockMatchEvent hash:%v, ..%v events received ", waitBlockMatchEvent.TxHash, len(waitBlockMatchEvent.ReceivedEvents))
	if waitBlockMatchEvent.AllEventReceived() {
		log.Debugf("waitBlockMatchEvent hash:%v all events received", waitBlockMatchEvent.TxHash)
		// push engine market execution record
		engineEvent := new(types.EngineEvent)
		engineEvent.EventID = types.UPDATE_ORDER_EVENT
		engineEvent.Data = types.UpdateLimitOrder{
			OrderList:       waitBlockMatchEvent.GetMatchedOrderIds(),
			ContractAddress: ev.Address,
		}.ToString()
		e.pushEngineEvent(*engineEvent)
	}
}

func (e *EventNotifier) pushEngineEvent(engineEvent types.EngineEvent) {
	b, _ := json.Marshal(engineEvent)
	e.pushRedis(string(b))
}

func (e *EventNotifier) pushRedis(data string) {
	key := types.ENGINE_KEY
	_, err := e.rdb.LPush(context.Background(), key, data).Result()
	if err != nil {
		log.Error(err)
	}
}

func parseRedisMsg(str string) *types.NotifierEvent {
	ev := new(types.NotifierEvent)
	err := json.Unmarshal([]byte(str), ev)
	if err != nil {
		log.Error(err)
		return nil
	}
	return ev
}

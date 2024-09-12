package mpq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"engine/common"
	"engine/notifier"
	util "engine/util/redis"

	"github.com/go-redis/redis/v8"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type EventNotifier struct {
	exchangeAddressList []string
	con                 string
	rdb                 *redis.Client
	eventId             string
	// order match event
	// key hash  v numbers of event needs to wait
	waitTxHashList map[string]*common.WaitOrderMatch
}

func NewEventNotifier(addressList []string, connection, eventId string) *EventNotifier {
	eventNotifier := new(EventNotifier)
	eventNotifier.exchangeAddressList = addressList
	eventNotifier.con = connection
	eventNotifier.rdb = util.NewRedis("127.0.0.1:6379")
	eventNotifier.eventId = eventId
	eventNotifier.waitTxHashList = make(map[string]*common.WaitOrderMatch)
	return eventNotifier
}

func (e *EventNotifier) Start() {
	go e.run()
	go e.drainMsg()
}

func (e *EventNotifier) run() {
	log.Debug("start EventNotifier")

	connection, err := amqp.Dial(e.con)
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	log.Debug("Successfully connected to RabbitMQ instance")

	ch, err := connection.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclarePassive(
		e.eventId, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		true,   // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	for d := range msgs {
		blockEvent := new(notifier.BlockchainEvent)

		err := json.Unmarshal(d.Body, blockEvent)
		if err != nil {
			log.Error(err)
		}

		for _, event := range blockEvent.Events {
			// log.Debugf("contract event: %#v", event)
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
	}
}

func (e *EventNotifier) isInterestedEvent(ev notifier.Event) bool {
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
			key := common.NOTIFIER_KEY
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

func (e *EventNotifier) handleEvent(ev *common.NotifierEvent) {
	if ev.EventID == common.WAIT_MATCH_EVENT {
		waitOrderMatch := common.NewWaitOrderMatch()
		err := json.Unmarshal([]byte(ev.Data), waitOrderMatch)
		if err != nil {
			log.Error(err)
			return
		}
		e.waitTxHashList[waitOrderMatch.TxHash] = waitOrderMatch
	}
}

func (e *EventNotifier) handleCreateOrder(ev notifier.Event, side string) {
	log.Debugf("=========> ")
	log.Debugf("  %#v", ev)
	id, creator, _, inputAmount, outputAmount := notifier.ParseOrderData(ev.Data)

	log.Debugf(" id:%v, creator:%v, inputAmount:%v, outputAmount:%v, side:%v", id, creator, inputAmount, outputAmount, side)
	if side == "buy" {
		event := common.EngineEvent{}
		newOrder := common.LimitOrder{
			Id:              fmt.Sprintf("%v", id),
			Qty:             outputAmount,
			Price:           inputAmount.Div(outputAmount),
			Side:            common.NEW_BUY_ORDER,
			ContractAddress: ev.Address,
			Creator:         creator,
			TimeStamp:       time.Now().Unix(),
			Filled:          decimal.RequireFromString("0"),
		}
		event.EventID = common.NEW_BUY_ORDER
		event.Data = newOrder.ToString()
		e.pushEngineEvent(event)
	} else {
		event := common.EngineEvent{}
		newOrder := common.LimitOrder{
			Id:              fmt.Sprintf("%v", id),
			Qty:             inputAmount,
			Price:           outputAmount.Div(inputAmount),
			Side:            common.NEW_SELL_ORDER,
			ContractAddress: ev.Address,
			Creator:         creator,
			TimeStamp:       time.Now().Unix(),
			Filled:          decimal.RequireFromString("0"),
		}
		event.EventID = common.NEW_SELL_ORDER
		event.Data = newOrder.ToString()
		e.pushEngineEvent(event)
	}
}

func (e *EventNotifier) handleMatch(ev notifier.Event) {
	matchOrdersExtEvent := notifier.ParseMatchOrdersExtEvent(ev.Topics)
	log.Debugf("caller:%v, epoch:%v, orderType:%v, orderId:%v, orderCreator:%v", matchOrdersExtEvent.Caller, matchOrdersExtEvent.Epoch, matchOrdersExtEvent.OrderType, matchOrdersExtEvent.OrderId, matchOrdersExtEvent.OrderCreator)
	waitBlockMatchEvent := e.waitTxHashList[ev.TxHash]

	waitBlockMatchEvent.ReceivedEvents = append(waitBlockMatchEvent.ReceivedEvents, matchOrdersExtEvent)
	log.Debugf("waitBlockMatchEvent hash:%v, ..%v events received ", waitBlockMatchEvent.TxHash, len(waitBlockMatchEvent.ReceivedEvents))
	if waitBlockMatchEvent.AllEventReceived() {
		log.Debugf("waitBlockMatchEvent hash:%v all events received", waitBlockMatchEvent.TxHash)
		// push engine market execution record
		engineEvent := new(common.EngineEvent)
		engineEvent.EventID = common.UPDATE_ORDER_EVENT
		engineEvent.Data = common.UpdateLimitOrder{
			OrderList:       waitBlockMatchEvent.GetMatchedOrderIds(),
			ContractAddress: ev.Address,
		}.ToString()
		e.pushEngineEvent(*engineEvent)
	}
}

func (e *EventNotifier) pushEngineEvent(engineEvent common.EngineEvent) {
	b, _ := json.Marshal(engineEvent)
	e.pushRedis(string(b))
}

func (e *EventNotifier) pushRedis(data string) {
	key := common.ENGINE_KEY
	_, err := e.rdb.LPush(context.Background(), key, data).Result()
	if err != nil {
		log.Error(err)
	}
}

func parseRedisMsg(str string) *common.NotifierEvent {
	ev := new(common.NotifierEvent)
	err := json.Unmarshal([]byte(str), ev)
	if err != nil {
		log.Error(err)
		return nil
	}
	return ev
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

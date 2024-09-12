package main

import (
	"engine/api"
	"engine/conf"
	"engine/launcher"

	// "engine/notifier/websocket"
	"engine/market"
	"engine/notifier/mpq"
	"io/ioutil"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	_ "engine/util/logger"

	"engine/model"
	"engine/util/mysql"

	log "github.com/sirupsen/logrus"
)

func initConf() {
	conf.LoadConfig()
}

func initDB() {
	mysql.InitMysql(
		conf.EngineConfig.Mysql.Host,
		conf.EngineConfig.Mysql.User,
		conf.EngineConfig.Mysql.Password,
		conf.EngineConfig.Mysql.Port,
		conf.EngineConfig.Mysql.DB,
	)
	sess := mysql.GetSession()
	defer sess.Close()

	err := sess.Sync2(model.Trade{}, model.MatchRecord{}, model.LimitOrder{})
	if err != nil {
		panic(err)
	}

	// db.EnterTest()
}

func initNotifier() {
	marketAddressList := make([]string, 0)
	for _, v := range conf.EngineConfig.Market {
		marketAddressList = append(marketAddressList, v.ExchangeAddress)
	}
	// notifier := websocket.NewEventNotifier(marketAddressList, conf.EngineConfig.EventNotifier.WebSocket)
	notifier := mpq.NewEventNotifier(marketAddressList, conf.EngineConfig.EventNotifier.Connection, conf.EngineConfig.EventNotifier.EventID)
	notifier.Start()
}

func initLauncher() {
	var gateway string
	var explorer string
	gateway = conf.EngineConfig.Launcher.GateWay
	explorer = conf.EngineConfig.Launcher.Explorer

	filePath := conf.EngineConfig.Launcher.WalletPem
	walletPem, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Error(err)
		return
	}
	l := launcher.NewLauncher(gateway, walletPem, explorer)
	l.Start()
}

func initEngine() {
	engine := market.NewEngine()
	engine.LoadMarket()
	engine.Start()
}

func main() {
	initConf()
	initDB()
	initNotifier()
	initLauncher()
	initEngine()

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	e.POST("/api/orderbook", api.OrderBook)
	e.POST("/api/openOrder", api.OpenOrder)
	e.POST("/api/trades", api.MyTrade)
	e.POST("/api/tradeHistory", api.MyTrade)

	listenPort := ":8000"
	log.Infof("start listen %v", listenPort)
	err := e.Start(listenPort)
	if err != nil {
		panic(err)
	}
}

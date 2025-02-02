// Generated by https://quicktype.io
package conf

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type engineConfig struct {
	Redis         string        `json:"redis"`
	Mysql         Mysql         `json:"mysql"`
	Market        []Market      `json:"market"`
	EventNotifier EventNotifier `json:"eventNotifier"`
	Launcher      Launcher      `json:"launcher"`
}

type EventNotifier struct {
	Connection string `json:"connection"`
	EventID    string `json:"eventId"`
	WebSocket  string `json:"websocket"`
}

type Launcher struct {
	GateWay   string `json:"gateway"`
	Explorer  string `json:"explorer"`
	WalletPem string `json:"walletPem"`
}

type Market struct {
	Market          string `json:"market"`
	ExchangeAddress string `json:"exchangeAddress"`
	BaseSymbol      string `json:"baseSymbol"`
	QuoteSymbol     string `json:"quoteSymbol"`
	BaseESDT        string `json:"baseESDT"`
	QuoteESDT       string `json:"quoteESDT"`
}

type Mysql struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	Port     int64  `json:"port"`
	DB       string `json:"db"`
}

var EngineConfig *engineConfig

func LoadConfig() {
	filePath := "engine.json"

	// Read the file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	EngineConfig = new(engineConfig)
	err = json.Unmarshal(content, EngineConfig)
	if err != nil {
		log.Fatal(err)
	}
}

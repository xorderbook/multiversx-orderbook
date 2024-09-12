package websocket

import (
	"encoding/json"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/multiversx/mx-chain-notifier-go/data"
	"github.com/multiversx/mx-chain-notifier-go/dispatcher"
	log "github.com/sirupsen/logrus"
)

type wsClient struct {
	wsConn    dispatcher.WSConnection
	mutWsConn sync.RWMutex
}

// NewWSClient creates a new websocket client
func NewWSClient(hostURL string) (*wsClient, error) {
	u := url.URL{
		Scheme: "ws",
		Host:   hostURL,
		Path:   wsPath,
	}

	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}
	log.Infof("Notifier websocket connected")
	return &wsClient{
		wsConn: ws,
	}, nil
}

// SendSubscribeMessage will send subscribe message
func (ws *wsClient) SendSubscribeMessage(subscribeEvent *data.SubscribeEvent) error {
	ws.mutWsConn.Lock()
	defer ws.mutWsConn.Unlock()

	m, err := json.Marshal(subscribeEvent)
	if err != nil {
		return err
	}

	return ws.wsConn.WriteMessage(websocket.BinaryMessage, m)
}

// ReadMessage will read the received message
func (ws *wsClient) ReadMessage() ([]byte, error) {
	ws.mutWsConn.Lock()
	defer ws.mutWsConn.Unlock()

	_, m, err := ws.wsConn.ReadMessage()
	if err != nil {
		return nil, err
	}

	return m, nil
}

// Close will close connection
func (ws *wsClient) Close() {
	ws.mutWsConn.Lock()
	defer ws.mutWsConn.Unlock()

	ws.wsConn.Close()
}

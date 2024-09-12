package launcher

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"engine/common"
	"engine/db"
	"fmt"
	"strconv"
	"strings"
	"time"

	util "engine/util/redis"

	"github.com/go-redis/redis/v8"
	"github.com/multiversx/mx-chain-crypto-go/signing"
	"github.com/multiversx/mx-chain-crypto-go/signing/ed25519"
	log "github.com/sirupsen/logrus"

	"github.com/multiversx/mx-sdk-go/blockchain"
	"github.com/multiversx/mx-sdk-go/blockchain/cryptoProvider"
	"github.com/multiversx/mx-sdk-go/builders"
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/interactors"
)

var (
	suite  = ed25519.NewEd25519()
	keyGen = signing.NewKeyGenerator(suite)
)

type launcher struct {
	ch          chan MatchOrderReq
	gateway     string
	walletPem   []byte
	explorer    string
	explorerURL string
	rdb         *redis.Client
}

var Launcher *launcher

func NewLauncher(gateway string, pem []byte, explorer string) *launcher {
	Launcher = new(launcher)
	Launcher.gateway = gateway
	Launcher.walletPem = pem
	Launcher.explorer = explorer
	Launcher.ch = make(chan MatchOrderReq, 100)
	Launcher.explorerURL = explorer
	Launcher.rdb = util.NewRedis("127.0.0.1:6379")

	return Launcher
}

func (l *launcher) Start() {
	go l.run()
}

func (l *launcher) run() {
	for {
		select {
		case txReq := <-l.ch:
			l.sendTransaction(txReq)
		}
	}
}

func (l *launcher) sendTransaction(txReq MatchOrderReq) {
	args := blockchain.ArgsProxy{
		ProxyURL:            l.gateway,
		Client:              nil,
		SameScState:         false,
		ShouldBeSynced:      false,
		FinalityCheck:       false,
		CacheExpirationTime: time.Minute,
		EntityType:          core.Proxy,
	}
	ep, err := blockchain.NewProxy(args)
	if err != nil {
		log.Error("error creating proxy", "error", err)
		return
	}

	w := interactors.NewWallet()

	privateKey, err := w.LoadPrivateKeyFromPemData(l.walletPem)
	if err != nil {
		log.Error("unable to load alice.pem", "error", err)
		return
	}
	// Generate address from private key
	address, err := w.GetAddressFromPrivateKey(privateKey)
	if err != nil {
		log.Error("unable to load the address from the private key", "error", err)
		return
	}

	// netConfigs can be used multiple times (for example when sending multiple transactions) as to improve the
	// responsiveness of the system
	netConfigs, err := ep.GetNetworkConfig(context.Background())
	if err != nil {
		log.Error("unable to get the network configs", "error", err)
		return
	}

	tx, _, err := ep.GetDefaultTransactionArguments(context.Background(), address, netConfigs)
	if err != nil {
		log.Error("unable to prepare the transaction creation arguments", "error", err)
		return
	}

	// set tx data
	tx.Receiver = txReq.contract // exchange contract address
	tx.Data = l.matchOrdersExtCallData(txReq.takerOrderId, txReq.makerOrderIdList)
	tx.GasLimit = 600000000
	tx.GasPrice = 10000000000
	tx.Value = "0"

	holder, _ := cryptoProvider.NewCryptoComponentsHolder(keyGen, privateKey)
	txBuilder, err := builders.NewTxBuilder(cryptoProvider.NewSigner())
	if err != nil {
		log.Error("unable to prepare the transaction creation arguments", "error", err)
		return
	}

	ti, err := interactors.NewTransactionInteractor(ep, txBuilder)
	if err != nil {
		log.Error("error creating transaction interactor", "error", err)
		return
	}

	err = ti.ApplySignature(holder, &tx)
	if err != nil {
		log.Error("error signing transaction", "error", err)
		return
	}

	ti.AddTransaction(&tx)
	hashes, err := ti.SendTransactionsAsBunch(context.Background(), 100)
	if err != nil {
		log.Error("error sending transaction", "error", err)
		return
	}

	explorerUrl := fmt.Sprintf("%s%v", l.explorerURL, hashes[0])
	db.UpdateMatchTxHash(txReq.matchId, hashes[0], explorerUrl)
	log.Info("transactions sent ", "hash: ", explorerUrl)
	waitOrderMatch := common.WaitOrderMatch{
		NumOfEventsNeedsToWait: len(txReq.makerOrderIdList) + 1,
		TxHash:                 hashes[0],
	}
	notifierEvent := common.NotifierEvent{
		EventID: common.WAIT_MATCH_EVENT,
		Data:    waitOrderMatch.ToString(),
	}
	l.pushNotifierEvent(notifierEvent)
}

func (l *launcher) matchOrdersExtCallData(takerOrderId uint64, makerOrderIds []uint64) []byte {
	log.Info("takerOrderId ", takerOrderId)
	log.Info("makerOrderIds ", makerOrderIds)
	_args := make([]string, 0)
	bytes064 := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes064, takerOrderId)

	compact := make([]byte, 0)
	for _, v := range bytes064 {
		if v == 0 {
			continue
		}
		compact = append(compact, v)
	}

	_args = append(_args, hex.EncodeToString(compact))
	for _, elem := range makerOrderIds {
		bytes164 := make([]byte, 8)
		binary.BigEndian.PutUint64(bytes164, elem)
		_args = append(_args, hex.EncodeToString(bytes164))
	}
	dataField := "matchOrdersExt" + "@" + _args[0] + "@" + strings.Join(_args[1:], "")
	log.Info("dataField ", dataField)

	return []byte(dataField)
}

func (l *launcher) pushNotifierEvent(waitEvent common.NotifierEvent) {
	data, _ := json.Marshal(waitEvent)
	key := common.NOTIFIER_KEY
	_, err := l.rdb.LPush(context.Background(), key, data).Result()
	if err != nil {
		log.Error(err)
	}
}

func SubmitBlockchainOrderMatch(matchId int64, receiver string, takerOrder string, makeOrderList []string) {
	makerOrderIdList := make([]uint64, 0)

	takerOrderId, _ := strconv.ParseInt(takerOrder, 10, 64)
	for _, makerOrder := range makeOrderList {
		makerOrderId, _ := strconv.ParseInt(makerOrder, 10, 64)
		makerOrderIdList = append(makerOrderIdList, uint64(makerOrderId))
	}

	Launcher.ch <- MatchOrderReq{matchId: matchId, contract: receiver, takerOrderId: uint64(takerOrderId), makerOrderIdList: makerOrderIdList}
}

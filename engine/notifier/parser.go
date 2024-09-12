package notifier

import (
	"encoding/base64"
	"encoding/hex"
	"engine/common"
	"engine/util/parser"

	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

func ParseOrderData(base64Data *string) (uint64, string, string, decimal.Decimal, decimal.Decimal) {
	data, _ := base64.StdEncoding.DecodeString(*base64Data)
	idx := 0
	ok, allOk := true, true
	id, idx, ok := parser.ParseUint64(data, idx)
	allOk = allOk && ok
	creator, idx, ok := parser.ParsePubkey(data, idx)
	allOk = allOk && ok
	matchProvider, idx, ok := parser.ParsePubkey(data, idx)
	allOk = allOk && ok
	_inputAmount, idx, ok := parser.ParseBigInt(data, idx)
	allOk = allOk && ok
	_outputAmount, idx, ok := parser.ParseBigInt(data, idx)
	allOk = allOk && ok

	inputAmount := decimal.RequireFromString(_inputAmount.String())
	outputAmount := decimal.RequireFromString(_outputAmount.String())

	return id, hex.EncodeToString(creator), hex.EncodeToString(matchProvider), inputAmount, outputAmount
}

func ParseOrderDataByte(data []byte) (uint64, string, string, decimal.Decimal, decimal.Decimal) {
	idx := 0
	ok, allOk := true, true
	id, idx, ok := parser.ParseUint64(data, idx)
	allOk = allOk && ok
	creator, idx, ok := parser.ParsePubkey(data, idx)
	allOk = allOk && ok
	matchProvider, idx, ok := parser.ParsePubkey(data, idx)
	allOk = allOk && ok
	_inputAmount, idx, ok := parser.ParseBigInt(data, idx)
	allOk = allOk && ok
	_outputAmount, idx, ok := parser.ParseBigInt(data, idx)
	allOk = allOk && ok

	inputAmount := decimal.RequireFromString(_inputAmount.String())
	outputAmount := decimal.RequireFromString(_outputAmount.String())

	return id, hex.EncodeToString(creator), hex.EncodeToString(matchProvider), inputAmount, outputAmount
}

// #[event("match_order")]
// fn match_order_event(
//
//	&self,
//	#[indexed] caller: &ManagedAddress,
//	#[indexed] epoch: u64,
//	#[indexed] order_type: OrderType,
//	#[indexed] order_id: u64,
//	#[indexed] order_creator: ManagedAddress,
//
// );
func ParseMatchOrdersExtEvent(topics []string) common.MatchOrdersExt {
	// for _, v := range topics {
	// 	log.Debug(v)
	// }
	// order_id, _ := ParseAmount
	// order_creator, _ := base64.StdEncoding.DecodeString(topics[5])
	return common.MatchOrdersExt{
		Caller:       "",
		Epoch:        0,
		OrderType:    0,
		OrderId:      ParseU64(topics[4]),
		OrderCreator: "",
	}
}

func ParseMatchOrdersExtEventByte(topics [][]byte) common.MatchOrdersExt {
	for _, v := range topics {
		log.Debug(v)
	}
	return common.MatchOrdersExt{
		Caller:       "",
		Epoch:        0,
		OrderType:    0,
		OrderId:      ParseU64Byte(topics[4]),
		OrderCreator: hex.EncodeToString(topics[5]),
	}
}

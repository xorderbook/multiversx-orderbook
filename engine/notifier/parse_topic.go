package notifier

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func ParseU64(str string) uint64 {
	decodedBytes, _ := base64.StdEncoding.DecodeString(str)
	// Encode to hexadecimal
	hexStr := hex.EncodeToString(decodedBytes)
	decimal, _ := hexToDecimal(hexStr)
	return toUint64(decimal)
}

func ParseU64Byte(bytes []byte) uint64 {
	// Encode to hexadecimal
	hexStr := hex.EncodeToString(bytes)
	decimal, _ := hexToDecimal(hexStr)
	return toUint64(decimal)
}

func hexToDecimal(hex string) (string, error) {
	dec := new(big.Int)
	_, success := dec.SetString(hex, 16)
	if !success {
		log.Error(fmt.Errorf("Invalid hexadecimal string: %s", hex))
		return "", fmt.Errorf("Invalid hexadecimal string: %s", hex)
	}

	return dec.String(), nil
}

func toUint64(hex string) uint64 {
	i, err := strconv.ParseInt(hex, 10, 64)
	if err != nil {
		log.Error(err)
		return 0
	}

	return uint64(i)
}

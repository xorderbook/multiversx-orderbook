package convert

import (
	"encoding/hex"

	"github.com/multiversx/mx-sdk-go/data"
)

func HexAddress(addressAsBech32 string) string {
	address, err := data.NewAddressFromBech32String(addressAsBech32)
	if err != nil {
		return ""
	}
	addressHex := hex.EncodeToString(address.AddressBytes())
	return addressHex
}

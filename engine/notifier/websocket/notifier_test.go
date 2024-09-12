package websocket

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
)

func TestParse(t *testing.T) {
	inputAmount, _ := decimal.NewFromString("1000000000000000000")
	outputAmount, _ := decimal.NewFromString("1000000000000000000")
	price := inputAmount.Div(outputAmount)
	fmt.Println(price)
}

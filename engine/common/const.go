package common

const (
	NEW_BUY_ORDER = iota
	NEW_SELL_ORDER
	ORDER_MATCH
	WAIT_MATCH_EVENT
	UPDATE_ORDER_EVENT
)

const (
	ENGINE_KEY   = "engine"
	NOTIFIER_KEY = "notifier"
)

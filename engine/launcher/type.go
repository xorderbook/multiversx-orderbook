package launcher

type MatchOrderReq struct {
	matchId          int64
	contract         string
	takerOrderId     uint64
	makerOrderIdList []uint64
}

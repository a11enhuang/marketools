package marketdata

import "com.reopenai/marketool/dal"

type apiResult struct {
	Code   int      `json:"code"`
	Msg    string   `json:"msg"`
	Data   RankData `json:"data"`
	Offset int      `json:"offset"`
}

type RankData struct {
	RankList []*dal.StockPrice `json:"rank_list"`
}
